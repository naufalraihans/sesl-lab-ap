import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { User } from '$lib/types';

function load(): User | null {
	if (!browser) return null;
	const raw = localStorage.getItem('user');
	return raw ? (JSON.parse(raw) as User) : null;
}

export const user = writable<User | null>(load());

export function setAuth(token: string, u: User) {
	if (browser) {
		localStorage.setItem('token', token);
		localStorage.setItem('user', JSON.stringify(u));
	}
	user.set(u);
}

export function clearAuth() {
	if (browser) {
		localStorage.removeItem('token');
		localStorage.removeItem('user');
	}
	user.set(null);
}

export function hasToken(): boolean {
	return browser && !!localStorage.getItem('token');
}
