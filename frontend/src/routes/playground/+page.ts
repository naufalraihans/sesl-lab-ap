// Playground butuh runtime browser (Pyodide/WASM) + tidak boleh di-prerender
// supaya header COOP/COEP dari hooks.server.ts ikut terpasang.
export const ssr = false;
export const prerender = false;
