<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Play, Square } from 'lucide-svelte';

	type Status = 'init' | 'loading' | 'ready' | 'running' | 'unsupported';

	let status = $state<Status>('init');
	let output = $state('');
	let needInput = $state(false);
	let inputValue = $state('');
	let inputEl: HTMLInputElement | undefined = $state();

	let worker: Worker | null = null;
	let meta: Int32Array | null = null;
	let dataArr: Uint8Array | null = null;

	const DATA_SIZE = 4096;

	function append(text: string) {
		output += text;
	}

	function setupWorker() {
		// SharedArrayBuffer butuh cross-origin isolation (header COOP/COEP).
		if (typeof SharedArrayBuffer === 'undefined' || !self.crossOriginIsolated) {
			status = 'unsupported';
			return;
		}
		const sab = new SharedArrayBuffer(8 + DATA_SIZE);
		meta = new Int32Array(sab, 0, 2);
		dataArr = new Uint8Array(sab, 8);

		status = 'loading';
		worker = new Worker('/c-worker.js', { type: 'module' });
		worker.onmessage = (e: MessageEvent) => {
			const m = e.data;
			if (m.type === 'ready') status = 'ready';
			else if (m.type === 'out' || m.type === 'err') append(m.data);
			else if (m.type === 'need-input') {
				needInput = true;
				queueMicrotask(() => inputEl?.focus());
			} else if (m.type === 'done') {
				status = 'ready';
				needInput = false;
			} else if (m.type === 'fatal') {
				append('\n[Gagal memuat runtime C: ' + m.data + ']\n');
				status = 'unsupported';
			}
		};
		worker.postMessage({ type: 'init', sab });
	}

	function run() {
		if (status !== 'ready' || !worker) return;
		output = '';
		needInput = false;
		status = 'running';
		worker.postMessage({ type: 'run' }); // POC: worker memuat /scanf.wasm
	}

	function submitInput() {
		if (!needInput || !meta || !dataArr) return;
		const line = inputValue;
		append(line + '\n'); // echo ke terminal
		const enc = new TextEncoder().encode(line + '\n');
		const n = Math.min(enc.length, DATA_SIZE);
		dataArr.set(enc.subarray(0, n));
		Atomics.store(meta, 1, n);
		Atomics.store(meta, 0, 1);
		Atomics.notify(meta, 0);
		inputValue = '';
		needInput = false;
	}

	function stop() {
		// Hentikan run yang mungkin sedang Atomics.wait di worker.
		worker?.terminate();
		worker = null;
		needInput = false;
		status = 'init';
		setupWorker();
	}

	onMount(setupWorker);
	onDestroy(() => worker?.terminate());
</script>

<div class="mx-auto max-w-4xl">
	<h1 class="mb-1 text-2xl font-bold text-ink-heading">Playground C (Live) — POC</h1>
	<p class="mb-5 text-sm text-ink-caption">
		Menjalankan program C (WebAssembly/WASI) langsung di browser dengan <code>scanf</code> interaktif.
		Program contoh masih tetap (<code>Nama? </code> lalu echo); editor + compiler menyusul setelah
		endpoint compile siap. (Eksperimental / Fase 1)
	</p>

	{#if status === 'unsupported'}
		<div class="rounded-lg bg-state-warning-bg p-4 text-sm text-state-warning">
			Browser ini belum mendukung mode interaktif (butuh <em>cross-origin isolation</em> /
			SharedArrayBuffer). Coba Chrome/Edge/Firefox terbaru via HTTPS dengan COOP/COEP aktif.
		</div>
	{:else}
		<div class="grid gap-4">
			<div class="flex items-center gap-2">
				<button class="btn-primary" onclick={run} disabled={status !== 'ready'}>
					<Play size={16} />
					{status === 'loading' ? 'Memuat runtime…' : status === 'running' ? 'Berjalan…' : 'Run contoh'}
				</button>
				<button class="btn-outline" onclick={stop} disabled={status === 'loading' || status === 'init'}>
					<Square size={14} /> Stop / Reset
				</button>
			</div>

			<div>
				<p class="mb-1 text-xs font-medium text-ink-caption">Output / Terminal:</p>
				<pre class="min-h-40 overflow-auto whitespace-pre-wrap rounded-lg bg-gray-900 p-3 text-sm text-gray-100">{output || '(klik Run contoh)'}</pre>
				{#if needInput}
					<form class="mt-2 flex items-center gap-2" onsubmit={(e) => { e.preventDefault(); submitInput(); }}>
						<span class="font-mono text-sm text-ink-caption">&gt;</span>
						<input
							bind:this={inputEl}
							bind:value={inputValue}
							class="input flex-1 font-mono text-sm"
							placeholder="ketik input lalu Enter…"
							autocomplete="off"
						/>
						<button class="btn-primary py-1.5" type="submit">Kirim</button>
					</form>
				{/if}
			</div>
		</div>
	{/if}
</div>
