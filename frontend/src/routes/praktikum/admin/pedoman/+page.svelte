<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	interface Pedoman { id: number; nama_dokumen: string; file_url: string; }

	let list = $state<Pedoman[]>([]);
	let err = $state(''); let msg = $state('');
	let editId = $state<number | null>(null);
	let form = $state({ nama_dokumen: '', file_url: '' });

	async function load() {
		try { list = (await api.get<Pedoman[]>('/api/admin/pedoman')) ?? []; }
		catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	function reset() { editId = null; form = { nama_dokumen: '', file_url: '' }; }
	function edit(p: Pedoman) {
		editId = p.id;
		form = { nama_dokumen: p.nama_dokumen, file_url: p.file_url };
	}

	async function uploadFile(ev: Event) {
		const input = ev.target as HTMLInputElement;
		if (!input.files?.[0]) return;
		const fd = new FormData();
		fd.append('file', input.files[0]); fd.append('folder', 'pedoman');
		try {
			const res = await api.upload<{ url: string }>('/api/admin/upload', fd);
			form.file_url = res.url;
			msg = 'File berhasil diunggah.';
		} catch (e) { err = (e as Error).message; }
	}

	async function save() {
		err = ''; msg = '';
		try {
			if (editId) await api.put(`/api/admin/pedoman/${editId}`, form);
			else await api.post('/api/admin/pedoman', form);
			msg = 'Tersimpan.'; reset(); await load();
		} catch (e) { err = (e as Error).message; }
	}

	async function del(id: number) {
		if (!confirm('Hapus pedoman ini?')) return;
		try { await api.del(`/api/admin/pedoman/${id}`); await load(); }
		catch (e) { err = (e as Error).message; }
	}
</script>

<h1 class="mb-4 text-2xl">Manajemen Pedoman Laporan</h1>
<p class="mb-4 text-sm text-ink-caption">File yang diunggah di sini akan tampil sebagai tombol download di halaman /info/laporan.</p>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="card">
		<h2 class="mb-3 text-lg">{editId ? 'Edit' : 'Tambah'} Pedoman</h2>
		<label class="label" for="nd">Nama Dokumen</label>
		<input id="nd" class="input" bind:value={form.nama_dokumen} placeholder="mis. Template Laporan Akhir" />
		<label class="label mt-2" for="fu">File</label>
		<input id="fu" type="file" onchange={uploadFile} />
		{#if form.file_url}
			<p class="mt-1 text-xs text-state-success">✓ File: <a href={form.file_url} target="_blank" rel="noopener">{form.file_url.split('/').pop()}</a></p>
		{/if}
		<div class="mt-3 flex gap-2">
			<button class="btn-primary" onclick={save}>Simpan</button>
			{#if editId}<button class="btn-outline" onclick={reset}>Batal</button>{/if}
		</div>
	</div>

	<div class="lg:col-span-2">
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr><th>Nama Dokumen</th><th>File</th><th>Aksi</th></tr></thead>
				<tbody>
					{#each list as p}
						<tr>
							<td>{p.nama_dokumen}</td>
							<td><a href={p.file_url} target="_blank" rel="noopener" class="text-primary hover:underline">Download</a></td>
							<td class="space-x-2">
								<button class="text-primary hover:underline" onclick={() => edit(p)}>Edit</button>
								<button class="text-state-error hover:underline" onclick={() => del(p.id)}>Hapus</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
