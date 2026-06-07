<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { user, clearAuth } from '$lib/stores/auth';

	let { children } = $props();

	interface NavLink { href: string; label: string; icon: string; }

	const userLinks: NavLink[] = [
		{ href: '/praktikum/dashboard', label: 'Dashboard', icon: '🏠' },
		{ href: '/praktikum/sesi', label: 'Daftar Sesi', icon: '📚' },
		{ href: '/praktikum/profil', label: 'Profil', icon: '👤' }
	];
	const adminLinks: NavLink[] = [
		{ href: '/praktikum/admin', label: 'Dashboard', icon: '📊' },
		{ href: '/praktikum/admin/users', label: 'Data User', icon: '👥' },
		{ href: '/praktikum/admin/asisten', label: 'Asisten', icon: '🧑‍🏫' },
		{ href: '/praktikum/admin/kelas', label: 'Kelas', icon: '🏫' },
		{ href: '/praktikum/admin/jadwal', label: 'Jadwal', icon: '🗓️' },
		{ href: '/praktikum/admin/pedoman', label: 'Pedoman', icon: '📄' },
		{ href: '/praktikum/admin/modul', label: 'Modul', icon: '📕' },
		{ href: '/praktikum/admin/sesi', label: 'Sesi & Soal', icon: '🧩' },
		{ href: '/praktikum/admin/aktivasi', label: 'Aktivasi Sesi', icon: '⚡' },
		{ href: '/praktikum/admin/penilaian', label: 'Penilaian', icon: '✅' },
		{ href: '/praktikum/admin/rekap-jawaban', label: 'Rekap Jawaban', icon: '📋' },
		{ href: '/praktikum/admin/rekap-nilai', label: 'Rekap Nilai', icon: '🏆' }
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
	<aside class="fixed inset-y-0 left-0 z-30 w-64 -translate-x-full bg-sidebar text-white transition md:static md:translate-x-0 {open ? 'translate-x-0' : ''}">
		<div class="flex items-center gap-2 px-5 py-4">
			<span class="grid h-8 w-8 place-items-center rounded-lg bg-primary font-bold">L</span>
			<span class="font-semibold">Lab AP</span>
		</div>
		<nav class="px-3">
			{#each links as l}
				<a
					href={l.href}
					onclick={() => (open = false)}
					class="mb-1 flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-white/90 hover:bg-primary-hover {active(l.href) ? 'bg-primary' : ''}"
				>
					<span>{l.icon}</span>{l.label}
				</a>
			{/each}
		</nav>
	</aside>

	{#if open}
		<button class="fixed inset-0 z-20 bg-black/40 md:hidden" onclick={() => (open = false)} aria-label="Tutup"></button>
	{/if}

	<!-- Main -->
	<div class="flex min-w-0 flex-1 flex-col">
		<header class="flex items-center justify-between border-b border-gray-200 bg-white px-4 py-3">
			<button class="md:hidden" onclick={() => (open = true)} aria-label="Menu">
				<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" /></svg>
			</button>
			<div class="ml-auto flex items-center gap-3">
				<span class="text-sm text-ink-body">{$user?.nama} <span class="text-ink-caption">({$user?.role})</span></span>
				<button class="btn-outline py-1.5" onclick={logout}>Logout</button>
			</div>
		</header>
		<main class="min-w-0 flex-1 p-4 md:p-6">
			{@render children()}
		</main>
	</div>
</div>
