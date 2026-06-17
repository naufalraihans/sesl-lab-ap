// Worker untuk menjalankan Python interaktif via Pyodide (WASM) di browser.
// stdin di-blok pakai Atomics.wait pada SharedArrayBuffer sampai main-thread
// mengirim baris input dari terminal. (Fase 1 POC — Update 7)

/* global importScripts, loadPyodide */
importScripts('https://cdn.jsdelivr.net/pyodide/v0.26.4/full/pyodide.js');

let pyodide = null;
let meta = null; // Int32Array: [0]=flag siap, [1]=panjang data
let dataArr = null; // Uint8Array: byte input

async function init(sab) {
	meta = new Int32Array(sab, 0, 2);
	dataArr = new Uint8Array(sab, 8);
	try {
		pyodide = await loadPyodide();
		pyodide.setStdout({ batched: (s) => postMessage({ type: 'out', data: s }) });
		pyodide.setStderr({ batched: (s) => postMessage({ type: 'err', data: s }) });
		pyodide.setStdin({ stdin: readLine });
		postMessage({ type: 'ready' });
	} catch (e) {
		postMessage({ type: 'fatal', data: String((e && e.message) || e) });
	}
}

// Dipanggil Pyodide saat Python butuh input. Blok sampai main mengirim baris.
function readLine() {
	postMessage({ type: 'need-input' });
	// Tunggu sampai meta[0] menjadi 1 (data tersedia).
	Atomics.wait(meta, 0, 0);
	const len = Atomics.load(meta, 1);
	const bytes = dataArr.slice(0, len);
	Atomics.store(meta, 0, 0); // reset untuk input berikutnya
	// Sertakan newline agar input() menutup baris.
	return new TextDecoder().decode(bytes);
}

async function run(code) {
	try {
		await pyodide.runPythonAsync(code);
	} catch (e) {
		postMessage({ type: 'err', data: String((e && e.message) || e) });
	} finally {
		postMessage({ type: 'done' });
	}
}

onmessage = (e) => {
	const m = e.data;
	if (m.type === 'init') init(m.sab);
	else if (m.type === 'run') run(m.code);
};
