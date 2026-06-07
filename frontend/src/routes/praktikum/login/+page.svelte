<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { setAuth } from '$lib/stores/auth';
	import type { AuthResponse, CekNIMResponse } from '$lib/types';
	import { ArrowLeft, User, Lock, ArrowRight } from 'lucide-svelte';

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

<div class="relative flex min-h-screen items-center justify-center overflow-hidden bg-gray-900">
	<!-- Animated Background Elements -->
	<div class="absolute inset-0 z-0 bg-[linear-gradient(to_right,#4f4f4f2e_1px,transparent_1px),linear-gradient(to_bottom,#4f4f4f2e_1px,transparent_1px)] bg-[size:14px_24px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)]"></div>
	<div class="animate-float absolute top-1/4 left-1/4 h-96 w-96 rounded-full bg-primary/20 blur-3xl filter"></div>
	<div class="animate-float absolute bottom-1/4 right-1/4 h-[30rem] w-[30rem] rounded-full bg-blue-500/10 blur-3xl filter" style="animation-delay: -2s;"></div>

	<!-- Login Card -->
	<div class="animate-fade-in-up relative z-10 w-full max-w-md px-4">
		<div class="glass-dark rounded-3xl p-10 text-white shadow-2xl">
			<div class="mb-8 text-center">
				<div class="mx-auto mb-4 grid h-14 w-14 place-items-center rounded-2xl bg-gradient-to-br from-primary to-primary-dark shadow-lg shadow-primary/30">
					<span class="text-2xl font-bold tracking-wider text-white">L</span>
				</div>
				<h1 class="text-2xl font-bold tracking-tight">Login Praktikum</h1>
				<p class="mt-2 text-sm text-gray-400">Lab Algoritma &amp; Pemrograman</p>
			</div>

		{#if err}
			<p class="mb-4 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>
		{/if}

		{#if step === 'nim'}
			<form onsubmit={(e) => { e.preventDefault(); cekNim(); }} class="space-y-5">
				<div>
					<label class="mb-1.5 block text-sm font-medium text-gray-300" for="nim">Nomor Induk Mahasiswa (NIM)</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
							<User size={18} />
						</div>
						<input id="nim" class="w-full rounded-xl border border-gray-700 bg-gray-800/50 py-3 pl-10 pr-4 text-white placeholder-gray-500 outline-none transition-all focus:border-primary focus:bg-gray-800 focus:ring-2 focus:ring-primary/20" bind:value={nim} placeholder="Masukkan NIM Anda" required />
					</div>
				</div>
				<button class="btn-primary w-full py-3 text-base shadow-primary/20 hover:shadow-primary/40" disabled={loading}>
					{loading ? 'Memeriksa…' : 'Lanjutkan'} 
					{#if !loading}<ArrowRight size={18} />{/if}
				</button>
			</form>
		{:else if step === 'login'}
			<form onsubmit={(e) => { e.preventDefault(); doLogin(); }} class="space-y-5 animate-fade-in">
				<div class="rounded-xl bg-gray-800/50 p-4 border border-gray-700">
					<p class="text-sm text-gray-300">Selamat datang kembali,</p>
					<p class="text-lg font-semibold text-white">{nama}</p>
				</div>
				<div>
					<label class="mb-1.5 block text-sm font-medium text-gray-300" for="pw">Password</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
							<Lock size={18} />
						</div>
						<input id="pw" type="password" class="w-full rounded-xl border border-gray-700 bg-gray-800/50 py-3 pl-10 pr-4 text-white placeholder-gray-500 outline-none transition-all focus:border-primary focus:bg-gray-800 focus:ring-2 focus:ring-primary/20" bind:value={password} placeholder="••••••••" required />
					</div>
				</div>
				<button class="btn-primary w-full py-3 text-base" disabled={loading}>{loading ? 'Autentikasi…' : 'Login ke Praktikum'}</button>
				<button type="button" class="mt-4 flex w-full items-center justify-center gap-2 text-sm text-gray-400 transition-colors hover:text-white" onclick={reset}><ArrowLeft size={16} /> Ganti Akun (NIM)</button>
			</form>
		{:else if step === 'register'}
			<form onsubmit={(e) => { e.preventDefault(); doRegister(); }} class="space-y-4 animate-fade-in">
				<div class="rounded-xl border border-blue-500/30 bg-blue-500/10 p-4 text-sm text-blue-200 shadow-inner">
					{pesan}
				</div>
				<div>
					<label class="mb-1.5 block text-sm font-medium text-gray-300" for="pw1">Buat Password</label>
					<input id="pw1" type="password" class="w-full rounded-xl border border-gray-700 bg-gray-800/50 py-3 px-4 text-white placeholder-gray-500 outline-none transition-all focus:border-primary focus:bg-gray-800 focus:ring-2 focus:ring-primary/20" bind:value={password} required minlength="6" />
				</div>
				<div>
					<label class="mb-1.5 block text-sm font-medium text-gray-300" for="pw2">Konfirmasi Password</label>
					<input id="pw2" type="password" class="w-full rounded-xl border border-gray-700 bg-gray-800/50 py-3 px-4 text-white placeholder-gray-500 outline-none transition-all focus:border-primary focus:bg-gray-800 focus:ring-2 focus:ring-primary/20" bind:value={passwordConfirm} required />
				</div>
				<button class="btn-primary w-full py-3 text-base" disabled={loading}>{loading ? 'Mendaftarkan Akun…' : 'Daftar & Masuk'}</button>
				<button type="button" class="mt-4 flex w-full items-center justify-center gap-2 text-sm text-gray-400 transition-colors hover:text-white" onclick={reset}><ArrowLeft size={16} /> Ganti Akun (NIM)</button>
			</form>
		{:else}
			<div class="animate-fade-in text-center space-y-6">
				<div class="rounded-xl border border-yellow-500/30 bg-yellow-500/10 p-5 text-sm text-yellow-200 shadow-inner">
					{pesan}
				</div>
				<button class="btn-outline w-full border-gray-600 text-gray-300 hover:border-gray-400 hover:text-white" onclick={reset}>
					<ArrowLeft size={16} /> Kembali
				</button>
			</div>
		{/if}

		<div class="mt-8 text-center">
			<a href="/info" class="inline-flex items-center gap-2 text-sm text-gray-500 transition-colors hover:text-gray-300">
				<ArrowLeft size={14} /> Kembali ke portal utama
			</a>
		</div>
	</div>
</div>
