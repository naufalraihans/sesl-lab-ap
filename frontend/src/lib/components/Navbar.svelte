<script lang="ts">
	import { page } from '$app/stores';

	let open = $state(false);

	const links = [
		{ href: '/info', label: 'Beranda' },
		{ href: '/info/jadwal', label: 'Jadwal' },
		{ href: '/info/asisten', label: 'Asisten Lab' },
		{ href: '/info/laporan', label: 'Pedoman Laporan' },
		{ href: '/info/modul', label: 'Modul' }
	];

	function active(href: string): boolean {
		return $page.url.pathname === href;
	}
</script>

<nav class="bg-navbar text-white">
	<div class="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
		<a href="/info" class="flex items-center gap-2 text-white hover:text-white">
			<img src="/logoLab.webp" alt="Logo Lab AP" class="h-9 w-9 object-contain" />
			<span class="font-semibold">Lab Algoritma &amp; Pemrograman</span>
		</a>

		<button class="md:hidden" onclick={() => (open = !open)} aria-label="Menu">
			<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
			</svg>
		</button>

		<div class="hidden items-center gap-1 md:flex">
			{#each links as l}
				<a
					href={l.href}
					class="rounded-lg px-3 py-2 text-sm text-white hover:bg-primary-hover {active(l.href) ? 'bg-primary' : ''}"
				>{l.label}</a>
			{/each}
			<a href="/praktikum/login" class="ml-2 rounded-lg bg-white px-3 py-2 text-sm font-medium text-primary hover:bg-surface-soft">Login Praktikum</a>
		</div>
	</div>

	{#if open}
		<div class="border-t border-white/20 px-4 pb-3 md:hidden">
			{#each links as l}
				<a href={l.href} class="block rounded-lg px-3 py-2 text-sm text-white hover:bg-primary-hover" onclick={() => (open = false)}>{l.label}</a>
			{/each}
			<a href="/praktikum/login" class="mt-1 block rounded-lg bg-white px-3 py-2 text-sm font-medium text-primary">Login Praktikum</a>
		</div>
	{/if}
</nav>
