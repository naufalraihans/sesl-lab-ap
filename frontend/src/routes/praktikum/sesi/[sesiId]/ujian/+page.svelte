<script lang="ts">
	import { page } from '$app/stores';
	import Ruang from '$lib/components/Ruang.svelte';
	let aktivasi = $derived(Number($page.url.searchParams.get('aktivasi')));
	let course = $derived(Number($page.url.searchParams.get('course')));
</script>

<a href={`/praktikum/sesi/${$page.params.sesiId}`} class="text-sm">← Kembali ke sesi</a>
<div class="mt-3">
	{#if aktivasi && course}
		<Ruang aktivasiSesiId={aktivasi} courseId={course} />
	{:else}
		<p class="text-state-error">Parameter tidak lengkap.</p>
	{/if}
</div>
