<script lang="ts">
	import { onMount } from 'svelte';
	import { swrGet } from '$lib/cache';
	import { labelJenis, labelStatus, statusBadgeClass } from '$lib/utils';
	import type { SesiUserItem } from '$lib/types';
	import { Calendar, User, GraduationCap, ArrowRight, Trophy, Bell, BookOpen } from 'lucide-svelte';

	interface Dashboard {
		profil: { nama: string; nim: string; nama_kelas?: string; shift?: number };
		jadwal?: { hari: string; jam_mulai: string; jam_selesai: string; keterangan: string } | null;
		sesi_aktif: SesiUserItem[];
		riwayat_nilai: { sesi_judul: string; jenis: string; status: string; total_nilai?: number | null }[];
	}

	let data = $state<Dashboard | null>(null);
	let err = $state('');

	// Tautan CTA utama: menuju sesi aktif pertama bila ada, jika tidak ke daftar sesi.
	let ctaHref = $derived(
		data && data.sesi_aktif.length > 0
			? `/praktikum/sesi/${data.sesi_aktif[0].sesi_id}`
			: '/praktikum/sesi'
	);

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

<div class="relative z-10 mx-auto max-w-6xl">
	{#if err}
		<p class="glass-card p-4 text-state-error">{err}</p>
	{:else if !data}
		<p class="text-ink-caption">Memuat…</p>
	{:else}
		<!-- Bento grid -->
		<div class="grid auto-rows-[minmax(158px,auto)] grid-cols-1 gap-[18px] sm:grid-cols-2 lg:grid-cols-4">
			<!-- Hero (2x2) -->
			<section
				class="glass-hero fade-up flex flex-col justify-between p-8 sm:col-span-2 lg:row-span-2"
			>
				<div>
					<span class="badge bg-primary/10 text-primary">Dashboard</span>
					<h1 class="mt-4 font-display text-3xl font-extrabold leading-tight tracking-tight text-ink-heading text-balance md:text-4xl">
						Halo, {data.profil.nama.split(' ')[0]}
					</h1>
					<p class="mt-3 max-w-md text-base leading-relaxed text-ink-body">
						Selamat datang kembali di Lab Algoritma &amp; Pemrograman. Pantau sesi aktif, jadwal, dan progres nilaimu di satu tempat.
					</p>
					<div class="mt-5 flex flex-wrap gap-2">
						<span class="badge border border-white/60 bg-white/50 text-ink-body">NIM {data.profil.nim}</span>
						<span class="badge border border-white/60 bg-white/50 text-ink-body">{data.profil.nama_kelas ?? 'Kelas -'}</span>
						<span class="badge border border-white/60 bg-white/50 text-ink-body">Shift {data.profil.shift ?? '-'}</span>
					</div>
				</div>
				<div class="mt-7 flex flex-wrap gap-3">
					<a href={ctaHref} class="btn-primary">
						{data.sesi_aktif.length > 0 ? 'Masuk ke Sesi' : 'Lihat Daftar Sesi'}
						<ArrowRight size={17} />
					</a>
					<a href="/praktikum/profil" class="btn-outline border-white/60 bg-white/40">Profil Saya</a>
				</div>
			</section>

			<!-- Jadwal -->
			<div class="glass-card glass-hover fade-up flex flex-col p-6" style="animation-delay:.06s">
				<span class="glass-badge h-12 w-12"><Calendar size={22} /></span>
				<h3 class="mt-4 font-display text-base font-bold text-ink-heading">Jadwal Praktikum</h3>
				{#if data.jadwal}
					<p class="mt-1 text-sm font-semibold text-ink-heading">{data.jadwal.hari}, {data.jadwal.jam_mulai}–{data.jadwal.jam_selesai}</p>
					{#if data.jadwal.keterangan}<p class="mt-1 text-xs leading-relaxed text-ink-caption">{data.jadwal.keterangan}</p>{/if}
				{:else}
					<p class="mt-1 text-sm text-ink-caption">Belum ada jadwal.</p>
				{/if}
			</div>

			<!-- Ringkasan sesi/nilai -->
			<div class="glass-card glass-hover fade-up flex flex-col p-6" style="animation-delay:.12s">
				<span class="glass-badge h-12 w-12"><BookOpen size={22} /></span>
				<h3 class="mt-4 font-display text-base font-bold text-ink-heading">Sesi Aktif</h3>
				<p class="mt-1 text-3xl font-extrabold text-primary">{data.sesi_aktif.length}</p>
				<p class="text-xs text-ink-caption">sesi tersedia untuk kelas &amp; shift kamu</p>
			</div>

			<!-- Profil (wide) -->
			<section class="glass-card fade-up b-wide flex items-start gap-4 p-6 sm:col-span-2 lg:col-span-4" style="animation-delay:.18s">
				<span class="glass-badge-soft h-11 w-11 flex-shrink-0"><User size={20} /></span>
				<div class="min-w-0">
					<h2 class="font-display text-base font-bold text-ink-heading">Profil Kamu</h2>
					<div class="mt-3 grid grid-cols-2 gap-x-6 gap-y-3 text-sm sm:grid-cols-4">
						<div><dt class="text-xs text-ink-caption">Nama</dt><dd class="font-medium text-ink-heading">{data.profil.nama}</dd></div>
						<div><dt class="text-xs text-ink-caption">NIM</dt><dd class="font-medium text-ink-heading">{data.profil.nim}</dd></div>
						<div><dt class="text-xs text-ink-caption">Kelas</dt><dd class="font-medium text-ink-heading">{data.profil.nama_kelas ?? '-'}</dd></div>
						<div><dt class="text-xs text-ink-caption">Shift</dt><dd class="font-medium text-ink-heading">{data.profil.shift ?? '-'}</dd></div>
					</div>
				</div>
			</section>
		</div>

		<!-- Sesi Aktif Sekarang -->
		<section class="mt-12 fade-up">
			<h2 class="font-display text-2xl font-bold tracking-tight text-ink-heading">Sesi Aktif Sekarang</h2>
			<p class="mt-1 text-sm text-ink-caption">Sesi yang sedang dibuka untuk kelas &amp; shift kamu.</p>
			{#if data.sesi_aktif.length === 0}
				<div class="glass-card mt-5 flex items-center gap-3 p-6">
					<span class="glass-badge-soft h-10 w-10"><Bell size={18} /></span>
					<p class="text-sm text-ink-caption">Tidak ada sesi aktif saat ini.</p>
				</div>
			{:else}
				<div class="mt-5 grid gap-[18px] md:grid-cols-2">
					{#each data.sesi_aktif as s}
						<div class="glass-card glass-hover flex flex-col p-6">
							<div class="flex items-center justify-between gap-2">
								<h3 class="font-display text-lg font-bold text-ink-heading">{s.judul}</h3>
								{#if s.susulan}<span class="badge bg-state-warning-bg text-state-warning">Susulan</span>{/if}
							</div>
							<div class="mt-4 space-y-2">
								{#each s.courses as c}
									<div class="flex items-center justify-between rounded-xl border border-white/50 bg-white/40 px-3 py-2 text-sm">
										<span class="text-ink-body">{labelJenis(c.jenis)}</span>
										<span class="badge {statusBadgeClass(c.status)}">{labelStatus(c.status)}{c.is_open ? '' : ' · terkunci'}</span>
									</div>
								{/each}
							</div>
							<a href={`/praktikum/sesi/${s.sesi_id}`} class="btn-primary mt-5 w-full">Masuk Sesi <ArrowRight size={16} /></a>
						</div>
					{/each}
				</div>
			{/if}
		</section>

		<!-- Riwayat Nilai -->
		<section class="mt-12 fade-up">
			<h2 class="font-display text-2xl font-bold tracking-tight text-ink-heading">Riwayat Nilai</h2>
			<p class="mt-1 text-sm text-ink-caption">Rekap nilai dari sesi yang sudah kamu ikuti.</p>
			{#if data.riwayat_nilai.length === 0}
				<div class="glass-card mt-5 flex items-center gap-3 p-6">
					<span class="glass-badge-soft h-10 w-10"><Trophy size={18} /></span>
					<p class="text-sm text-ink-caption">Belum ada nilai.</p>
				</div>
			{:else}
				<div class="glass-card mt-5 overflow-x-auto p-0">
					<table class="tbl">
						<thead><tr><th>Sesi</th><th>Course</th><th>Status</th><th>Nilai</th></tr></thead>
						<tbody>
							{#each data.riwayat_nilai as r}
								<tr>
									<td class="font-medium text-ink-heading">{r.sesi_judul}</td>
									<td>{labelJenis(r.jenis)}</td>
									<td>{labelStatus(r.status)}</td>
									<td class="font-semibold text-primary">{r.total_nilai ?? '-'}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</section>

		<!-- CTA band -->
		<section class="cta-band fade-up mt-12 p-10 text-center md:p-12">
			<h2 class="font-display text-2xl font-extrabold tracking-tight text-white md:text-3xl">Siap lanjut praktikum?</h2>
			<p class="mx-auto mt-3 max-w-md text-sm leading-relaxed text-white/85">
				Buka sesi aktifmu, kerjakan pretest, keterampilan, dan posttest langsung dari browser.
			</p>
			<div class="mt-6 flex flex-wrap justify-center gap-3">
				<a href={ctaHref} class="btn inline-flex bg-white font-bold text-primary hover:-translate-y-0.5 hover:shadow-lg">
					{data.sesi_aktif.length > 0 ? 'Masuk ke Sesi' : 'Lihat Daftar Sesi'}
					<ArrowRight size={16} />
				</a>
				<a href="/praktikum/profil" class="btn inline-flex border border-white/40 bg-white/15 font-semibold text-white hover:-translate-y-0.5">
					Profil Saya
				</a>
			</div>
		</section>
	{/if}
</div>
