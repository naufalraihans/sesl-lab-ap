// Package cwasm mengompilasi sumber C ke wasm32-wasi via clang, untuk dijalankan
// di browser mahasiswa (lewat WASI shim). Server HANYA compile sekali per request;
// eksekusi program + stdin interaktif tetap di browser. Bukan sandbox eksekusi.
package cwasm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Compiler struct {
	clang   string        // path/binary clang (mis. "clang")
	sysroot string        // wasi-sysroot (mis. /opt/wasi-sdk/share/wasi-sysroot)
	timeout time.Duration // batas waktu satu kompilasi
}

// New membuat compiler. Nonaktif bila clang/sysroot kosong.
func New(clang, sysroot string) *Compiler {
	return &Compiler{clang: clang, sysroot: sysroot, timeout: 10 * time.Second}
}

// Enabled menandakan compiler dikonfigurasi.
func (c *Compiler) Enabled() bool { return c.clang != "" && c.sysroot != "" }

// Compile mengembalikan (wasm, "", nil) bila sukses; (nil, stderr, nil) bila kode
// gagal dikompilasi (data, bukan error); (nil, "", err) untuk kegagalan infra.
//
// ponytail: compile-only di temp dir + timeout. Untrusted source -> risiko sebatas
// CPU/RAM/disk kompilasi (bukan RCE). Isolasi lebih (cgroup/container) urusan deploy,
// tambah kalau throughput/abuse jadi masalah.
func (c *Compiler) Compile(source string) ([]byte, string, error) {
	if !c.Enabled() {
		return nil, "", fmt.Errorf("compiler C belum dikonfigurasi")
	}
	dir, err := os.MkdirTemp("", "cwasm-*")
	if err != nil {
		return nil, "", err
	}
	defer os.RemoveAll(dir)

	src := filepath.Join(dir, "main.c")
	out := filepath.Join(dir, "main.wasm")
	if err := os.WriteFile(src, []byte(source), 0o600); err != nil {
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, c.clang,
		"--target=wasm32-wasi",
		"--sysroot="+c.sysroot,
		"-O2",
		"-o", out, src,
	)
	cmd.Stderr = &stderr

	if runErr := cmd.Run(); runErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, "Kompilasi melebihi batas waktu.", nil
		}
		// Exit non-zero dari clang = gagal kompilasi -> kembalikan stderr sebagai data.
		var ee *exec.ExitError
		if errors.As(runErr, &ee) {
			return nil, stderr.String(), nil
		}
		return nil, "", runErr // clang tak ditemukan / kegagalan infra
	}

	wasm, err := os.ReadFile(out)
	if err != nil {
		return nil, "", err
	}
	return wasm, "", nil
}
