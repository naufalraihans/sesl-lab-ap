<script lang="ts">
	import { page } from '$app/stores';
	import Ruang from '$lib/components/Ruang.svelte';
	import { ArrowLeft } from 'lucide-svelte';
	let aktivasi = $derived(Number($page.url.searchParams.get('aktivasi')));
	let course = $derived(Number($page.url.searchParams.get('course')));
</script>

<div class="mb-4">
	<a href={`/praktikum/sesi/${$page.params.sesiId}`} class="inline-flex items-center gap-1.5 text-xs font-bold text-primary border border-primary/20 bg-primary/5 hover:bg-primary hover:text-white px-3 py-1.5 rounded-xl transition-all shadow-sm">
		<ArrowLeft size={13} />
		Kembali ke Sesi
	</a>
</div>
<div class="mt-3">
	{#if aktivasi && course}
		<Ruang aktivasiSesiId={aktivasi} courseId={course} />
	{:else}
		<p class="text-state-error">Parameter tidak lengkap.</p>
	{/if}
</div>
