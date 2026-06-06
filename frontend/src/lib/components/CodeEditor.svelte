<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	let {
		value = $bindable(''),
		language = 'c',
		readonly = false,
		height = '320px'
	}: { value?: string; language?: string; readonly?: boolean; height?: string } = $props();

	let el: HTMLDivElement;
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	let editor: any = null;
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	let monacoRef: any = null;

	onMount(async () => {
		const loader = (await import('@monaco-editor/loader')).default;
		const monaco = await loader.init();
		monacoRef = monaco;
		editor = monaco.editor.create(el, {
			value,
			language,
			readOnly: readonly,
			automaticLayout: true,
			minimap: { enabled: false },
			fontSize: 14,
			scrollBeyondLastLine: false,
			theme: 'vs'
		});
		editor.onDidChangeModelContent(() => {
			value = editor.getValue();
		});
	});

	// Sinkronkan jika value diubah dari luar (mis. load awal).
	$effect(() => {
		if (editor && value !== editor.getValue()) {
			editor.setValue(value ?? '');
		}
	});

	onDestroy(() => {
		editor?.dispose();
	});
</script>

<div bind:this={el} style="height: {height};" class="overflow-hidden rounded-lg border border-gray-300"></div>
