<script lang="ts">
	import { onMount } from 'svelte';
	import Papa from 'papaparse';
	import { api } from '$lib/api';
	import { KeyRound, Trash2 } from 'lucide-svelte';
	import type { Kelas, User } from '$lib/types';

	let users = $state<User[]>([]);
	let kelas = $state<Kelas[]>([]);
	let err = $state('');
	let msg = $state('');

	let editId = $state<number | null>(null);
	let form = $state({ nim: '', nama: '', kelas_id: null as number | null, shift: 1 as number, kelompok: '' });

	let showImport = $state(false);
	let importRows = $state<any[]>([]);
	let importErrors = $state<string[]>([]);
	let isImporting = $state(false);

	// Bulk select
	let selectedIds = $state<Set<number>>(new Set());
	let bulkBusy = $state(false);

	async function load() {
		try {
			users = (await api.get<User[]>('/api/admin/users')) ?? [];
			kelas = (await api.get<Kelas[]>('/api/admin/kelas')) ?? [];
		} catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	function toggleSelect(id: number) {
		const s = new Set(selectedIds);
		s.has(id) ? s.delete(id) : s.add(id);
		selectedIds = s;
	}
	function toggleSelectAll() {
		selectedIds = selectedIds.size === users.length ? new Set() : new Set(users.map((u) => u.id));
	}

	async function bulkAction(action: 'delete' | 'reset_pw') {
		if (selectedIds.size === 0) return;
		const label = action === 'delete' ? 'Hapus' : 'Reset password';
		const extra = action === 'reset_pw' ? ' Mereka harus register ulang.' : '';
		if (!confirm(`${label} ${selectedIds.size} mahasiswa terpilih?${extra}`)) return;
		bulkBusy = true; err = ''; msg = '';
		let ok = 0, fail = 0;
		for (const id of selectedIds) {
			try {
				if (action === 'delete') await api.del(`/api/admin/users/${id}`);
				else await api.post(`/api/admin/users/${id}/reset-password`);
				ok++;
			} catch { fail++; }
		}
		bulkBusy = false;
		selectedIds = new Set();
		msg = `${label} selesai: ${ok} berhasil${fail ? `, ${fail} gagal` : ''}.`;
		await load();
	}

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

	function onFileChange(e: Event) {
		const file = (e.target as HTMLInputElement).files?.[0];
		if (!file) return;
		Papa.parse(file, {
			header: true,
			skipEmptyLines: true,
			complete: (results) => {
				importErrors = [];
				let valid = [];
				for (let i = 0; i < results.data.length; i++) {
					const row = results.data[i] as any;
					if (!row.NIM || !row.Nama || !row.Kelas) {
						importErrors.push(`Baris ${i + 1}: Format tidak lengkap. Harus ada header NIM, Nama, Kelas.`);
						continue;
					}
					const k = kelas.find((kl) => kl.nama_kelas === row.Kelas);
					if (!k) {
						importErrors.push(`Baris ${i + 1} (${row.NIM}): Kelas "${row.Kelas}" tidak ditemukan di database.`);
						continue;
					}
					valid.push({
						nim: row.NIM,
						nama: row.Nama,
						kelas_id: k.id,
						shift: Number(row.Shift) || 1,
						kelompok: row.Kelompok || ''
					});
				}
				if (importErrors.length === 0) importRows = valid;
				else importRows = [];
			}
		});
	}

	async function importBulk() {
		if (importRows.length === 0) return;
		isImporting = true;
		err = ''; msg = '';
		try {
			const res = await api.post<any>('/api/admin/users/bulk', { users: importRows });
			msg = `Berhasil memproses ${res?.total_processed ?? importRows.length} mahasiswa.`;
			showImport = false;
			importRows = [];
			await load();
		} catch (e) { err = (e as Error).message; }
		finally { isImporting = false; }
	}

	function downloadTemplate() {
		const csvContent = "NIM,Nama,Kelas,Shift,Kelompok\n12345678,Budi Santoso,4IA01,1,A\n";
		const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.setAttribute('href', url);
		link.setAttribute('download', 'template_import_mahasiswa.csv');
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}
</script>

<div class="mb-4 flex items-center justify-between">
	<h1 class="text-2xl">Manajemen Data User</h1>
	<div class="space-x-2 flex">
		<button class="btn-outline" onclick={downloadTemplate}>Unduh Template</button>
		<button class="btn-primary" onclick={() => (showImport = !showImport)}>
			{showImport ? 'Tutup Import' : 'Import CSV'}
		</button>
	</div>
</div>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="mb-6 grid gap-4 lg:grid-cols-3">
	{#if showImport}
		<div class="card border-primary border-dashed border-2 lg:col-span-3">
			<h2 class="mb-3 text-lg font-bold">Import Mahasiswa via CSV</h2>
			<p class="mb-3 text-sm text-ink-caption">
				Pastikan file CSV memiliki header (baris pertama) persis: <strong>NIM, Nama, Kelas, Shift, Kelompok</strong>.<br/>
				Jika NIM sudah ada di sistem, data namanya, kelas, dan shift akan di-_update_ tanpa mereset password.
			</p>
			<input type="file" accept=".csv" class="mb-3 block w-full text-sm file:mr-4 file:rounded-lg file:border-0 file:bg-primary file:px-4 file:py-2 file:text-sm file:font-semibold file:text-white hover:file:bg-primary/90" onchange={onFileChange} />
			
			{#if importErrors.length > 0}
				<div class="mb-3 max-h-48 overflow-y-auto rounded-lg bg-state-error-bg p-3 text-sm text-state-error">
					<p class="font-bold">Terdapat Error pada file CSV:</p>
					<ul class="list-inside list-disc">
						{#each importErrors as ie}<li>{ie}</li>{/each}
					</ul>
					<p class="mt-2 font-bold">Harap perbaiki file CSV dan unggah ulang.</p>
				</div>
			{/if}

			{#if importRows.length > 0}
				<p class="mb-2 text-sm font-bold text-state-success">Valid! {importRows.length} data siap diimpor.</p>
				<button class="btn-primary w-full max-w-sm" onclick={importBulk} disabled={isImporting}>
					{isImporting ? 'Memproses...' : 'Mulai Import'}
				</button>
			{/if}
		</div>
	{/if}
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
		{#if selectedIds.size > 0}
			<div class="mb-3 flex flex-wrap items-center gap-3 rounded-xl border border-primary/20 bg-primary/5 p-3 shadow-sm">
				<span class="text-sm font-medium text-primary"><strong>{selectedIds.size}</strong> dipilih</span>
				<div class="ml-auto flex flex-wrap gap-2">
					<button class="btn-outline inline-flex items-center gap-1 border-state-warning py-1.5 text-state-warning hover:bg-state-warning hover:text-white" disabled={bulkBusy} onclick={() => bulkAction('reset_pw')}>
						<KeyRound size={14} /> Reset PW
					</button>
					<button class="btn-outline inline-flex items-center gap-1 border-state-error py-1.5 text-state-error hover:bg-state-error hover:text-white" disabled={bulkBusy} onclick={() => bulkAction('delete')}>
						<Trash2 size={14} /> Hapus
					</button>
					<button class="btn-outline py-1.5" disabled={bulkBusy} onclick={() => (selectedIds = new Set())}>Batal</button>
				</div>
			</div>
		{/if}
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr>
					<th class="w-10 text-center">
						<input type="checkbox" class="rounded border-gray-300 text-primary focus:ring-primary"
							checked={users.length > 0 && selectedIds.size === users.length}
							onchange={toggleSelectAll} aria-label="Pilih semua" />
					</th>
					<th>NIM</th><th>Nama</th><th>Kelas</th><th>Shift</th><th>Kelompok</th><th>Status</th><th>Aksi</th>
				</tr></thead>
				<tbody>
					{#each users as u}
						<tr class={selectedIds.has(u.id) ? 'bg-primary/5' : ''}>
							<td class="text-center">
								<input type="checkbox" class="rounded border-gray-300 text-primary focus:ring-primary"
									checked={selectedIds.has(u.id)} onchange={() => toggleSelect(u.id)} aria-label={`Pilih ${u.nim}`} />
							</td>
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
