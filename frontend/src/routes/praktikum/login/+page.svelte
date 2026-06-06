<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { setAuth } from '$lib/stores/auth';
	import type { AuthResponse, CekNIMResponse } from '$lib/types';

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

<div class="grid min-h-screen place-items-center bg-surface-soft px-4">
	<div class="w-full max-w-md rounded-lg border border-gray-200 bg-white p-8 shadow-sm">
		<div class="mb-6 text-center">
			<div class="mx-auto mb-3 grid h-12 w-12 place-items-center rounded-lg bg-primary text-xl font-bold text-white">L</div>
			<h1 class="text-xl">Login Praktikum</h1>
			<p class="text-sm text-ink-caption">Lab Algoritma &amp; Pemrograman</p>
		</div>

		{#if err}
			<p class="mb-4 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>
		{/if}

		{#if step === 'nim'}
			<form onsubmit={(e) => { e.preventDefault(); cekNim(); }}>
				<label class="label" for="nim">NIM</label>
				<input id="nim" class="input" bind:value={nim} placeholder="Masukkan NIM" required />
				<button class="btn-primary mt-4 w-full" disabled={loading}>{loading ? 'Memeriksa…' : 'Lanjut'}</button>
			</form>
		{:else if step === 'login'}
			<form onsubmit={(e) => { e.preventDefault(); doLogin(); }}>
				<p class="mb-3 text-sm text-ink-body">Halo <strong>{nama}</strong>, masukkan password Anda.</p>
				<label class="label" for="pw">Password</label>
				<input id="pw" type="password" class="input" bind:value={password} required />
				<button class="btn-primary mt-4 w-full" disabled={loading}>{loading ? 'Masuk…' : 'Login'}</button>
				<button type="button" class="mt-2 w-full text-sm text-ink-caption" onclick={reset}>← Ganti NIM</button>
			</form>
		{:else if step === 'register'}
			<form onsubmit={(e) => { e.preventDefault(); doRegister(); }}>
				<p class="mb-3 rounded-lg bg-state-info-bg p-3 text-sm text-state-info">{pesan}</p>
				<label class="label" for="pw1">Buat Password</label>
				<input id="pw1" type="password" class="input" bind:value={password} required minlength="6" />
				<label class="label mt-3" for="pw2">Konfirmasi Password</label>
				<input id="pw2" type="password" class="input" bind:value={passwordConfirm} required />
				<button class="btn-primary mt-4 w-full" disabled={loading}>{loading ? 'Mendaftar…' : 'Daftar & Masuk'}</button>
				<button type="button" class="mt-2 w-full text-sm text-ink-caption" onclick={reset}>← Ganti NIM</button>
			</form>
		{:else}
			<p class="rounded-lg bg-state-warning-bg p-3 text-sm text-state-warning">{pesan}</p>
			<button class="btn-outline mt-4 w-full" onclick={reset}>← Kembali</button>
		{/if}

		<p class="mt-6 text-center text-sm"><a href="/info">← Kembali ke beranda</a></p>
	</div>
</div>
