<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	let { children } = $props();
</script>

<div class="flex min-h-screen flex-col">
	<Navbar />
	<main class="mx-auto w-full max-w-6xl flex-1 px-4 py-8">
		{@render children()}
	</main>
	<footer class="border-t border-gray-200 bg-surface-muted py-6 text-center text-sm text-ink-caption">
		© {new Date().getFullYear()} Laboratorium Algoritma &amp; Pemrograman
	</footer>
</div>
