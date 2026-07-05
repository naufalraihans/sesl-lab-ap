<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { setAuth } from '$lib/stores/auth';
	import type { AuthResponse, CekNIMResponse } from '$lib/types';
	import { ArrowLeft, User, Lock, ArrowRight, AlertCircle } from 'lucide-svelte';

	let step = $state<'nim' | 'login' | 'register' | 'blocked'>('nim');
	let nim = $state('');
	let password = $state('');
	let passwordConfirm = $state('');
	let nama = $state('');
	let pesan = $state('');
	let err = $state('');
	let loading = $state(false);

	async function cekNim() {
		err = '';
		loading = true;
		try {
			const res = await api.post<CekNIMResponse>('/api/auth/cek-nim', { nim });
			nama = res.nama ?? '';
			pesan = res.pesan;
			if (!res.ditemukan) {
				step = 'blocked';
			} else if (res.is_registered) {
				step = 'login';
			} else if (res.is_register_open) {
				step = 'register';
			} else {
				step = 'blocked';
			}
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	function redirectByRole(role: string) {
		goto(role === 'admin' ? '/praktikum/admin' : '/praktikum/dashboard');
	}

	async function doLogin() {
		err = '';
		loading = true;
		try {
			const res = await api.post<AuthResponse>('/api/auth/login', { nim, password });
			setAuth(res.token, res.user);
			redirectByRole(res.user.role);
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	async function doRegister() {
		err = '';
		if (password !== passwordConfirm) {
			err = 'Konfirmasi password tidak cocok.';
			return;
		}
		loading = true;
		try {
			const res = await api.post<AuthResponse>('/api/auth/register', { nim, password });
			setAuth(res.token, res.user);
			redirectByRole(res.user.role);
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	function reset() {
		step = 'nim';
		password = '';
		passwordConfirm = '';
		err = '';
	}
</script>

<div class="relative flex min-h-screen items-center justify-center p-4 overflow-hidden"
	style="background: linear-gradient(135deg, #fff5f5 0%, #fffbfb 50%, #fff5f5 100%);">
	
	<!-- Tech background Grid Pattern (very soft rose using logo RGB 138, 21, 56) -->
	<div class="absolute inset-0 opacity-[0.02] pointer-events-none z-0" 
		style="background-size: 30px 30px; background-image: linear-gradient(to right, rgba(138, 21, 56, 0.3) 1px, transparent 1px), linear-gradient(to bottom, rgba(138, 21, 56, 0.3) 1px, transparent 1px);">
	</div>

	<!-- Glowing Mesh Gradients (soft rose/maroon) -->
	<div class="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[400px] h-[400px] rounded-full bg-[#8A1538]/5 blur-[80px] pointer-events-none z-0"></div>
	<div class="absolute bottom-10 right-10 w-72 h-72 rounded-full bg-rose-200/20 blur-[80px] pointer-events-none z-0"></div>

	<!-- Login Card Container (Centered Light Mode Premium Card) -->
	<div class="relative z-10 w-full max-w-md bg-white border border-rose-100/80 rounded-[2.5rem] p-8 sm:p-10 shadow-[0_20px_50px_rgba(138,21,56,0.06)] text-center">
		
		<!-- Top border gradient highlight matching the logo maroon -->
		<div class="absolute top-0 left-12 right-12 h-1 bg-gradient-to-r from-transparent via-[#8A1538] to-transparent rounded-t-[2.5rem]"></div>

		<!-- Logo & Brand -->
		<div class="mb-8">
			<div class="mx-auto mb-4 w-16 h-16 bg-gradient-to-br from-rose-50 to-rose-100 border border-rose-200/40 rounded-2xl flex items-center justify-center shadow-sm group">
				<img src="/logo.png" alt="Logo Lab AP" class="h-10 w-10 object-contain transition-transform duration-500 group-hover:scale-110" />
			</div>
			<h1 class="text-2xl font-black tracking-tight text-slate-900">Login Portal</h1>
			<p class="mt-1 text-xs font-bold text-[#8A1538] uppercase tracking-wider">Lab Algoritma &amp; Pemrograman</p>
		</div>

		{#if err}
			<div class="mb-5 rounded-2xl border border-red-100 bg-red-50 p-4 text-xs font-semibold text-red-700 shadow-sm flex items-start gap-2.5 text-left animate-shake">
				<AlertCircle class="shrink-0 text-red-500 mt-0.5" size={15} />
				<span>{err}</span>
			</div>
		{/if}

		{#if step === 'nim'}
			<form onsubmit={(e) => { e.preventDefault(); cekNim(); }} class="space-y-5 text-left">
				<div>
					<label class="block text-xs font-bold text-slate-600 uppercase tracking-wider mb-2" for="nim">Nomor Induk Mahasiswa (NIM)</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-slate-400">
							<User size={15} />
						</div>
						<input id="nim" class="w-full h-12 pl-10 pr-4 bg-slate-50/80 border border-slate-200 rounded-2xl text-sm font-semibold text-slate-800 placeholder-slate-400/50 focus:outline-none focus:border-[#8A1538]/50 focus:bg-white transition-all" bind:value={nim} placeholder="Contoh: 202431001" required />
					</div>
				</div>
				<button class="w-full h-12 bg-[#8A1538] hover:bg-[#730d2d] text-white rounded-2xl text-sm font-black flex items-center justify-center gap-2 shadow-lg shadow-[#8A1538]/10 transition-all active:scale-[0.99]" disabled={loading}>
					{loading ? 'Memeriksa…' : 'Lanjutkan'} 
					{#if !loading}<ArrowRight size={14} />{/if}
				</button>
			</form>
		{:else if step === 'login'}
			<form onsubmit={(e) => { e.preventDefault(); doLogin(); }} class="space-y-5 animate-fade-in text-left">
				<div class="rounded-2xl bg-rose-50/40 p-4 border border-rose-100/50">
					<p class="text-[9px] font-black text-[#8A1538] uppercase tracking-wider">Praktikan Terdaftar</p>
					<p class="text-base font-extrabold text-slate-900 mt-1">{nama}</p>
					<p class="text-[10px] font-bold text-slate-500 mt-0.5">NIM · {nim}</p>
				</div>
				<div>
					<label class="block text-xs font-bold text-slate-600 uppercase tracking-wider mb-2" for="pw">Password</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-slate-400">
							<Lock size={15} />
						</div>
						<input id="pw" type="password" class="w-full h-12 pl-10 pr-4 bg-slate-50/80 border border-slate-200 rounded-2xl text-sm font-semibold text-slate-800 placeholder-slate-400/50 focus:outline-none focus:border-[#8A1538]/50 focus:bg-white transition-all" bind:value={password} placeholder="Masukkan password Anda" required />
					</div>
				</div>
				<button class="w-full h-12 bg-[#8A1538] hover:bg-[#730d2d] text-white rounded-2xl text-sm font-black flex items-center justify-center shadow-lg shadow-[#8A1538]/10 transition-all active:scale-[0.99]" disabled={loading}>
					{loading ? 'Autentikasi…' : 'Masuk ke Portal'}
				</button>
				<button type="button" class="w-full py-2.5 text-xs font-extrabold text-[#8A1538] hover:text-[#610a24] transition-colors flex items-center justify-center gap-1" onclick={reset}>
					<ArrowLeft size={13} /> Ganti Akun NIM
				</button>
			</form>
		{:else if step === 'register'}
			<form onsubmit={(e) => { e.preventDefault(); doRegister(); }} class="space-y-4 animate-fade-in text-left">
				<div class="rounded-2xl border border-blue-100 bg-blue-50/40 p-4 text-xs font-medium text-blue-800 leading-relaxed">
					{pesan}
				</div>
				<div>
					<label class="block text-xs font-bold text-slate-600 uppercase tracking-wider mb-2" for="pw1">Buat Password Baru</label>
					<input id="pw1" type="password" class="w-full h-12 px-4 bg-slate-50/80 border border-slate-200 rounded-2xl text-sm font-semibold text-slate-800 focus:outline-none focus:border-[#8A1538]/50 focus:bg-white" bind:value={password} placeholder="Minimal 6 karakter" required minlength="6" />
				</div>
				<div>
					<label class="block text-xs font-bold text-slate-600 uppercase tracking-wider mb-2" for="pw2">Konfirmasi Password</label>
					<input id="pw2" type="password" class="w-full h-12 px-4 bg-slate-50/80 border border-slate-200 rounded-2xl text-sm font-semibold text-slate-800 focus:outline-none focus:border-[#8A1538]/50 focus:bg-white" bind:value={passwordConfirm} placeholder="Ketik ulang password baru Anda" required />
				</div>
				<button class="w-full h-12 bg-[#8A1538] hover:bg-[#730d2d] text-white rounded-2xl text-sm font-black flex items-center justify-center shadow-lg shadow-[#8A1538]/10 transition-all active:scale-[0.99]" disabled={loading}>
					{loading ? 'Mendaftarkan Akun…' : 'Daftar & Masuk'}
				</button>
				<button type="button" class="w-full py-2.5 text-xs font-extrabold text-[#8A1538] hover:text-[#610a24] transition-colors flex items-center justify-center gap-1" onclick={reset}>
					<ArrowLeft size={13} /> Ganti Akun NIM
				</button>
			</form>
		{:else if step === 'blocked'}
			<div class="space-y-5 animate-fade-in text-center">
				<div class="mx-auto w-12 h-12 rounded-full bg-red-50 text-red-600 flex items-center justify-center border border-red-100">
					<AlertCircle size={20} />
				</div>
				<div>
					<h3 class="text-lg font-black text-slate-900">Akses Masuk Ditutup</h3>
					<p class="text-xs text-slate-500 font-semibold leading-relaxed mt-2">
						NIM Anda tidak terdaftar atau registrasi kelas Anda belum dibuka oleh asisten. Silakan hubungi asisten penanggung jawab untuk info lebih lanjut.
					</p>
				</div>
				<button type="button" class="w-full h-12 bg-slate-50 hover:bg-slate-100 text-slate-700 border border-slate-200 rounded-2xl text-sm font-bold flex items-center justify-center gap-1.5 transition-all" onclick={reset}>
					<ArrowLeft size={14} /> Ganti NIM Lain
				</button>
			</div>
		{/if}

		<!-- Bottom Back Link -->
		<div class="pt-6 border-t border-slate-100 mt-6 text-center">
			<a href="/info" class="inline-flex items-center gap-1.5 text-xs font-bold text-slate-400 hover:text-primary transition-colors">
				<ArrowLeft size={12} /> Kembali ke portal utama
			</a>
		</div>

	</div>
</div>
