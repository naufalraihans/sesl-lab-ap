<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { Download, FileText, Eye, EyeOff } from 'lucide-svelte';

	interface Pedoman { id: number; nama_dokumen: string; file_url: string; }
	let items = $state<Pedoman[]>([]);
	let loading = $state(true);
	let err = $state('');

	let selectedItem = $state<Pedoman | null>(null);
	let showPreview = $state(false);

	onMount(async () => {
		try {
			items = (await api.get<Pedoman[]>('/api/info/laporan')) ?? [];
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});

	function togglePreview(item: Pedoman) {
		if (selectedItem?.id === item.id && showPreview) {
			showPreview = false;
		} else {
			selectedItem = item;
			showPreview = true;
		}
	}

	function getFormatBadge(url: string): string {
		const ext = url.split('.').pop()?.toLowerCase() ?? '';
		const map: Record<string, string> = {
			pdf: 'PDF', docx: 'DOCX', doc: 'DOC',
			xlsx: 'XLSX', xls: 'XLS', pptx: 'PPTX', ppt: 'PPT',
			png: 'PNG', jpg: 'JPG', jpeg: 'JPEG', webp: 'WEBP',
		};
		return map[ext] ?? ext.toUpperCase();
	}
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex items-center gap-3 text-left">
		<div class="w-11 h-11 rounded-xl flex items-center justify-center shadow-md" style="background: linear-gradient(135deg, #f97316, #f59e0b);">
			<FileText size={22} color="white" />
		</div>
		<div>
			<h1 class="text-2xl font-extrabold text-slate-900 leading-none">Pedoman Laporan</h1>
			<p class="mt-1 text-xs font-bold text-slate-600">Unduh &amp; pratinjau berkas panduan laporan praktikum.</p>
		</div>
	</div>

	{#if loading}
		<p class="text-slate-700 font-semibold">Memuat…</p>
	{:else if err}
		<p class="rounded-lg bg-red-50 p-3 text-red-700 font-bold border border-red-200">{err}</p>
	{:else if items.length === 0}
		<p class="text-slate-600 font-bold">Belum ada dokumen pedoman.</p>
	{:else}
		<div class="grid gap-5 sm:grid-cols-2">
			{#each items as it}
				{@const badge = getFormatBadge(it.file_url)}
				<div class="group bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-lg hover:-translate-y-1 transition-all duration-300 overflow-hidden flex flex-col text-left">
					<!-- Card body -->
					<div class="flex items-stretch gap-0 flex-1">
						<!-- Gradient icon strip -->
						<div class="flex-shrink-0 w-16 flex items-center justify-center" style="background: linear-gradient(135deg, #7c3aed, #4f46e5); border-radius: 1rem 0 0 0;">
							<FileText size={28} color="white" />
						</div>
						<!-- Content -->
						<div class="flex-1 p-4">
							<div class="flex items-start justify-between gap-2 mb-1">
								<h2 class="font-extrabold text-slate-900 text-sm leading-snug">{it.nama_dokumen}</h2>
								{#if badge}
									<span class="flex-shrink-0 text-[9px] font-black tracking-wider uppercase px-2 py-0.5 rounded-full border" style="background-color:#ede9fe; color:#5b21b6; border-color:#ddd6fe;">{badge}</span>
								{/if}
							</div>
							<span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Berkas Panduan</span>
						</div>
					</div>
					<!-- Action buttons -->
					<div class="flex gap-2 px-4 pb-4 pt-3 border-t border-slate-100">
						<a href={it.file_url} target="_blank" rel="noopener"
							class="btn-primary flex-1 py-2 text-xs font-bold flex items-center justify-center gap-1.5">
							<Download size={13} /> Download
						</a>
						<button onclick={() => togglePreview(it)}
							class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-xl text-xs font-bold border border-slate-200 bg-slate-50 hover:bg-slate-100 text-slate-700 transition-all">
							{#if selectedItem?.id === it.id && showPreview}
								<EyeOff size={13} /> Tutup Preview
							{:else}
								<Eye size={13} /> Preview
							{/if}
						</button>
					</div>
				</div>
			{/each}
		</div>

		{#if showPreview && selectedItem}
			<div class="bg-white rounded-2xl border border-slate-200 shadow-lg text-left p-5 space-y-4">
				<div class="flex items-center justify-between border-b border-slate-100 pb-3">
					<div>
						<h3 class="font-extrabold text-slate-900 text-base">Preview: {selectedItem.nama_dokumen}</h3>
						<p class="text-xs text-slate-600 font-bold mt-0.5">Jika dokumen tidak tampil, silakan download langsung.</p>
					</div>
					<button onclick={() => showPreview = false} class="btn-outline px-4 py-1.5 text-xs font-bold">Tutup Preview</button>
				</div>
				<div class="w-full bg-slate-50 border border-slate-200 rounded-2xl overflow-hidden min-h-[500px] flex items-center justify-center">
					{#if selectedItem.file_url.toLowerCase().endsWith('.pdf')}
						<iframe title="PDF Reader" src={selectedItem.file_url} class="w-full h-[600px] border-0" allowfullscreen={true}></iframe>
					{:else if selectedItem.file_url.toLowerCase().match(/\.(png|jpe?g|webp|gif)$/)}
						<img src={selectedItem.file_url} alt={selectedItem.nama_dokumen} class="max-w-full max-h-[600px] object-contain p-2" />
					{:else}
						<iframe title="Doc Reader" src="https://docs.google.com/gview?url={encodeURIComponent(selectedItem.file_url)}&embedded=true" class="w-full h-[600px] border-0" allowfullscreen={true}></iframe>
					{/if}
				</div>
			</div>
		{/if}
	{/if}
</div>
