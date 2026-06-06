<script lang="ts">
	const quickLinks = [
		{ href: '/info/jadwal', title: 'Jadwal Praktikum', desc: 'Lihat jadwal per kelas & shift.' },
		{ href: '/info/asisten', title: 'Daftar Asisten', desc: 'Profil & kontak asisten lab.' },
		{ href: '/info/laporan', title: 'Pedoman Laporan', desc: 'Unduh template & pedoman.' },
		{ href: '/info/modul', title: 'Modul Praktikum', desc: 'Unduh modul praktikum (PDF).' }
	];
</script>

<section class="rounded-lg bg-surface-soft p-8">
	<h1 class="text-3xl">Selamat Datang di Laboratorium Algoritma &amp; Pemrograman</h1>
	<p class="mt-3 max-w-2xl text-ink-body">
		Portal informasi dan praktikum Lab Algoritma &amp; Pemrograman. Temukan jadwal, daftar asisten,
		pedoman laporan, dan modul praktikum. Mahasiswa dapat mengerjakan pre-test, post-test,
		keterampilan, dan ujian praktik melalui halaman praktikum.
	</p>
	<div class="mt-5 flex gap-3">
		<a href="/praktikum/login" class="btn-primary">Masuk ke Praktikum</a>
		<a href="/info/jadwal" class="btn-outline">Lihat Jadwal</a>
	</div>
</section>

<section class="mt-8">
	<h2 class="mb-4 text-xl">Quick Links</h2>
	<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
		{#each quickLinks as q}
			<a href={q.href} class="card block">
				<h3 class="text-lg">{q.title}</h3>
				<p class="mt-1 text-sm text-ink-caption">{q.desc}</p>
			</a>
		{/each}
	</div>
</section>

<section class="mt-8 rounded-lg border border-state-info-bg bg-state-info-bg/40 p-5">
	<h2 class="text-lg text-state-info">📢 Pengumuman</h2>
	<p class="mt-1 text-sm text-ink-body">
		Pastikan Anda sudah melakukan registrasi akun praktikum menggunakan NIM. Jika akses register
		belum dibuka, hubungi asisten lab.
	</p>
</section>
