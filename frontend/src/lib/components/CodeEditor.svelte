<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { api } from '$lib/api';
	import { Play } from 'lucide-svelte';

	let {
		value = $bindable(''),
		language = 'c',
		readonly = false,
		height = '420px',
		runnable = false,
		oninput
	}: {
		value?: string;
		language?: string;
		readonly?: boolean;
		height?: string;
		runnable?: boolean;
		oninput?: () => void;
	} = $props();

	let el: HTMLDivElement;
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	let editor: any = null;
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	let monacoRef: any = null;

	// Run state
	let runLang = $state('c');
	let stdin = $state('');
	let running = $state(false);
	let runErr = $state('');
	let output = $state<{ stdout: string; stderr: string; error: string } | null>(null);

	onMount(async () => {
		if (language === 'python') runLang = 'python';
		const loader = (await import('@monaco-editor/loader')).default;
		const monaco = await loader.init();
		monacoRef = monaco;
		editor = monaco.editor.create(el, {
			value,
			language: runLang,
			readOnly: readonly,
			automaticLayout: true,
			minimap: { enabled: false },
			fontSize: 14,
			scrollBeyondLastLine: false,
			theme: 'vs'
		});
		editor.onDidChangeModelContent(() => {
			value = editor.getValue();
			oninput?.();
		});
	});

	// Sinkronkan jika value diubah dari luar (mis. load awal).
	$effect(() => {
		if (editor && value !== editor.getValue()) {
			editor.setValue(value ?? '');
		}
	});

	function onLangChange() {
		if (monacoRef && editor) {
			monacoRef.editor.setModelLanguage(editor.getModel(), runLang);
		}
	}

	async function runCode() {
		running = true;
		runErr = '';
		output = null;
		try {
			output = await api.post<{ stdout: string; stderr: string; error: string }>(
				'/api/praktikum/run',
				{ language: runLang, source: value, stdin }
			);
		} catch (e) {
			runErr = (e as Error).message;
		} finally {
			running = false;
		}
	}

	function formatOutput(o: { stdout: string; stderr: string; error: string }): string {
		let s = o.stdout ?? '';
		if (o.stderr) s += (s ? '\n' : '') + o.stderr;
		if (o.error) s += (s ? '\n' : '') + `[${o.error}]`;
		return s.trim() || '(tidak ada output)';
	}

	onDestroy(() => {
		editor?.dispose();
	});
</script>

<div bind:this={el} style="height: {height};" class="overflow-hidden rounded-lg border border-gray-300"></div>

{#if runnable}
	<div class="mt-2 rounded-lg border border-gray-200 bg-surface-soft p-3">
		<div class="flex flex-wrap items-center gap-2">
			<select bind:value={runLang} onchange={onLangChange} class="input h-9 w-auto py-1 text-sm">
				<option value="c">C</option>
				<option value="python">Python</option>
			</select>
			<button class="btn-primary py-1.5 text-sm" onclick={runCode} disabled={running}>
				<Play size={14} /> {running ? 'Menjalankan…' : 'Run'}
			</button>
			<span class="text-xs text-ink-caption">Coba jalankan kode (input di bawah, opsional)</span>
		</div>
		<textarea
			bind:value={stdin}
			class="input mt-2 min-h-12 font-mono text-sm"
			placeholder="stdin (opsional) — input untuk program"
		></textarea>
		{#if runErr}<p class="mt-2 text-sm text-state-error">{runErr}</p>{/if}
		{#if output}
			<div class="mt-2">
				<p class="mb-1 text-xs font-medium text-ink-caption">Output:</p>
				<pre class="max-h-60 overflow-auto whitespace-pre-wrap rounded-lg bg-gray-900 p-3 text-sm text-gray-100">{formatOutput(output)}</pre>
			</div>
		{/if}
	</div>
{/if}
