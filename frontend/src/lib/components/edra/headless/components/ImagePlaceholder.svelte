<script lang="ts">
	import MediaPlaceHolder from '../../components/MediaPlaceHolder.svelte';
	import type { NodeViewProps } from '@tiptap/core';
	import Image from '@lucide/svelte/icons/image';
	import { api } from '$lib/api';

	const { editor }: NodeViewProps = $props();

	let fileInput = $state<HTMLInputElement>();

	function handleClick() {
		fileInput?.click();
	}

	async function uploadImage(ev: Event) {
		const input = ev.target as HTMLInputElement;
		if (!input.files?.[0]) return;
		
		const fd = new FormData();
		fd.append('file', input.files[0]);
		fd.append('folder', 'soal-images');
		
		try {
			const res = await api.upload<{ url: string }>('/api/admin/upload', fd);
			editor.chain().focus().setImage({ src: res.url }).run();
		} catch (e) {
			alert('Gagal mengunggah gambar: ' + (e as Error).message);
		} finally {
			if (fileInput) fileInput.value = '';
		}
	}
</script>

<input
	type="file"
	accept="image/*"
	bind:this={fileInput}
	onchange={uploadImage}
	style="display: none;"
/>

<MediaPlaceHolder
	class="edra-media-placeholder-wrapper"
	icon={Image}
	title="Upload an Image"
	onClick={handleClick}
/>
