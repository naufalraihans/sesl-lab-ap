<script lang="ts">
	import { onMount } from 'svelte';
	import { swrGet } from '$lib/cache';
	import { labelJenis, labelStatus, statusBadgeClass } from '$lib/utils';
	import { Lock, Search, ChevronDown, ChevronUp, ArrowRight, Filter, BookOpen, Layers } from 'lucide-svelte';
	import type { SesiUserItem } from '$lib/types';

	let sesi = $state<SesiUserItem[]>([]);
	let err = $state('');
	let loading = $state(true);

	// Filter state
	let searchQuery = $state('');
	let filterMode = $state<'semua' | 'aktif' | 'terkunci'>('semua');

	// Expandable state
	let expandedIds = $state<Set<number>>(new Set());

	onMount(() => {
		swrGet<SesiUserItem[]>('/api/praktikum/sesi', (v) => {
			sesi = v ?? [];
			loading = false;
		}).catch((e) => {
			err = (e as Error).message;
			loading = false;
		});
	});

	let filteredSesi = $derived(
		sesi.filter((s) => {
			const q = searchQuery.toLowerCase().trim();
			const matchesSearch = q === '' || s.judul.toLowerCase().includes(q) || s.deskripsi.toLowerCase().includes(q);
			const matchesFilter =
				filterMode === 'semua' ||
				(filterMode === 'aktif' && s.aktif) ||
				(filterMode === 'terkunci' && !s.aktif);
			return matchesSearch && matchesFilter;
		})
	);

	function toggleExpand(id: number) {
		const s = new Set(expandedIds);
		s.has(id) ? s.delete(id) : s.add(id);
		expandedIds = s;
	}

	let countAktif = $derived(sesi.filter((s) => s.aktif).length);
	let countTerkunci = $derived(sesi.filter((s) => !s.aktif).length);
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex flex-col gap-1">
		<div class="flex items-center gap-3">
			<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
				<Layers size={20} />
			</div>
			<div>
				<h1 class="text-2xl font-bold text-ink-heading">Daftar Sesi Praktikum</h1>
				<p class="text-sm text-ink-caption">{sesi.length} sesi tersedia · {countAktif} aktif · {countTerkunci} terkunci</p>
			</div>
		</div>
	</div>

	{#if loading}
		<p class="text-ink-caption">Memuat…</p>
	{:else if err}
		<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
	{:else}
		<!-- Filter Controls -->
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
			<div class="relative flex-1 max-w-sm">
				<Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-ink-caption" />
				<input
					type="text"
					placeholder="Cari sesi..."
					bind:value={searchQuery}
					class="input pl-9 w-full"
				/>
			</div>
			<div class="flex items-center gap-2">
				<Filter size={14} class="text-ink-caption" />
				<div class="flex rounded-lg border border-gray-200 overflow-hidden">
					<button
						class="px-3 py-1.5 text-xs font-semibold transition-colors {filterMode === 'semua' ? 'bg-primary text-white' : 'bg-white text-ink-body hover:bg-gray-50'}"
						onclick={() => (filterMode = 'semua')}
					>Semua ({sesi.length})</button>
					<button
						class="px-3 py-1.5 text-xs font-semibold transition-colors border-l border-gray-200 {filterMode === 'aktif' ? 'bg-state-success text-white' : 'bg-white text-ink-body hover:bg-gray-50'}"
						onclick={() => (filterMode = 'aktif')}
					>Aktif ({countAktif})</button>
					<button
						class="px-3 py-1.5 text-xs font-semibold transition-colors border-l border-gray-200 {filterMode === 'terkunci' ? 'bg-ink-caption text-white' : 'bg-white text-ink-body hover:bg-gray-50'}"
						onclick={() => (filterMode = 'terkunci')}
					>Terkunci ({countTerkunci})</button>
				</div>
			</div>
		</div>

		<!-- Results -->
		{#if filteredSesi.length === 0}
			<div class="rounded-xl border border-dashed border-gray-300 bg-gray-50 p-10 text-center">
				<p class="text-sm font-semibold text-ink-caption">Tidak ada sesi yang sesuai filter.</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each filteredSesi as s}
					<div class="rounded-2xl border bg-white shadow-sm flex flex-col justify-between overflow-hidden transition-all duration-300 hover:-translate-y-1 hover:shadow-md hover:border-primary/20 {s.aktif ? 'border-slate-200 border-l-4 border-l-primary' : 'border-slate-200 border-l-4 border-l-slate-300 bg-slate-50/40 opacity-80' }">
						<div class="p-6 flex flex-col flex-grow min-w-0">
							<!-- Top header info -->
							<div class="flex items-start justify-between gap-3 mb-4">
								<div class="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-xl {s.aktif ? 'bg-primary/10 text-primary border border-primary/20' : 'bg-slate-100 text-slate-400 border border-slate-200'}">
									{#if s.aktif}
										<BookOpen size={20} />
									{:else}
										<Lock size={20} />
									{/if}
								</div>
								
								<div class="flex flex-wrap gap-1.5 justify-end">
									{#if s.aktif}
										<span class="badge bg-primary/10 text-primary border-primary/20 text-[10px] px-2 py-0.5">Aktif</span>
									{:else}
										<span class="badge inline-flex items-center gap-1 bg-slate-100 text-slate-400 border border-slate-200 text-[10px] px-2 py-0.5"><Lock size={10} /> Terkunci</span>
									{/if}
									{#if s.susulan}
										<span class="badge bg-state-warning-bg text-state-warning text-[10px] px-2 py-0.5">Susulan</span>
									{/if}
								</div>
							</div>

							<!-- Title & Description -->
							<h3 class="text-base font-bold text-ink-heading mb-1.5 line-clamp-2">{s.judul}</h3>
							<p class="text-xs text-ink-caption mb-5 line-clamp-3 leading-relaxed">{s.deskripsi}</p>

							<!-- Courses visual listing -->
							{#if s.courses && s.courses.length > 0}
								<div class="mt-auto pt-4 border-t border-slate-100 space-y-2">
									<h4 class="text-[10px] font-bold text-ink-caption uppercase tracking-wider mb-2">Daftar Modul</h4>
									{#each s.courses as c}
										<div class="flex items-center justify-between rounded-xl border border-slate-100/80 bg-slate-50/50 px-3 py-2 text-xs">
											<div class="flex items-center gap-2.5 min-w-0">
												<span class="text-sm flex-shrink-0">
													{#if c.jenis === 'pretest'}📝{:else if c.jenis === 'posttest'}✅{:else if c.jenis === 'keterampilan'}🧪{:else}📋{/if}
												</span>
												<span class="font-semibold text-ink-heading truncate">{labelJenis(c.jenis)}</span>
											</div>
											<span class="badge text-[9px] px-1.5 py-0.5 flex-shrink-0 {statusBadgeClass(c.status)}">{labelStatus(c.status)}{c.is_open ? '' : ' · terkunci'}</span>
										</div>
									{/each}
								</div>
							{/if}
						</div>

						<!-- Card Footer CTA Action -->
						<div class="px-6 pb-6 pt-2 flex-shrink-0">
							{#if s.aktif}
								<a href={`/praktikum/sesi/${s.sesi_id}`} class="btn-primary w-full text-xs py-2 shadow-sm">
									Masuk Sesi <ArrowRight size={14} />
								</a>
							{:else}
								<button class="btn w-full text-xs py-2 bg-slate-100 text-slate-400 border border-slate-200 cursor-not-allowed inline-flex items-center justify-center gap-1.5" disabled>
									<Lock size={12} /> Belum Diaktifkan
								</button>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	{/if}
</div>
