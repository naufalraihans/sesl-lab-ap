<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { labelJenis } from '$lib/utils';
	import { KeyRound } from 'lucide-svelte';
	import type { Sesi, Kelas, User } from '$lib/types';

	interface AktivasiSesi {
		id: number; sesi_praktikum_id: number; kelas_id: number; shift: number;
		is_active: boolean; activated_at: string; token?: string;
		sesi?: { judul_sesi: string };
		kelas?: { nama_kelas: string };
		aktivasi_courses?: AktivasiCourse[];
	}
	interface AktivasiCourse {
		id: number; course_id: number; is_open: boolean; urutan: number;
		course?: { jenis: string; judul: string };
	}
	interface Susulan { id: number; mahasiswa_id: number; alasan: string; mahasiswa?: { nama: string; nim: string } }

	let aktivasiList = $state<AktivasiSesi[]>([]);
	let sesiList = $state<Sesi[]>([]);
	let kelasList = $state<Kelas[]>([]);
	let users = $state<User[]>([]);
	let err = $state(''); let msg = $state('');

	let form = $state({ sesi_praktikum_id: 0, kelas_id: 0, shift: 1, gacha_pilihan: 'pretest' });

	let selected = $state<AktivasiSesi | null>(null);
	let susulanList = $state<Susulan[]>([]);
	let susulanForm = $state({ mahasiswa_id: 0, alasan: '' });

	async function load() {
		try {
			aktivasiList = (await api.get<AktivasiSesi[]>('/api/admin/aktivasi')) ?? [];
			sesiList = (await api.get<Sesi[]>('/api/admin/sesi')) ?? [];
			kelasList = (await api.get<Kelas[]>('/api/admin/kelas')) ?? [];
			users = (await api.get<User[]>('/api/admin/users')) ?? [];
		} catch (e) { err = (e as Error).message; }
	}
	onMount(load);

	async function aktivasi() {
		err = ''; msg = '';
		try {
			await api.post('/api/admin/aktivasi', form);
			msg = 'Sesi berhasil diaktifkan.'; await load();
		} catch (e) { err = (e as Error).message; }
	}

	async function selectAktivasi(a: AktivasiSesi) {
		selected = a;
		try {
			const detail = await api.get<AktivasiSesi>(`/api/admin/aktivasi/${a.id}`);
			if (detail) selected = detail;
			susulanList = (await api.get<Susulan[]>(`/api/admin/aktivasi/${a.id}/susulan`)) ?? [];
		} catch (e) { err = (e as Error).message; }
	}

	async function toggleCourse(ac: AktivasiCourse) {
		err = ''; msg = '';
		const action = ac.is_open ? 'MENUTUP (auto-submit massal)' : 'MEMBUKA';
		if (!confirm(`${action} course ini?`)) return;
		try {
			await api.post('/api/admin/aktivasi-course/buka-tutup', {
				aktivasi_course_id: ac.id,
				is_open: !ac.is_open
			});
			msg = ac.is_open ? 'Course ditutup. Auto-submit massal dijalankan.' : 'Course dibuka.';
			if (selected) await selectAktivasi(selected);
		} catch (e) { err = (e as Error).message; }
	}

	async function addSusulan() {
		if (!selected) return;
		err = ''; msg = '';
		try {
			await api.post(`/api/admin/aktivasi/${selected.id}/susulan`, susulanForm);
			msg = 'Mahasiswa susulan didaftarkan.';
			susulanForm = { mahasiswa_id: 0, alasan: '' };
			await selectAktivasi(selected);
		} catch (e) { err = (e as Error).message; }
	}

	async function removeSusulan(mhsId: number) {
		if (!selected || !confirm('Hapus akses susulan?')) return;
		try {
			await api.del(`/api/admin/aktivasi/${selected.id}/susulan/${mhsId}`);
			await selectAktivasi(selected);
		} catch (e) { err = (e as Error).message; }
	}

	async function generateToken(a: AktivasiSesi) {
		err = ''; msg = '';
		if (!confirm('Generate/Reset PIN ujian untuk kelas ini?')) return;
		try {
			await api.post(`/api/admin/aktivasi/${a.id}/token`);
			msg = 'PIN Ujian berhasil dibuat/direset.';
			await load();
			if (selected && selected.id === a.id) {
				await selectAktivasi(selected);
			}
		} catch (e) { err = (e as Error).message; }
	}
</script>

<h1 class="mb-4 text-2xl">Aktivasi Sesi</h1>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="card">
		<h2 class="mb-3 text-lg">Aktifkan Sesi</h2>
		<label class="label" for="as">Sesi</label>
		<select id="as" class="input" bind:value={form.sesi_praktikum_id}>
			<option value={0}>— Pilih Sesi —</option>
			{#each sesiList as s}<option value={s.id}>{s.judul_sesi}</option>{/each}
		</select>
		<label class="label mt-2" for="ak">Kelas</label>
		<select id="ak" class="input" bind:value={form.kelas_id}>
			<option value={0}>— Pilih Kelas —</option>
			{#each kelasList as k}<option value={k.id}>{k.nama_kelas}</option>{/each}
		</select>
		<label class="label mt-2" for="ash">Shift</label>
		<select id="ash" class="input" bind:value={form.shift}>
			<option value={1}>Shift 1</option>
			<option value={2}>Shift 2</option>
		</select>
		<label class="label mt-2" for="ag">Gacha (Pre/Post-test)</label>
		<select id="ag" class="input" bind:value={form.gacha_pilihan}>
			<option value="pretest">Pre-test</option>
			<option value="posttest">Post-test</option>
		</select>
		<p class="mt-1 text-xs text-ink-caption">Diabaikan untuk sesi ujian praktik.</p>
		<button class="btn-primary mt-3 w-full" onclick={aktivasi}>Aktifkan Sesi</button>
	</div>

	<div class="lg:col-span-2">
		<div class="table-wrap">
			<table class="tbl">
				<thead><tr><th>Sesi</th><th>Kelas</th><th>Shift</th><th>Aksi</th></tr></thead>
				<tbody>
					{#each aktivasiList as a}
						<tr class={selected?.id === a.id ? 'ring-2 ring-primary' : ''}>
							<td>{a.sesi?.judul_sesi ?? a.sesi_praktikum_id}</td>
							<td>{a.kelas?.nama_kelas ?? a.kelas_id}</td>
							<td>{a.shift}</td>
							<td><button class="text-state-info hover:underline" onclick={() => selectAktivasi(a)}>Detail</button></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>

{#if selected}
	<hr class="my-6 border-gray-200" />
	<div class="mb-5 flex flex-col md:flex-row md:items-center justify-between gap-4">
		<h2 class="text-xl">
			{selected.sesi?.judul_sesi} — {selected.kelas?.nama_kelas} Shift {selected.shift}
		</h2>
		<div class="flex items-center gap-3">
			{#if selected.token}
				<div class="bg-primary/10 border-2 border-primary text-primary px-4 py-1.5 rounded-lg font-mono text-xl font-bold tracking-widest shadow-inner">
					{selected.token}
				</div>
			{/if}
			<button class="btn-primary" onclick={() => generateToken(selected!)}>
				<KeyRound size={16} /> {selected.token ? 'Reset PIN' : 'Generate PIN'}
			</button>
		</div>
	</div>

	<h3 class="mb-2 text-lg">Buka / Tutup Course</h3>
	<p class="mb-3 text-xs text-ink-caption">Menutup course = auto-submit massal untuk semua mahasiswa yang belum submit.</p>
	<div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
		{#each selected.aktivasi_courses ?? [] as ac}
			<div class="card flex items-center justify-between">
				<div>
					<p class="font-medium">{ac.course?.judul ?? labelJenis(ac.course?.jenis ?? '')}</p>
					<span class="badge mt-1 {ac.is_open ? 'bg-state-success-bg text-state-success' : 'bg-gray-100 text-ink-caption'}">
						{ac.is_open ? 'Terbuka' : 'Tertutup'}
					</span>
				</div>
				<button
					class="btn-outline py-1.5 text-xs {ac.is_open ? 'text-state-error' : 'text-state-success'}"
					onclick={() => toggleCourse(ac)}
				>{ac.is_open ? 'Tutup' : 'Buka'}</button>
			</div>
		{/each}
	</div>

	<hr class="my-6 border-gray-200" />
	<h3 class="mb-2 text-lg">Peserta Susulan</h3>
	<div class="grid gap-4 lg:grid-cols-3">
		<div class="card">
			<label class="label" for="sm">Mahasiswa</label>
			<select id="sm" class="input" bind:value={susulanForm.mahasiswa_id}>
				<option value={0}>— Pilih —</option>
				{#each users as u}<option value={u.id}>{u.nama} ({u.nim})</option>{/each}
			</select>
			<label class="label mt-2" for="sa">Alasan</label>
			<input id="sa" class="input" bind:value={susulanForm.alasan} />
			<button class="btn-primary mt-3" onclick={addSusulan}>Tambah Susulan</button>
		</div>
		<div class="lg:col-span-2">
			{#if susulanList.length === 0}
				<p class="text-sm text-ink-caption">Tidak ada peserta susulan.</p>
			{:else}
				<div class="table-wrap">
					<table class="tbl">
						<thead><tr><th>Nama</th><th>NIM</th><th>Alasan</th><th>Aksi</th></tr></thead>
						<tbody>
							{#each susulanList as s}
								<tr>
									<td>{s.mahasiswa?.nama ?? s.mahasiswa_id}</td>
									<td>{s.mahasiswa?.nim ?? '-'}</td>
									<td>{s.alasan}</td>
									<td><button class="text-state-error hover:underline" onclick={() => removeSusulan(s.mahasiswa_id)}>Hapus</button></td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	</div>
{/if}
