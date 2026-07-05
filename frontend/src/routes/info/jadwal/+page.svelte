<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { FileText, Calendar, Search, SlidersHorizontal, ChevronDown, ChevronUp, Users } from 'lucide-svelte';
	import type { Jadwal, User, AmpuanKelompok, Kelas } from '$lib/types';

	let jadwal = $state<Jadwal[]>([]);
	let config = $state<{ mode: string; gdrive_url: string }>({ mode: 'internal', gdrive_url: '' });
	let loading = $state(true);
	let err = $state('');

	let selectedKelas = $state<Kelas | null>(null);
	let mhsList = $state<User[]>([]);
	let ampuanList = $state<AmpuanKelompok[]>([]);
	let loadingMhs = $state(false);

	// Collapsible state: kelompok name → open/closed
	let collapsedGroups = $state<Record<string, boolean>>({});

	let filterQuery = $state('');
	let filterHari = $state('');
	let filterShift = $state('');

	onMount(async () => {
		try {
			config = await api.get('/api/info/jadwal/config');
			if (config.mode !== 'gdrive') {
				jadwal = (await api.get<Jadwal[]>('/api/info/jadwal')) ?? [];
			}
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});

	let filteredJadwal = $derived(
		jadwal.filter((j) => {
			const q = filterQuery.toLowerCase().trim();
			const matchesQuery = q === '' ||
				(j.kelas?.nama_kelas ?? '').toLowerCase().includes(q) ||
				(j.keterangan ?? '').toLowerCase().includes(q) ||
				(j.hari ?? '').toLowerCase().includes(q);
			const matchesHari = filterHari === '' || j.hari === filterHari;
			const matchesShift = filterShift === '' || String(j.shift) === filterShift;
			return matchesQuery && matchesHari && matchesShift;
		})
	);

	async function showKelas(k: Kelas) {
		selectedKelas = k;
		loadingMhs = true;
		collapsedGroups = {}; // reset collapse state
		try {
			const res = await api.get<{ mahasiswa: User[]; ampuan: AmpuanKelompok[] }>(
				`/api/info/kelas/${k.id}/mahasiswa`
			);
			mhsList = res?.mahasiswa ?? [];
			ampuanList = res?.ampuan ?? [];
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loadingMhs = false;
		}
	}

	function ampuanForKelompok(kel: string): string {
		return ampuanList.find(x => x.kelompok === kel)?.asisten?.nama ?? '-';
	}

	let kelompokGroups = $derived(
		mhsList.reduce<Record<string, User[]>>((acc, m) => {
			const kel = m.kelompok ?? '(Belum ada)';
			if (!acc[kel]) acc[kel] = [];
			acc[kel].push(m);
			return acc;
		}, {})
	);

	function toggleGroup(kel: string) {
		collapsedGroups = { ...collapsedGroups, [kel]: !collapsedGroups[kel] };
	}
</script>

<div class="space-y-6">
	<div class="flex items-center gap-3 text-left">
		<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-fun-blue to-blue-600 text-white flex items-center justify-center shadow-md shadow-blue-200">
			<Calendar size={20} />
		</div>
		<div>
			<h1 class="text-2xl font-extrabold text-slate-900 leading-none">Jadwal Praktikum</h1>
			<p class="text-xs font-bold text-slate-600 mt-1">Jadwal resmi kelas &amp; shift per semester.</p>
		</div>
	</div>

	{#if loading}
		<p class="text-slate-700 font-semibold">Memuat…</p>
	{:else if err}
		<p class="rounded-lg bg-state-error-bg p-3 text-state-error font-semibold">{err}</p>
	{:else if config.mode === 'gdrive' && config.gdrive_url}
		<a href={config.gdrive_url} target="_blank" rel="noopener" class="btn-primary"><FileText size={16}/> Buka Jadwal (Google Drive)</a>
	{:else if jadwal.length === 0}
		<p class="text-slate-600 font-semibold">Belum ada jadwal yang dipublikasikan.</p>
	{:else}
		<!-- Filters -->
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-3 bg-white p-5 rounded-2xl border border-slate-200 shadow-sm text-left">
			<div>
				<label for="search-kelas" class="block text-xs font-black text-slate-700 uppercase tracking-wider mb-1.5 flex items-center gap-1.5">
					<Search size={12} /> Cari Kelas
				</label>
				<input id="search-kelas" type="text" placeholder="Ketik nama kelas..." bind:value={filterQuery} class="input w-full" />
			</div>
			<div>
				<label for="filter-hari" class="block text-xs font-black text-slate-700 uppercase tracking-wider mb-1.5 flex items-center gap-1.5">
					<SlidersHorizontal size={12} /> Filter Hari
				</label>
				<select id="filter-hari" bind:value={filterHari} class="input w-full bg-white text-slate-800 font-semibold">
					<option value="">Semua Hari</option>
					{#each ['Senin','Selasa','Rabu','Kamis','Jumat','Sabtu','Minggu'] as h}
						<option value={h}>{h}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="filter-shift" class="block text-xs font-black text-slate-700 uppercase tracking-wider mb-1.5 flex items-center gap-1.5">
					<SlidersHorizontal size={12} /> Filter Shift
				</label>
				<select id="filter-shift" bind:value={filterShift} class="input w-full bg-white text-slate-800 font-semibold">
					<option value="">Semua Shift</option>
					{#each ['1','2','3','4','5'] as s}
						<option value={s}>Shift {s}</option>
					{/each}
				</select>
			</div>
		</div>

		{#if filteredJadwal.length === 0}
			<p class="text-slate-600 py-8 text-center bg-white border border-slate-200 rounded-2xl font-bold text-sm">Tidak ada jadwal yang sesuai dengan filter.</p>
		{:else}
			<div class="table-wrap">
				<table class="tbl">
					<thead>
						<tr><th>Kelas</th><th>Shift</th><th>Hari</th><th>Jam</th><th>Keterangan</th><th></th></tr>
					</thead>
					<tbody>
						{#each filteredJadwal as j}
							<tr>
								<td class="font-extrabold text-slate-900">{j.kelas?.nama_kelas ?? j.kelas_id}</td>
								<td><span class="badge bg-blue-50 text-blue-800 border-blue-200 font-bold">Shift {j.shift}</span></td>
								<td class="font-bold text-slate-800">{j.hari}</td>
								<td class="font-mono text-xs font-bold text-slate-800">{j.jam_mulai} – {j.jam_selesai}</td>
								<td class="text-slate-700 font-semibold">{j.keterangan}</td>
								<td>
									{#if j.kelas}
										<button
											class="text-sm font-extrabold text-primary hover:text-primary-hover hover:underline transition-colors focus:outline-none"
											onclick={() => showKelas(j.kelas!)}>
											Lihat Mahasiswa
										</button>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}

		{#if selectedKelas}
			<hr class="my-6 border-slate-200" />
			<div class="bg-slate-50 rounded-2xl border border-slate-200 p-6 text-left">
				<div class="flex items-center gap-3 mb-5">
					<div class="w-9 h-9 rounded-xl bg-gradient-to-br from-primary to-fun-purple text-white flex items-center justify-center shadow-sm">
						<Users size={16} />
					</div>
					<h2 class="text-xl font-extrabold text-slate-900">Daftar Mahasiswa — {selectedKelas.nama_kelas}</h2>
				</div>

				{#if loadingMhs}
					<p class="text-slate-700 font-semibold">Memuat…</p>
				{:else if mhsList.length === 0}
					<p class="text-slate-600 font-semibold">Belum ada mahasiswa di kelas ini.</p>
				{:else}
					<div class="grid gap-4">
						{#each Object.entries(kelompokGroups) as [kel, members]}
							{@const isOpen = collapsedGroups[kel] !== true}
							<div class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
								<!-- Collapsible header -->
								<button
									class="w-full flex flex-wrap items-center justify-between gap-3 px-5 py-4 border-b border-slate-100 hover:bg-slate-50 transition-colors text-left"
									onclick={() => toggleGroup(kel)}
								>
									<div class="flex items-center gap-3">
										<h3 class="text-base font-extrabold text-slate-900">Kelompok {kel}</h3>
										<span class="badge bg-primary/5 text-primary border-primary/20 font-bold text-xs">
											Asisten: {ampuanForKelompok(kel)}
										</span>
										<span class="text-xs font-bold text-slate-600 bg-slate-100 px-2 py-0.5 rounded-full">{members.length} mahasiswa</span>
									</div>
									<span class="text-slate-500">
										{#if isOpen}
											<ChevronUp size={18} />
										{:else}
											<ChevronDown size={18} />
										{/if}
									</span>
								</button>
								<!-- Collapsible body -->
								{#if isOpen}
									<div class="table-wrap rounded-t-none border-0">
										<table class="tbl">
											<thead><tr><th>No</th><th>NIM</th><th>Nama</th><th>Shift</th></tr></thead>
											<tbody>
												{#each members as m, i}
													<tr>
														<td class="font-mono text-xs text-slate-600">{i + 1}</td>
														<td class="font-mono text-slate-900 font-bold">{m.nim}</td>
														<td class="font-semibold text-slate-800">{m.nama}</td>
														<td><span class="badge bg-slate-50 border-slate-200 text-slate-700 font-bold">{m.shift ?? '-'}</span></td>
													</tr>
												{/each}
											</tbody>
										</table>
									</div>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{/if}
	{/if}
</div>
