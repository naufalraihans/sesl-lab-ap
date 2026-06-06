<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	onMount(() => goto('/info'));
</script>

<p class="p-8 text-ink-caption">Mengalihkan ke halaman utama…</p>
