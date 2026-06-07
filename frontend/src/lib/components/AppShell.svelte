<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { user, clearAuth } from '$lib/stores/auth';

	import { 
		Home, BookOpen, User, BarChart2, Users, GraduationCap, School, 
		Calendar, FileText, Book, Puzzle, Zap, CheckCircle, ClipboardList, Trophy 
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
		{ href: '/praktikum/admin/rekap-nilai', label: 'Rekap Nilai', icon: Trophy }
	];

	let links = $derived($user?.role === 'admin' ? adminLinks : userLinks);
	let open = $state(false);

	function active(href: string): boolean {
		return $page.url.pathname === href || $page.url.pathname.startsWith(href + '/');
	}

	async function logout() {
		try { await api.post('/api/auth/logout'); } catch { /* ignore */ }
		clearAuth();
		goto('/praktikum/login');
	}
</script>

<div class="flex min-h-screen bg-surface-muted">
	<!-- Sidebar -->
	<aside class="fixed inset-y-0 left-0 z-30 w-64 -translate-x-full bg-sidebar/95 backdrop-blur-xl border-r border-white/10 text-white transition-all duration-300 md:static md:translate-x-0 {open ? 'translate-x-0' : ''} shadow-2xl md:shadow-none">
		<div class="flex items-center gap-3 px-6 py-5">
			<span class="grid h-10 w-10 place-items-center rounded-xl bg-white/10 shadow-inner font-bold text-lg border border-white/20">L</span>
			<span class="font-bold tracking-wide text-lg text-white/90">Lab AP</span>
		</div>
		<nav class="px-4 py-2 flex flex-col gap-1 overflow-y-auto max-h-[calc(100vh-5rem)]">
			{#each links as l}
				{@const Icon = l.icon}
				<a
					href={l.href}
					onclick={() => (open = false)}
					class="group flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition-all duration-300
					{active(l.href) ? 'bg-primary/80 text-white shadow-lg shadow-black/20' : 'text-white/70 hover:bg-white/10 hover:text-white'}"
				>
					<Icon size={18} class="transition-transform duration-300 group-hover:scale-110 {active(l.href) ? 'text-white' : 'text-white/60 group-hover:text-white'}" />
					{l.label}
				</a>
			{/each}
		</nav>
	</aside>

	{#if open}
		<button class="fixed inset-0 z-20 bg-black/40 md:hidden" onclick={() => (open = false)} aria-label="Tutup"></button>
	{/if}

	<!-- Main -->
	<div class="flex min-w-0 flex-1 flex-col">
		<header class="sticky top-0 z-20 flex items-center justify-between glass border-b px-6 py-4 transition-all">
			<button class="md:hidden rounded-lg p-2 hover:bg-surface-soft transition-colors" onclick={() => (open = true)} aria-label="Menu">
				<svg class="h-6 w-6 text-ink-body" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" /></svg>
			</button>
			<div class="ml-auto flex items-center gap-4">
				<div class="flex flex-col items-end">
					<span class="text-sm font-semibold text-ink-heading leading-tight">{$user?.nama}</span>
					<span class="text-xs font-medium text-primary bg-primary/10 px-2 py-0.5 rounded-full mt-1 uppercase tracking-wide">{$user?.role}</span>
				</div>
				<div class="h-8 w-px bg-gray-200"></div>
				<button class="btn-outline border-transparent hover:border-gray-200 px-3 py-1.5 text-sm font-semibold" onclick={logout}>Logout</button>
			</div>
		</header>
		<main class="min-w-0 flex-1 p-4 md:p-6">
			{@render children()}
		</main>
	</div>
</div>
