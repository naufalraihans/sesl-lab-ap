<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { Download, BookOpen, Calendar, FileText, Users, ArrowRight, Lightbulb, Eye, EyeOff, Book, Code, Cpu, Terminal } from 'lucide-svelte';

	let fileUrl = $state('');
	let loading = $state(true);
	let showPreview = $state(false);

	const modulCards = [
		{ no: 1, judul: 'Dasar Logika Algoritma & Pemrograman Bahasa C', topik: 'Pengenalan algoritma, flowchart, dasar C', fromColor: '#48CAE4', toColor: '#0284C7', borderColor: '#bae6fd', icon: Terminal },
		{ no: 2, judul: 'Struktur Kontrol dalam Bahasa C', topik: 'if/else, switch, for, while, do-while', fromColor: '#9D4EDD', toColor: '#4F46E5', borderColor: '#ddd6fe', icon: Code },
		{ no: 3, judul: 'Array, Struct dan Operasi File', topik: 'Array 1D/2D, struct, fopen/fclose', fromColor: '#06D6A0', toColor: '#059669', borderColor: '#a7f3d0', icon: Cpu },
		{ no: 4, judul: 'Pengenalan Dasar Bahasa Python', topik: 'Sintaks Python, variabel, tipe data, I/O', fromColor: '#FFC300', toColor: '#EA580C', borderColor: '#fed7aa', icon: BookOpen },
		{ no: 5, judul: 'Percabangan dan Perulangan pada Bahasa Python', topik: 'if/elif/else, for, while, list comprehension', fromColor: '#F472B6', toColor: '#E11D48', borderColor: '#fecdd3', icon: Book },
		{ no: 6, judul: 'List, Dictionary dan Operasi File', topik: 'List, dict, tuple, baca/tulis file', fromColor: '#14B8A6', toColor: '#0E7490', borderColor: '#99f6e4', icon: FileText },
	];

	const tautan = [
		{ href: '/info/jadwal', title: 'Jadwal Praktikum', desc: 'Lihat jadwal per kelas & shift.', icon: Calendar, fromColor: '#60A5FA', toColor: '#2563EB', borderColor: '#bfdbfe' },
		{ href: '/info/laporan', title: 'Pedoman Laporan', desc: 'Unduh template & pedoman laporan.', icon: FileText, fromColor: '#9D4EDD', toColor: '#4F46E5', borderColor: '#ddd6fe' },
		{ href: '/info/asisten', title: 'Asisten Lab', desc: 'Kontak & profil asisten.', icon: Users, fromColor: '#06D6A0', toColor: '#059669', borderColor: '#a7f3d0' },
	];

	onMount(async () => {
		try {
			const res = await api.get<{ file_url: string }>('/api/info/modul');
			fileUrl = res?.file_url ?? '';
		} finally {
			loading = false;
		}
	});
</script>

<div class="space-y-8">
	<div class="flex items-center gap-3 text-left">
		<div class="w-10 h-10 rounded-xl text-white flex items-center justify-center shadow-md" style="background: linear-gradient(135deg, #06D6A0, #059669);">
			<BookOpen size={20} />
		</div>
		<div>
			<h1 class="text-2xl font-extrabold text-slate-900 leading-none">Modul Praktikum</h1>
			<p class="mt-1 text-xs font-bold text-slate-600">Materi resmi Algoritma &amp; Pemrograman. Baca sebelum sesi dimulai.</p>
		</div>
	</div>

	{#if loading}
		<p class="text-slate-700 font-semibold">Memuat…</p>
	{:else}
		<!-- Download card + tips -->
		<div class="grid gap-6 lg:grid-cols-3">
			<div class="card lg:col-span-2 flex flex-col gap-5 sm:flex-row sm:items-center text-left" style="border-color: #a7f3d0; background: linear-gradient(to bottom right, white, rgba(236,253,245,0.4));">
				<div class="grid h-20 w-20 flex-shrink-0 place-items-center rounded-2xl text-white shadow-lg" style="background: linear-gradient(135deg, #06D6A0, #059669);">
					<BookOpen size={36} />
				</div>
				<div class="flex-1">
					<h2 class="text-lg font-extrabold text-slate-900">Modul Praktikum (PDF Lengkap)</h2>
					{#if fileUrl}
						<p class="mt-1 text-sm text-slate-700 font-bold leading-relaxed">Unduh modul lengkap dalam format PDF untuk dipelajari sebelum praktikum.</p>
						<div class="flex flex-wrap gap-2.5 mt-4">
							<a href={fileUrl} target="_blank" rel="noopener" class="btn-primary py-2 text-sm font-bold flex items-center gap-1.5"><Download size={14}/> Download Modul (PDF)</a>
							<button onclick={() => showPreview = !showPreview}
								class="btn bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold px-4 py-2 rounded-xl text-sm transition-all border border-slate-200 flex items-center gap-1.5">
								{#if showPreview}<EyeOff size={14}/> Tutup Preview{:else}<Eye size={14}/> Tampilkan Preview{/if}
							</button>
						</div>
					{:else}
						<p class="mt-1 text-sm text-slate-700 font-bold">Modul belum diunggah oleh admin. Silakan cek kembali nanti.</p>
					{/if}
				</div>
			</div>
			<div class="card text-left" style="background: linear-gradient(to bottom right, #eff6ff, #eef2ff); border-color: #bfdbfe;">
				<div class="mb-3 flex items-center gap-2" style="color: #1e40af;">
					<Lightbulb size={18}/>
					<h3 class="font-extrabold" style="color: #1e3a8a;">Petunjuk Belajar</h3>
				</div>
				<ul class="space-y-2 text-xs font-bold text-slate-700 leading-relaxed">
					<li>• Baca modul sebelum praktikum dimulai.</li>
					<li>• Siapkan pertanyaan untuk asisten.</li>
					<li>• Modul jadi acuan pre-test &amp; post-test.</li>
				</ul>
			</div>
		</div>

		{#if showPreview && fileUrl}
			<div class="card p-5 bg-white border border-slate-200 shadow-lg text-left space-y-4">
				<div class="flex items-center justify-between border-b border-slate-100 pb-3">
					<div>
						<h3 class="font-extrabold text-slate-900 text-base">Preview Modul Praktikum</h3>
						<p class="text-xs text-slate-600 font-bold mt-0.5">Membaca dokumen secara real-time.</p>
					</div>
					<button onclick={() => showPreview = false} class="btn-outline px-4 py-1.5 text-xs font-bold">Tutup Preview</button>
				</div>
				<div class="w-full bg-slate-50 border border-slate-200 rounded-2xl overflow-hidden min-h-[600px]">
					{#if fileUrl.toLowerCase().endsWith('.pdf')}
						<iframe title="Modul PDF Reader" src={fileUrl} class="w-full h-[650px] border-0" allowfullscreen={true}></iframe>
					{:else}
						<iframe title="Modul Doc Reader" src="https://docs.google.com/gview?url={encodeURIComponent(fileUrl)}&embedded=true" class="w-full h-[650px] border-0" allowfullscreen={true}></iframe>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Module cards -->
		<section>
			<h2 class="text-lg font-extrabold text-slate-900 mb-5">Daftar Modul</h2>
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each modulCards as m}
					{@const Icon = m.icon}
					<div class="group relative bg-white rounded-2xl p-5 shadow-sm hover:shadow-lg hover:-translate-y-1 transition-all duration-300 overflow-hidden border"
						style="border-color: {m.borderColor};">
						<div class="absolute top-0 left-0 right-0 h-1 rounded-t-2xl"
							style="background: linear-gradient(90deg, {m.fromColor}, {m.toColor});"></div>
						<div class="absolute inset-0 rounded-2xl pointer-events-none opacity-0 group-hover:opacity-5 transition-opacity"
							style="background: linear-gradient(135deg, {m.fromColor}, {m.toColor});"></div>
						<div class="relative z-10">
							<div class="w-9 h-9 rounded-xl text-white flex items-center justify-center shadow-md mb-3"
								style="background: linear-gradient(135deg, {m.fromColor}, {m.toColor});">
								<Icon size={16}/>
							</div>
							<p class="text-[10px] font-black text-slate-500 uppercase tracking-wider mb-1">Modul {m.no}</p>
							<h3 class="font-extrabold text-slate-900 text-sm leading-snug mb-1">{m.judul}</h3>
							<p class="text-[11px] font-bold text-slate-600 leading-relaxed">{m.topik}</p>
						</div>
					</div>
				{/each}
			</div>
		</section>

		<!-- Related links -->
		<section>
			<h2 class="mb-4 text-lg font-extrabold text-slate-900">Informasi Lainnya</h2>
			<div class="grid gap-5 sm:grid-cols-3">
				{#each tautan as t}
					{@const Icon = t.icon}
					<a href={t.href} class="group block text-left bg-white rounded-2xl p-6 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all border"
						style="border-color: {t.borderColor};">
						<div class="mb-3 inline-flex rounded-xl p-3 text-white shadow-sm"
							style="background: linear-gradient(135deg, {t.fromColor}, {t.toColor});">
							<Icon size={22}/>
						</div>
						<h3 class="flex items-center gap-1 text-base font-extrabold text-slate-900">
							{t.title} <ArrowRight size={15} class="opacity-0 group-hover:opacity-100 transition-opacity"/>
						</h3>
						<p class="mt-1 text-xs font-bold leading-relaxed text-slate-600">{t.desc}</p>
					</a>
				{/each}
			</div>
		</section>
	{/if}
</div>
