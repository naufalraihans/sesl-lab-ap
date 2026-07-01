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

<nav class="nav-glass sticky top-0 z-40 text-ink-body">
	<div class="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
		<a href="/info" class="flex items-center gap-3 text-ink-heading hover:text-ink-heading">
			<span class="glass-badge h-10 w-10">
				<img src="/logoLab.webp" alt="Logo Lab AP" class="h-6 w-6 object-contain" />
			</span>
			<span class="font-semibold tracking-tight">Lab Algoritma &amp; Pemrograman</span>
		</a>

		<button class="md:hidden text-ink-heading" onclick={() => (open = !open)} aria-label="Menu">
			<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
			</svg>
		</button>

		<div class="hidden items-center gap-1 md:flex">
			{#each links as l}
				<a
					href={l.href}
					class="rounded-xl px-3.5 py-2 text-sm font-medium transition-colors {active(l.href) ? 'bg-primary text-white shadow-sm' : 'text-ink-body hover:bg-primary/10'}"
				>{l.label}</a>
			{/each}
			<a href="/praktikum/login" class="ml-2 rounded-xl border border-white/70 bg-white/60 px-4 py-2 text-sm font-semibold text-primary backdrop-blur transition-transform hover:-translate-y-0.5">Login Praktikum</a>
		</div>
	</div>

	{#if open}
		<div class="border-t border-white/40 px-4 pb-3 md:hidden">
			{#each links as l}
				<a href={l.href} class="block rounded-xl px-3 py-2 text-sm font-medium {active(l.href) ? 'bg-primary text-white' : 'text-ink-body hover:bg-primary/10'}" onclick={() => (open = false)}>{l.label}</a>
			{/each}
			<a href="/praktikum/login" class="mt-1 block rounded-xl border border-white/70 bg-white/60 px-3 py-2 text-sm font-semibold text-primary" onclick={() => (open = false)}>Login Praktikum</a>
		</div>
	{/if}
</nav>
