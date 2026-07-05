<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Timer } from 'lucide-svelte';
	import { fmtCountdown, secondsUntil } from '$lib/utils';

	let { deadline, onExpire }: { deadline?: string | null; onExpire?: () => void } = $props();

	let remaining = $state(secondsUntil(deadline));
	let timer: ReturnType<typeof setInterval>;
	let fired = false;

	onMount(() => {
		timer = setInterval(() => {
			remaining = secondsUntil(deadline);
			if (remaining <= 0 && !fired) {
				fired = true;
				onExpire?.();
			}
		}, 1000);
	});

	onDestroy(() => clearInterval(timer));

	// Green (> 15 minutes), Yellow/Orange (5-15 minutes), Red (<= 5 minutes)
	let colorClass = $derived(
		remaining > 900
			? 'bg-emerald-50 text-emerald-600 border border-emerald-200'
			: remaining > 300
			? 'bg-amber-50 text-amber-600 border border-amber-200'
			: 'bg-red-50 text-red-600 border border-red-300 animate-tick shadow-lg shadow-red-500/10'
	);
</script>

<span class="inline-flex items-center gap-2 font-mono text-lg md:text-xl px-4 py-2.5 rounded-2xl transition-all duration-300 shadow-sm font-bold {colorClass}">
	<Timer size={18} class="flex-shrink-0" /> 
	<span>{fmtCountdown(remaining)}</span>
</span>

<style>
	@keyframes tick-alert {
		0%, 100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.06);
			background-color: #fef2f2; /* slightly lighter red background on pulse */
		}
	}
	
	:global(.animate-tick) {
		animation: tick-alert 0.8s infinite ease-in-out;
		border-color: #fca5a5 !important;
	}
</style>
