<script lang="ts">
	import { Calendar, Users, FileText, Book, ArrowRight, Bell } from 'lucide-svelte';

	const quickLinks = [
		{ href: '/info/jadwal', title: 'Jadwal Praktikum', desc: 'Lihat jadwal per kelas & shift.', icon: Calendar },
		{ href: '/info/asisten', title: 'Daftar Asisten', desc: 'Profil & kontak asisten lab.', icon: Users },
		{ href: '/info/laporan', title: 'Pedoman Laporan', desc: 'Unduh template & pedoman.', icon: FileText },
		{ href: '/info/modul', title: 'Modul Praktikum', desc: 'Unduh modul praktikum (PDF).', icon: Book }
	];
</script>

<div class="bento">
	<!-- Hero (2x2) -->
	<section class="b-hero glass-hero fade-up flex flex-col justify-between p-9 md:p-10">
		<div>
			<h1 class="font-sans text-4xl font-extrabold leading-[1.08] tracking-tight text-ink-heading text-balance">
				Laboratorium Algoritma<br /><span class="text-primary">&amp; Pemrograman</span>
			</h1>
			<p class="mt-4 max-w-md text-base leading-relaxed text-ink-body">
				Satu tempat untuk jadwal kelas, profil asisten, pedoman laporan, modul praktikum, hingga fasilitas live coding. Semua rapi dalam satu portal.
			</p>
		</div>
		<div class="mt-7 flex flex-wrap gap-3">
			<a href="/praktikum/login" class="btn-primary px-6 py-3 text-base font-bold">
				Masuk ke Praktikum <ArrowRight size={17} />
			</a>
			<a href="/info/jadwal" class="inline-flex items-center gap-2 rounded-2xl border border-white/70 bg-white/55 px-5 py-3 text-base font-semibold text-ink-heading backdrop-blur transition-transform hover:-translate-y-0.5">
				Lihat Jadwal
			</a>
		</div>
	</section>

	<!-- Quick link cards -->
	{#each quickLinks as q, i}
		{@const Icon = q.icon}
		<a
			href={q.href}
			class="glass-card glass-hover fade-up flex flex-col p-6"
			style="animation-delay: {60 + i * 60}ms"
		>
			<span class="glass-badge h-12 w-12">
				<Icon size={23} />
			</span>
			<h3 class="mt-4 font-sans text-base font-bold text-ink-heading">{q.title}</h3>
			<p class="mt-1.5 text-sm leading-relaxed text-ink-caption">{q.desc}</p>
		</a>
	{/each}

	<!-- Pengumuman (wide) -->
	<section class="b-wide glass-card fade-up flex items-start gap-4 p-7" style="animation-delay: 300ms">
		<span class="glass-badge-soft h-12 w-12 flex-shrink-0">
			<Bell size={22} />
		</span>
		<div>
			<h2 class="font-sans text-lg font-bold text-ink-heading">Pengumuman Penting</h2>
			<p class="mt-2 text-sm leading-relaxed text-ink-body">
				Pastikan Anda sudah melakukan registrasi akun praktikum menggunakan <strong>NIM</strong> yang valid. Jika akses registrasi untuk kelas Anda belum dibuka, silakan hubungi asisten laboratorium yang bertugas.
			</p>
		</div>
	</section>
</div>
