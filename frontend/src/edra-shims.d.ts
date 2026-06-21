// Augmentasi tipe untuk komponen editor `edra` (third-party, di-vendor ke
// src/lib/components/edra). edra mengakses storage & API ProseMirror yang belum
// dideklarasikan tipenya di @tiptap v3 — shim ini membersihkan error type-check
// tanpa mengubah kode edra. Bukan logika runtime, hanya tipe.
import '@tiptap/core';

/* eslint-disable @typescript-eslint/no-explicit-any */
declare module '@tiptap/core' {
	interface Storage {
		searchAndReplace: any;
		slashCommand: any;
	}
}

declare module '@tiptap/pm/view' {
	export function __serializeForClipboard(view: any, slice: any): { dom: HTMLElement; text: string };
}
