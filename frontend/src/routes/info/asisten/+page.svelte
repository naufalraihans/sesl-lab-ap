<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { User } from '$lib/types';
	import { MessageCircle, Link2 } from 'lucide-svelte';

	let asisten = $state<User[]>([]);
	let loading = $state(true);
	let err = $state('');

	// Colour pairs [from, to] used as inline CSS gradients so Tailwind purge doesn't remove them
	const palette = [
		['#48CAE4', '#06B6D4'], // fun-blue → cyan
		['#9D4EDD', '#4F46E5'], // fun-purple → indigo
		['#06D6A0', '#059669'], // fun-green → emerald
		['#FFC300', '#F97316'], // fun-yellow → orange
		['#EC4899', '#E11D48'], // pink → rose
		['#8A1538', '#B21F47'], // primary → primary-hover
		['#14B8A6', '#0891B2'], // teal → cyan
		['#48CAE4', '#8A1538'], // blue → maroon
	];

	function cardGradient(i: number): string {
		const [a, b] = palette[i % palette.length];
		return `background: linear-gradient(135deg, ${a}, ${b})`;
	}

	function avatarGradient(i: number): string {
		const [a, b] = palette[i % palette.length];
		return `background: linear-gradient(135deg, ${a}, ${b})`;
	}

	onMount(async () => {
		try {
			asisten = (await api.get<User[]>('/api/info/asisten')) ?? [];
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-extrabold text-slate-900">Daftar Asisten Lab</h1>
		<p class="mt-1 text-sm font-bold text-slate-600">Hubungi asisten lab yang bertugas untuk pertanyaan seputar praktikum.</p>
	</div>

	{#if loading}
		<p class="text-slate-700 font-semibold">Memuat…</p>
	{:else if err}
		<p class="rounded-lg bg-state-error-bg p-3 text-state-error font-semibold">{err}</p>
	{:else if asisten.length === 0}
		<p class="text-slate-600 font-semibold">Belum ada data asisten.</p>
	{:else}
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each asisten as a, i}
				<div class="bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-lg hover:-translate-y-1 transition-all duration-300 overflow-hidden text-center group">
					<!-- Gradient header band -->
					<div class="h-24 relative" style={cardGradient(i)}>
						<div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_1px_1px,#fff_1px,transparent_0)] [background-size:16px_16px]"></div>
					</div>

					<div class="px-6 pb-6">
						<!-- Avatar -->
						<div class="-mt-16 flex justify-center mb-3">
							{#if a.foto_url}
								<img
									src={a.foto_url}
									alt={a.nama}
									class="h-28 w-28 rounded-full object-cover ring-4 ring-white shadow-lg group-hover:scale-105 transition-transform"
								/>
							{:else}
								<div
									class="h-28 w-28 rounded-full ring-4 ring-white shadow-lg flex items-center justify-center group-hover:scale-105 transition-transform"
									style={avatarGradient(i)}
								>
									<span class="text-4xl font-black text-white">{a.nama?.charAt(0).toUpperCase()}</span>
								</div>
							{/if}
						</div>

						<h3 class="text-lg font-extrabold text-slate-900 leading-tight">{a.nama}</h3>
						<p class="mt-0.5 text-sm font-bold text-slate-600">{a.nim}</p>
						<span class="inline-block mt-2 px-3 py-1 rounded-full text-xs font-black bg-primary/10 text-primary border border-primary/20 uppercase tracking-wide">
							Asisten Lab
						</span>

						<!-- Contact buttons -->
						<div class="mt-4 flex flex-wrap justify-center gap-2">
							{#if a.nomor_hp}
								<a
									href={`https://wa.me/${a.nomor_hp.replace(/^0/, '62')}`}
									target="_blank" rel="noopener"
									class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-emerald-50 text-emerald-800 border border-emerald-200 hover:bg-emerald-100 transition-colors"
								>
									<MessageCircle size={12}/> WhatsApp
								</a>
							{/if}
							{#if a.medsos_link}
								<a
									href={a.medsos_link}
									target="_blank" rel="noopener"
									class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-blue-50 text-blue-800 border border-blue-200 hover:bg-blue-100 transition-colors"
								>
									<Link2 size={12}/> Media Sosial
								</a>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
