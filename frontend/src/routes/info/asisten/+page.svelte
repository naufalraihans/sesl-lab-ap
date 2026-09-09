<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import type { User } from '$lib/types';

	let asisten = $state<User[]>([]);
	let loading = $state(true);
	let err = $state('');

	onMount(async () => {
		try {
			asisten = (await api.get<User[]>('/api/info/asisten')) ?? [];
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});

	function getInitials(name: string): string {
		return name.split(' ').map((n) => n[0]).join('').substring(0, 2).toUpperCase();
	}
</script>

<div class="page">
	<div class="page-header">
		<h1 class="page-title">Tim Asisten Laboratorium</h1>
		<p class="page-desc">Kenali para asisten lab dan temukan kontak mereka untuk bantuan praktikum.</p>
	</div>

	{#if loading}
		<div class="loading-state">
			<div class="spinner"></div>
			<p>Memuat data asisten…</p>
		</div>
	{:else if err}
		<div class="error-state">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="48" height="48"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
			<p>{err}</p>
		</div>
	{:else if asisten.length === 0}
		<div class="empty-state">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="64" height="64"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
			<p>Belum ada data asisten</p>
			<span class="empty-hint">Data asisten akan ditampilkan setelah admin menambahkannya.</span>
		</div>
	{:else}
		<div class="asisten-grid">
			{#each asisten as a}
				{@const nm = (a.nama ?? '').trim()}
				{@const cut = nm.lastIndexOf(' ')}
				{@const head = cut > 0 ? nm.slice(0, cut + 1) : ''}
				{@const last = cut > 0 ? nm.slice(cut + 1) : nm}
				<div class="asisten-card">
					{#if a.foto_url}
						<img src={a.foto_url} alt={a.nama} class="photo" loading="lazy" />
					{:else}
						<div class="photo photo-fallback" aria-hidden="true">
							{getInitials(a.nama ?? '?')}
						</div>
					{/if}
					<div class="card-body">
						<h3 class="asisten-name">
							{head}<span class="name-last">{last}<svg class="badge" viewBox="0 0 24 24" fill="currentColor" role="img" aria-label="Asisten terverifikasi"><path d="M12 2l2.4 1.8 3-.3 1 2.8 2.6 1.5-.9 2.9.9 2.9-2.6 1.5-1 2.8-3-.3L12 22l-2.4-1.8-3 .3-1-2.8L3 16.2l.9-2.9L3 10.4l2.6-1.5 1-2.8 3 .3z"/><path d="M10.6 15.2l-2.8-2.8 1.2-1.2 1.6 1.6 4-4 1.2 1.2z" fill="#fff"/></svg></span>
						</h3>
						<span class="asisten-nim">{a.nim}</span>

						<div class="contact-links">
							{#if a.nomor_hp}
								<a href={`https://wa.me/${a.nomor_hp.replace(/^0/, '62')}`} target="_blank" rel="noopener noreferrer" class="contact-btn wa-btn">
									<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.362 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.338 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
									WhatsApp
								</a>
							{/if}
							{#if a.medsos_link}
								<a href={a.medsos_link} target="_blank" rel="noopener noreferrer" class="contact-btn medsos-btn">
									<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>
									Profil
								</a>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	/* Palette Lab-AP dipertahankan: maroon #8A1538 / hover #B21F47 */
	.page { min-height: 60vh; }
	.page-header { padding: 0 0 2rem; }
	.page-title {
		font-size: 1.75rem;
		font-weight: 800;
		color: #0f172a;
		margin-bottom: 0.5rem;
	}
	.page-desc { color: #475569; font-size: 0.95rem; font-weight: 600; }

	.loading-state, .error-state, .empty-state {
		display: flex; flex-direction: column; align-items: center; justify-content: center;
		padding: 4rem 2rem; text-align: center; color: #64748b;
	}
	.spinner {
		width: 2.5rem; height: 2.5rem; border: 3px solid #e2e8f0;
		border-top-color: #8A1538; border-radius: 50%;
		animation: spin 0.8s linear infinite; margin-bottom: 1rem;
	}
	@keyframes spin { to { transform: rotate(360deg); } }
	.error-state { color: #dc2626; }
	.error-state svg { margin-bottom: 1rem; opacity: 0.6; }
	.empty-state svg { margin-bottom: 1rem; color: #cbd5e1; }
	.empty-hint { font-size: 0.85rem; color: #64748b; margin-top: 0.25rem; }

	/* Grid kartu profil — foto fade ke body (adaptasi layout mikon) */
	.asisten-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
		gap: 3rem 2.5rem;
		padding: 0.5rem 0 1rem;
	}
	.asisten-card {
		background: #ffffff;
		border-radius: 1.75rem;
		box-shadow: 0 0 0 3px #fff, 0 20px 40px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.05);
		overflow: hidden;
		transition: transform 0.25s ease, box-shadow 0.25s ease;
	}
	.asisten-card:hover, .asisten-card:focus-within {
		transform: translateY(-6px);
		box-shadow: 0 0 0 3px #fff, 0 30px 60px rgba(0,0,0,0.16), 0 3px 8px rgba(0,0,0,0.06);
	}
	.photo {
		display: block;
		width: 100%;
		aspect-ratio: 1 / 1.05;
		object-fit: cover;
		object-position: 50% 25%;
		background: #C9D2D4;
		-webkit-mask-image: linear-gradient(#000 72%, transparent);
		mask-image: linear-gradient(#000 72%, transparent);
	}
	.photo-fallback {
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 3rem;
		font-weight: 700;
		letter-spacing: 0.02em;
		color: #fff;
		background: linear-gradient(135deg, #8A1538 0%, #5C0E25 100%);
	}
	.card-body {
		padding: 0 1.25rem 1.25rem;
		margin-top: -0.75rem;
		text-align: left;
	}
	.asisten-name {
		font-size: 1.3rem;
		font-weight: 650;
		letter-spacing: -0.01em;
		line-height: 1.25;
		color: #0f172a;
		margin-bottom: 0.25rem;
	}
	.name-last { white-space: nowrap; }
	.badge {
		display: inline;
		width: 1.125rem;
		height: 1.125rem;
		margin-left: 0.375rem;
		vertical-align: -0.15em;
		color: #8A1538;
	}
	.asisten-nim {
		font-size: 0.9rem;
		color: #3D3D3D;
		font-weight: 400;
	}
	.contact-links {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-top: 1.5rem;
	}
	.contact-btn {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.7rem 1.1rem;
		border-radius: 999px;
		font-size: 0.9rem;
		font-weight: 500;
		text-decoration: none;
		background: #fff;
		box-shadow: 0 2px 8px rgba(0,0,0,0.12);
		transition: background 0.2s ease, transform 0.15s ease;
	}
	.contact-btn:active { transform: scale(0.97); }
	.contact-btn:focus-visible { outline: 2px solid #8A1538; outline-offset: 2px; }
	.wa-btn { color: #059669; }
	.wa-btn:hover, .wa-btn:visited { color: #059669; }
	.wa-btn:hover { background: #ecfdf5; }
	.medsos-btn { color: #8A1538; }
	.medsos-btn:hover, .medsos-btn:visited { color: #8A1538; }
	.medsos-btn:hover { background: #FDF2F4; }

	@media (prefers-reduced-motion: reduce) {
		.asisten-card, .contact-btn { transition: none; }
	}
</style>
