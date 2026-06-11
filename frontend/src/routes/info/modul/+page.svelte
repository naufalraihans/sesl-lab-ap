<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { Download } from 'lucide-svelte';

	let fileUrl = $state('');
	let loading = $state(true);

	onMount(async () => {
		try {
			const res = await api.get<{ file_url: string }>('/api/info/modul');
			fileUrl = res?.file_url ?? '';
		} finally {
			loading = false;
		}
	});
</script>

<h1 class="mb-4 text-2xl">Modul Praktikum</h1>

{#if loading}
	<p class="text-ink-caption">Memuat…</p>
{:else if fileUrl}
	<div class="card max-w-md">
		<p class="mb-3 text-ink-body">Modul praktikum tersedia dalam format PDF.</p>
		<a href={fileUrl} target="_blank" rel="noopener" class="btn-primary"><Download size={16} /> Download Modul (PDF)</a>
	</div>
{:else}
	<p class="text-ink-caption">Modul belum diunggah oleh admin.</p>
{/if}
