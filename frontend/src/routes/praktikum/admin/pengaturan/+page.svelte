<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Calendar, UserCheck, Copy, AlertTriangle, Plus, Trash2,
		Download, Upload, Save
	} from 'lucide-svelte';

	// ─── Types ────────────────────────────────────────────────────────────────
	interface RescheduleCard { id: string; kelas: string; hari: string; jam: string; catatan: string; }
	interface RecruitCard    { id: string; judul: string; deskripsi: string; link?: string; }
	interface SusulanCard    { id: string; kelas: string; deskripsi: string; }
	interface PlagiarismRow  { nim: string; nama: string; keterangan: string; }
	interface PlagiarismData { rows: PlagiarismRow[]; updatedAt: string; catatan?: string; }

	let reschedules  = $state<RescheduleCard[]>([]);
	let recruitCards = $state<RecruitCard[]>([]);
	let susulanCards = $state<SusulanCard[]>([]);
	let plagiarismData = $state<PlagiarismData>({ rows: [], updatedAt: '' });
	let annSaved = $state(false);

	let newRs = $state<Omit<RescheduleCard,'id'>>({ kelas:'', hari:'', jam:'', catatan:'' });
	let newRc = $state<Omit<RecruitCard,'id'>>({ judul:'', deskripsi:'', link:'' });
	let newSu = $state<Omit<SusulanCard,'id'>>({ kelas:'', deskripsi:'' });
	let newPl = $state<PlagiarismRow>({ nim:'', nama:'', keterangan:'' });
	let importCsv = $state('');

	function uid() { return Math.random().toString(36).slice(2); }

	onMount(async () => {
		try {
			// Fetch configurations from backend API
			const res = await api.get<any[]>('/api/admin/konfigurasi') || [];
			const configMap: Record<string, string> = {};
			res.forEach(item => {
				if (item && item.key) configMap[item.key] = item.value;
			});
			if (configMap['ann_reschedules']) reschedules = JSON.parse(configMap['ann_reschedules']);
			if (configMap['ann_recruit']) recruitCards = JSON.parse(configMap['ann_recruit']);
			if (configMap['ann_susulan']) susulanCards = JSON.parse(configMap['ann_susulan']);
			if (configMap['ann_plagiarism']) plagiarismData = JSON.parse(configMap['ann_plagiarism']);
		} catch {
			// Fallback to local storage if API fails
			const r = localStorage.getItem('ann_reschedules'); if (r) reschedules = JSON.parse(r);
			const rc = localStorage.getItem('ann_recruit');    if (rc) recruitCards = JSON.parse(rc);
			const s = localStorage.getItem('ann_susulan');     if (s) susulanCards = JSON.parse(s);
			const p = localStorage.getItem('ann_plagiarism');  if (p) plagiarismData = JSON.parse(p);
		}
	});

	async function saveAnnouncements() {
		try {
			// Save to backend database configurations table via API
			await api.post('/api/admin/konfigurasi', { key: 'ann_reschedules', value: JSON.stringify(reschedules) });
			await api.post('/api/admin/konfigurasi', { key: 'ann_recruit', value: JSON.stringify(recruitCards) });
			await api.post('/api/admin/konfigurasi', { key: 'ann_susulan', value: JSON.stringify(susulanCards) });
			await api.post('/api/admin/konfigurasi', { key: 'ann_plagiarism', value: JSON.stringify(plagiarismData) });
		} catch { /* ignore */ }

		// Sync to local storage as robust fallback
		localStorage.setItem('ann_reschedules', JSON.stringify(reschedules));
		localStorage.setItem('ann_recruit',     JSON.stringify(recruitCards));
		localStorage.setItem('ann_susulan',     JSON.stringify(susulanCards));
		localStorage.setItem('ann_plagiarism',  JSON.stringify(plagiarismData));
		annSaved = true; setTimeout(() => annSaved = false, 3000);
	}

	function addReschedule() {
		if (!newRs.kelas || !newRs.hari) return;
		reschedules = [...reschedules, { ...newRs, id: uid() }];
		newRs = { kelas:'', hari:'', jam:'', catatan:'' };
	}
	function removeReschedule(id: string) { reschedules = reschedules.filter(r => r.id !== id); }

	function addRecruit() {
		if (!newRc.judul) return;
		recruitCards = [...recruitCards, { ...newRc, id: uid() }];
		newRc = { judul:'', deskripsi:'', link:'' };
	}
	function removeRecruit(id: string) { recruitCards = recruitCards.filter(r => r.id !== id); }

	function addSusulan() {
		if (!newSu.kelas) return;
		susulanCards = [...susulanCards, { ...newSu, id: uid() }];
		newSu = { kelas:'', deskripsi:'' };
	}
	function removeSusulan(id: string) { susulanCards = susulanCards.filter(s => s.id !== id); }

	function addPlagiarismRow() {
		if (!newPl.nim || !newPl.nama) return;
		plagiarismData = { ...plagiarismData, rows: [...plagiarismData.rows, { ...newPl }], updatedAt: new Date().toLocaleDateString('id-ID') };
		newPl = { nim:'', nama:'', keterangan:'' };
	}
	function removePlagiarismRow(nim: string) {
		plagiarismData = { ...plagiarismData, rows: plagiarismData.rows.filter(r => r.nim !== nim), updatedAt: new Date().toLocaleDateString('id-ID') };
	}
	function importFromCsv() {
		try {
			const rows: PlagiarismRow[] = importCsv.trim().split('\n').filter(Boolean).map(l => {
				const [nim, nama, keterangan] = l.split(',').map(s => s.trim());
				return { nim, nama, keterangan: keterangan || '' };
			});
			plagiarismData = { ...plagiarismData, rows, updatedAt: new Date().toLocaleDateString('id-ID') };
			importCsv = '';
		} catch { /* ignore */ }
	}
	function exportCsv() {
		const csv = plagiarismData.rows.map(r => `${r.nim},${r.nama},${r.keterangan}`).join('\n');
		const a = Object.assign(document.createElement('a'), {
			href: URL.createObjectURL(new Blob([csv], { type:'text/csv' })),
			download: 'plagiarisme.csv'
		});
		a.click(); URL.revokeObjectURL(a.href);
	}
</script>

<div class="space-y-8 pb-16 max-w-5xl">

	<!-- Header -->
	<div class="flex flex-wrap items-center justify-between gap-4">
		<div>
			<h1 class="text-2xl font-extrabold text-ink-heading">Pengumuman</h1>
			<p class="text-sm text-ink-body mt-1">Kelola kartu pengumuman yang tampil di papan informasi halaman /info.</p>
		</div>
		<div class="flex items-center gap-3">
			{#if annSaved}
				<span class="text-xs font-bold text-emerald-700 flex items-center gap-1.5 bg-emerald-50 border border-emerald-200 px-3 py-1.5 rounded-lg">
					<Save size={13}/> Tersimpan!
				</span>
			{/if}
			<button onclick={saveAnnouncements} class="btn-primary text-sm flex items-center gap-2">
				<Save size={15}/> Simpan Pengumuman
			</button>
		</div>
	</div>

	<!-- ── 1. RESCHEDULE ── -->
	<section id="reschedule" class="bg-white rounded-2xl border border-amber-100 shadow-sm overflow-hidden">
		<div class="px-6 py-4 border-b border-amber-50 bg-amber-50/50 flex items-center gap-3">
			<div class="w-8 h-8 rounded-lg flex items-center justify-center text-white shrink-0" style="background:linear-gradient(135deg,#f97316,#ea580c)"><Calendar size={15}/></div>
			<div class="flex-1">
				<h2 class="text-base font-extrabold text-slate-900">Reschedule Praktikum</h2>
				<p class="text-xs text-slate-500 font-semibold">Kartu ini tampil di sisi kanan papan informasi.</p>
			</div>
			<span class="text-xs font-bold text-amber-700 bg-amber-50 border border-amber-200 px-2.5 py-1 rounded-full shrink-0">{reschedules.length} kartu</span>
		</div>
		<div class="p-6 space-y-4">
			{#each reschedules as rs, i}
				<div class="bg-amber-50 border border-amber-200 rounded-xl p-4 flex gap-3 items-start">
					<div class="flex-1 grid grid-cols-1 sm:grid-cols-2 gap-3">
						<div><label class="label text-xs">Kelas</label><input type="text" class="input" bind:value={reschedules[i].kelas}/></div>
						<div><label class="label text-xs">Hari Baru</label><input type="text" class="input" bind:value={reschedules[i].hari} placeholder="Sabtu"/></div>
						<div><label class="label text-xs">Jam</label><input type="text" class="input" bind:value={reschedules[i].jam} placeholder="07:20 WIB"/></div>
						<div><label class="label text-xs">Catatan</label><input type="text" class="input" bind:value={reschedules[i].catatan}/></div>
					</div>
					<button onclick={() => removeReschedule(rs.id)} class="p-2 rounded-lg bg-red-50 text-red-500 hover:bg-red-100 transition-colors shrink-0 mt-6"><Trash2 size={15}/></button>
				</div>
			{/each}
			<div class="border border-dashed border-slate-300 rounded-xl p-4 space-y-3 bg-slate-50">
				<p class="text-xs font-extrabold text-slate-600 uppercase tracking-wider">+ Tambah Kartu</p>
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
					<div><label class="label text-xs">Kelas</label><input type="text" class="input" bind:value={newRs.kelas} placeholder="mis: TL A"/></div>
					<div><label class="label text-xs">Hari Baru</label><input type="text" class="input" bind:value={newRs.hari} placeholder="mis: Sabtu"/></div>
					<div><label class="label text-xs">Jam</label><input type="text" class="input" bind:value={newRs.jam} placeholder="mis: 07:20 WIB"/></div>
					<div><label class="label text-xs">Catatan</label><input type="text" class="input" bind:value={newRs.catatan} placeholder="Harap hadir tepat waktu"/></div>
				</div>
				<button onclick={addReschedule} class="btn-primary text-sm flex items-center gap-1.5"><Plus size={14}/> Tambah</button>
			</div>
		</div>
	</section>

	<!-- ── 2. OPEN RECRUITMENT ── -->
	<section id="rekrutmen" class="bg-white rounded-2xl border border-blue-100 shadow-sm overflow-hidden">
		<div class="px-6 py-4 border-b border-blue-50 bg-blue-50/50 flex items-center gap-3">
			<div class="w-8 h-8 rounded-lg flex items-center justify-center text-white shrink-0" style="background:linear-gradient(135deg,#48CAE4,#2563EB)"><UserCheck size={15}/></div>
			<div class="flex-1">
				<h2 class="text-base font-extrabold text-slate-900">Open Recruitment</h2>
				<p class="text-xs text-slate-500 font-semibold">Pengumuman rekrutmen asisten/kepanitiaan.</p>
			</div>
			<span class="text-xs font-bold text-blue-700 bg-blue-50 border border-blue-200 px-2.5 py-1 rounded-full shrink-0">{recruitCards.length} kartu</span>
		</div>
		<div class="p-6 space-y-4">
			{#each recruitCards as rc, i}
				<div class="bg-blue-50 border border-blue-200 rounded-xl p-4 flex gap-3 items-start">
					<div class="flex-1 space-y-2">
						<div><label class="label text-xs">Judul</label><input type="text" class="input" bind:value={recruitCards[i].judul}/></div>
						<div><label class="label text-xs">Deskripsi</label><textarea class="input h-16 resize-none" bind:value={recruitCards[i].deskripsi}></textarea></div>
						<div><label class="label text-xs">Link (opsional)</label><input type="url" class="input" bind:value={recruitCards[i].link} placeholder="https://forms.gle/..."/></div>
					</div>
					<button onclick={() => removeRecruit(rc.id)} class="p-2 rounded-lg bg-red-50 text-red-500 hover:bg-red-100 transition-colors shrink-0 mt-6"><Trash2 size={15}/></button>
				</div>
			{/each}
			<div class="border border-dashed border-slate-300 rounded-xl p-4 space-y-3 bg-slate-50">
				<p class="text-xs font-extrabold text-slate-600 uppercase tracking-wider">+ Tambah Kartu</p>
				<div><label class="label text-xs">Judul</label><input type="text" class="input" bind:value={newRc.judul} placeholder="Open Recruitment Asisten 2025"/></div>
				<div><label class="label text-xs">Deskripsi</label><textarea class="input h-16 resize-none" bind:value={newRc.deskripsi} placeholder="Deskripsi singkat…"></textarea></div>
				<div><label class="label text-xs">Link (opsional)</label><input type="url" class="input" bind:value={newRc.link} placeholder="https://..."/></div>
				<button onclick={addRecruit} class="btn-primary text-sm flex items-center gap-1.5"><Plus size={14}/> Tambah</button>
			</div>
		</div>
	</section>

	<!-- ── 3. SUSULAN ── -->
	<section id="susulan" class="bg-white rounded-2xl border border-purple-100 shadow-sm overflow-hidden">
		<div class="px-6 py-4 border-b border-purple-50 bg-purple-50/50 flex items-center gap-3">
			<div class="w-8 h-8 rounded-lg flex items-center justify-center text-white shrink-0" style="background:linear-gradient(135deg,#9D4EDD,#4F46E5)"><Copy size={15}/></div>
			<div class="flex-1">
				<h2 class="text-base font-extrabold text-slate-900">Pengumuman Susulan</h2>
				<p class="text-xs text-slate-500 font-semibold">Susulan laporan analisis &amp; jurnal per kelas.</p>
			</div>
			<span class="text-xs font-bold text-purple-700 bg-purple-50 border border-purple-200 px-2.5 py-1 rounded-full shrink-0">{susulanCards.length} kartu</span>
		</div>
		<div class="p-6 space-y-4">
			{#each susulanCards as su, i}
				<div class="bg-purple-50 border border-purple-200 rounded-xl p-4 flex gap-3 items-start">
					<div class="flex-1 space-y-2">
						<div><label class="label text-xs">Kelas</label><input type="text" class="input" bind:value={susulanCards[i].kelas}/></div>
						<div><label class="label text-xs">Deskripsi</label><textarea class="input h-16 resize-none" bind:value={susulanCards[i].deskripsi}></textarea></div>
					</div>
					<button onclick={() => removeSusulan(su.id)} class="p-2 rounded-lg bg-red-50 text-red-500 hover:bg-red-100 transition-colors shrink-0 mt-6"><Trash2 size={15}/></button>
				</div>
			{/each}
			<div class="border border-dashed border-slate-300 rounded-xl p-4 space-y-3 bg-slate-50">
				<p class="text-xs font-extrabold text-slate-600 uppercase tracking-wider">+ Tambah Kartu</p>
				<div><label class="label text-xs">Kelas</label><input type="text" class="input" bind:value={newSu.kelas} placeholder="mis: TSE B"/></div>
				<div><label class="label text-xs">Deskripsi</label><textarea class="input h-16 resize-none" bind:value={newSu.deskripsi} placeholder="Detail susulan laporan/jurnal…"></textarea></div>
				<button onclick={addSusulan} class="btn-primary text-sm flex items-center gap-1.5"><Plus size={14}/> Tambah</button>
			</div>
		</div>
	</section>

	<!-- ── 4. DATA PLAGIARISME ── -->
	<section id="plagiarisme" class="bg-white rounded-2xl border border-red-100 shadow-sm overflow-hidden">
		<div class="px-6 py-4 border-b border-red-50 bg-red-50/50 flex items-center gap-3">
			<div class="w-8 h-8 rounded-lg flex items-center justify-center text-white shrink-0" style="background:linear-gradient(135deg,#ef4444,#e11d48)"><AlertTriangle size={15}/></div>
			<div class="flex-1">
				<h2 class="text-base font-extrabold text-slate-900">Data Plagiarisme</h2>
				<p class="text-xs text-slate-500 font-semibold">Publik hanya bisa lihat &amp; unduh. Import/edit hanya di sini.</p>
			</div>
			<span class="text-xs font-bold text-red-700 bg-red-50 border border-red-200 px-2.5 py-1 rounded-full shrink-0">{plagiarismData.rows.length} mahasiswa</span>
		</div>
		<div class="p-6 space-y-5">
			<div class="bg-slate-50 border border-dashed border-slate-300 rounded-xl p-4 space-y-3">
				<h3 class="text-sm font-extrabold text-slate-800 flex items-center gap-2"><Upload size={14}/> Import dari CSV</h3>
				<p class="text-xs text-slate-600 font-semibold">Format: <code class="bg-white border border-slate-200 px-1.5 py-0.5 rounded text-xs font-mono">NIM,Nama,Keterangan</code> — satu baris per mahasiswa.</p>
				<textarea class="input h-24 resize-none font-mono text-xs" bind:value={importCsv}
					placeholder="2211001,Budi Santoso,Posttest Modul 3&#10;2211002,Ani Wati,Jurnal Modul 2"></textarea>
				<button onclick={importFromCsv} class="btn-primary text-sm flex items-center gap-1.5"><Upload size={14}/> Terapkan Import</button>
			</div>

			<div class="bg-slate-50 border border-slate-200 rounded-xl p-4 space-y-2">
				<label class="block text-xs font-black text-slate-600 uppercase tracking-wider">Catatan / Instruksi Panggilan Kehadiran</label>
				<p class="text-[10px] text-slate-500 font-semibold leading-relaxed">Catatan/instruksi ini akan muncul di bagian bawah daftar pelanggar pada modal plagiarisme halaman depan mahasiswa.</p>
				<textarea class="input h-16 resize-none font-semibold text-xs" bind:value={plagiarismData.catatan}
					placeholder="Bagi nama-nama yang disebutkan di atas, silakan hadir di Lab. Algoritma & Pemrograman pada Rabu, 08 Juli 2026 jam 08:00 Pagi"></textarea>
			</div>

			{#if plagiarismData.rows.length > 0}
				<div class="flex items-center justify-between">
					<p class="text-xs font-bold text-slate-500">Diperbarui: {plagiarismData.updatedAt || '—'}</p>
					<button onclick={exportCsv} class="btn-outline text-xs flex items-center gap-1.5"><Download size={12}/> Ekspor CSV</button>
				</div>
				<div class="table-wrap">
					<table class="tbl">
						<thead><tr><th>NIM</th><th>Nama</th><th>Keterangan</th><th></th></tr></thead>
						<tbody>
							{#each plagiarismData.rows as row}
								<tr>
									<td class="font-mono font-bold text-slate-900">{row.nim}</td>
									<td class="font-semibold text-slate-800">{row.nama}</td>
									<td class="text-slate-700">{row.keterangan}</td>
									<td><button onclick={() => removePlagiarismRow(row.nim)} class="p-1.5 rounded text-red-500 hover:bg-red-50 transition-colors"><Trash2 size={13}/></button></td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}

			<div class="border border-dashed border-slate-300 rounded-xl p-4 space-y-3 bg-slate-50">
				<p class="text-xs font-extrabold text-slate-600 uppercase tracking-wider">+ Tambah Baris Manual</p>
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
					<div><label class="label text-xs">NIM</label><input type="text" class="input font-mono" bind:value={newPl.nim} placeholder="2211001"/></div>
					<div><label class="label text-xs">Nama</label><input type="text" class="input" bind:value={newPl.nama} placeholder="Nama Mahasiswa"/></div>
					<div><label class="label text-xs">Keterangan</label><input type="text" class="input" bind:value={newPl.keterangan} placeholder="Posttest Modul 3"/></div>
				</div>
				<button onclick={addPlagiarismRow} class="btn-primary text-sm flex items-center gap-1.5"><Plus size={14}/> Tambah Baris</button>
			</div>
		</div>
	</section>

	<div class="flex justify-end">
		<button onclick={saveAnnouncements} class="btn-primary text-sm flex items-center gap-2"><Save size={15}/> Simpan Pengumuman</button>
	</div>

</div>
