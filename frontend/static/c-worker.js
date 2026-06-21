// Worker C interaktif (Fase 1 POC). Menjalankan wasm32-wasi di browser; stdin
// di-blok pakai Atomics.wait pada SharedArrayBuffer — protokol SAB SAMA PERSIS
// dengan playground-worker.js, jadi UI terminal bisa dipakai ulang.
import { WASI, Fd, WASIProcExit } from 'https://esm.sh/@bjorn3/browser_wasi_shim@0.4.2';

let meta = null; // Int32Array: [0]=flag siap, [1]=panjang data
let dataArr = null; // Uint8Array: byte input dari terminal

// stdin (fd 0): blok sampai main-thread mengirim satu baris dari terminal.
class SabStdin extends Fd {
	fd_read(size) {
		postMessage({ type: 'need-input' });
		Atomics.wait(meta, 0, 0); // tunggu meta[0] -> 1
		const len = Atomics.load(meta, 1);
		const bytes = dataArr.slice(0, len);
		Atomics.store(meta, 0, 0); // reset untuk input berikutnya
		// ponytail: sisa byte > size diabaikan; cukup untuk input per-baris (scanf/gets)
		return { ret: 0, data: bytes.subarray(0, size) };
	}
}

// stdout(1)/stderr(2): teruskan ke terminal.
class OutFd extends Fd {
	constructor(kind) { super(); this.kind = kind; }
	fd_write(data) {
		postMessage({ type: this.kind, data: new TextDecoder().decode(data) });
		return { ret: 0, nwritten: data.byteLength };
	}
}

async function run(wasmBytes) {
	try {
		// wasmBytes: hasil compile dari server (Fase 2). POC: fallback ke /scanf.wasm.
		const bytes = wasmBytes ?? (await (await fetch('/scanf.wasm')).arrayBuffer());
		const wasi = new WASI([], [], [new SabStdin(), new OutFd('out'), new OutFd('err')]);
		const { instance } = await WebAssembly.instantiate(bytes, {
			wasi_snapshot_preview1: wasi.wasiImport,
		});
		try { wasi.start(instance); }
		catch (e) { if (!(e instanceof WASIProcExit)) throw e; }
	} catch (e) {
		postMessage({ type: 'err', data: String((e && e.message) || e) });
	} finally {
		postMessage({ type: 'done' });
	}
}

onmessage = (e) => {
	const m = e.data;
	if (m.type === 'init') {
		meta = new Int32Array(m.sab, 0, 2);
		dataArr = new Uint8Array(m.sab, 8);
		postMessage({ type: 'ready' });
	} else if (m.type === 'run') {
		run(m.wasm); // m.wasm = ArrayBuffer (Fase 2) atau undefined (POC)
	}
};
