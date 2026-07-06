<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import { api } from '$lib/api';
	import { user, clearAuth } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import ConfirmModal from '$lib/components/ConfirmModal.svelte';

	import { 
		Home, BookOpen, User, BarChart2, Users, GraduationCap, School, 
		Calendar, FileText, Book, Puzzle, Zap, CheckCircle, ClipboardList, Trophy,
		LogOut, Menu, X, ArrowUp, Settings, Terminal
	} from 'lucide-svelte';

	let { children } = $props();

	interface NavLink { href: string; label: string; icon: any; }

	const userLinks: NavLink[] = [
		{ href: '/praktikum/dashboard', label: 'Dashboard', icon: Home },
		{ href: '/praktikum/sesi', label: 'Daftar Sesi', icon: BookOpen },
		{ href: '/praktikum/profil', label: 'Profil', icon: User }
	];
	const adminLinks: NavLink[] = [
		{ href: '/praktikum/admin', label: 'Dashboard', icon: BarChart2 },
		{ href: '/praktikum/admin/users', label: 'Data User', icon: Users },
		{ href: '/praktikum/admin/asisten', label: 'Asisten', icon: GraduationCap },
		{ href: '/praktikum/admin/kelas', label: 'Kelas', icon: School },
		{ href: '/praktikum/admin/jadwal', label: 'Jadwal', icon: Calendar },
		{ href: '/praktikum/admin/pedoman', label: 'Pedoman', icon: FileText },
		{ href: '/praktikum/admin/modul', label: 'Modul', icon: Book },
		{ href: '/praktikum/admin/sesi', label: 'Sesi & Soal', icon: Puzzle },
		{ href: '/praktikum/admin/aktivasi', label: 'Aktivasi Sesi', icon: Zap },
		{ href: '/praktikum/admin/penilaian', label: 'Penilaian', icon: CheckCircle },
		{ href: '/praktikum/admin/rekap-jawaban', label: 'Rekap Jawaban', icon: ClipboardList },
		{ href: '/praktikum/admin/rekap-nilai', label: 'Rekap Nilai', icon: Trophy },
		{ href: '/praktikum/admin/log', label: 'Log Aktivitas', icon: Terminal },
		{ href: '/praktikum/admin/pengaturan', label: 'Pengaturan Lobby', icon: Settings }
	];

	let links = $derived($user?.role === 'admin' ? adminLinks : userLinks);
	let open = $state(false);

	// Saat user di halaman pengerjaan soal, sembunyikan sidebar agar fokus.
	let inSession = $derived(
		/\/praktikum\/sesi\/\d+\/(pretest|posttest|keterampilan|ujian)/.test($page.url.pathname)
	);

	function active(href: string): boolean {
		if (href === '/praktikum/admin' || href === '/praktikum/dashboard') {
			return $page.url.pathname === href;
		}
		return $page.url.pathname === href || $page.url.pathname.startsWith(href + '/');
	}

	async function logout() {
		try { await api.post('/api/auth/logout'); } catch { /* ignore */ }
		clearAuth();
		goto('/praktikum/login');
	}
</script>

<div class="flex min-h-screen bg-slate-50/50 selection:bg-maroon-100 selection:text-primary">
	<!-- Sidebar (hidden saat sesi pengerjaan) -->
	{#if !inSession}
	<aside class="fixed inset-y-0 left-0 z-35 w-64 -translate-x-full bg-white border-r border-slate-200 text-slate-800 transition-all duration-300 md:static md:translate-x-0 {open ? 'translate-x-0' : ''} shadow-[4px_0_24px_rgba(0,0,0,0.02)] flex flex-col shrink-0">
		<!-- Logo Area -->
		<div class="h-16 flex items-center px-6 border-b border-slate-100 shrink-0 justify-between">
			<a href="/praktikum/dashboard" class="flex items-center">
				<img src="/logo_new.png" alt="Logo Lab AP" class="h-9 object-contain" />
			</a>
			<!-- Close button for mobile -->
			<button class="text-slate-400 hover:text-slate-650 md:hidden" onclick={() => (open = false)} aria-label="Tutup Menu">
				<X size={18} />
			</button>
		</div>

		<!-- Scrollable Navigation -->
		<nav class="flex-1 overflow-y-auto py-5 px-3 space-y-1 custom-scrollbar">
			{#each links as l}
				{@const Icon = l.icon}
				{@const isActive = active(l.href)}
				<a
					href={l.href}
					onclick={() => (open = false)}
					class="flex items-center px-4 py-2.5 rounded-xl group relative transition-all duration-300 overflow-hidden text-sm font-semibold
					{isActive ? 'bg-maroon-bg text-primary border border-primary/10 shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-50 hover:text-slate-800'}"
				>
					{#if isActive}
						<div class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-6 bg-primary rounded-r-full"></div>
					{/if}
					<Icon size={16} class="mr-3 transition-colors shrink-0 {isActive ? 'text-primary' : 'text-slate-400 group-hover:text-primary'}" />
					<span class="transition-transform duration-300 {isActive ? '' : 'group-hover:translate-x-1'}">{l.label}</span>
				</a>
			{/each}
		</nav>

	</aside>

	{#if open}
		<button class="fixed inset-0 z-20 bg-slate-900/50 backdrop-blur-sm md:hidden" onclick={() => (open = false)} aria-label="Tutup"></button>
	{/if}
	{/if}

	<!-- Main -->
	<div class="flex min-w-0 flex-1 flex-col overflow-hidden h-screen">
		<header class="bg-white/95 backdrop-blur-md border-b border-slate-200 h-16 flex items-center justify-between px-6 shrink-0 z-10 shadow-sm">
			{#if !inSession}
			<button class="md:hidden text-slate-500 hover:text-primary transition-colors focus:outline-none p-1" onclick={() => (open = true)} aria-label="Menu">
				<Menu size={20} />
			</button>
			{/if}
			
			<div class="ml-auto flex items-center gap-4 shrink-0">
				<!-- Profile Capsule -->
				<a href="/praktikum/profil" class="flex items-center gap-3 bg-white pl-1.5 pr-4 py-1.5 rounded-full border border-slate-200 shadow-sm hover:border-primary/20 hover:shadow-md transition-all cursor-pointer group">
					<img src={`https://ui-avatars.com/api/?name=${encodeURIComponent($user?.nama ?? '')}&background=FBE5E9&color=8A1538&bold=true`} alt="Avatar" class="w-8 h-8 rounded-full border border-white shadow-sm group-hover:scale-105 transition-transform shrink-0">
					<div class="text-right hidden sm:block">
						<p class="text-sm font-extrabold text-slate-900 leading-none group-hover:text-primary transition-colors">{$user?.nama}</p>
						<p class="text-[10px] font-black text-primary uppercase tracking-widest mt-1">{$user?.role}</p>
					</div>
				</a>

				<!-- Divider -->
				<div class="h-6 w-px bg-slate-200 hidden sm:block"></div>

				<!-- Logout Button (Prominent Red Pill) -->
				<button class="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-extrabold text-red-650 hover:text-white bg-red-50 hover:bg-red-600 border border-red-200 hover:border-transparent rounded-full shadow-sm transition-all active:scale-95 focus:outline-none shrink-0" onclick={logout} title="Logout">
					<LogOut size={13} /> Logout
				</button>
			</div>
		</header>
		<main class="flex-1 overflow-x-hidden overflow-y-auto bg-slate-50/30 p-4 md:p-6 lg:p-8">
			{#key $page.url.pathname}
				<div in:fly={{ y: 8, duration: 240, easing: cubicOut }}>
					{@render children()}
				</div>
			{/key}
		</main>
	</div>
	<ConfirmModal />
</div>
