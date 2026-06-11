<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { labelJenis, labelStatus, statusBadgeClass } from '$lib/utils';
	import { Lock, CalendarClock } from 'lucide-svelte';
	import type { SesiUserItem } from '$lib/types';

	let sesiId = $derived(Number($page.params.sesiId));
	let sesi = $state<SesiUserItem | null>(null);
	let err = $state('');
	let loading = $state(true);

	function jenisPath(jenis: string): string {
		if (jenis === 'ujian_praktik') return 'ujian';
		return jenis; // pretest | posttest | keterampilan
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
</script>

<a href="/praktikum/sesi" class="text-sm">← Kembali ke daftar sesi</a>

{#if loading}
	<p class="mt-4 text-ink-caption">Memuat…</p>
{:else if err}
	<p class="mt-4 rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
{:else if !sesi || !sesi.aktif}
	<div class="mt-6 flex flex-col items-center rounded-2xl border border-dashed border-gray-300 bg-surface-soft p-10 text-center">
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
	<h1 class="mb-1 mt-3 text-2xl">{sesi.judul}</h1>
	<p class="mb-5 text-ink-caption">{sesi.deskripsi}</p>

	<div class="space-y-3">
		{#each sesi.courses as c}
			<div class="card flex items-center justify-between">
				<div>
					<h3 class="text-lg">{labelJenis(c.jenis)}</h3>
					<p class="text-sm text-ink-caption">Durasi {c.durasi_menit} menit</p>
					<span class="badge mt-1 {statusBadgeClass(c.status)}">{labelStatus(c.status)}</span>
				</div>
				<div class="text-right">
					{#if c.is_open && c.status !== 'selesai'}
						<a
							href={`/praktikum/sesi/${sesiId}/${jenisPath(c.jenis)}?aktivasi=${sesi.aktivasi_sesi_id}&course=${c.course_id}`}
							class="btn-primary"
						>Kerjakan</a>
					{:else if c.status === 'selesai'}
						<a
							href={`/praktikum/sesi/${sesiId}/${jenisPath(c.jenis)}?aktivasi=${sesi.aktivasi_sesi_id}&course=${c.course_id}`}
							class="btn-outline"
						>Lihat</a>
					{:else}
						<span class="badge inline-flex items-center gap-1 bg-gray-100 text-ink-caption"><Lock size={12} /> Terkunci</span>
					{/if}
				</div>
			</div>
		{/each}
	</div>
{/if}
