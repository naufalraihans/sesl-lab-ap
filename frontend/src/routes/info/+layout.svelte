<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import { page } from '$app/stores';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	let { children } = $props();
</script>

<!-- Latar aurora glassmorphism -->
<div class="aurora" aria-hidden="true">
	<div class="aurora-blob aurora-blob-1"></div>
	<div class="aurora-blob aurora-blob-2"></div>
	<div class="aurora-blob aurora-blob-3"></div>
</div>

<div class="relative z-10 flex min-h-screen flex-col">
	<Navbar />
	<main class="mx-auto w-full max-w-6xl flex-1 px-4 py-9">
		{#key $page.url.pathname}
			<div in:fly={{ y: 12, duration: 280, easing: cubicOut }}>
				{@render children()}
			</div>
		{/key}
	</main>
	<footer class="border-t border-white/40 py-6 text-center text-sm text-ink-caption">
		© {new Date().getFullYear()} Laboratorium Algoritma &amp; Pemrograman
	</footer>
</div>
