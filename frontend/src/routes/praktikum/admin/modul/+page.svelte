<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let fileUrl = $state('');
	let err = $state(''); let msg = $state('');
	let uploading = $state(false);

	async function load() {
		try {
			const konf = (await api.get<{ key: string; value: string }[]>('/api/admin/konfigurasi')) ?? [];
			fileUrl = konf.find((k) => k.key === 'modul_file_url')?.value ?? '';
		} catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	async function uploadModul(ev: Event) {
		const input = ev.target as HTMLInputElement;
		if (!input.files?.[0]) return;
		uploading = true; err = ''; msg = '';
		const fd = new FormData();
		fd.append('file', input.files[0]); fd.append('folder', 'modul');
		try {
			const res = await api.upload<{ url: string }>('/api/admin/upload', fd);
			await api.post('/api/admin/konfigurasi', { key: 'modul_file_url', value: res.url });
			fileUrl = res.url;
			msg = 'Modul berhasil diunggah.';
		} catch (e) { err = (e as Error).message; }
		finally { uploading = false; }
	}
</script>

<h1 class="mb-4 text-2xl">Manajemen Modul Praktikum</h1>
<p class="mb-4 text-sm text-ink-caption">Modul bersifat global (1 file PDF). File ini akan tampil di halaman /info/modul.</p>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="card max-w-lg">
	<h2 class="mb-3 text-lg">File Modul Saat Ini</h2>
	{#if fileUrl}
		<div class="mb-4 flex items-center gap-3 rounded-lg border border-state-success-bg bg-state-success-bg/40 p-3">
			<span class="text-state-success">📄</span>
			<a href={fileUrl} target="_blank" rel="noopener" class="flex-1 text-sm text-primary hover:underline">{fileUrl.split('/').pop()}</a>
			<a href={fileUrl} target="_blank" rel="noopener" class="btn-outline py-1.5 text-xs">Download</a>
		</div>
	{:else}
		<p class="mb-4 text-sm text-ink-caption">Belum ada modul yang diunggah.</p>
	{/if}

	<label class="label" for="modul">Upload / Ganti Modul (PDF)</label>
	<input id="modul" type="file" accept=".pdf,application/pdf" onchange={uploadModul} disabled={uploading} />
	{#if uploading}<p class="mt-2 text-sm text-ink-caption">Mengunggah…</p>{/if}
</div>
