<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { User, Kelas, AmpuanKelompok } from '$lib/types';

	let list = $state<User[]>([]);
	let kelasList = $state<Kelas[]>([]);
	let ampuanList = $state<AmpuanKelompok[]>([]);
	let err = $state(''); let msg = $state('');
	let editId = $state<number | null>(null);
	let form = $state({ nim: '', nama: '', nomor_hp: '', medsos_link: '', foto_url: '', password: '' });
	let ampuanForm = $state({ asisten_id: 0, kelas_id: 0, kelompok: '' });

	async function load() {
		try {
			list = (await api.get<User[]>('/api/admin/asisten')) ?? [];
			kelasList = (await api.get<Kelas[]>('/api/admin/kelas')) ?? [];
			ampuanList = (await api.get<AmpuanKelompok[]>('/api/admin/ampuan')) ?? [];
		} catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	function reset() { editId = null; form = { nim: '', nama: '', nomor_hp: '', medsos_link: '', foto_url: '', password: '' }; }
	function edit(a: User) {
		editId = a.id;
		form = { nim: a.nim, nama: a.nama, nomor_hp: a.nomor_hp ?? '', medsos_link: a.medsos_link ?? '', foto_url: a.foto_url ?? '', password: '' };
	}

	async function uploadFoto(ev: Event) {
		const input = ev.target as HTMLInputElement;
		if (!input.files?.[0]) return;
		const fd = new FormData();
		fd.append('file', input.files[0]); fd.append('folder', 'asisten');
		try { const res = await api.upload<{ url: string }>('/api/admin/upload', fd); form.foto_url = res.url; }
		catch (e) { err = (e as Error).message; }
	}

	async function save() {
		err = ''; msg = '';
		const body: Record<string, unknown> = {
			nim: form.nim, nama: form.nama,
			nomor_hp: form.nomor_hp || null, medsos_link: form.medsos_link || null, foto_url: form.foto_url || null
		};
		if (form.password) body.password = form.password;
		try {
			if (editId) await api.put(`/api/admin/asisten/${editId}`, body);
			else await api.post('/api/admin/asisten', body);
			msg = 'Tersimpan.'; reset(); await load();
		} catch (e) { err = (e as Error).message; }
	}

	async function addAmpuan() {
		err = ''; msg = '';
		try {
			await api.post('/api/admin/ampuan', ampuanForm);
			msg = 'Ampuan ditambahkan.';
			ampuanForm = { asisten_id: 0, kelas_id: 0, kelompok: '' };
			await load();
		} catch (e) { err = (e as Error).message; }
	}

	async function delAmpuan(id: number) {
		if (!confirm('Hapus ampuan ini?')) return;
		try { await api.del(`/api/admin/ampuan/${id}`); await load(); }
		catch (e) { err = (e as Error).message; }
	}
</script>

<h1 class="mb-4 text-2xl">Manajemen Asisten Lab</h1>
<p class="mb-4 text-sm text-ink-caption">Asisten = admin. Data ini juga tampil di halaman publik /info/asisten.</p>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="card lg:col-span-1">
		<h2 class="mb-3 text-lg">{editId ? 'Edit' : 'Tambah'} Asisten</h2>
		<label class="label" for="nim">NIM</label>
		<input id="nim" class="input" bind:value={form.nim} />
		<label class="label mt-2" for="nama">Nama</label>
		<input id="nama" class="input" bind:value={form.nama} />
		<label class="label mt-2" for="hp">Nomor HP</label>
		<input id="hp" class="input" bind:value={form.nomor_hp} />
		<label class="label mt-2" for="ms">Link Medsos</label>
		<input id="ms" class="input" bind:value={form.medsos_link} />
		<label class="label mt-2" for="foto">Foto</label>
		{#if form.foto_url}<img src={form.foto_url} alt="foto" class="mb-2 h-16 w-16 rounded-full object-cover" />{/if}
		<input id="foto" type="file" accept="image/*" onchange={uploadFoto} />
		<label class="label mt-2" for="pw">Password {editId ? '(kosongkan jika tidak diubah)' : ''}</label>
		<input id="pw" type="password" class="input" bind:value={form.password} />
		<div class="mt-3 flex gap-2">
			<button class="btn-primary" onclick={save}>Simpan</button>
			{#if editId}<button class="btn-outline" onclick={reset}>Batal</button>{/if}
		</div>
	</div>

	<div class="lg:col-span-2">
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr><th>Foto</th><th>Nama</th><th>NIM</th><th>HP</th><th>Aksi</th></tr></thead>
				<tbody>
					{#each list as a}
						<tr>
							<td>{#if a.foto_url}<img src={a.foto_url} alt="" class="h-10 w-10 rounded-full object-cover" />{:else}-{/if}</td>
							<td>{a.nama}</td>
							<td>{a.nim}</td>
							<td>{a.nomor_hp ?? '-'}</td>
							<td><button class="text-primary hover:underline" onclick={() => edit(a)}>Edit</button></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>

<hr class="my-6 border-gray-200" />
<h2 class="mb-3 text-xl">Ampuan Kelompok</h2>
<p class="mb-4 text-sm text-ink-caption">Assign asisten sebagai pengampu kelompok di setiap kelas.</p>

<div class="grid gap-4 lg:grid-cols-3">
	<div class="card">
		<h3 class="mb-3 text-lg">Tambah Ampuan</h3>
		<label class="label" for="amp-as">Asisten</label>
		<select id="amp-as" class="input" bind:value={ampuanForm.asisten_id}>
			<option value={0}>— Pilih Asisten —</option>
			{#each list as a}<option value={a.id}>{a.nama}</option>{/each}
		</select>
		<label class="label mt-2" for="amp-kl">Kelas</label>
		<select id="amp-kl" class="input" bind:value={ampuanForm.kelas_id}>
			<option value={0}>— Pilih Kelas —</option>
			{#each kelasList as k}<option value={k.id}>{k.nama_kelas}</option>{/each}
		</select>
		<label class="label mt-2" for="amp-kel">Kelompok</label>
		<input id="amp-kel" class="input" bind:value={ampuanForm.kelompok} placeholder="mis. A, B, 1" />
		<button class="btn-primary mt-3 w-full" onclick={addAmpuan}>Tambah</button>
	</div>

	<div class="lg:col-span-2">
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr><th>Asisten</th><th>Kelas</th><th>Kelompok</th><th>Aksi</th></tr></thead>
				<tbody>
					{#each ampuanList as a}
						<tr>
							<td>{a.asisten?.nama ?? a.asisten_id}</td>
							<td>{a.kelas?.nama_kelas ?? a.kelas_id}</td>
							<td>{a.kelompok}</td>
							<td><button class="text-state-error hover:underline" onclick={() => delAmpuan(a.id)}>Hapus</button></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>

