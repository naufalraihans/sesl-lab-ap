<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { api } from '$lib/api';
	import { labelJenis } from '$lib/utils';
	import CodeEditor from './CodeEditor.svelte';
	import Countdown from './Countdown.svelte';
	import type { RuangCourse } from '$lib/types';

	let { aktivasiSesiId, courseId }: { aktivasiSesiId: number; courseId: number } = $props();

	let ruang = $state<RuangCourse | null>(null);
	let answers = $state<Record<number, string>>({});
	let err = $state('');
	let info = $state('');
	let loading = $state(true);
	let submitting = $state(false);
	let dirty = new Set<number>();
	let saveTimer: ReturnType<typeof setInterval>;
	
	let inputToken = $state('');
	let starting = $state(false);

	let locked = $derived(!ruang || !ruang.is_open || ruang.status === 'selesai');

	async function load(start = false) {
		try {
			if (start) {
				ruang = await api.post<RuangCourse>('/api/praktikum/mulai', {
					aktivasi_sesi_id: aktivasiSesiId,
					course_id: courseId,
					token: inputToken ? inputToken : null
				});
			} else {
				ruang = await api.get<RuangCourse>(
					`/api/praktikum/ruang?aktivasi_sesi_id=${aktivasiSesiId}&course_id=${courseId}`
				);
			}
			const map: Record<number, string> = {};
			for (const s of ruang.soal) map[s.soal_terpilih_id] = s.jawaban_teks ?? '';
			answers = map;
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	async function saveOne(soalTerpilihId: number) {
		if (locked) return;
		try {
			await api.post('/api/praktikum/autosave', {
				soal_terpilih_id: soalTerpilihId,
				jawaban_teks: answers[soalTerpilihId] ?? ''
			});
			dirty.delete(soalTerpilihId);
			info = 'Tersimpan otomatis ' + new Date().toLocaleTimeString();
		} catch (e) {
			err = (e as Error).message;
		}
	}

	function markDirty(id: number) {
		dirty.add(id);
	}

	async function flush() {
		for (const id of Array.from(dirty)) await saveOne(id);
	}

	async function submit() {
		if (!confirm('Submit jawaban? Anda tidak dapat mengubah setelah submit.')) return;
		submitting = true;
		try {
			await flush();
			await api.post('/api/praktikum/submit', {
				aktivasi_sesi_id: aktivasiSesiId,
				course_id: courseId
			});
			info = 'Jawaban berhasil di-submit.';
			await load();
		} catch (e) {
			err = (e as Error).message;
		} finally {
			submitting = false;
		}
	}

	async function onExpire() {
		info = 'Waktu habis. Jawaban otomatis ter-submit.';
		await flush().catch(() => {});
		await load();
	}

	onMount(async () => {
		await load();
		// Mulai otomatis jika tidak butuh token, belum mulai & masih terbuka.
		if (ruang && ruang.is_open && ruang.status !== 'selesai' && !ruang.waktu_mulai) {
			if (!ruang.require_token) {
				await load(true);
			}
		}
		// Auto-save berkala tiap 15 detik.
		saveTimer = setInterval(flush, 15000);
	});

	onDestroy(() => clearInterval(saveTimer));
</script>

{#if loading}
	<p class="text-ink-caption">Memuat ruang pengerjaan…</p>
{:else if err && !ruang}
	<p class="rounded-lg bg-state-error-bg p-3 text-state-error">{err}</p>
{:else if ruang}
	{#if !ruang.waktu_mulai && ruang.is_open && ruang.status !== 'selesai' && ruang.require_token}
		<!-- LOCK SCREEN -->
		<div class="card max-w-md mx-auto mt-12 text-center p-8 border-t-4 border-t-primary">
			<h2 class="text-2xl font-bold mb-2">Sesi Ujian Terkunci</h2>
			<p class="text-ink-caption mb-6">Masukkan PIN Ujian yang diberikan oleh Asisten untuk memulai pengerjaan.</p>
			
			<div class="mb-4">
				<input type="text" class="input text-center text-2xl font-mono tracking-widest uppercase" placeholder="PIN 6 DIGIT" maxlength="6" bind:value={inputToken} />
			</div>
			
			{#if err}<p class="mb-4 text-sm text-state-error">{err}</p>{/if}
			
			<button class="btn-primary w-full" disabled={!inputToken || starting} onclick={async () => {
				starting = true;
				err = '';
				await load(true);
				starting = false;
			}}>
				{starting ? 'Memverifikasi...' : 'Mulai Ujian'}
			</button>
		</div>
	{:else}
		<div class="mb-4 flex flex-wrap items-center justify-between gap-3">
		<div>
			<h1 class="text-2xl">{labelJenis(ruang.jenis)}</h1>
			<p class="text-sm text-ink-caption">Durasi {ruang.durasi_menit} menit · {ruang.soal.length} soal</p>
		</div>
		<div class="flex items-center gap-3">
			{#if ruang.deadline && !locked}
				<Countdown deadline={ruang.deadline} {onExpire} />
			{/if}
			{#if ruang.status === 'selesai'}
				<span class="badge bg-state-success-bg text-state-success">Selesai</span>
			{:else if !ruang.is_open}
				<span class="badge bg-gray-100 text-ink-caption">Ditutup admin</span>
			{/if}
		</div>
	</div>

	{#if err}<p class="mb-3 rounded-lg bg-state-error-bg p-3 text-sm text-state-error">{err}</p>{/if}
	{#if info}<p class="mb-3 rounded-lg bg-state-info-bg p-3 text-sm text-state-info">{info}</p>{/if}

	<div class="space-y-5">
		{#each ruang.soal as s, i}
			<div class="card">
				<div class="mb-2 flex items-center justify-between">
					<h3 class="text-base">Soal {i + 1}
						{#if s.kategori_ujian}<span class="badge ml-2 bg-surface-soft text-primary">{s.kategori_ujian}</span>{/if}
					</h3>
					<span class="text-sm text-ink-caption">{s.poin} poin · {s.jenis_soal}</span>
				</div>
				<div class="prose prose-sm max-w-none text-ink-body">
					{@html s.teks_soal}
				</div>
				{#if s.gambar_url}
					<img src={s.gambar_url} alt="Flowchart" class="mt-3 max-w-full rounded-lg border border-gray-200" />
				{/if}

				<div class="mt-3">
					{#if s.jenis_soal === 'coding'}
						<CodeEditor
							bind:value={answers[s.soal_terpilih_id]}
							readonly={locked}
							language="c"
						/>
					{:else}
						<textarea
							class="input min-h-32"
							bind:value={answers[s.soal_terpilih_id]}
							readonly={locked}
							placeholder="Tulis jawaban Anda…"
							oninput={() => markDirty(s.soal_terpilih_id)}
							onblur={() => saveOne(s.soal_terpilih_id)}
						></textarea>
					{/if}
					{#if s.jenis_soal === 'coding' && !locked}
						<button class="btn-outline mt-2 py-1.5 text-xs" onclick={() => saveOne(s.soal_terpilih_id)}>💾 Simpan soal ini</button>
					{/if}
				</div>
			</div>
		{/each}
	</div>

	{#if !locked}
		<div class="sticky bottom-4 mt-6 flex justify-end">
			<button class="btn-primary px-8" onclick={submit} disabled={submitting}>
				{submitting ? 'Mengirim…' : 'Submit Jawaban'}
			</button>
		</div>
	{/if}
	{/if}
{/if}
