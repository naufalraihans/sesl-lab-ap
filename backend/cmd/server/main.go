package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lab-ap/config"
	"lab-ap/database"
	_ "lab-ap/docs"
	"lab-ap/internal/delivery/http/handler"
	"lab-ap/internal/delivery/http/route"
	"lab-ap/internal/repository"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/jwt"
	"lab-ap/pkg/ollama"
	"lab-ap/pkg/online"
	"lab-ap/pkg/supabase"
)

// @title           Lab-AP API
// @version         1.0
// @description     REST API Documentation for Web Lab-AP v3
// @termsOfService  http://swagger.io/terms/

// @contact.name   Developer
// @contact.email  naufalraihans@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DSN(), cfg.AppEnv == "development")
	if err != nil {
		log.Fatalf("Gagal konek database: %v", err)
	}

	// Infrastruktur stateful: registry online in-memory.
	reg := online.NewRegistry(cfg.OnlineTTL)
	stop := make(chan struct{})
	reg.StartJanitor(cfg.OnlineTTL, stop)

	jm := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpireHours)
	sb := supabase.New(cfg.SupabaseURL, cfg.SupabaseServiceKey, cfg.SupabaseBucket)
	oc := ollama.NewClient(cfg)

	// ---- Repository ----
	userRepo := repository.NewUserRepository(db)
	kelasRepo := repository.NewKelasRepository(db)
	jadwalRepo := repository.NewJadwalRepository(db)
	konfRepo := repository.NewKonfigurasiRepository(db)
	pedomanRepo := repository.NewPedomanRepository(db)
	sesiRepo := repository.NewSesiRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	soalRepo := repository.NewSoalRepository(db)
	aktivasiRepo := repository.NewAktivasiRepository(db)
	terpilihRepo := repository.NewSoalTerpilihRepository(db)
	jawabanRepo := repository.NewJawabanRepository(db)
	pengerjaanRepo := repository.NewPengerjaanRepository(db)
	ampuanRepo := repository.NewAmpuanRepository(db)
	rekapRepo := repository.NewRekapRepository(db)
	penilaianTxRepo := repository.NewPenilaianTxRepo(db)
	aktivasiTxRepo := repository.NewAktivasiTxRepo(db)

	// ---- Usecase ----
	authUC := usecase.NewAuthUsecase(userRepo, kelasRepo, jm, reg)
	profileUC := usecase.NewProfileUsecase(userRepo)
	dashboardUC := usecase.NewDashboardUsecase(userRepo, aktivasiRepo, pengerjaanRepo, reg)
	sesiUC := usecase.NewSesiUsecase(sesiRepo, courseRepo)
	soalUC := usecase.NewSoalUsecase(soalRepo, courseRepo, terpilihRepo)
	aktivasiUC := usecase.NewAktivasiUsecase(aktivasiRepo, sesiRepo, courseRepo, kelasRepo, jawabanRepo, pengerjaanRepo, soalUC, aktivasiTxRepo)
	jawabanUC := usecase.NewJawabanUsecase(aktivasiRepo, terpilihRepo, jawabanRepo, pengerjaanRepo, userRepo, courseRepo)
	penilaianUC := usecase.NewPenilaianUsecase(jawabanRepo, pengerjaanRepo, userRepo, penilaianTxRepo)
	konfUC := usecase.NewKonfigurasiUsecase(konfRepo)
	userUC := usecase.NewUserUsecase(userRepo, kelasRepo)
	kelasUC := usecase.NewKelasUsecase(kelasRepo)
	jadwalUC := usecase.NewJadwalUsecase(jadwalRepo)
	pedomanUC := usecase.NewPedomanUsecase(pedomanRepo)
	praktikumUC := usecase.NewPraktikumUsecase(sesiRepo, courseRepo, aktivasiRepo, pengerjaanRepo, userRepo, jadwalRepo)
	ampuanUC := usecase.NewAmpuanUsecase(ampuanRepo)
	rekapUC := usecase.NewRekapUsecase(rekapRepo, kelasRepo)
	aiGradingUC := usecase.NewAIGradingUsecase(jawabanRepo, penilaianUC, oc)

	// ---- Handler ----
	h := route.Handlers{
		Auth:         handler.NewAuthHandler(authUC, profileUC),
		Dashboard:    handler.NewDashboardHandler(dashboardUC),
		Sesi:         handler.NewSesiHandler(sesiUC),
		Aktivasi:     handler.NewAktivasiHandler(aktivasiUC),
		Soal:         handler.NewSoalHandler(soalUC),
		Jawaban:      handler.NewJawabanHandler(jawabanUC),
		Penilaian:    handler.NewPenilaianHandler(penilaianUC),
		Konfigurasi:  handler.NewKonfigurasiHandler(konfUC),
		Profile:      handler.NewProfileHandler(profileUC),
		Jadwal:       handler.NewJadwalHandler(jadwalUC),
		User:         handler.NewUserHandler(userUC),
		Kelas:        handler.NewKelasHandler(kelasUC),
		Pedoman:      handler.NewPedomanHandler(pedomanUC),
		Praktikum:    handler.NewPraktikumHandler(praktikumUC),
		Upload:       handler.NewUploadHandler(sb),
		Ampuan:       handler.NewAmpuanHandler(ampuanUC, userUC),
		Rekap:        handler.NewRekapHandler(rekapUC),
		AIGrading:    handler.NewAIGradingHandler(aiGradingUC),
		RekapJawaban: handler.NewRekapJawabanHandler(penilaianUC),
	}

	r := route.Setup(cfg, jm, reg, h)

	// ---- Background sweeper: auto-submit pengerjaan yang lewat deadline ----
	startSweeper(jawabanUC, cfg.SweeperInterval, stop)

	srv := &http.Server{Addr: ":" + cfg.AppPort, Handler: r}

	go func() {
		log.Printf("🚀 Server berjalan di http://localhost:%s (env=%s)", cfg.AppPort, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// ---- Graceful shutdown ----
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Mematikan server...")
	close(stop)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown paksa: %v", err)
	}
	log.Println("Server berhenti dengan rapi.")
}

// startSweeper menjalankan job periodik auto-submit (server stateful).
func startSweeper(uc *usecase.JawabanUsecase, interval time.Duration, stop <-chan struct{}) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				if n, err := uc.AutoSubmitExpired(); err != nil {
					log.Printf("[sweeper] error: %v", err)
				} else if n > 0 {
					log.Printf("[sweeper] auto-submit %d pengerjaan kedaluwarsa", n)
				}
			}
		}
	}()
}
