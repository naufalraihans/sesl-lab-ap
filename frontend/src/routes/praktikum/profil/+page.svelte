<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { user } from '$lib/stores/auth';

	let isAdmin = $derived($user?.role === 'admin');
	let nama = $state('');
	let nomorHp = $state('');
	let medsos = $state('');
	let fotoUrl = $state('');
	let passwordLama = $state('');
	let passwordBaru = $state('');
	let msg = $state('');
	let err = $state('');

	onMount(async () => {
		try {
			const u = await api.get<any>('/api/profile');
			nama = u.nama ?? '';
			nomorHp = u.nomor_hp ?? '';
			medsos = u.medsos_link ?? '';
			fotoUrl = u.foto_url ?? '';
		} catch (e) {
			err = (e as Error).message;
		}
	});

	async function uploadFoto(ev: Event) {
		const input = ev.target as HTMLInputElement;
		if (!input.files?.[0]) return;
		const fd = new FormData();
		fd.append('file', input.files[0]);
		fd.append('folder', 'asisten');
		try {
			const res = await api.upload<{ url: string }>('/api/admin/upload', fd);
			fotoUrl = res.url;
			msg = 'Foto terunggah.';
		} catch (e) {
			err = (e as Error).message;
		}
	}

	async function save() {
		msg = ''; err = '';
		const body: Record<string, unknown> = { nama };
		if (isAdmin) {
			body.nomor_hp = nomorHp;
			body.medsos_link = medsos;
			body.foto_url = fotoUrl;
		}
		if (passwordBaru) {
			body.password_lama = passwordLama;
			body.password_baru = passwordBaru;
		}
		try {
			await api.put('/api/profile', body);
			msg = 'Profil diperbarui.';
			passwordLama = ''; passwordBaru = '';
		} catch (e) {
			err = (e as Error).message;
		}
	}
</script>

<h1 class="mb-4 text-2xl">Profil Saya</h1>

{#if msg}<p class="mb-3 rounded-lg bg-state-success-bg p-3 text-sm text-state-success">{msg}</p>{/if}
{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}

<div class="card max-w-lg space-y-3">
	<div>
		<label class="label" for="nama">Nama</label>
		<input id="nama" class="input" bind:value={nama} />
	</div>

	{#if isAdmin}
		<div>
			<label class="label" for="hp">Nomor HP/WhatsApp</label>
			<input id="hp" class="input" bind:value={nomorHp} placeholder="08xxx" />
		</div>
		<div>
			<label class="label" for="ms">Link Media Sosial / LinkedIn</label>
			<input id="ms" class="input" bind:value={medsos} />
		</div>
		<div>
			<label class="label" for="foto">Foto Profil</label>
			{#if fotoUrl}<img src={fotoUrl} alt="foto" class="mb-2 h-20 w-20 rounded-full object-cover" />{/if}
			<input id="foto" type="file" accept="image/*" onchange={uploadFoto} />
		</div>
	{/if}

	<hr class="my-2" />
	<p class="text-sm text-ink-caption">Ganti password (opsional)</p>
	<div>
		<label class="label" for="pl">Password Lama</label>
		<input id="pl" type="password" class="input" bind:value={passwordLama} />
	</div>
	<div>
		<label class="label" for="pb">Password Baru</label>
		<input id="pb" type="password" class="input" bind:value={passwordBaru} minlength="6" />
	</div>

	<button class="btn-primary" onclick={save}>Simpan Perubahan</button>
</div>
