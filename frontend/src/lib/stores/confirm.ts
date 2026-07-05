import { writable } from 'svelte/store';

export interface ConfirmConfig {
	show: boolean;
	title: string;
	message: string;
	confirmText: string;
	cancelText: string;
	danger: boolean;
	resolve?: (value: boolean) => void;
}

export const confirmStore = writable<ConfirmConfig>({
	show: false,
	title: '',
	message: '',
	confirmText: 'Ya, Lanjutkan',
	cancelText: 'Batalkan',
	danger: true
});

export function confirmAction(config: Partial<Omit<ConfirmConfig, 'show' | 'resolve'>>): Promise<boolean> {
	return new Promise((resolve) => {
		confirmStore.set({
			show: true,
			title: config.title ?? 'Konfirmasi Tindakan',
			message: config.message ?? 'Apakah Anda yakin ingin melanjutkan?',
			confirmText: config.confirmText ?? 'Ya, Lanjutkan',
			cancelText: config.cancelText ?? 'Batalkan',
			danger: config.danger ?? true,
			resolve
		});
	});
}
