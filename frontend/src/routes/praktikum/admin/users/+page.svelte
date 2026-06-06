<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { Kelas, User } from '$lib/types';

	let users = $state<User[]>([]);
	let kelas = $state<Kelas[]>([]);
	let err = $state('');
	let msg = $state('');

	let editId = $state<number | null>(null);
	let form = $state({ nim: '', nama: '', kelas_id: null as number | null, shift: 1 as number, kelompok: '' });

	async function load() {
		try {
			users = (await api.get<User[]>('/api/admin/users')) ?? [];
			kelas = (await api.get<Kelas[]>('/api/admin/kelas')) ?? [];
		} catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	function resetForm() {
		editId = null;
		form = { nim: '', nama: '', kelas_id: kelas[0]?.id ?? null, shift: 1, kelompok: '' };
	}

	function edit(u: User) {
		editId = u.id;
		form = { nim: u.nim, nama: u.nama, kelas_id: u.kelas_id ?? null, shift: u.shift ?? 1, kelompok: u.kelompok ?? '' };
	}

	async function save() {
		err = ''; msg = '';
		try {
			const body = { nim: form.nim, nama: form.nama, kelas_id: form.kelas_id, shift: Number(form.shift), kelompok: form.kelompok || null };
			if (editId) await api.put(`/api/admin/users/${editId}`, body);
			else await api.post('/api/admin/users', body);
			msg = 'Tersimpan.';
			resetForm();
			await load();
		} catch (e) { err = (e as Error).message; }
	}

	async function del(id: number) {
		if (!confirm('Hapus mahasiswa ini?')) return;
		try { await api.del(`/api/admin/users/${id}`); await load(); }
		catch (e) { err = (e as Error).message; }
	}

	async function resetPw(id: number) {
		if (!confirm('Reset password? Mahasiswa harus register ulang.')) return;
		try { await api.post(`/api/admin/users/${id}/reset-password`); msg = 'Password direset.'; await load(); }
		catch (e) { err = (e as Error).message; }
	}

	async function toggleRegister(k: Kelas) {
		try {
			await api.post('/api/admin/kelas-register', { kelas_id: k.id, open: !k.is_register_open });
			await load();
		} catch (e) { err = (e as Error).message; }
	}
</script>

<h1 class="mb-4 text-2xl">Manajemen Data User</h1>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="mb-6 grid gap-4 lg:grid-cols-3">
	<div class="card lg:col-span-2">
		<h2 class="mb-2 text-lg">Akses Register per Kelas</h2>
		<div class="flex flex-wrap gap-2">
			{#each kelas as k}
				<button
					class="badge {k.is_register_open ? 'bg-state-success-bg text-state-success' : 'bg-gray-100 text-ink-caption'} cursor-pointer px-3 py-1"
					onclick={() => toggleRegister(k)}
				>{k.nama_kelas}: {k.is_register_open ? 'DIBUKA' : 'ditutup'}</button>
			{/each}
		</div>
	</div>
</div>

<div class="grid gap-4 lg:grid-cols-3">
	<div class="card lg:col-span-1">
		<h2 class="mb-3 text-lg">{editId ? 'Edit' : 'Tambah'} Mahasiswa</h2>
		<label class="label" for="nim">NIM</label>
		<input id="nim" class="input" bind:value={form.nim} />
		<label class="label mt-2" for="nama">Nama</label>
		<input id="nama" class="input" bind:value={form.nama} />
		<label class="label mt-2" for="kelas">Kelas</label>
		<select id="kelas" class="input" bind:value={form.kelas_id}>
			{#each kelas as k}<option value={k.id}>{k.nama_kelas}</option>{/each}
		</select>
		<label class="label mt-2" for="shift">Shift</label>
		<select id="shift" class="input" bind:value={form.shift}>
			<option value={1}>Shift 1</option>
			<option value={2}>Shift 2</option>
		</select>
		<label class="label mt-2" for="kelompok">Kelompok</label>
		<input id="kelompok" class="input" bind:value={form.kelompok} placeholder="mis. A, B, 1" />
		<div class="mt-3 flex gap-2">
			<button class="btn-primary" onclick={save}>Simpan</button>
			{#if editId}<button class="btn-outline" onclick={resetForm}>Batal</button>{/if}
		</div>
	</div>

	<div class="lg:col-span-2">
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr><th>NIM</th><th>Nama</th><th>Kelas</th><th>Shift</th><th>Kelompok</th><th>Status</th><th>Aksi</th></tr></thead>
				<tbody>
					{#each users as u}
						<tr>
							<td>{u.nim}</td>
							<td>{u.nama}</td>
							<td>{u.nama_kelas ?? u.kelas?.nama_kelas ?? '-'}</td>
							<td>{u.shift ?? '-'}</td>
							<td>{u.kelompok ?? '-'}</td>
							<td>
								{#if u.is_registered}
									<span class="badge bg-state-success-bg text-state-success">Terdaftar</span>
								{:else}
									<span class="badge bg-gray-100 text-ink-caption">Belum</span>
								{/if}
							</td>
							<td class="space-x-1 whitespace-nowrap">
								<button class="text-primary hover:underline" onclick={() => edit(u)}>Edit</button>
								<button class="text-state-warning hover:underline" onclick={() => resetPw(u.id)}>Reset PW</button>
								<button class="text-state-error hover:underline" onclick={() => del(u.id)}>Hapus</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
