<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { Kelas } from '$lib/types';

	let list = $state<Kelas[]>([]);
	let err = $state(''); let msg = $state('');
	let editId = $state<number | null>(null);
	let nama = $state('');

	async function load() {
		try { list = (await api.get<Kelas[]>('/api/admin/kelas')) ?? []; }
		catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	function reset() { editId = null; nama = ''; }
	function edit(k: Kelas) { editId = k.id; nama = k.nama_kelas; }

	async function save() {
		err = ''; msg = '';
		try {
			if (editId) await api.put(`/api/admin/kelas/${editId}`, { nama_kelas: nama });
			else await api.post('/api/admin/kelas', { nama_kelas: nama });
			msg = 'Tersimpan.'; reset(); await load();
		} catch (e) { err = (e as Error).message; }
	}
	async function del(id: number) {
		if (!confirm('Hapus kelas?')) return;
		try { await api.del(`/api/admin/kelas/${id}`); await load(); }
		catch (e) { err = (e as Error).message; }
	}
</script>

<h1 class="mb-4 text-2xl">Manajemen Kelas</h1>
{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="card">
		<h2 class="mb-3 text-lg">{editId ? 'Edit' : 'Tambah'} Kelas</h2>
		<label class="label" for="n">Nama Kelas</label>
		<input id="n" class="input" bind:value={nama} placeholder="mis. TTL A" />
		<div class="mt-3 flex gap-2">
			<button class="btn-primary" onclick={save}>Simpan</button>
			{#if editId}<button class="btn-outline" onclick={reset}>Batal</button>{/if}
		</div>
	</div>
	<div class="lg:col-span-2">
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr><th>Nama</th><th>Register</th><th>Aksi</th></tr></thead>
				<tbody>
					{#each list as k}
						<tr>
							<td>{k.nama_kelas}</td>
							<td>{k.is_register_open ? 'Dibuka' : 'Ditutup'}</td>
							<td class="space-x-2">
								<button class="text-primary hover:underline" onclick={() => edit(k)}>Edit</button>
								<button class="text-state-error hover:underline" onclick={() => del(k.id)}>Hapus</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
