<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { api } from '$lib/api';
	import { labelJenis } from '$lib/utils';

	interface Stat {
		online_sekarang: number;
		total_mahasiswa: number;
		total_asisten: number;
		sudah_register: number;
		belum_register: number;
		per_kelas_shift: { nama_kelas: string; shift: number; jumlah: number }[];
		sesi_aktif: {
			aktivasi_sesi_id: number; judul_sesi: string; nama_kelas: string; shift: number;
			courses: { course_id: number; jenis: string; is_open: boolean; selesai: number; sedang: number; belum: number }[];
		}[];
	}

	let stat = $state<Stat | null>(null);
	let online = $state(0);
	let err = $state('');
	let poll: ReturnType<typeof setInterval>;

	async function loadStat() {
		try { stat = await api.get<Stat>('/api/admin/dashboard'); }
		catch (e) { err = (e as Error).message; }
	}
	async function loadOnline() {
		try {
			const res = await api.get<{ total: number }>('/api/admin/dashboard/online');
			online = res.total;
		} catch { /* ignore */ }
	}

	onMount(async () => {
		await loadStat();
		await loadOnline();
		poll = setInterval(loadOnline, 10000);
	});
	onDestroy(() => clearInterval(poll));
</script>

<h1 class="mb-4 text-2xl">Dashboard Admin</h1>

{#if err}<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>{/if}

{#if stat}
	<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
		<div class="card">
			<p class="text-sm text-ink-caption">🟢 Online Sekarang</p>
			<p class="mt-1 text-3xl font-bold text-state-success">{online}</p>
			<p class="text-xs text-ink-caption">real-time (registry server)</p>
		</div>
		<div class="card">
			<p class="text-sm text-ink-caption">Total Mahasiswa</p>
			<p class="mt-1 text-3xl font-bold text-primary">{stat.total_mahasiswa}</p>
		</div>
		<div class="card">
			<p class="text-sm text-ink-caption">Sudah Register</p>
			<p class="mt-1 text-3xl font-bold text-state-info">{stat.sudah_register}</p>
			<p class="text-xs text-ink-caption">Belum: {stat.belum_register}</p>
		</div>
		<div class="card">
			<p class="text-sm text-ink-caption">Total Asisten</p>
			<p class="mt-1 text-3xl font-bold text-ink-heading">{stat.total_asisten}</p>
		</div>
	</div>

	<h2 class="mb-3 mt-6 text-xl">Mahasiswa per Kelas &amp; Shift</h2>
	<div class="table-wrap max-w-xl">
		<table class="tbl">
			<thead><tr><th>Kelas</th><th>Shift</th><th>Jumlah</th></tr></thead>
			<tbody>
				{#each stat.per_kelas_shift as r}
					<tr><td>{r.nama_kelas}</td><td>{r.shift}</td><td>{r.jumlah}</td></tr>
				{/each}
			</tbody>
		</table>
	</div>

	<h2 class="mb-3 mt-6 text-xl">Sesi Aktif &amp; Progress</h2>
	{#if stat.sesi_aktif.length === 0}
		<p class="text-ink-caption">Tidak ada sesi aktif.</p>
	{:else}
		<div class="grid gap-4 md:grid-cols-2">
			{#each stat.sesi_aktif as s}
				<div class="card">
					<h3 class="text-lg">{s.judul_sesi}</h3>
					<p class="text-sm text-ink-caption">{s.nama_kelas} · Shift {s.shift}</p>
					<div class="mt-3 space-y-2">
						{#each s.courses as c}
							<div class="rounded-lg border border-gray-100 p-2 text-sm">
								<div class="flex justify-between">
									<span>{labelJenis(c.jenis)}</span>
									<span class="badge {c.is_open ? 'bg-state-success-bg text-state-success' : 'bg-gray-100 text-ink-caption'}">{c.is_open ? 'Terbuka' : 'Tertutup'}</span>
								</div>
								<div class="mt-1 text-xs text-ink-caption">Selesai {c.selesai} · Sedang {c.sedang} · Belum {c.belum}</div>
							</div>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{/if}
{:else if !err}
	<p class="text-ink-caption">Memuat…</p>
{/if}
