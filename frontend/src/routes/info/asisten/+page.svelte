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

<h1 class="mb-1 text-2xl font-bold text-ink-heading">Daftar Asisten Lab</h1>
<p class="mb-6 text-sm text-ink-caption">Hubungi asisten lab yang bertugas untuk pertanyaan seputar praktikum.</p>

{#if loading}
	<p class="text-ink-caption">Memuat…</p>
{:else if err}
	<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
{:else if asisten.length === 0}
	<p class="text-ink-caption">Belum ada data asisten.</p>
{:else}
	<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
		{#each asisten as a}
			<div class="card hover-card overflow-hidden p-0 text-center">
				<div class="h-24 bg-gradient-to-r from-primary to-primary-dark"></div>
				<div class="px-6 pb-6">
					<div class="-mt-16 flex justify-center">
						{#if a.foto_url}
							<img src={a.foto_url} alt={a.nama} class="h-36 w-36 rounded-full object-cover ring-4 ring-white shadow-lg" />
						{:else}
							<div class="grid h-36 w-36 place-items-center rounded-full bg-surface-soft text-5xl font-bold text-primary ring-4 ring-white shadow-lg">
								{a.nama?.charAt(0)}
							</div>
						{/if}
					</div>
					<h3 class="mt-4 text-xl font-bold text-ink-heading">{a.nama}</h3>
					<p class="mt-0.5 text-sm text-ink-caption">{a.nim}</p>
					<span class="badge mt-2 inline-block bg-primary/10 text-primary">Asisten Lab</span>
					<div class="mt-4 flex flex-wrap justify-center gap-2 text-sm">
						{#if a.nomor_hp}
							<a href={`https://wa.me/${a.nomor_hp.replace(/^0/, '62')}`} target="_blank" rel="noopener" class="badge bg-state-success-bg text-state-success">WhatsApp</a>
						{/if}
						{#if a.medsos_link}
							<a href={a.medsos_link} target="_blank" rel="noopener" class="badge bg-state-info-bg text-state-info">Media Sosial</a>
						{/if}
					</div>
				</div>
			</div>
		{/each}
	</div>
{/if}
