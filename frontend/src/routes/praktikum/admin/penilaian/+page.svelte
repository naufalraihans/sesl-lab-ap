<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { labelJenis, labelShift, renderMath } from '$lib/utils';
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

	// AI grading sinkron 1-per-1: frontend yang me-loop.
	let aiRunning = $state(false);
	let aiProcessed = $state(0);
	let aiTotal = $state(0);
	let aiFailed = $state(0);
	let aiFailedItems = $state<{ jawaban_id: number; nama: string; nim: string }[]>([]);
	let aiStopRequested = false;

	let nilaiEdits = $state<Record<number, { nilai: number; feedback: string }>>({});
	// Checklist bulk-select untuk grading sebagian.
	let selected = $state<Record<number, boolean>>({});
	let selectedCount = $derived(rekap.filter((r) => selected[r.jawaban_id]).length);

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

	// Inti: nilai sederet jawaban_id satu per satu (dipakai "nilai semua" & "nilai terpilih").
	async function gradeIds(ids: number[]) {
		if (ids.length === 0) { msg = 'Tidak ada jawaban untuk dinilai.'; return; }
		err = ''; msg = '';
		aiProcessed = 0; aiFailed = 0; aiFailedItems = []; aiStopRequested = false;
		aiTotal = ids.length; aiRunning = true;
		try {
			for (const id of ids) {
				if (aiStopRequested) break;
				try { await api.post('/api/admin/penilaian/ai-grade/one', { jawaban_id: id }); }
				catch {
					aiFailed++; // 1 jawaban gagal/timeout → lanjut yang lain, catat siapa
					const r = rekap.find((x) => x.jawaban_id === id);
					aiFailedItems = [...aiFailedItems, { jawaban_id: id, nama: r?.nama_mahasiswa ?? '?', nim: r?.nim ?? '-' }];
				}
				aiProcessed++;
			}
		} finally {
			aiRunning = false;
			const sukses = aiProcessed - aiFailed;
			msg = `AI grading selesai: ${sukses} berhasil`
				+ (aiFailed ? `, ${aiFailed} gagal (lihat daftar di bawah)` : '')
				+ (aiStopRequested ? ' — dihentikan' : '') + '.';
			if (selectedCourseId) await loadRekap(selectedCourseId);
		}
	}

	// Nilai SEMUA jawaban yang belum dinilai (perilaku lama).
	async function startAIGrading() {
		if (!selectedAktivasi || !selectedCourseId) return;
		if (!confirm('Mulai penilaian AI untuk SEMUA jawaban yang belum dinilai?')) return;
		try {
			const res = await api.get<{ jawaban_ids: number[]; total: number }>(
				`/api/admin/penilaian/ai-grade/targets?aktivasi_sesi_id=${selectedAktivasi.id}&course_id=${selectedCourseId}`
			);
			await gradeIds(res?.jawaban_ids ?? []);
		} catch (e) { err = (e as Error).message; }
	}

	// Ulangi grading HANYA untuk jawaban yang tadi gagal.
	async function retryFailed() {
		const ids = aiFailedItems.map((x) => x.jawaban_id);
		if (ids.length === 0) return;
		await gradeIds(ids);
	}

	// Nilai HANYA jawaban yang dicentang.
	async function gradeSelected() {
		const ids = rekap.filter((r) => selected[r.jawaban_id]).map((r) => r.jawaban_id);
		if (ids.length === 0) { msg = 'Belum ada jawaban yang dipilih.'; return; }
		if (!confirm(`Nilai ${ids.length} jawaban terpilih dengan AI?`)) return;
		await gradeIds(ids);
		selected = {};
	}

	// Reset nilai (& feedback) jawaban terpilih jadi kosong → bisa dinilai ulang.
	async function resetSelected() {
		const ids = rekap.filter((r) => selected[r.jawaban_id]).map((r) => r.jawaban_id);
		if (ids.length === 0) { msg = 'Belum ada jawaban yang dipilih.'; return; }
		if (!confirm(`Reset nilai ${ids.length} jawaban terpilih jadi kosong? Feedback juga dihapus.`)) return;
		err = ''; msg = '';
		try {
			await api.post('/api/admin/penilaian/bulk-action', { action: 'reset_nilai', jawaban_ids: ids });
			msg = `${ids.length} nilai direset.`;
			selected = {};
			if (selectedCourseId) await loadRekap(selectedCourseId);
		} catch (e) { err = (e as Error).message; }
	}

	// Genosida: reset SEMUA nilai di course/aktivasi yang sedang dibuka (jawaban tetap).
	async function resetSemua() {
		const ids = rekap.map((r) => r.jawaban_id);
		if (ids.length === 0) { msg = 'Tidak ada jawaban untuk direset.'; return; }
		if (!confirm(`GENOSIDA NILAI: reset SEMUA ${ids.length} nilai di course ini jadi kosong (jawaban TIDAK dihapus, feedback ikut terhapus)? Tidak bisa di-undo.`)) return;
		if (!confirm(`Yakin? Ini menimpa nilai ${ids.length} jawaban (termasuk yang sudah dinilai).`)) return;
		err = ''; msg = '';
		try {
			await api.post('/api/admin/penilaian/bulk-action', { action: 'reset_nilai', jawaban_ids: ids });
			msg = `${ids.length} nilai direset (genosida).`;
			selected = {};
			if (selectedCourseId) await loadRekap(selectedCourseId);
		} catch (e) { err = (e as Error).message; }
	}

	function pilihSemua() {
		const s: Record<number, boolean> = {};
		for (const r of rekap) s[r.jawaban_id] = true;
		selected = s;
	}
	function pilihBelumDinilai() {
		const s: Record<number, boolean> = {};
		for (const r of rekap) if (r.is_submitted && r.nilai == null && r.jawaban_teks?.trim()) s[r.jawaban_id] = true;
		selected = s;
	}
	function kosongkanPilihan() { selected = {}; }

	function stopAIGrading() { aiStopRequested = true; }
</script>

<h1 class="mb-4 text-2xl">Penilaian Mahasiswa</h1>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

{#if !aiRunning && aiFailedItems.length > 0}
	<div class="mb-3 rounded-lg border border-state-error/30 bg-state-error-bg p-3 text-sm">
		<div class="mb-2 flex flex-wrap items-center justify-between gap-2">
			<span class="font-medium text-state-error">{aiFailedItems.length} jawaban gagal dinilai:</span>
			<button class="btn-outline border-state-error px-3 py-1 text-xs text-state-error hover:bg-state-error hover:text-white" onclick={retryFailed}>Ulangi yang gagal ({aiFailedItems.length})</button>
		</div>
		<ul class="ml-4 list-disc space-y-0.5 text-ink-body">
			{#each aiFailedItems as f}
				<li>{f.nama} <span class="text-ink-caption">({f.nim})</span></li>
			{/each}
		</ul>
	</div>
{/if}

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
					<p class="text-xs text-ink-caption">{a.kelas?.nama_kelas ?? a.kelas_id} · {labelShift(a.shift)}</p>
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
				{#if aiRunning}
					<div class="mb-4 rounded-lg border border-primary/20 bg-primary/5 p-4">
						<div class="mb-2 flex items-center justify-between">
							<h3 class="font-medium text-primary">AI Grading sedang berjalan…</h3>
							<button class="btn-outline border-state-error px-3 py-1 text-xs text-state-error hover:bg-state-error hover:text-white" onclick={stopAIGrading}>Hentikan</button>
						</div>
						<div class="mb-2 h-2 w-full overflow-hidden rounded-full bg-gray-200">
							<div class="h-full bg-primary transition-all duration-300" style="width: {aiTotal > 0 ? (aiProcessed / aiTotal) * 100 : 0}%"></div>
						</div>
						<p class="text-sm text-ink-body">
							Diproses: {aiProcessed} / {aiTotal}{aiFailed ? ` · gagal: ${aiFailed}` : ''}
						</p>
					</div>
				{:else}
					<div class="mb-4 rounded-lg border border-gray-200 bg-surface-muted p-4">
						<div class="flex flex-wrap items-center justify-between gap-3">
							<div>
								<h3 class="font-medium text-ink-body">Otomatisasi Penilaian (AI)</h3>
								<p class="text-sm text-ink-caption">Nilai semua yang belum dinilai, atau centang beberapa lalu "Nilai Terpilih". Satu per satu, bisa dihentikan.</p>
							</div>
							<div class="flex flex-wrap gap-2">
								<button class="btn-outline py-2" onclick={gradeSelected} disabled={selectedCount === 0}>Nilai Terpilih ({selectedCount})</button>
								<button class="btn-outline border-state-error py-2 text-state-error hover:bg-state-error hover:text-white" onclick={resetSelected} disabled={selectedCount === 0}>Reset Terpilih ({selectedCount})</button>
								<button class="btn-outline border-state-error bg-state-error/5 py-2 text-state-error hover:bg-state-error hover:text-white" onclick={resetSemua} disabled={rekap.length === 0}>Reset Semua ({rekap.length})</button>
								<button class="btn-primary py-2" onclick={startAIGrading}>
									<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2"><path d="m2 16 3-3 3 3"></path><path d="m2 16 3 3 3-3"></path><path d="M14 6h-4a4 4 0 0 0-4 4v10"></path><path d="M18 10a4 4 0 0 1 4 4v6"></path><path d="m22 20-3 3-3-3"></path></svg>
									Nilai Semua
								</button>
							</div>
						</div>
						<div class="mt-3 flex flex-wrap gap-4 text-xs">
							<button class="text-state-info hover:underline" onclick={pilihSemua}>Pilih semua</button>
							<button class="text-state-info hover:underline" onclick={pilihBelumDinilai}>Pilih yang belum dinilai</button>
							<button class="text-ink-caption hover:underline" onclick={kosongkanPilihan}>Kosongkan</button>
						</div>
					</div>
				{/if}

				<div class="space-y-4">
					{#each rekap as r}
						<div class="card">
							<div class="flex flex-wrap items-center gap-3 text-sm">
								<input
									type="checkbox" class="h-4 w-4"
									checked={!!selected[r.jawaban_id]}
									onchange={(e) => (selected[r.jawaban_id] = (e.target as HTMLInputElement).checked)}
									title="Pilih untuk dinilai AI"
								/>
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
