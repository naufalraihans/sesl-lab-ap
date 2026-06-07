<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { labelJenis, renderMath } from '$lib/utils';
	import type { Kelas } from '$lib/types';

	interface AktivasiSesi {
		id: number; sesi_praktikum_id: number; kelas_id: number; shift: number;
		sesi?: { judul_sesi: string };
		kelas?: { nama_kelas: string };
		aktivasi_courses?: { id: number; course_id: number; course?: { jenis: string; judul: string } }[];
	}
	interface RekapItem {
		jawaban_id: number; mahasiswa_id: number; nama_mahasiswa: string; nim: string;
		soal_id: number; teks_soal: string; poin: number;
		jawaban_teks: string; is_submitted: boolean;
		nilai: number | null; feedback: string | null;
	}

	let aktivasiList = $state<AktivasiSesi[]>([]);
	let err = $state(''); let msg = $state('');

	let selectedAktivasi = $state<AktivasiSesi | null>(null);
	let selectedCourseId = $state<number | null>(null);
	let rekap = $state<RekapItem[]>([]);
	let loading = $state(false);

	let aiJobId = $state<string | null>(null);
	let aiJobStatus = $state<string>('');
	let aiProcessed = $state(0);
	let aiTotal = $state(0);
	let pollInterval: ReturnType<typeof setInterval>;

	let nilaiEdits = $state<Record<number, { nilai: number; feedback: string }>>({});

	async function loadAktivasi() {
		try { aktivasiList = (await api.get<AktivasiSesi[]>('/api/admin/aktivasi')) ?? []; }
		catch (e) { err = (e as Error).message; }
	}
	onMount(loadAktivasi);

	async function selectAktivasi(a: AktivasiSesi) {
		selectedAktivasi = a; selectedCourseId = null; rekap = [];
		try {
			const detail = await api.get<AktivasiSesi>(`/api/admin/aktivasi/${a.id}`);
			if (detail) selectedAktivasi = detail;
		} catch (e) { err = (e as Error).message; }
	}

	async function loadRekap(courseId: number) {
		if (!selectedAktivasi) return;
		selectedCourseId = courseId; loading = true; err = '';
		try {
			const res = await api.get<{ items: RekapItem[] }>(
				`/api/admin/penilaian/rekap?aktivasi_sesi_id=${selectedAktivasi.id}&course_id=${courseId}`
			);
			rekap = res?.items ?? [];
			const edits: Record<number, { nilai: number; feedback: string }> = {};
			for (const r of rekap) {
				edits[r.jawaban_id] = { nilai: r.nilai ?? 0, feedback: r.feedback ?? '' };
			}
			nilaiEdits = edits;
		} catch (e) { err = (e as Error).message; }
		finally { loading = false; }
	}

	async function simpanNilai(jawabanId: number) {
		err = ''; msg = '';
		const edit = nilaiEdits[jawabanId];
		if (!edit) return;
		try {
			await api.post('/api/admin/penilaian', {
				jawaban_id: jawabanId,
				nilai: Number(edit.nilai),
				feedback: edit.feedback || null
			});
			msg = 'Nilai disimpan.';
			if (selectedCourseId) await loadRekap(selectedCourseId);
		} catch (e) { err = (e as Error).message; }
	}

	async function startAIGrading() {
		if (!selectedAktivasi || !selectedCourseId) return;
		if (!confirm('Apakah Anda yakin ingin memulai penilaian AI untuk jawaban yang belum dinilai? Proses ini membutuhkan waktu beberapa saat.')) return;
		err = ''; msg = '';
		try {
			const res = await api.post<{ job_id: string; status: string; total: number; processed: number }>(
				'/api/admin/penilaian/ai-grade/bulk',
				{ aktivasi_sesi_id: selectedAktivasi.id, course_id: selectedCourseId }
			);
			if (res && res.job_id) {
				aiJobId = res.job_id;
				aiJobStatus = res.status;
				aiProcessed = res.processed;
				aiTotal = res.total;
				pollAIGrading();
			}
		} catch (e) { err = (e as Error).message; }
	}

	function pollAIGrading() {
		if (pollInterval) clearInterval(pollInterval);
		pollInterval = setInterval(async () => {
			if (!aiJobId) {
				clearInterval(pollInterval);
				return;
			}
			try {
				const res = await api.get<{ job_id: string; status: string; total: number; processed: number; message: string }>(
					`/api/admin/jobs/${aiJobId}`
				);
				if (res) {
					aiJobStatus = res.status;
					aiProcessed = res.processed;
					aiTotal = res.total;
					
					if (res.status === 'completed' || res.status === 'failed') {
						clearInterval(pollInterval);
						msg = res.message;
						aiJobId = null; // Sembunyikan progress
						if (selectedCourseId) await loadRekap(selectedCourseId); // Refresh hasil
					}
				}
			} catch (e) {
				err = (e as Error).message;
				clearInterval(pollInterval);
				aiJobId = null;
			}
		}, 2000);
	}
</script>

<h1 class="mb-4 text-2xl">Penilaian Mahasiswa</h1>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="grid gap-4 lg:grid-cols-4">
	<div class="card">
		<h2 class="mb-3 text-lg">Pilih Aktivasi</h2>
		<div class="space-y-2">
			{#each aktivasiList as a}
				<button
					class="w-full rounded-lg border p-3 text-left text-sm transition hover:bg-surface-soft {selectedAktivasi?.id === a.id ? 'border-primary bg-surface-soft' : 'border-gray-200'}"
					onclick={() => selectAktivasi(a)}
				>
					<p class="font-medium">{a.sesi?.judul_sesi ?? `Sesi #${a.sesi_praktikum_id}`}</p>
					<p class="text-xs text-ink-caption">{a.kelas?.nama_kelas ?? a.kelas_id} · Shift {a.shift}</p>
				</button>
			{/each}
		</div>
	</div>

	<div class="lg:col-span-3">
		{#if selectedAktivasi}
			<div class="mb-4 flex flex-wrap gap-2">
				{#each selectedAktivasi.aktivasi_courses ?? [] as ac}
					<button
						class="badge cursor-pointer px-3 py-1.5 {selectedCourseId === ac.course_id ? 'bg-primary text-white' : 'bg-surface-soft text-ink-body'}"
						onclick={() => loadRekap(ac.course_id)}
					>{ac.course?.judul ?? labelJenis(ac.course?.jenis ?? '')}</button>
				{/each}
			</div>

			{#if loading}
				<p class="text-ink-caption">Memuat rekap…</p>
			{:else if selectedCourseId && rekap.length === 0}
				<p class="text-ink-caption">Belum ada jawaban yang ter-submit.</p>
			{:else if rekap.length > 0}
				{#if aiJobId}
					<div class="mb-4 rounded-lg border border-primary/20 bg-primary/5 p-4">
						<h3 class="mb-2 font-medium text-primary">AI Grading sedang berjalan...</h3>
						<div class="mb-2 h-2 w-full overflow-hidden rounded-full bg-gray-200">
							<div class="h-full bg-primary transition-all duration-500" style="width: {aiTotal > 0 ? (aiProcessed / aiTotal) * 100 : 0}%"></div>
						</div>
						<p class="text-sm text-ink-body">Status: {aiJobStatus} | Diproses: {aiProcessed} / {aiTotal} jawaban</p>
					</div>
				{:else}
					<div class="mb-4 flex items-center justify-between rounded-lg border border-gray-200 bg-surface-muted p-4">
						<div>
							<h3 class="font-medium text-ink-body">Otomatisasi Penilaian</h3>
							<p class="text-sm text-ink-caption">Gunakan AI untuk menilai otomatis jawaban yang belum memiliki nilai.</p>
						</div>
						<button class="btn-primary" onclick={startAIGrading}>
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2"><path d="m2 16 3-3 3 3"></path><path d="m2 16 3 3 3-3"></path><path d="M14 6h-4a4 4 0 0 0-4 4v10"></path><path d="M18 10a4 4 0 0 1 4 4v6"></path><path d="m22 20-3 3-3-3"></path></svg>
							Mulai AI Grading
						</button>
					</div>
				{/if}

				<div class="space-y-4">
					{#each rekap as r}
						<div class="card">
							<div class="flex flex-wrap items-center gap-3 text-sm">
								<span class="font-medium">{r.nama_mahasiswa}</span>
								<span class="text-ink-caption">{r.nim}</span>
								<span class="badge {r.is_submitted ? 'bg-state-success-bg text-state-success' : 'bg-state-warning-bg text-state-warning'}">
									{r.is_submitted ? 'Submitted' : 'Belum Submit'}
								</span>
								<span class="text-ink-caption">Maks: {r.poin} poin</span>
							</div>
							<div class="mt-3 rounded-lg border border-gray-100 bg-surface-muted p-3">
								<p class="mb-1 text-xs font-medium text-ink-caption">Soal:</p>
								<div class="prose prose-sm max-w-none text-ink-body" use:renderMath>
									{@html r.teks_soal}
								</div>
							</div>
							<div class="mt-2 rounded-lg border border-gray-100 bg-surface-muted p-3">
								<p class="mb-1 text-xs font-medium text-ink-caption">Jawaban:</p>
								<pre class="whitespace-pre-wrap text-sm text-ink-body">{r.jawaban_teks || '(kosong)'}</pre>
							</div>
							{#if nilaiEdits[r.jawaban_id]}
								<div class="mt-3 flex flex-wrap items-end gap-3">
									<div class="w-24">
										<label class="label" for={`n${r.jawaban_id}`}>Nilai</label>
										<input id={`n${r.jawaban_id}`} type="number" class="input" bind:value={nilaiEdits[r.jawaban_id].nilai} min="0" max={r.poin} />
									</div>
									<div class="flex-1">
										<label class="label" for={`f${r.jawaban_id}`}>Feedback</label>
										<input id={`f${r.jawaban_id}`} class="input" bind:value={nilaiEdits[r.jawaban_id].feedback} />
									</div>
									<button class="btn-primary py-2" onclick={() => simpanNilai(r.jawaban_id)}>Simpan</button>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{:else}
				<p class="text-ink-caption">Pilih course untuk melihat rekap jawaban.</p>
			{/if}
		{:else}
			<p class="text-ink-caption">Pilih aktivasi di sebelah kiri untuk memulai penilaian.</p>
		{/if}
	</div>
</div>
