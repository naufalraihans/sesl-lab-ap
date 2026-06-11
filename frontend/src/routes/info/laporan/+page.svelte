<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { Download } from 'lucide-svelte';

	interface Pedoman { id: number; nama_dokumen: string; file_url: string; }
	let items = $state<Pedoman[]>([]);
	let loading = $state(true);
	let err = $state('');

	onMount(async () => {
		try {
			items = (await api.get<Pedoman[]>('/api/info/laporan')) ?? [];
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});
</script>

<h1 class="mb-4 text-2xl">Pedoman Laporan</h1>

{#if loading}
	<p class="text-ink-caption">Memuat…</p>
{:else if err}
	<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
{:else if items.length === 0}
	<p class="text-ink-caption">Belum ada dokumen pedoman.</p>
{:else}
	<div class="grid gap-3 sm:grid-cols-2">
		{#each items as it}
			<div class="card flex items-center justify-between">
				<span class="font-medium">{it.nama_dokumen}</span>
				<a href={it.file_url} target="_blank" rel="noopener" class="btn-primary"><Download size={16} /> Download</a>
			</div>
		{/each}
	</div>
{/if}
