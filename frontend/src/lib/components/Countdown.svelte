<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
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

	let warn = $derived(remaining <= 60);
</script>

<span
	class="badge font-mono text-base {warn ? 'bg-state-error-bg text-state-error' : 'bg-state-info-bg text-state-info'}"
>
	⏱ {fmtCountdown(remaining)}
</span>
