<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { User } from '$lib/types';

	let asisten = $state<User[]>([]);
	let loading = $state(true);
	let err = $state('');

	onMount(async () => {
		try {
			asisten = (await api.get<User[]>('/api/info/asisten')) ?? [];
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});
</script>

<h1 class="mb-4 text-2xl">Daftar Asisten Lab</h1>

{#if loading}
	<p class="text-ink-caption">Memuat…</p>
{:else if err}
	<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
{:else if asisten.length === 0}
	<p class="text-ink-caption">Belum ada data asisten.</p>
{:else}
	<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
		{#each asisten as a}
			<div class="card">
				<div class="flex items-center gap-4">
					{#if a.foto_url}
						<img src={a.foto_url} alt={a.nama} class="h-16 w-16 rounded-full object-cover" />
					{:else}
						<div class="grid h-16 w-16 place-items-center rounded-full bg-surface-soft text-xl font-bold text-primary">
							{a.nama?.charAt(0)}
						</div>
					{/if}
					<div>
						<h3 class="text-lg">{a.nama}</h3>
						<p class="text-sm text-ink-caption">{a.nim}</p>
					</div>
				</div>
				<div class="mt-3 flex flex-wrap gap-2 text-sm">
					{#if a.nomor_hp}
						<a href={`https://wa.me/${a.nomor_hp.replace(/^0/, '62')}`} target="_blank" rel="noopener" class="badge bg-state-success-bg text-state-success">WhatsApp</a>
					{/if}
					{#if a.medsos_link}
						<a href={a.medsos_link} target="_blank" rel="noopener" class="badge bg-state-info-bg text-state-info">Media Sosial</a>
					{/if}
				</div>
			</div>
		{/each}
	</div>
{/if}
