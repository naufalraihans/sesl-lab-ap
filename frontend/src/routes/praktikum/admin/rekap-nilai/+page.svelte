<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	interface RekapKolom {
		key: string;
		label: string;
	}

	interface RekapMahasiswa {
		nim: string;
		nama: string;
		scores: Record<string, number>;
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
				row.push((m.scores[c.key] ?? 0).toString());
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
	<div class="flex-none">
		<button class="btn-primary" onclick={exportCSV} disabled={!rekap || filteredData.length === 0}>
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" x2="12" y1="15" y2="3"/></svg>
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
								<td class="text-center text-ink-caption">
									{row.scores[col.key] ?? '-'}
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
