<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { labelJenis } from '$lib/utils';
	import type { Jadwal, Kelas } from '$lib/types';

	let list = $state<Jadwal[]>([]);
	let kelas = $state<Kelas[]>([]);
	let err = $state(''); let msg = $state('');
	let editId = $state<number | null>(null);
	let form = $state({ kelas_id: null as number | null, shift: 1, hari: 'Senin', jam_mulai: '08:00', jam_selesai: '10:00', keterangan: '' });

	// Konfigurasi mode jadwal publik (gdrive / internal)
	let mode = $state('internal');
	let gdriveUrl = $state('');

	async function load() {
		try {
			list = (await api.get<Jadwal[]>('/api/admin/jadwal')) ?? [];
			kelas = (await api.get<Kelas[]>('/api/admin/kelas')) ?? [];
			if (!form.kelas_id) form.kelas_id = kelas[0]?.id ?? null;
			const konf = (await api.get<{ key: string; value: string }[]>('/api/admin/konfigurasi')) ?? [];
			mode = konf.find((k) => k.key === 'jadwal_mode')?.value ?? 'internal';
			gdriveUrl = konf.find((k) => k.key === 'gdrive_jadwal_url')?.value ?? '';
		} catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	function reset() { editId = null; form = { kelas_id: kelas[0]?.id ?? null, shift: 1, hari: 'Senin', jam_mulai: '08:00', jam_selesai: '10:00', keterangan: '' }; }
	function edit(j: Jadwal) {
		editId = j.id;
		form = { kelas_id: j.kelas_id, shift: j.shift, hari: j.hari, jam_mulai: j.jam_mulai?.slice(0,5), jam_selesai: j.jam_selesai?.slice(0,5), keterangan: j.keterangan };
	}

	async function save() {
		err = ''; msg = '';
		const body = { ...form, shift: Number(form.shift), jam_mulai: form.jam_mulai + ':00', jam_selesai: form.jam_selesai + ':00' };
		try {
			if (editId) await api.put(`/api/admin/jadwal/${editId}`, body);
			else await api.post('/api/admin/jadwal', body);
			msg = 'Tersimpan.'; reset(); await load();
		} catch (e) { err = (e as Error).message; }
	}
	async function del(id: number) {
		if (!confirm('Hapus jadwal?')) return;
		try { await api.del(`/api/admin/jadwal/${id}`); await load(); } catch (e) { err = (e as Error).message; }
	}

	async function saveKonfig() {
		try {
			await api.post('/api/admin/konfigurasi', { key: 'jadwal_mode', value: mode });
			await api.post('/api/admin/konfigurasi', { key: 'gdrive_jadwal_url', value: gdriveUrl });
			msg = 'Konfigurasi jadwal disimpan.';
		} catch (e) { err = (e as Error).message; }
	}
</script>

<h1 class="mb-4 text-2xl">Manajemen Jadwal</h1>
{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="card mb-6 max-w-2xl">
	<h2 class="mb-2 text-lg">Mode Tampilan Jadwal Publik</h2>
	<div class="flex flex-wrap items-end gap-3">
		<div>
			<label class="label" for="mode">Mode</label>
			<select id="mode" class="input" bind:value={mode}>
				<option value="internal">Tabel Internal</option>
				<option value="gdrive">Link Google Drive</option>
			</select>
		</div>
		<div class="flex-1">
			<label class="label" for="gd">URL Google Drive</label>
			<input id="gd" class="input" bind:value={gdriveUrl} placeholder="https://drive.google.com/..." />
		</div>
		<button class="btn-primary" onclick={saveKonfig}>Simpan</button>
	</div>
</div>

<div class="grid gap-4 lg:grid-cols-3">
	<div class="card">
		<h2 class="mb-3 text-lg">{editId ? 'Edit' : 'Tambah'} Jadwal</h2>
		<label class="label" for="k">Kelas</label>
		<select id="k" class="input" bind:value={form.kelas_id}>{#each kelas as k}<option value={k.id}>{k.nama_kelas}</option>{/each}</select>
		<label class="label mt-2" for="s">Shift</label>
		<select id="s" class="input" bind:value={form.shift}><option value={1}>1</option><option value={2}>2</option></select>
		<label class="label mt-2" for="h">Hari</label>
		<input id="h" class="input" bind:value={form.hari} />
		<div class="mt-2 flex gap-2">
			<div class="flex-1"><label class="label" for="jm">Mulai</label><input id="jm" type="time" class="input" bind:value={form.jam_mulai} /></div>
			<div class="flex-1"><label class="label" for="js">Selesai</label><input id="js" type="time" class="input" bind:value={form.jam_selesai} /></div>
		</div>
		<label class="label mt-2" for="ket">Keterangan</label>
		<input id="ket" class="input" bind:value={form.keterangan} placeholder="mis. Minggu 1-4" />
		<div class="mt-3 flex gap-2">
			<button class="btn-primary" onclick={save}>Simpan</button>
			{#if editId}<button class="btn-outline" onclick={reset}>Batal</button>{/if}
		</div>
	</div>

	<div class="lg:col-span-2">
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr><th>Kelas</th><th>Shift</th><th>Hari</th><th>Jam</th><th>Ket.</th><th>Aksi</th></tr></thead>
				<tbody>
					{#each list as j}
						<tr>
							<td>{j.kelas?.nama_kelas ?? j.kelas_id}</td>
							<td>{j.shift}</td>
							<td>{j.hari}</td>
							<td>{j.jam_mulai?.slice(0,5)}–{j.jam_selesai?.slice(0,5)}</td>
							<td>{j.keterangan}</td>
							<td class="space-x-2">
								<button class="text-primary hover:underline" onclick={() => edit(j)}>Edit</button>
								<button class="text-state-error hover:underline" onclick={() => del(j.id)}>Hapus</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
