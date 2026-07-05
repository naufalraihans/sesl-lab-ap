<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { labelJenis, labelStatus, statusBadgeClass } from '$lib/utils';
	import { Lock, CalendarClock, ChevronDown, ArrowLeft, Clock, CheckCircle2, BookOpen, FlaskConical, PenLine, ClipboardCheck } from 'lucide-svelte';
	import type { SesiUserItem } from '$lib/types';

	let sesiId = $derived(Number($page.params.sesiId));
	let sesi = $state<SesiUserItem | null>(null);
	let err = $state('');
	let loading = $state(true);

	// Expandable state
	let expandedCourses = $state<Set<number>>(new Set());

	function jenisPath(jenis: string): string {
		if (jenis === 'ujian_praktik') return 'ujian';
		return jenis;
	}

	function jenisIcon(jenis: string) {
		switch (jenis) {
			case 'pretest': return PenLine;
			case 'posttest': return ClipboardCheck;
			case 'keterampilan': return FlaskConical;
			default: return BookOpen;
		}
	}

	function toggleCourse(id: number) {
		const s = new Set(expandedCourses);
		s.has(id) ? s.delete(id) : s.add(id);
		expandedCourses = s;
	}

	onMount(async () => {
		try {
			const list = (await api.get<SesiUserItem[]>('/api/praktikum/sesi')) ?? [];
			sesi = list.find((s) => s.sesi_id === sesiId) ?? null;
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});

	let completedCount = $derived(sesi ? sesi.courses.filter((c) => c.status === 'selesai').length : 0);
	let totalCount = $derived(sesi ? sesi.courses.length : 0);
	let progressPercent = $derived(totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0);
</script>

<div class="space-y-6">
	<!-- Back button -->
	<a href="/praktikum/sesi" class="inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-semibold text-ink-caption hover:bg-gray-100 hover:text-ink-heading transition-colors">
		<ArrowLeft size={16} /> Kembali ke daftar sesi
	</a>

	{#if loading}
		<p class="text-ink-caption">Memuat…</p>
	{:else if err}
		<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
	{:else if !sesi || !sesi.aktif}
		<div class="flex flex-col items-center rounded-2xl border border-dashed border-gray-300 bg-surface-soft p-10 text-center">
			<div class="mb-3 grid h-14 w-14 place-items-center rounded-full bg-state-warning-bg text-state-warning">
				<CalendarClock size={28} />
			</div>
			<h2 class="text-lg font-bold text-ink-heading">Sesi belum dibuka</h2>
			<p class="mt-1 max-w-sm text-sm text-ink-caption">
				Sesi ini belum diaktifkan oleh asisten untuk kelas/shift Anda. Silakan cek lagi nanti.
			</p>
			<a href="/praktikum/sesi" class="btn-outline mt-5">Kembali ke daftar sesi</a>
		</div>
	{:else}
		<!-- Sesi Header -->
		<div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm">
			<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div>
					<h1 class="text-2xl font-bold text-ink-heading">{sesi.judul}</h1>
					<p class="mt-1 text-sm text-ink-caption">{sesi.deskripsi}</p>
					{#if sesi.susulan}<span class="badge mt-2 bg-state-warning-bg text-state-warning">Sesi Susulan</span>{/if}
				</div>
				<div class="flex items-center gap-2">
					<span class="badge bg-state-success-bg text-state-success">Aktif</span>
				</div>
			</div>

			<!-- Progress Bar -->
			<div class="mt-5">
				<div class="flex items-center justify-between mb-2">
					<span class="text-xs font-bold text-ink-caption uppercase tracking-wider">Progress Course</span>
					<span class="text-xs font-bold text-primary">{completedCount}/{totalCount} selesai ({progressPercent}%)</span>
				</div>
				<div class="h-2 w-full rounded-full bg-gray-100 overflow-hidden">
					<div
						class="h-full rounded-full bg-primary transition-all duration-500 ease-out"
						style="width: {progressPercent}%"
					></div>
				</div>
			</div>
		</div>

		<!-- Course Cards -->
		<div class="space-y-3">
			{#each sesi.courses as c}
				{@const Icon = jenisIcon(c.jenis)}
				{@const isExpanded = expandedCourses.has(c.course_id)}
				<div class="rounded-2xl border border-slate-200 bg-white shadow-sm transition-all duration-300 hover:shadow-md hover:border-primary/20 {
					c.status === 'selesai' ? 'border-l-4 border-l-emerald-500' :
					c.status === 'sedang_dikerjakan' ? 'border-l-4 border-l-amber-500 bg-amber-50/[0.01]' :
					c.is_open ? 'border-l-4 border-l-primary bg-primary/[0.01]' :
					'border-l-4 border-l-slate-300 bg-slate-50/40 opacity-70'
				}">
					<!-- Course Header -->
					<button
						class="flex w-full items-center justify-between gap-4 p-5 text-left transition-colors hover:bg-gray-50/50 rounded-2xl"
						onclick={() => toggleCourse(c.course_id)}
					>
						<div class="flex items-center gap-4">
							<div class="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-xl {
								c.status === 'selesai' ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' :
								c.status === 'sedang_dikerjakan' ? 'bg-amber-50 text-amber-600 border border-amber-100' :
								c.is_open ? 'bg-primary/10 text-primary border border-primary/20' :
								'bg-slate-100 text-slate-400 border border-slate-200'
							}">
								{#if c.status === 'selesai'}
									<CheckCircle2 size={20} />
								{:else}
									<Icon size={20} />
								{/if}
							</div>
							<div>
								<h3 class="text-base font-bold text-ink-heading">{labelJenis(c.jenis)}</h3>
								<div class="mt-0.5 flex items-center gap-2 text-xs text-ink-caption">
									<Clock size={12} />
									<span>{c.durasi_menit} menit</span>
								</div>
							</div>
						</div>
						<div class="flex items-center gap-3 flex-shrink-0">
							<span class="badge {statusBadgeClass(c.status)}">{labelStatus(c.status)}</span>
							{#if !c.is_open && c.status !== 'selesai'}
								<span class="badge inline-flex items-center gap-1 bg-slate-100 text-slate-400 border border-slate-200"><Lock size={10} /> Terkunci</span>
							{/if}
							<div class="text-ink-caption transition-transform duration-200 {isExpanded ? 'rotate-180' : ''}">
								<ChevronDown size={18} />
							</div>
						</div>
					</button>

					<!-- Expandable Content -->
					{#if isExpanded}
						<div class="border-t border-gray-100 px-5 pb-5 pt-4">
							<div class="grid grid-cols-2 gap-4 sm:grid-cols-4 mb-4">
								<div class="rounded-xl bg-slate-50 border border-slate-100 p-3">
									<p class="text-[10px] font-bold text-ink-caption uppercase tracking-wider">Durasi</p>
									<p class="mt-1 text-sm font-bold text-ink-heading">{c.durasi_menit} menit</p>
								</div>
								<div class="rounded-xl bg-slate-50 border border-slate-100 p-3">
									<p class="text-[10px] font-bold text-ink-caption uppercase tracking-wider">Status</p>
									<p class="mt-1 text-sm font-bold text-ink-heading">{labelStatus(c.status)}</p>
								</div>
								<div class="rounded-xl bg-slate-50 border border-slate-100 p-3">
									<p class="text-[10px] font-bold text-ink-caption uppercase tracking-wider">Akses</p>
									<p class="mt-1 text-sm font-bold {c.is_open ? 'text-emerald-600' : 'text-red-500'}">{c.is_open ? 'Terbuka' : 'Tertutup'}</p>
								</div>
								<div class="rounded-xl bg-slate-50 border border-slate-100 p-3">
									<p class="text-[10px] font-bold text-ink-caption uppercase tracking-wider">Jenis</p>
									<p class="mt-1 text-sm font-bold text-ink-heading">{labelJenis(c.jenis)}</p>
								</div>
							</div>

							<div class="flex justify-end gap-2">
								{#if c.is_open && c.status !== 'selesai'}
									<a
										href={`/praktikum/sesi/${sesiId}/${jenisPath(c.jenis)}?aktivasi=${sesi.aktivasi_sesi_id}&course=${c.course_id}`}
										class="btn-primary"
									>Kerjakan</a>
								{:else if c.status === 'selesai'}
									<a
										href={`/praktikum/sesi/${sesiId}/${jenisPath(c.jenis)}?aktivasi=${sesi.aktivasi_sesi_id}&course=${c.course_id}`}
										class="btn-outline"
									>Lihat Hasil</a>
								{:else}
									<span class="badge inline-flex items-center gap-1 bg-gray-100 text-ink-caption px-4 py-2"><Lock size={12} /> Belum dapat diakses</span>
								{/if}
							</div>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
