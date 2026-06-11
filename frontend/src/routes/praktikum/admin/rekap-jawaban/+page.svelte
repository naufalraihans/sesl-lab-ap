<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { labelJenis, renderMath } from '$lib/utils';
	import { RotateCcw, Trash2, X } from 'lucide-svelte';
	import type { Kelas } from '$lib/types';

	interface Sesi {
		id: number;
		judul_sesi: string;
	}

	interface RekapJawabanItem {
		jawaban_id: number;
		nim: string;
		nama_mahasiswa: string;
		kelas_id: number;
		nama_kelas: string;
		sesi_praktikum_id: number;
		judul_sesi: string;
		course_id: number;
		judul_course: string;
		jenis_course: string;
		jenis_soal: string;
		teks_soal: string;
		poin_maksimal: number;
		jawaban_teks: string;
		is_submitted: boolean;
		waktu_submit: string | null;
		nilai: number | null;
		feedback: string | null;
	}

	interface RekapResponse {
		items: RekapJawabanItem[];
		total: number;
	}

	let kelasList = $state<Kelas[]>([]);
	let sesiList = $state<Sesi[]>([]);

	let selectedKelas = $state<number>(0);
	let selectedSesi = $state<number>(0);
	let selectedJenis = $state<string>('');
	let searchQuery = $state<string>('');

	let rekap = $state<RekapResponse | null>(null);
	let loading = $state(false);
	let errorMsg = $state('');
	let successMsg = $state('');

	// Bulk Selection
	let selectedIds = $state<Set<number>>(new Set());

	// Detail Modal
	let showDetailModal = $state(false);
	let detailData = $state<RekapJawabanItem | null>(null);

	// Edit Nilai Modal
	let showEditModal = $state(false);
	let editData = $state<RekapJawabanItem | null>(null);
	let editNilai = $state<number | null>(null);
	let editFeedback = $state<string>('');

	onMount(async () => {
		try {
			kelasList = await api.get<Kelas[]>('/api/admin/kelas');
			sesiList = await api.get<Sesi[]>('/api/admin/sesi');
			await fetchRekap();
		} catch (e) {
			errorMsg = (e as Error).message;
		}
	});

	async function fetchRekap() {
		loading = true;
		errorMsg = '';
		rekap = null;
		selectedIds = new Set();

		const params = new URLSearchParams();
		if (selectedKelas > 0) params.append('kelas_id', selectedKelas.toString());
		if (selectedSesi > 0) params.append('sesi_id', selectedSesi.toString());
		if (selectedJenis) params.append('jenis', selectedJenis);
		if (searchQuery) params.append('search', searchQuery);

		try {
			rekap = await api.get<RekapResponse>(`/api/admin/rekap-jawaban?${params.toString()}`);
		} catch (e) {
			errorMsg = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	function handleSearch(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			fetchRekap();
		}
	}

	function toggleSelectAll() {
		if (!rekap || rekap.items.length === 0) return;
		if (selectedIds.size === rekap.items.length) {
			selectedIds.clear();
		} else {
			rekap.items.forEach(i => selectedIds.add(i.jawaban_id));
		}
		selectedIds = new Set(selectedIds);
	}

	function toggleSelect(id: number) {
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = new Set(selectedIds);
	}

	async function doBulkAction(action: 'reset_nilai' | 'delete') {
		if (selectedIds.size === 0) return;
		const msgConfirm = action === 'delete' 
			? `Hapus permanen ${selectedIds.size} jawaban terpilih?` 
			: `Reset nilai (jadi null) untuk ${selectedIds.size} jawaban terpilih?`;
		
		if (!confirm(msgConfirm)) return;

		errorMsg = '';
		successMsg = '';
		loading = true;

		try {
			await api.post('/api/admin/penilaian/bulk-action', {
				action: action,
				jawaban_ids: Array.from(selectedIds)
			});
			successMsg = `Berhasil melakukan bulk ${action} pada ${selectedIds.size} data.`;
			await fetchRekap();
		} catch (e) {
			errorMsg = (e as Error).message;
			loading = false; // set false here since fetchRekap won't be called if error
		}
	}

	function openDetail(item: RekapJawabanItem) {
		detailData = item;
		showDetailModal = true;
	}

	function openEdit(item: RekapJawabanItem) {
		editData = item;
		editNilai = item.nilai;
		editFeedback = item.feedback || '';
		showEditModal = true;
	}

	async function saveNilai() {
		if (!editData) return;
		errorMsg = '';
		successMsg = '';
		try {
			await api.post('/api/admin/penilaian', {
				jawaban_id: editData.jawaban_id,
				nilai: Number(editNilai),
				feedback: editFeedback || null
			});
			successMsg = 'Nilai berhasil disimpan.';
			showEditModal = false;
			await fetchRekap();
		} catch (e) {
			errorMsg = (e as Error).message;
		}
	}

	function formatDate(ts: string | null) {
		if (!ts) return '-';
		return new Date(ts).toLocaleString('id-ID', { dateStyle: 'short', timeStyle: 'short' });
	}

	function truncate(str: string, len: number) {
		if (!str) return '-';
		if (str.length <= len) return str;
		return str.slice(0, len) + '...';
	}
</script>

<h1 class="mb-4 text-2xl font-bold">Rekap Jawaban Global</h1>

{#if errorMsg}
	<div class="mb-4 rounded-lg bg-state-error-bg p-3 text-state-error">{errorMsg}</div>
{/if}
{#if successMsg}
	<div class="mb-4 rounded-lg bg-state-success-bg p-3 text-state-success">{successMsg}</div>
{/if}

<!-- Filter Bar -->
<div class="mb-6 rounded-xl border border-gray-100 bg-white p-5 shadow-sm flex flex-wrap items-end gap-4">
	<div class="flex-1 min-w-[150px]">
		<label for="kelas" class="mb-1 block text-sm font-medium text-ink-caption">Kelas</label>
		<select id="kelas" class="input w-full" bind:value={selectedKelas} onchange={fetchRekap}>
			<option value={0}>Semua Kelas</option>
			{#each kelasList as k}
				<option value={k.id}>{k.nama_kelas}</option>
			{/each}
		</select>
	</div>
	<div class="flex-1 min-w-[150px]">
		<label for="sesi" class="mb-1 block text-sm font-medium text-ink-caption">Sesi Praktikum</label>
		<select id="sesi" class="input w-full" bind:value={selectedSesi} onchange={fetchRekap}>
			<option value={0}>Semua Sesi</option>
			{#each sesiList as s}
				<option value={s.id}>{s.judul_sesi}</option>
			{/each}
		</select>
	</div>
	<div class="flex-1 min-w-[150px]">
		<label for="jenis" class="mb-1 block text-sm font-medium text-ink-caption">Jenis Tes</label>
		<select id="jenis" class="input w-full" bind:value={selectedJenis} onchange={fetchRekap}>
			<option value="">Semua Jenis</option>
			<option value="pretest">Pre-test</option>
			<option value="posttest">Post-test</option>
			<option value="keterampilan">Keterampilan</option>
			<option value="ujian_praktik">Ujian Praktik</option>
		</select>
	</div>
	<div class="flex-1 min-w-[200px]">
		<label for="search" class="mb-1 block text-sm font-medium text-ink-caption">Cari (NIM/Nama)</label>
		<input type="text" id="search" placeholder="Tekan Enter untuk cari..." class="input w-full" bind:value={searchQuery} onkeydown={handleSearch} />
	</div>
	<div class="flex-none">
		<button class="btn-primary" onclick={fetchRekap} disabled={loading}>
			{#if loading} Memuat... {:else} Refresh {/if}
		</button>
	</div>
</div>

<!-- Bulk Action Bar -->
{#if selectedIds.size > 0}
	<div class="mb-4 flex flex-wrap items-center gap-4 rounded-xl border border-primary/20 bg-primary/5 p-4 shadow-sm animate-fade-in">
		<div class="font-medium text-primary"><strong>{selectedIds.size}</strong> jawaban dipilih</div>
		<div class="ml-auto flex gap-2">
			<button class="btn-outline inline-flex items-center gap-1 border-state-warning text-state-warning hover:bg-state-warning hover:text-white" onclick={() => doBulkAction('reset_nilai')}>
				<RotateCcw size={14} /> Reset Nilai
			</button>
			<button class="btn-outline inline-flex items-center gap-1 border-state-error text-state-error hover:bg-state-error hover:text-white" onclick={() => doBulkAction('delete')}>
				<Trash2 size={14} /> Hapus
			</button>
			<button class="btn-outline" onclick={() => (selectedIds = new Set())}>Batal</button>
		</div>
	</div>
{/if}

<!-- Data Table -->
{#if loading && !rekap}
	<div class="py-10 text-center text-ink-caption">Memuat data jawaban...</div>
{:else if rekap}
	<div class="mb-2 text-sm text-ink-caption">
		Total: <strong>{rekap.total}</strong> jawaban
	</div>
	
	{#if rekap.items.length === 0}
		<div class="py-10 text-center text-ink-caption">Tidak ada data jawaban yang cocok.</div>
	{:else}
		<div class="table-wrap rounded-xl border border-gray-100 bg-white shadow-sm overflow-x-auto">
			<table class="tbl w-full min-w-max">
				<thead>
					<tr>
						<th class="w-12 text-center">
							<input type="checkbox" 
								checked={selectedIds.size === rekap.items.length && rekap.items.length > 0}
								onchange={toggleSelectAll} 
								class="rounded border-gray-300 text-primary focus:ring-primary" />
						</th>
						<th>NIM / Nama</th>
						<th>Kelas / Sesi</th>
						<th>Soal & Modul</th>
						<th>Waktu</th>
						<th class="text-center">Nilai</th>
						<th class="text-center">Aksi</th>
					</tr>
				</thead>
				<tbody>
					{#each rekap.items as row}
						<tr class="hover:bg-gray-50/50 transition-colors {selectedIds.has(row.jawaban_id) ? 'bg-primary/5' : ''}">
							<td class="text-center">
								<input type="checkbox" 
									checked={selectedIds.has(row.jawaban_id)}
									onchange={() => toggleSelect(row.jawaban_id)}
									class="rounded border-gray-300 text-primary focus:ring-primary" />
							</td>
							<td>
								<div class="font-medium text-ink-body">{row.nim}</div>
								<div class="text-xs text-ink-caption truncate max-w-[150px]">{row.nama_mahasiswa}</div>
							</td>
							<td>
								<div class="text-sm text-ink-body">{row.nama_kelas}</div>
								<div class="text-xs text-ink-caption">{row.judul_sesi}</div>
							</td>
							<td>
								<div class="flex items-center gap-2 mb-1">
									<span class="badge {row.jenis_course === 'keterampilan' || row.jenis_course === 'ujian_praktik' ? 'bg-amber-100 text-amber-800' : 'bg-primary/10 text-primary'}">
										{labelJenis(row.jenis_course)}
									</span>
									<span class="text-xs font-medium text-ink-caption">({row.jenis_soal})</span>
								</div>
								<div class="text-xs font-medium text-ink-body mb-1">{row.judul_course}</div>
								<div class="text-xs text-ink-body truncate max-w-[200px]" title={row.teks_soal}>
									{truncate(row.teks_soal, 40)}
								</div>
							</td>
							<td class="text-sm text-ink-caption">
								{#if row.is_submitted}
									<span class="block text-state-success">Submitted</span>
									<span class="text-xs">{formatDate(row.waktu_submit)}</span>
								{:else}
									<span class="text-state-warning">Belum Submit</span>
								{/if}
							</td>
							<td class="text-center">
								{#if row.nilai === null}
									<span class="text-ink-caption">-</span>
								{:else}
									<span class="font-bold {row.nilai < (row.poin_maksimal * 0.5) ? 'text-state-error' : 'text-state-success'}">
										{row.nilai}
									</span>
									<span class="text-xs text-ink-caption block">/{row.poin_maksimal}</span>
								{/if}
							</td>
							<td class="text-center space-x-1">
								<button class="btn-outline px-2 py-1 text-xs" onclick={() => openDetail(row)}>Detail</button>
								<button class="btn-outline px-2 py-1 text-xs border-amber-500 text-amber-600 hover:bg-amber-50" onclick={() => openEdit(row)}>Nilai</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
{/if}

<!-- Modal Detail Jawaban -->
{#if showDetailModal && detailData}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6">
		<button class="absolute inset-0 bg-black/60 backdrop-blur-sm" onclick={() => showDetailModal = false} aria-label="Tutup Modal"></button>
		<div class="relative flex w-full max-w-4xl flex-col rounded-2xl bg-white shadow-2xl max-h-[90vh]">
			<div class="flex items-center justify-between border-b px-6 py-4">
				<h3 class="text-lg font-bold text-ink-body">Detail Jawaban - {detailData.nama_mahasiswa}</h3>
				<button class="text-ink-caption hover:text-ink-body" onclick={() => showDetailModal = false}><X size={20} /></button>
			</div>
			
			<div class="overflow-y-auto p-6">
				<div class="grid grid-cols-2 gap-4 mb-6">
					<div class="rounded-lg bg-surface-soft p-4 text-sm">
						<div class="text-ink-caption mb-1">Soal / Modul</div>
						<div class="font-medium text-ink-body">{detailData.judul_course} ({detailData.jenis_soal})</div>
						<div class="mt-2 prose prose-sm max-w-none bg-white p-3 rounded border" use:renderMath>
							{@html detailData.teks_soal}
						</div>
					</div>
					<div class="rounded-lg bg-surface-soft p-4 text-sm">
						<div class="grid grid-cols-2 gap-y-3">
							<div>
								<span class="text-ink-caption block text-xs">NIM</span>
								<span class="font-medium">{detailData.nim}</span>
							</div>
							<div>
								<span class="text-ink-caption block text-xs">Kelas</span>
								<span class="font-medium">{detailData.nama_kelas}</span>
							</div>
							<div>
								<span class="text-ink-caption block text-xs">Status</span>
								<span class="font-medium {detailData.is_submitted ? 'text-state-success' : 'text-state-warning'}">
									{detailData.is_submitted ? 'Submitted' : 'Belum Submit'}
								</span>
							</div>
							<div>
								<span class="text-ink-caption block text-xs">Nilai Akhir</span>
								<span class="font-bold text-lg">{detailData.nilai ?? '-'} <span class="text-sm font-normal text-ink-caption">/ {detailData.poin_maksimal}</span></span>
							</div>
						</div>
					</div>
				</div>

				<div class="rounded-lg border border-gray-200">
					<div class="bg-gray-50 px-4 py-2 border-b font-medium text-ink-body">Jawaban Mahasiswa</div>
					<div class="p-4 bg-white">
						{#if detailData.jenis_soal === 'coding'}
							<pre class="bg-[#1e1e1e] text-white p-4 rounded-lg overflow-x-auto text-sm font-mono">{detailData.jawaban_teks || '(Kosong)'}</pre>
						{:else}
							<div class="prose max-w-none text-ink-body bg-gray-50/50 p-4 rounded-lg min-h-[100px] whitespace-pre-wrap">{detailData.jawaban_teks || '(Kosong)'}</div>
						{/if}
					</div>
				</div>

				{#if detailData.feedback}
					<div class="mt-4 rounded-lg border border-amber-200 bg-amber-50 p-4">
						<div class="font-medium text-amber-800 mb-1">Feedback:</div>
						<p class="text-sm text-amber-900">{detailData.feedback}</p>
					</div>
				{/if}
			</div>
			
			<div class="border-t bg-gray-50 px-6 py-4 flex justify-end gap-3 rounded-b-2xl">
				<button class="btn-outline" onclick={() => showDetailModal = false}>Tutup</button>
				<button class="btn-primary" onclick={() => { showDetailModal = false; openEdit(detailData!); }}>Edit Nilai</button>
			</div>
		</div>
	</div>
{/if}

<!-- Modal Edit Nilai -->
{#if showEditModal && editData}
	<div class="fixed inset-0 z-[60] flex items-center justify-center p-4">
		<button class="absolute inset-0 bg-black/60 backdrop-blur-sm" onclick={() => showEditModal = false} aria-label="Tutup Modal"></button>
		<div class="relative w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl">
			<h3 class="mb-4 text-lg font-bold text-ink-body">Beri Nilai - {editData.nama_mahasiswa}</h3>
			
			<div class="mb-4">
				<label class="label block mb-1" for="input_nilai">Nilai (Maks: {editData.poin_maksimal})</label>
				<input id="input_nilai" type="number" class="input w-full" bind:value={editNilai} min="0" max={editData.poin_maksimal} />
			</div>
			<div class="mb-6">
				<label class="label block mb-1" for="input_feedback">Feedback (Opsional)</label>
				<textarea id="input_feedback" class="input w-full min-h-[80px]" bind:value={editFeedback}></textarea>
			</div>
			
			<div class="flex justify-end gap-3">
				<button class="btn-outline" onclick={() => showEditModal = false}>Batal</button>
				<button class="btn-primary" onclick={saveNilai}>Simpan Nilai</button>
			</div>
		</div>
	</div>
{/if}
