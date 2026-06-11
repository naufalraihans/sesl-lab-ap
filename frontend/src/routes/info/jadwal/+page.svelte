<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { FileText } from 'lucide-svelte';
	import type { Jadwal, User, AmpuanKelompok, Kelas } from '$lib/types';

	let jadwal = $state<Jadwal[]>([]);
	let config = $state<{ mode: string; gdrive_url: string }>({ mode: 'internal', gdrive_url: '' });
	let loading = $state(true);
	let err = $state('');

	let selectedKelas = $state<Kelas | null>(null);
	let mhsList = $state<User[]>([]);
	let ampuanList = $state<AmpuanKelompok[]>([]);
	let loadingMhs = $state(false);

	onMount(async () => {
		try {
			config = await api.get('/api/info/jadwal/config');
			if (config.mode !== 'gdrive') {
				jadwal = (await api.get<Jadwal[]>('/api/info/jadwal')) ?? [];
			}
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});

	let kelasList = $derived(
		jadwal.reduce<Kelas[]>((acc, j) => {
			if (j.kelas && !acc.find((k) => k.id === j.kelas!.id)) acc.push(j.kelas);
			return acc;
		}, [])
	);

	async function showKelas(k: Kelas) {
		selectedKelas = k;
		loadingMhs = true;
		try {
			const res = await api.get<{ mahasiswa: User[]; ampuan: AmpuanKelompok[] }>(
				`/api/info/kelas/${k.id}/mahasiswa`
			);
			mhsList = res?.mahasiswa ?? [];
			ampuanList = res?.ampuan ?? [];
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loadingMhs = false;
		}
	}

	function ampuanForKelompok(kel: string): string {
		const a = ampuanList.find((x) => x.kelompok === kel);
		return a?.asisten?.nama ?? '-';
	}

	let kelompokGroups = $derived(
		mhsList.reduce<Record<string, User[]>>((acc, m) => {
			const kel = m.kelompok ?? '(Belum ada)';
			if (!acc[kel]) acc[kel] = [];
			acc[kel].push(m);
			return acc;
		}, {})
	);
</script>

<h1 class="mb-4 text-2xl">Jadwal Praktikum</h1>

{#if loading}
	<p class="text-ink-caption">Memuat…</p>
{:else if err}
	<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
{:else if config.mode === 'gdrive' && config.gdrive_url}
	<a href={config.gdrive_url} target="_blank" rel="noopener" class="btn-primary"><FileText size={16} /> Buka Jadwal (Google Drive)</a>
{:else if jadwal.length === 0}
	<p class="text-ink-caption">Belum ada jadwal yang dipublikasikan.</p>
{:else}
	<div class="table-wrap">
		<table class="tbl">
			<thead>
				<tr><th>Kelas</th><th>Shift</th><th>Hari</th><th>Jam</th><th>Keterangan</th><th></th></tr>
			</thead>
			<tbody>
				{#each jadwal as j}
					<tr>
						<td>{j.kelas?.nama_kelas ?? j.kelas_id}</td>
						<td>Shift {j.shift}</td>
						<td>{j.hari}</td>
						<td>{j.jam_mulai} – {j.jam_selesai}</td>
						<td>{j.keterangan}</td>
						<td>
							{#if j.kelas}
								<button
									class="text-sm text-primary hover:underline"
									onclick={() => showKelas(j.kelas!)}
								>Lihat Mahasiswa</button>
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	{#if selectedKelas}
		<hr class="my-6 border-gray-200" />
		<h2 class="mb-3 text-xl">Daftar Mahasiswa — {selectedKelas.nama_kelas}</h2>

		{#if loadingMhs}
			<p class="text-ink-caption">Memuat…</p>
		{:else if mhsList.length === 0}
			<p class="text-ink-caption">Belum ada mahasiswa di kelas ini.</p>
		{:else}
			{#each Object.entries(kelompokGroups) as [kel, members]}
				<div class="mb-4">
					<div class="flex items-center gap-3 mb-2">
						<h3 class="text-lg">Kelompok {kel}</h3>
						<span class="badge bg-surface-soft text-ink-caption">Asisten: {ampuanForKelompok(kel)}</span>
					</div>
					<div class="table-wrap">
						<table class="tbl">
							<thead><tr><th>No</th><th>NIM</th><th>Nama</th><th>Shift</th></tr></thead>
							<tbody>
								{#each members as m, i}
									<tr>
										<td>{i + 1}</td>
										<td>{m.nim}</td>
										<td>{m.nama}</td>
										<td>{m.shift ?? '-'}</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			{/each}
		{/if}
	{/if}
{/if}
