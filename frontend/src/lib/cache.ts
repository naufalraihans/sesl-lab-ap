import { api } from './api';

// Cache GET sederhana di memori (per-sesi SPA). Pola stale-while-revalidate:
// saat halaman dibuka lagi, data lama tampil INSTAN lalu disegarkan diam-diam.
const store = new Map<string, unknown>();

/**
 * GET dengan pola stale-while-revalidate.
 * - Bila ada cache, `onUpdate` dipanggil seketika (tampil instan, tanpa loading).
 * - Lalu data fresh di-fetch; `onUpdate` dipanggil lagi dengan hasil terbaru.
 * - Error hanya dilempar bila TIDAK ada cache yang bisa ditampilkan.
 */
export async function swrGet<T>(path: string, onUpdate: (value: T) => void): Promise<void> {
	const cached = store.get(path) as T | undefined;
	if (cached !== undefined) onUpdate(cached);
	try {
		const fresh = await api.get<T>(path);
		store.set(path, fresh);
		onUpdate(fresh);
	} catch (e) {
		if (cached === undefined) throw e; // tak ada fallback → surface error
	}
}

/** Kosongkan seluruh cache. Dipanggil saat logout agar data antar-user tak bocor. */
export function clearSwrCache(): void {
	store.clear();
}
