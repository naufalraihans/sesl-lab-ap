// Sama seperti /playground: butuh runtime browser (WASM) + tak boleh prerender
// supaya header COOP/COEP dari hooks.server.ts ikut terpasang (SharedArrayBuffer).
export const ssr = false;
export const prerender = false;
