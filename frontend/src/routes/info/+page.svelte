<script lang="ts">
	import { Calendar, Users, FileText, Book, ArrowRight, Bell, Code } from 'lucide-svelte';

	const quickLinks = [
		{ href: '/info/jadwal', title: 'Jadwal Praktikum', desc: 'Lihat jadwal per kelas & shift.', icon: Calendar },
		{ href: '/info/asisten', title: 'Daftar Asisten', desc: 'Profil & kontak asisten lab.', icon: Users },
		{ href: '/info/laporan', title: 'Pedoman Laporan', desc: 'Unduh template & pedoman.', icon: FileText },
		{ href: '/info/modul', title: 'Modul Praktikum', desc: 'Unduh modul praktikum (PDF).', icon: Book }
	];
</script>

<!-- Latar aurora glassmorphism -->
<div class="aurora" aria-hidden="true">
	<div class="aurora-blob aurora-blob-1"></div>
	<div class="aurora-blob aurora-blob-2"></div>
	<div class="aurora-blob aurora-blob-3"></div>
</div>

<div class="relative z-10">
	<!-- Hero kaca -->
	<section class="glass-hero fade-up flex flex-col items-start gap-8 p-8 md:flex-row md:items-center md:p-12">
		<div class="flex-1">
			<span class="badge bg-primary/10 text-primary">Portal Praktikum</span>
			<h1 class="mt-4 font-display text-4xl font-extrabold leading-tight tracking-tight text-ink-heading text-balance md:text-5xl">
				Laboratorium Algoritma <span class="text-primary">&amp; Pemrograman</span>
			</h1>
			<p class="mt-4 max-w-xl text-lg leading-relaxed text-ink-body">
				Portal informasi terpadu dan platform ujian praktikum. Temukan jadwal kelas, informasi asisten, pedoman laporan, hingga fasilitas live coding dalam satu tempat.
			</p>
			<div class="mt-8 flex flex-wrap gap-3">
				<a href="/praktikum/login" class="btn-primary px-8 py-3 text-base font-bold">
					Masuk ke Praktikum <ArrowRight size={18} />
				</a>
				<a href="/info/jadwal" class="btn-outline border-white/60 bg-white/40 px-6 py-3">Lihat Jadwal</a>
			</div>
		</div>
		<div class="hidden lg:block">
			<span class="glass-badge h-40 w-40 rotate-3 transition-transform duration-500 hover:rotate-0" style="border-radius:1.9rem">
				<Code size={64} />
			</span>
		</div>
	</section>

	<!-- Quick Links -->
	<section class="mt-12">
		<h2 class="font-display text-2xl font-bold tracking-tight text-ink-heading">Quick Links</h2>
		<div class="mt-6 grid gap-[18px] sm:grid-cols-2 lg:grid-cols-4">
			{#each quickLinks as q, i}
				{@const Icon = q.icon}
				<a href={q.href} class="glass-card glass-hover fade-up group flex flex-col p-6" style="animation-delay:{i * 0.06}s">
					<span class="glass-badge-soft h-12 w-12"><Icon size={22} /></span>
					<h3 class="mt-4 font-display text-base font-bold text-ink-heading transition-colors group-hover:text-primary">{q.title}</h3>
					<p class="mt-1 text-sm leading-relaxed text-ink-caption">{q.desc}</p>
				</a>
			{/each}
		</div>
	</section>

	<!-- Pengumuman -->
	<section class="glass-card fade-up mt-12 flex items-start gap-4 p-6" style="animation-delay:.28s">
		<span class="glass-badge-soft h-11 w-11 flex-shrink-0"><Bell size={20} /></span>
		<div>
			<h2 class="font-display text-lg font-bold text-ink-heading">Pengumuman Penting</h2>
			<p class="mt-2 text-sm leading-relaxed text-ink-body">
				Pastikan Anda sudah melakukan registrasi akun praktikum menggunakan <strong>NIM</strong> yang valid. Jika akses registrasi untuk kelas Anda belum dibuka, silakan hubungi asisten laboratorium yang bertugas.
			</p>
		</div>
	</section>
</div>
