<script lang="ts">
	import { onMount } from 'svelte';
	import { swrGet } from '$lib/cache';
	import { Lock } from 'lucide-svelte';
	import type { SesiUserItem } from '$lib/types';

	let sesi = $state<SesiUserItem[]>([]);
	let err = $state('');
	let loading = $state(true);

	onMount(() => {
		// SWR: kalau sudah pernah dibuka, tampil instan dari cache lalu disegarkan.
		swrGet<SesiUserItem[]>('/api/praktikum/sesi', (v) => {
			sesi = v ?? [];
			loading = false;
		}).catch((e) => {
			err = (e as Error).message;
			loading = false;
		});
	});
</script>

<h1 class="mb-4 text-2xl">Daftar Sesi Praktikum</h1>

{#if loading}
	<p class="text-ink-caption">Memuat…</p>
{:else if err}
	<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
{:else}
	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
		{#each sesi as s}
			<div class="card flex flex-col {s.aktif ? '' : 'opacity-70'}">
				<div class="flex items-center justify-between">
					<h3 class="text-lg">{s.judul}</h3>
					{#if s.aktif}
						<span class="badge bg-state-success-bg text-state-success">Aktif</span>
					{:else}
						<span class="badge inline-flex items-center gap-1 bg-gray-100 text-ink-caption"><Lock size={12} /> Terkunci</span>
					{/if}
				</div>
				<p class="mt-1 flex-1 text-sm text-ink-caption">{s.deskripsi}</p>
				{#if s.aktif}
					<a href={`/praktikum/sesi/${s.sesi_id}`} class="btn-primary mt-4">Buka Sesi</a>
				{:else}
					<button class="btn-outline mt-4 cursor-not-allowed" disabled>Belum Diaktifkan</button>
				{/if}
			</div>
		{/each}
	</div>
{/if}
