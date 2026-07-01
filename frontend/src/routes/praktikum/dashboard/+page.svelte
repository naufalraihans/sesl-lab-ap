<script lang="ts">
	import { onMount } from 'svelte';
	import { swrGet } from '$lib/cache';
	import { labelJenis, labelStatus, statusBadgeClass } from '$lib/utils';
	import type { SesiUserItem } from '$lib/types';

	interface Dashboard {
		profil: { nama: string; nim: string; nama_kelas?: string; shift?: number };
		jadwal?: { hari: string; jam_mulai: string; jam_selesai: string; keterangan: string } | null;
		sesi_aktif: SesiUserItem[];
		riwayat_nilai: { sesi_judul: string; jenis: string; status: string; total_nilai?: number | null }[];
	}

	let data = $state<Dashboard | null>(null);
	let err = $state('');

	onMount(() => {
		// SWR: tampil instan dari cache (bila ada) lalu disegarkan diam-diam.
		swrGet<Dashboard>('/api/praktikum/dashboard', (v) => {
			data = v;
		}).catch((e) => {
			err = (e as Error).message;
		});
	});
</script>

<!-- Latar aurora glassmorphism -->
<div class="aurora" aria-hidden="true">
	<div class="aurora-blob aurora-blob-1"></div>
	<div class="aurora-blob aurora-blob-2"></div>
	<div class="aurora-blob aurora-blob-3"></div>
</div>

<div class="relative z-10">
	<h1 class="mb-1 text-2xl">Dashboard</h1>
	<p class="mb-5 text-sm text-ink-caption">Selamat datang kembali di Lab Algoritma &amp; Pemrograman.</p>

	{#if err}
		<p class="glass-card p-4 text-state-error">{err}</p>
	{:else if !data}
		<p class="text-ink-caption">Memuat…</p>
	{:else}
		<div class="grid gap-5 md:grid-cols-3">
			<div class="glass-card glass-hover p-6 md:col-span-1">
				<h2 class="text-lg">Profil</h2>
				<dl class="mt-3 space-y-2 text-sm">
					<div class="flex justify-between border-b border-white/40 pb-2"><dt class="text-ink-caption">Nama</dt><dd class="font-medium text-ink-heading">{data.profil.nama}</dd></div>
					<div class="flex justify-between border-b border-white/40 pb-2"><dt class="text-ink-caption">NIM</dt><dd class="font-medium text-ink-heading">{data.profil.nim}</dd></div>
					<div class="flex justify-between border-b border-white/40 pb-2"><dt class="text-ink-caption">Kelas</dt><dd class="font-medium text-ink-heading">{data.profil.nama_kelas ?? '-'}</dd></div>
					<div class="flex justify-between"><dt class="text-ink-caption">Shift</dt><dd class="font-medium text-ink-heading">{data.profil.shift ?? '-'}</dd></div>
				</dl>
			</div>

			<div class="glass-card glass-hover p-6 md:col-span-2">
				<h2 class="text-lg">Jadwal Praktikum</h2>
				{#if data.jadwal}
					<p class="mt-3 text-lg font-semibold text-ink-heading">{data.jadwal.hari}, {data.jadwal.jam_mulai} – {data.jadwal.jam_selesai}</p>
					{#if data.jadwal.keterangan}<p class="mt-1 text-sm text-ink-caption">{data.jadwal.keterangan}</p>{/if}
				{:else}
					<p class="mt-3 text-sm text-ink-caption">Belum ada jadwal.</p>
				{/if}
			</div>
		</div>

		<h2 class="mb-3 mt-7 text-xl">Sesi Aktif Sekarang</h2>
		{#if data.sesi_aktif.length === 0}
			<div class="glass-card p-5"><p class="text-ink-caption">Tidak ada sesi aktif untuk kelas &amp; shift Anda.</p></div>
		{:else}
			<div class="grid gap-5 md:grid-cols-2">
				{#each data.sesi_aktif as s}
					<div class="glass-card glass-hover p-6">
						<div class="flex items-center justify-between">
							<h3 class="text-lg">{s.judul}</h3>
							{#if s.susulan}<span class="badge bg-state-warning-bg text-state-warning">Susulan</span>{/if}
						</div>
						<div class="mt-4 space-y-2">
							{#each s.courses as c}
								<div class="flex items-center justify-between rounded-xl border border-white/50 bg-white/40 px-3 py-2 text-sm">
									<span>{labelJenis(c.jenis)}</span>
									<span class="badge {statusBadgeClass(c.status)}">{labelStatus(c.status)}{c.is_open ? '' : ' · terkunci'}</span>
								</div>
							{/each}
						</div>
						<a href={`/praktikum/sesi/${s.sesi_id}`} class="btn-primary mt-5 w-full">Masuk Sesi</a>
					</div>
				{/each}
			</div>
		{/if}

		<h2 class="mb-3 mt-7 text-xl">Riwayat Nilai</h2>
		{#if data.riwayat_nilai.length === 0}
			<div class="glass-card p-5"><p class="text-ink-caption">Belum ada nilai.</p></div>
		{:else}
			<div class="glass-card overflow-x-auto">
				<table class="tbl">
					<thead><tr><th>Sesi</th><th>Course</th><th>Status</th><th>Nilai</th></tr></thead>
					<tbody>
						{#each data.riwayat_nilai as r}
							<tr>
								<td>{r.sesi_judul}</td>
								<td>{labelJenis(r.jenis)}</td>
								<td>{labelStatus(r.status)}</td>
								<td>{r.total_nilai ?? '-'}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	{/if}
</div>
