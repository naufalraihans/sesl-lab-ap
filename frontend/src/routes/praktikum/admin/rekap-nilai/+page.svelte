<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { FileSpreadsheet, Save } from 'lucide-svelte';

	interface RekapKolom {
		key: string;
		label: string;
	}

	interface RekapSel {
		pengerjaan_id: number;
		murni: number;
		keaktifan: number | null;
		total: number;
		edit_keaktifan: boolean;
		nilai_akhir: number | null;
	}

	interface RekapMahasiswa {
		nim: string;
		nama: string;
		scores: Record<string, RekapSel>;
		total_score: number;
	}

	interface RekapResponse {
		columns: RekapKolom[];
		data: RekapMahasiswa[];
	}

	interface Kelas {
		id: number;
		nama_kelas: string;
	}

	let kelasList = $state<Kelas[]>([]);
	let selectedKelas = $state<number>(0);
	let rekap = $state<RekapResponse | null>(null);
	let loading = $state(false);
	let errorMsg = $state('');
	let searchQuery = $state('');
	// Perubahan keaktifan yang belum disimpan: pengerjaan_id -> nilai baru.
	let edits = $state<Record<number, number>>({});
	let saving = $state(false);

	function liveKeaktifan(sel: RekapSel): number {
		return Number(edits[sel.pengerjaan_id] ?? sel.keaktifan ?? 0);
	}

	async function saveKeaktifan() {
		const items = Object.entries(edits).map(([pid, nilai]) => ({
			pengerjaan_id: Number(pid),
			nilai: Number(nilai)
		}));
		if (items.length === 0) return;
		saving = true;
		errorMsg = '';
		try {
			await api.post('/api/admin/keaktifan', { items });
			edits = {};
			await fetchRekap();
		} catch (e) {
			errorMsg = (e as Error).message;
		} finally {
			saving = false;
		}
	}

	onMount(async () => {
		try {
			kelasList = await api.get<Kelas[]>('/api/admin/kelas');
			if (kelasList.length > 0) {
				selectedKelas = kelasList[0].id;
				await fetchRekap();
			}
		} catch (e) {
			errorMsg = (e as Error).message;
		}
	});

	async function fetchRekap() {
		if (!selectedKelas) return;
		loading = true;
		errorMsg = '';
		rekap = null;
		try {
			rekap = await api.get<RekapResponse>(`/api/admin/rekap/kelas/${selectedKelas}`);
		} catch (e) {
			errorMsg = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	let filteredData = $derived(
		rekap?.data.filter((m) => {
			const q = searchQuery.toLowerCase();
			return m.nim.toLowerCase().includes(q) || m.nama.toLowerCase().includes(q);
		}) || []
	);

	function exportCSV() {
		if (!rekap || filteredData.length === 0) return;

		// Buat Header
		const header = ['NIM', 'Nama', ...rekap.columns.map(c => c.label), 'Total Nilai'];
		
		// Buat Baris Data
		const rows = filteredData.map(m => {
			const row = [m.nim, m.nama];
			rekap!.columns.forEach(c => {
				const sel = m.scores[c.key];
				row.push(sel ? sel.total.toString() : '0');
			});
			row.push(m.total_score.toString());
			return row;
		});

		// Gabungkan jadi CSV
		const csvContent = [
			header.join(','),
			...rows.map(r => r.map(v => `"${v}"`).join(','))
		].join('\n');

		// Buat blob dan trigger download
		const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.setAttribute('href', url);
		link.setAttribute('download', `rekap_nilai_kelas_${selectedKelas}.csv`);
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}
</script>

<h1 class="mb-4 text-2xl font-bold">Rekap Nilai</h1>

{#if errorMsg}
	<div class="mb-4 rounded-lg bg-state-error-bg p-3 text-state-error">{errorMsg}</div>
{/if}

<div class="mb-6 rounded-xl bg-white p-5 shadow-sm border border-gray-100 flex flex-wrap items-end gap-4">
	<div class="flex-1 min-w-[200px]">
		<label for="kelas" class="mb-1 block text-sm font-medium text-ink-caption">Pilih Kelas</label>
		<select id="kelas" class="input w-full" bind:value={selectedKelas} onchange={fetchRekap}>
			{#each kelasList as k}
				<option value={k.id}>{k.nama_kelas}</option>
			{/each}
		</select>
	</div>
	<div class="flex-1 min-w-[200px]">
		<label for="search" class="mb-1 block text-sm font-medium text-ink-caption">Cari Mahasiswa</label>
		<input type="text" id="search" placeholder="NIM atau Nama..." class="input w-full" bind:value={searchQuery} />
	</div>
	<div class="flex-none flex gap-2 items-end">
		<button class="btn-primary flex items-center gap-1.5 py-2.5 text-sm font-semibold" onclick={saveKeaktifan} disabled={saving || Object.keys(edits).length === 0}>
			<Save size={16} />
			{saving ? 'Menyimpan…' : `Simpan Keaktifan${Object.keys(edits).length ? ` (${Object.keys(edits).length})` : ''}`}
		</button>
		<button class="btn-outline flex items-center gap-1.5 py-2.5 text-sm font-semibold hover:bg-slate-50 transition-colors" onclick={exportCSV} disabled={!rekap || filteredData.length === 0}>
			<FileSpreadsheet size={16} class="text-emerald-600" />
			Export CSV
		</button>
	</div>
</div>

{#if loading}
	<div class="py-10 text-center text-ink-caption">Memuat data rekapitulasi...</div>
{:else if rekap}
	{#if filteredData.length === 0}
		<div class="py-10 text-center text-ink-caption">Tidak ada data yang cocok dengan pencarian.</div>
	{:else}
		<div class="table-wrap rounded-xl border border-gray-100 bg-white shadow-sm overflow-x-auto">
			<table class="tbl w-full min-w-max">
				<thead>
					<tr>
						<th class="sticky left-0 bg-gray-50 z-10 w-24 border-r border-gray-200">NIM</th>
						<th class="sticky left-24 bg-gray-50 z-10 w-48 border-r border-gray-200">Nama</th>
						{#each rekap.columns as col}
							<th class="text-center">{col.label}</th>
						{/each}
						<th class="text-center bg-gray-50 border-l border-gray-200 text-brand-blue">Total</th>
					</tr>
				</thead>
				<tbody>
					{#each filteredData as row}
						<tr class="hover:bg-gray-50/50 transition-colors">
							<td class="sticky left-0 bg-white z-10 font-medium text-ink-body border-r border-gray-100">{row.nim}</td>
							<td class="sticky left-24 bg-white z-10 text-ink-body border-r border-gray-100 truncate max-w-xs">{row.nama}</td>
							{#each rekap.columns as col}
								{@const sel = row.scores[col.key]}
								<td class="text-center text-ink-caption">
									{#if !sel}
										-
									{:else if sel.edit_keaktifan}
										<div class="flex items-center justify-center gap-1 text-xs">
											<span>{sel.murni}</span>
											<span class="text-ink-caption">+</span>
											<input
												type="number" min="0" max="100"
												class="w-12 rounded border border-gray-200 px-1 py-0.5 text-center"
												value={edits[sel.pengerjaan_id] ?? sel.keaktifan ?? ''}
												oninput={(e) => (edits[sel.pengerjaan_id] = Number((e.target as HTMLInputElement).value))}
											/>
											<span class="font-semibold text-ink-body">= {sel.murni + liveKeaktifan(sel)}</span>
										</div>
									{:else if sel.nilai_akhir != null}
										<span class="font-semibold text-brand-blue">{sel.nilai_akhir}</span>
										<span class="text-xs text-ink-caption">({sel.total})</span>
									{:else}
										{sel.total}
									{/if}
								</td>
							{/each}
							<td class="text-center bg-gray-50/50 border-l border-gray-100 font-semibold text-brand-blue">
								{row.total_score}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
{/if}
