<script lang="ts">
	import { untrack } from 'svelte';
	import { EdraEditor, EdraToolBar, EdraBubbleMenu } from './edra/headless/index.js';
	import type { Editor } from '@tiptap/core';

	let {
		value = $bindable(''),
		placeholder = 'Tulis di sini...',
		class: className = ''
	} = $props();

	let editor = $state<Editor>();
	let isUpdatingInternal = false;

	function handleUpdate() {
		if (editor) {
			isUpdatingInternal = true;
			value = editor.getHTML();
		}
	}

	$effect(() => {
		const html = value;
		untrack(() => {
			if (editor && html !== undefined) {
				if (isUpdatingInternal) {
					isUpdatingInternal = false;
					return;
				}
				if (html !== editor.getHTML()) {
					editor.commands.setContent(html);
				}
			}
		});
	});
</script>

<div class={`rounded-md border border-slate-300 bg-white shadow-sm overflow-hidden ${className}`}>
	<div class="border-b border-slate-200 bg-slate-50 p-2">
		{#if editor}
			<EdraToolBar {editor} />
		{/if}
	</div>
	
	{#if editor}
		<EdraBubbleMenu {editor} />
	{/if}
	
	<div class="p-4 min-h-[300px] prose prose-base max-w-none">
		<EdraEditor
			bind:editor
			content={value}
			onUpdate={handleUpdate}
		/>
	</div>
</div>

<style>
	:global(.edra-editor .ProseMirror) {
		min-height: 300px;
		outline: none;
	}
	:global(.edra-editor .ProseMirror p.is-editor-empty:first-child::before) {
		color: #9ca3af;
		content: attr(data-placeholder);
		float: left;
		height: 0;
		pointer-events: none;
	}
</style>
