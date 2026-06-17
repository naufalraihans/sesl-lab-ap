import type { Handle } from '@sveltejs/kit';

// Hanya /playground yang di-isolasi (COOP/COEP) agar SharedArrayBuffer aktif
// untuk runner Python interaktif. Halaman lain TIDAK ikut, supaya resource
// cross-origin (gambar Supabase, dll) tetap dimuat normal.
export const handle: Handle = async ({ event, resolve }) => {
	const response = await resolve(event);
	if (event.url.pathname.startsWith('/playground')) {
		response.headers.set('Cross-Origin-Opener-Policy', 'same-origin');
		response.headers.set('Cross-Origin-Embedder-Policy', 'credentialless');
	}
	return response;
};
