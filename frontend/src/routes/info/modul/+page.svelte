<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { Download, BookOpen, Calendar, FileText, Users, ArrowRight, Lightbulb } from 'lucide-svelte';

	let fileUrl = $state('');
	let loading = $state(true);

	const tautan = [
		{ href: '/info/jadwal', title: 'Jadwal Praktikum', desc: 'Lihat jadwal per kelas & shift.', icon: Calendar },
		{ href: '/info/laporan', title: 'Pedoman Laporan', desc: 'Unduh template & pedoman laporan.', icon: FileText },
		{ href: '/info/asisten', title: 'Asisten Lab', desc: 'Kontak & profil asisten.', icon: Users }
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

<h1 class="mb-1 text-2xl font-bold text-ink-heading">Modul Praktikum</h1>
<p class="mb-6 text-sm text-ink-caption">Materi resmi praktikum Algoritma &amp; Pemrograman. Baca sebelum sesi dimulai.</p>

{#if loading}
	<p class="text-ink-caption">Memuat…</p>
{:else}
	<div class="grid gap-6 lg:grid-cols-3">
		<!-- Kartu modul utama -->
		<div class="card lg:col-span-2 flex flex-col gap-5 sm:flex-row sm:items-center">
			<div class="grid h-20 w-20 flex-shrink-0 place-items-center rounded-2xl bg-primary/10 text-primary">
				<BookOpen size={36} />
			</div>
			<div class="flex-1">
				<h2 class="text-lg font-bold text-ink-heading">Modul Praktikum (PDF)</h2>
				{#if fileUrl}
					<p class="mt-1 text-sm text-ink-caption">Unduh modul lengkap dalam format PDF untuk dipelajari sebelum praktikum.</p>
					<a href={fileUrl} target="_blank" rel="noopener" class="btn-primary mt-4"><Download size={16} /> Download Modul (PDF)</a>
				{:else}
					<p class="mt-1 text-sm text-ink-caption">Modul belum diunggah oleh admin. Silakan cek kembali nanti.</p>
				{/if}
			</div>
		</div>

		<!-- Petunjuk -->
		<div class="card bg-blue-50/60">
			<div class="mb-2 flex items-center gap-2 text-blue-800">
				<Lightbulb size={18} />
				<h3 class="font-bold">Petunjuk</h3>
			</div>
			<ul class="space-y-1.5 text-sm leading-relaxed text-ink-body">
				<li>• Baca modul sebelum praktikum dimulai.</li>
				<li>• Siapkan pertanyaan untuk asisten.</li>
				<li>• Modul jadi acuan pre-test &amp; post-test.</li>
			</ul>
		</div>
	</div>

	<!-- Tautan informasi lain -->
	<section class="mt-10">
		<h2 class="mb-4 text-lg font-bold text-ink-heading">Informasi Lainnya</h2>
		<div class="grid gap-5 sm:grid-cols-3">
			{#each tautan as t}
				{@const Icon = t.icon}
				<a href={t.href} class="card hover-card group block">
					<div class="mb-3 inline-flex rounded-xl bg-primary/10 p-3 text-primary transition-colors group-hover:bg-primary group-hover:text-white">
						<Icon size={22} />
					</div>
					<h3 class="flex items-center gap-1 text-base font-bold text-ink-heading transition-colors group-hover:text-primary">
						{t.title} <ArrowRight size={15} class="opacity-0 transition-opacity group-hover:opacity-100" />
					</h3>
					<p class="mt-1 text-sm leading-relaxed text-ink-caption">{t.desc}</p>
				</a>
			{/each}
		</div>
	</section>
{/if}
