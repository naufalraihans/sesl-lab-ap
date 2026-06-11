<script lang="ts">
	import { Calendar, Users, FileText, Book, ArrowRight, Bell, Code } from 'lucide-svelte';

	const quickLinks = [
		{ href: '/info/jadwal', title: 'Jadwal Praktikum', desc: 'Lihat jadwal per kelas & shift.', icon: Calendar },
		{ href: '/info/asisten', title: 'Daftar Asisten', desc: 'Profil & kontak asisten lab.', icon: Users },
		{ href: '/info/laporan', title: 'Pedoman Laporan', desc: 'Unduh template & pedoman.', icon: FileText },
		{ href: '/info/modul', title: 'Modul Praktikum', desc: 'Unduh modul praktikum (PDF).', icon: Book }
	];
</script>

<!-- Hero Section -->
<section class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-gray-900 via-gray-800 to-primary-dark p-8 md:p-12 shadow-2xl animate-fade-in">
	<!-- Background decorations -->
	<div class="absolute -top-24 -right-24 h-64 w-64 rounded-full bg-primary/20 blur-3xl filter animate-float"></div>
	<div class="absolute -bottom-24 -left-24 h-64 w-64 rounded-full bg-blue-500/20 blur-3xl filter animate-float" style="animation-delay: -2s;"></div>
	
	<div class="relative z-10 flex flex-col md:flex-row items-center gap-8">
		<div class="flex-1">
			<h1 class="text-4xl md:text-5xl font-extrabold tracking-tight text-white leading-tight">
				Laboratorium Algoritma <br/><span class="text-transparent bg-clip-text bg-gradient-to-r from-primary-hover to-red-400">&amp; Pemrograman</span>
			</h1>
			<p class="mt-4 max-w-xl text-lg text-gray-300 leading-relaxed">
				Portal informasi terpadu dan platform ujian praktikum. Temukan jadwal kelas, informasi asisten, pedoman laporan, hingga fasilitas live coding dalam satu tempat.
			</p>
			<div class="mt-8 flex flex-wrap gap-4">
				<a href="/praktikum/login" class="btn-primary btn-glow px-8 py-3 text-base font-bold shadow-primary/40">
					Masuk ke Praktikum <ArrowRight size={18} />
				</a>
				<a href="/info/jadwal" class="btn bg-white/10 text-white hover:bg-white/20 backdrop-blur-md border border-white/10 px-6 py-3">
					Lihat Jadwal
				</a>
			</div>
		</div>
		<div class="hidden lg:block">
			<div class="grid h-40 w-40 place-items-center rounded-3xl bg-gradient-to-tr from-primary to-primary-dark shadow-[0_0_40px_rgba(208,49,83,0.4)] rotate-3 hover:rotate-0 transition-transform duration-500">
				<Code size={64} class="text-white" />
			</div>
		</div>
	</div>
</section>

<section class="mt-12">
	<div class="mb-6 flex items-center justify-between">
		<h2 class="text-2xl font-bold tracking-tight text-ink-heading">Quick Links</h2>
	</div>
	<div class="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
		{#each quickLinks as q, i}
			{@const Icon = q.icon}
			<a href={q.href} class="card hover-card group block animate-fade-in-up" style="animation-delay: {i * 100}ms">
				<div class="mb-4 inline-flex rounded-xl bg-primary/10 p-3 text-primary transition-colors group-hover:bg-primary group-hover:text-white">
					<Icon size={24} />
				</div>
				<h3 class="text-lg font-bold text-ink-heading transition-colors group-hover:text-primary">{q.title}</h3>
				<p class="mt-1.5 text-sm leading-relaxed text-ink-caption">{q.desc}</p>
			</a>
		{/each}
	</div>
</section>

<section class="mt-12 overflow-hidden rounded-2xl border border-blue-200 bg-gradient-to-r from-blue-50 to-white shadow-sm animate-fade-in-up" style="animation-delay: 400ms">
	<div class="flex items-start gap-4 p-6">
		<div class="mt-1 flex-shrink-0 rounded-full bg-blue-100 p-2 text-blue-600">
			<Bell size={20} />
		</div>
		<div>
			<h2 class="text-lg font-bold text-blue-900">Pengumuman Penting</h2>
			<p class="mt-2 text-sm leading-relaxed text-blue-800/80">
				Pastikan Anda sudah melakukan registrasi akun praktikum menggunakan <strong>NIM</strong> yang valid. Jika akses registrasi untuk kelas Anda belum dibuka, silakan hubungi asisten laboratorium yang bertugas.
			</p>
		</div>
	</div>
</section>
