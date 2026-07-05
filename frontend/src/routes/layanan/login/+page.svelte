<script lang="ts">
	import { 
		User, Lock, Phone, Info, HelpCircle, Calendar, ArrowRight, ArrowLeft, 
		X, CheckCircle2, AlertCircle, Sparkles
	} from 'lucide-svelte';
	import { slide, fade } from 'svelte/transition';

	// Tab State
	let activeTab = $state<'masuk' | 'daftar'>('masuk');
	let isLoading = $state(false);
	let errorMsg = $state('');
	let successMsg = $state('');

	// Input States
	let loginUsername = $state('');
	let loginPassword = $state('');
	let showLoginPassword = $state(false);

	let signupName = $state('');
	let signupUsername = $state('');
	let signupPassword = $state('');
	let showSignupPassword = $state(false);
	let signupPhone = $state('');
	let signupRole = $state<'praktikan' | 'penyewa'>('praktikan');

	// Modal Dialogs States
	let openPanduan = $state(false);
	let openProsedur = $state(false);
	let openReschedule = $state(false);

	// Mock/Static Announcements
	const announcements = [
		{ type: 'oprec', text: '📢 OPEN RECRUITMENT ASISTEN LAB AP SEDANG DIBUKA! DAFTAR SEKARANG!' },
		{ type: 'info', text: '✨ AKSES LAYANAN SEWA ALAT DAN JADWAL JAGA AKTIF' },
		{ type: 'reschedule', text: '⚠️ PERUBAHAN JADWAL WAJIB DIKONFIRMASI H-1 KEPADA ASISTEN JAGA' }
	];

	async function handleLogin(e: SubmitEvent) {
		e.preventDefault();
		isLoading = true;
		errorMsg = '';
		successMsg = '';

		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username: loginUsername, password: loginPassword })
			});

			const result = await response.json();
			if (!response.ok || !result.success) {
				throw new Error(result.message || 'Gagal login');
			}

			localStorage.setItem('lab_jwt_token', result.data.token);
			localStorage.setItem('lab_user', JSON.stringify(result.data.user));

			successMsg = `Login Berhasil! Selamat datang, ${result.data.user.full_name || result.data.user.username}`;
			
			setTimeout(() => {
				window.location.href = '/layanan/beranda';
			}, 1000);
		} catch (err: any) {
			errorMsg = err.message || 'Terjadi kesalahan sistem';
		} finally {
			isLoading = false;
		}
	}

	async function handleSignup(e: SubmitEvent) {
		e.preventDefault();
		isLoading = true;
		errorMsg = '';
		successMsg = '';

		if (signupPassword.length < 6) {
			errorMsg = 'Password minimal 6 karakter.';
			isLoading = false;
			return;
		}

		if (signupRole === 'penyewa' && !signupPhone) {
			errorMsg = 'Nomor WhatsApp wajib diisi untuk verifikasi sewa.';
			isLoading = false;
			return;
		}

		try {
			const response = await fetch('/api/auth/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					username: signupUsername,
					password: signupPassword,
					full_name: signupName,
					role: signupRole,
					phone_number: signupPhone
				})
			});

			const result = await response.json();
			if (!response.ok || !result.success) {
				throw new Error(result.message || 'Gagal mendaftar');
			}

			successMsg = 'Pendaftaran berhasil! Silakan masuk dengan akun Anda.';
			
			signupName = '';
			signupUsername = '';
			signupPassword = '';
			signupPhone = '';
			activeTab = 'masuk';
		} catch (err: any) {
			errorMsg = err.message || 'Gagal melakukan pendaftaran';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="relative flex min-h-screen items-center justify-center overflow-hidden bg-slate-50 px-4 py-20 animate-fade-in-up">
	
	<!-- Subtle tech bg pattern -->
	<div class="absolute inset-0 opacity-[0.1] pointer-events-none" 
		style="background-size: 40px 40px; background-image: linear-gradient(to right, rgba(138, 21, 56, 0.08) 1px, transparent 1px), linear-gradient(to bottom, rgba(138, 21, 56, 0.08) 1px, transparent 1px);">
	</div>

	<!-- Main Box Container -->
	<div class="relative z-10 w-full max-w-lg">
		
		<!-- Announcement Ticker (Sleek Clean style) -->
		<div class="mb-6 w-full overflow-hidden rounded-xl border border-slate-200 bg-white text-slate-700 px-4 py-2.5 shadow-sm font-semibold text-xs flex items-center gap-3">
			<div class="bg-primary/10 text-primary border border-primary/20 rounded-md px-2 py-0.5 font-bold uppercase tracking-wider">PENGUMUMAN</div>
			<div class="flex-1 whitespace-nowrap overflow-hidden relative h-4">
				<div class="absolute animate-[marquee_20s_linear_infinite] flex gap-8">
					{#each announcements as a}
						<span>{a.text}</span>
					{/each}
				</div>
			</div>
		</div>

		<!-- Login Card (Formal style) -->
		<div class="card bg-white p-8 sm:p-10 border border-slate-200 shadow-md">
			
			<!-- Logo Header -->
			<div class="mb-8 text-center group">
				<a href="/info" class="inline-block">
					<div class="mx-auto mb-4 w-14 h-14 bg-white border border-slate-200 rounded-xl flex items-center justify-center shadow-sm">
						<img src="/logo.png" alt="Logo Lab AP" class="h-9 w-9 object-contain" />
					</div>
				</a>
				<h1 class="text-2xl font-bold tracking-tight text-slate-900 leading-tight">Portal Layanan Lab AP</h1>
				<p class="mt-1 text-sm font-medium text-slate-400">Peminjaman Alat, Absensi, & Administrasi</p>
			</div>

			<!-- Tab buttons -->
			<div class="grid grid-cols-2 gap-2 mb-6 p-1 border border-slate-200 rounded-xl bg-slate-50">
				<button 
					type="button" 
					class="py-2 rounded-lg text-sm font-semibold transition-all duration-150 {activeTab === 'masuk' ? 'bg-white text-slate-900 border border-slate-200 shadow-sm' : 'text-slate-500 hover:text-slate-800'}"
					onclick={() => { activeTab = 'masuk'; errorMsg = ''; successMsg = ''; }}
				>
					Masuk
				</button>
				<button 
					type="button" 
					class="py-2 rounded-lg text-sm font-semibold transition-all duration-150 {activeTab === 'daftar' ? 'bg-white text-slate-900 border border-slate-200 shadow-sm' : 'text-slate-500 hover:text-slate-800'}"
					onclick={() => { activeTab = 'daftar'; errorMsg = ''; successMsg = ''; }}
				>
					Daftar Akun
				</button>
			</div>

			<!-- Alerts -->
			{#if errorMsg}
				<div class="mb-5 rounded-xl border border-red-200 bg-red-50 p-4 text-sm font-semibold text-red-600 shadow-sm flex items-center gap-3">
					<AlertCircle class="shrink-0 text-red-500" />
					<span>{errorMsg}</span>
				</div>
			{/if}

			{#if successMsg}
				<div class="mb-5 rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-sm font-semibold text-emerald-600 shadow-sm flex items-center gap-3">
					<CheckCircle2 class="shrink-0 text-emerald-500" />
					<span>{successMsg}</span>
				</div>
			{/if}

			<!-- Form Section -->
			{#if activeTab === 'masuk'}
				<form onsubmit={handleLogin} class="space-y-5">
					<div class="text-left">
						<label class="label" for="username">Username / NIM</label>
						<div class="relative">
							<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-slate-400">
								<User size={16} />
							</div>
							<input id="username" class="input pl-10" bind:value={loginUsername} placeholder="Masukkan ID / NIM Anda" required />
						</div>
					</div>

					<div class="text-left">
						<div class="flex items-center justify-between mb-1.5">
							<label class="label m-0" for="pw-login">Password</label>
						</div>
						<div class="relative">
							<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-slate-400">
								<Lock size={16} />
							</div>
							<input id="pw-login" type={showLoginPassword ? "text" : "password"} class="input pl-10" bind:value={loginPassword} placeholder="••••••••" required />
						</div>
					</div>

					<button class="btn-primary w-full py-3 text-sm mt-2 font-semibold" disabled={isLoading}>
						{isLoading ? 'Menghubungkan…' : 'Masuk Layanan'} 
						{#if !isLoading}<ArrowRight size={16} />{/if}
					</button>
				</form>
			{:else}
				<form onsubmit={handleSignup} class="space-y-4">
					<div class="text-left">
						<label class="label" for="reg-name">Nama Lengkap</label>
						<input id="reg-name" class="input" bind:value={signupName} placeholder="Masukkan Nama Lengkap Anda" required />
					</div>

					<div class="text-left">
						<label class="label" for="reg-username">ID / NIM (Untuk Username)</label>
						<input id="reg-username" class="input" bind:value={signupUsername} placeholder="Contoh: 202411001" required />
					</div>

					<div class="text-left">
						<label class="label" for="reg-pw">Password Baru</label>
						<input id="reg-pw" type="password" class="input" bind:value={signupPassword} placeholder="Minimal 6 karakter" required minlength="6" />
					</div>

					<div class="text-left">
						<label class="label" for="reg-role">Pilih Peran</label>
						<select id="reg-role" class="input py-2 bg-white" bind:value={signupRole}>
							<option value="praktikan">Praktikan / Mahasiswa Jaga</option>
							<option value="penyewa">Penyewa Alat (Umum)</option>
						</select>
					</div>

					{#if signupRole === 'penyewa'}
						<div transition:slide={{ duration: 200 }} class="text-left">
							<label class="label" for="reg-phone">No. WhatsApp (WA wajib aktif)</label>
							<div class="relative">
								<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-slate-400">
									<Phone size={16} />
								</div>
								<input id="reg-phone" class="input pl-10" bind:value={signupPhone} placeholder="Contoh: 08123456789" required />
							</div>
						</div>
					{/if}

					<button class="btn-primary w-full py-3 text-sm mt-2 font-semibold" disabled={isLoading}>
						{isLoading ? 'Mendaftarkan…' : 'Daftar & Buat Akun'}
					</button>
				</form>
			{/if}

			<!-- Help & Procedures Links Grid -->
			<div class="mt-8 pt-6 border-t border-slate-100 grid grid-cols-3 gap-3 text-center">
				<button type="button" class="flex flex-col items-center gap-1.5 hover:text-primary transition-colors focus:outline-none" onclick={() => openPanduan = true}>
					<span class="w-10 h-10 rounded-xl bg-slate-50 text-slate-600 border border-slate-200 flex items-center justify-center shadow-sm hover:bg-slate-100 transition-colors"><HelpCircle size={18} /></span>
					<span class="text-[11px] font-semibold text-slate-500 leading-tight">Panduan</span>
				</button>
				
				<button type="button" class="flex flex-col items-center gap-1.5 hover:text-primary transition-colors focus:outline-none" onclick={() => openProsedur = true}>
					<span class="w-10 h-10 rounded-xl bg-slate-50 text-slate-600 border border-slate-200 flex items-center justify-center shadow-sm hover:bg-slate-100 transition-colors"><Info size={18} /></span>
					<span class="text-[11px] font-semibold text-slate-500 leading-tight">Sewa Alat</span>
				</button>

				<button type="button" class="flex flex-col items-center gap-1.5 hover:text-primary transition-colors focus:outline-none" onclick={() => openReschedule = true}>
					<span class="w-10 h-10 rounded-xl bg-slate-50 text-slate-600 border border-slate-200 flex items-center justify-center shadow-sm hover:bg-slate-100 transition-colors"><Calendar size={18} /></span>
					<span class="text-[11px] font-semibold text-slate-500 leading-tight">Jadwal Jaga</span>
				</button>
			</div>

		</div>

		<!-- Back link -->
		<div class="mt-8 text-center">
			<a href="/info" class="inline-flex items-center gap-2 text-xs font-semibold text-slate-400 hover:text-primary transition-colors">
				<ArrowLeft size={12} /> Kembali ke portal utama
			</a>
		</div>
	</div>
</div>

<!-- Modal Dialog 1: Panduan -->
{#if openPanduan}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm" transition:fade={{ duration: 150 }}>
		<div class="w-full max-w-md bg-white border border-slate-200 rounded-2xl p-6 shadow-xl relative animate-fade-in" role="dialog">
			<button class="absolute top-4 right-4 w-7 h-7 rounded-full border border-slate-200 bg-white flex items-center justify-center text-slate-400 hover:text-slate-700 hover:bg-slate-50 shadow-sm transition-all" onclick={() => openPanduan = false} aria-label="Tutup">
				<X size={14} />
			</button>
			<h2 class="text-xl font-bold text-slate-900 mb-4 flex items-center gap-2">
				<HelpCircle class="text-primary" size={20} /> Panduan Layanan
			</h2>
			<div class="space-y-3 text-sm text-slate-500 leading-relaxed overflow-y-auto max-h-[60vh] pr-2 text-left">
				<p>Selamat datang di sistem manajemen terpadu Lab AP. Di sini Anda dapat:</p>
				<ul class="list-disc pl-5 space-y-1">
					<li>Mengakses jadwal jaga asisten dan kehadiran.</li>
					<li>Melakukan peminjaman alat laboratorium (mahasiswa & umum).</li>
					<li>Melihat statistik dan laporan keuangan lab (asisten/admin).</li>
				</ul>
				<p class="font-bold text-slate-700 mt-4">1. Bagaimana cara meminjam alat?</p>
				<p>Daftar sebagai <strong>Penyewa</strong>, lalu login untuk memilih barang, mengisi tanggal peminjaman, dan mengajukan persetujuan.</p>
				<p class="font-bold text-slate-700 mt-2">2. Berapa lama persetujuan sewa?</p>
				<p>Asisten akan memverifikasi dalam 1x24 jam dan menghubungi Anda lewat nomor WhatsApp yang didaftarkan.</p>
			</div>
		</div>
	</div>
{/if}

<!-- Modal Dialog 2: Prosedur Sewa -->
{#if openProsedur}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm" transition:fade={{ duration: 150 }}>
		<div class="w-full max-w-md bg-white border border-slate-200 rounded-2xl p-6 shadow-xl relative animate-fade-in" role="dialog">
			<button class="absolute top-4 right-4 w-7 h-7 rounded-full border border-slate-200 bg-white flex items-center justify-center text-slate-400 hover:text-slate-700 hover:bg-slate-50 shadow-sm transition-all" onclick={() => openProsedur = false} aria-label="Tutup">
				<X size={14} />
			</button>
			<h2 class="text-xl font-bold text-slate-900 mb-4 flex items-center gap-2">
				<Info class="text-primary" size={20} /> Prosedur Sewa Alat
			</h2>
			<div class="space-y-4 text-sm text-slate-500 leading-relaxed overflow-y-auto max-h-[60vh] pr-2 text-left">
				<div class="flex gap-3">
					<span class="w-5 h-5 rounded-full bg-slate-100 text-slate-700 flex items-center justify-center font-bold text-[10px] shrink-0 border border-slate-200">1</span>
					<p><strong>Pendaftaran</strong>: Buat akun di tab "Daftar" dan pilih peran "Penyewa Alat (Umum)". Lengkapi nomor WhatsApp aktif.</p>
				</div>
				<div class="flex gap-3">
					<span class="w-5 h-5 rounded-full bg-slate-100 text-slate-700 flex items-center justify-center font-bold text-[10px] shrink-0 border border-slate-200">2</span>
					<p><strong>Pilih Alat</strong>: Login ke dasbor, temukan barang yang dicari di menu "Sewa Barang", masukkan kuantitas dan durasi pinjam.</p>
				</div>
				<div class="flex gap-3">
					<span class="w-5 h-5 rounded-full bg-slate-100 text-slate-700 flex items-center justify-center font-bold text-[10px] shrink-0 border border-slate-200">3</span>
					<p><strong>Verifikasi</strong>: Asisten laboratorium akan melakukan approval dan memvalidasi kelayakan penyewaan barang.</p>
				</div>
				<div class="flex gap-3">
					<span class="w-5 h-5 rounded-full bg-slate-100 text-slate-700 flex items-center justify-center font-bold text-[10px] shrink-0 border border-slate-200">4</span>
					<p><strong>Pengambilan</strong>: Datang ke lab AP pada waktu yang ditentukan untuk serah terima barang secara fisik dengan membawa KTP/KTM.</p>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Modal Dialog 3: Reschedule/Jadwal Jaga -->
{#if openReschedule}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm" transition:fade={{ duration: 150 }}>
		<div class="w-full max-w-md bg-white border border-slate-200 rounded-2xl p-6 shadow-xl relative animate-fade-in" role="dialog">
			<button class="absolute top-4 right-4 w-7 h-7 rounded-full border border-slate-200 bg-white flex items-center justify-center text-slate-400 hover:text-slate-700 hover:bg-slate-50 shadow-sm transition-all" onclick={() => openReschedule = false} aria-label="Tutup">
				<X size={14} />
			</button>
			<h2 class="text-xl font-bold text-slate-900 mb-4 flex items-center gap-2">
				<Calendar class="text-primary" size={20} /> Jadwal Jaga Asisten
			</h2>
			<div class="space-y-3 text-sm text-slate-500 leading-relaxed overflow-y-auto max-h-[60vh] pr-2 text-left">
				<p>Halaman reschedule ini khusus untuk **asisten yang ingin bertukar shift** jaga di Lab AP:</p>
				<ul class="list-disc pl-5 space-y-1">
					<li>Pertukaran shift maksimal diajukan <strong>H-1</strong> sebelum tugas jaga dimulai.</li>
					<li>Pastikan sudah ada asisten pengganti yang menyetujui pertukaran.</li>
					<li>Status reschedule dapat dipantau di menu "Manajemen Jadwal".</li>
				</ul>
				<p class="mt-4 p-3 bg-slate-50 border border-slate-200 text-slate-600 rounded-xl text-xs">
					💡 *Harap login terlebih dahulu sebagai asisten laboratorium untuk mengajukan pertukaran jadwal jaga resmi.*
				</p>
			</div>
		</div>
	</div>
{/if}

<style>
	@keyframes marquee {
		0% { transform: translateX(100%); }
		100% { transform: translateX(-100%); }
	}
</style>
