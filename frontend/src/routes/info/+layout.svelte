<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import { page } from '$app/stores';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import { onMount } from 'svelte';

	let { children } = $props();

	// Preloader state
	let preloading = $state(true);
	let progress = $state(0);
	let codeLines = [
		'#include <stdio.h>',
		'',
		'int main() {',
		'  printf("Login Success");',
		'  return 0;',
		'}',
	];
	let visibleLines = $state(0);

	onMount(() => {
		// Animate code lines appearing slowly (synchronized with loading progress)
		const lineInterval = setInterval(() => {
			visibleLines = Math.min(visibleLines + 1, codeLines.length);
		}, 420);

		// Progress bar animation to complete in ~2.5 to 3 seconds
		const progressInterval = setInterval(() => {
			progress = Math.min(progress + Math.random() * 3.5 + 2, 100);
			if (progress >= 100) {
				clearInterval(progressInterval);
				setTimeout(() => { preloading = false; }, 400);
			}
		}, 100);

		return () => {
			clearInterval(lineInterval);
			clearInterval(progressInterval);
		};
	});
</script>






<!-- Preloader -->
{#if preloading}
	<div class="fixed inset-0 z-[200] flex flex-col items-center justify-center"
		style="background: linear-gradient(135deg, #fff5f5 0%, #fffbfb 50%, #fff5f5 100%);">

		<!-- Background grid (very soft maroon/pink using logo RGB 138,21,56) -->
		<div class="absolute inset-0 opacity-[0.015]"
			style="background-size: 40px 40px; background-image: linear-gradient(to right, rgba(138,21,56,0.3) 1px, transparent 1px), linear-gradient(to bottom, rgba(138,21,56,0.3) 1px, transparent 1px);">
		</div>
		<!-- Glow orbs (extremely soft/faded rose using logo RGB 138,21,56) -->
		<div class="absolute top-1/4 left-1/4 w-64 h-64 rounded-full blur-[120px]" style="background: rgba(138,21,56,0.05);"></div>
		<div class="absolute bottom-1/4 right-1/4 w-64 h-64 rounded-full blur-[120px]" style="background: rgba(138,21,56,0.02);"></div>

		<div class="relative z-10 flex flex-col items-center gap-8 px-6 max-w-sm w-full">
			<!-- Logo -->
			<img src="/logo_new.png" alt="Lab AP" class="h-14 w-auto object-contain opacity-90" />

			<!-- Animated code window (clean light rose card, soft contrast) -->
			<div class="w-full bg-[#fffcfc] rounded-2xl border border-rose-100/80 shadow-xl overflow-hidden">
				<div class="bg-[#fff1f2] px-4 py-2.5 flex items-center gap-2 border-b border-rose-100/80">
					<div class="w-2.5 h-2.5 rounded-full bg-red-400/80"></div>
					<div class="w-2.5 h-2.5 rounded-full bg-amber-400/80"></div>
					<div class="w-2.5 h-2.5 rounded-full bg-green-400/80"></div>
					<span class="ml-2 text-[10px] font-mono text-[#8A1538]/70 font-extrabold">main.c</span>
				</div>
				<div class="p-4 font-mono text-[11px] leading-relaxed min-h-[120px] text-left">
					{#each codeLines.slice(0, visibleLines) as line, i}
						<div class="flex gap-3 items-start">
							<span class="text-rose-300 w-4 shrink-0 select-none text-right">{i + 1}</span>
							<span class={line.startsWith('#') ? 'text-[#8A1538] font-bold' : line.includes('int') || line.includes('return') ? 'text-purple-600 font-bold' : line.includes('"') ? 'text-amber-700 font-medium' : line.includes('printf') ? 'text-[#8A1538]' : 'text-slate-600'}>
								{line || '\u00A0'}
							</span>
						</div>
					{/each}
					<!-- blinking cursor -->
					{#if visibleLines <= codeLines.length}
						<div class="flex gap-3 items-start">
							<span class="text-rose-300 w-4 shrink-0 select-none">&nbsp;</span>
							<span class="inline-block w-2 h-4 bg-[#8A1538] animate-pulse rounded-sm"></span>
						</div>
					{/if}
				</div>
			</div>

			<!-- Progress bar -->
			<div class="w-full space-y-2">
				<div class="flex justify-between items-center">
					<span class="text-xs font-bold text-slate-500">Memuat sistem…</span>
					<span class="text-xs font-black text-[#8A1538]">{Math.round(progress)}%</span>
				</div>
				<div class="h-1.5 w-full bg-rose-100/60 rounded-full overflow-hidden border border-rose-200/30">
					<div class="h-full rounded-full transition-all duration-100"
						style="width: {progress}%; background: linear-gradient(90deg, #ffccd5, #8A1538);">
					</div>
				</div>
			</div>

			<p class="text-xs font-bold text-slate-400 tracking-widest uppercase">Laboratorium Algoritma &amp; Pemrograman</p>
		</div>
	</div>
{/if}

<div class="flex min-h-screen flex-col" class:opacity-0={preloading} style="transition: opacity 0.3s ease;">
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
