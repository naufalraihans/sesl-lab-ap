import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { browser } from '$app/environment';

const BASE = PUBLIC_API_BASE_URL || 'http://localhost:8080';

export interface Envelope<T = unknown> {
	success: boolean;
	message: string;
	data?: T;
	error?: unknown;
}

export class ApiError extends Error {
	status: number;
	constructor(message: string, status: number) {
		super(message);
		this.status = status;
	}
}

function token(): string | null {
	if (!browser) return null;
	return localStorage.getItem('token');
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
	const headers: Record<string, string> = {};
	const t = token();
	if (t) headers['Authorization'] = `Bearer ${t}`;

	let payload: BodyInit | undefined;
	if (body instanceof FormData) {
		payload = body;
	} else if (body !== undefined) {
		headers['Content-Type'] = 'application/json';
		payload = JSON.stringify(body);
	}

	const res = await fetch(`${BASE}${path}`, { method, headers, body: payload });
	let json: Envelope<T>;
	try {
		json = await res.json();
	} catch {
		throw new ApiError(`HTTP ${res.status}`, res.status);
	}
	if (!res.ok || !json.success) {
		throw new ApiError(json.message || `HTTP ${res.status}`, res.status);
	}
	return json.data as T;
}

export const api = {
	get: <T>(path: string) => request<T>('GET', path),
	post: <T>(path: string, body?: unknown) => request<T>('POST', path, body),
	put: <T>(path: string, body?: unknown) => request<T>('PUT', path, body),
	del: <T>(path: string) => request<T>('DELETE', path),
	upload: <T>(path: string, form: FormData) => request<T>('POST', path, form),
	base: BASE
};
