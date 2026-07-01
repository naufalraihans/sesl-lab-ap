<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import { page } from '$app/stores';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	let { children } = $props();
</script>

<div class="flex min-h-screen flex-col">
	<Navbar />
	<main class="mx-auto w-full max-w-6xl flex-1 px-4 py-8">
		{#key $page.url.pathname}
			<div in:fly={{ y: 12, duration: 280, easing: cubicOut }}>
				{@render children()}
			</div>
		{/key}
	</main>
	<footer class="relative z-10 border-t border-gray-200 bg-surface-muted py-6 text-center text-sm text-ink-caption">
		© {new Date().getFullYear()} Laboratorium Algoritma &amp; Pemrograman
	</footer>
</div>
