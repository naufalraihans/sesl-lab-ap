// Package handler adalah entrypoint serverless untuk Vercel.
// Berbeda dari cmd/server/main.go (server persisten dengan sweeper goroutine),
// file ini hanya membungkus engine Gin jadi satu fungsi serverless yang
// dipanggil Vercel tiap request. TIDAK menjalankan goroutine latar belakang —
// auto-submit ditangani endpoint /api/cron/auto-submit + cron eksternal.
package handler

import (
	"net/http"
	"strconv"
	"sync"

	"lab-ap/config"
	_ "lab-ap/docs"
	"lab-ap/internal/app"

	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine
	initMu sync.Mutex
)

// ensureInit merakit engine sekali per instance hangat. Bila gagal (mis. DB
// belum siap), error dikembalikan agar bisa dilaporkan ke klien, lalu dicoba
// lagi pada request berikutnya (bukan crash).
func ensureInit() error {
	initMu.Lock()
	defer initMu.Unlock()

	if engine != nil {
		return nil
	}

	gin.SetMode(gin.ReleaseMode)
	cfg := config.Load()
	r, _, err := app.Build(cfg)
	if err != nil {
		return err
	}
	engine = r
	return nil
}

// Handler adalah entrypoint yang dipanggil Vercel untuk tiap request.
func Handler(w http.ResponseWriter, r *http.Request) {
	if err := ensureInit(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"success":false,"message":"startup error","error":` + strconv.Quote(err.Error()) + `}`))
		return
	}
	engine.ServeHTTP(w, r)
}
