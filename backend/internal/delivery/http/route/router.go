package route

import (
	"net/http"
	"time"

	"lab-ap/config"
	"lab-ap/internal/delivery/http/handler"
	"lab-ap/internal/delivery/http/middleware"
	"lab-ap/internal/entity"
	"lab-ap/pkg/jwt"
	"lab-ap/pkg/online"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers mengumpulkan seluruh handler untuk diregistrasi.
type Handlers struct {
	Auth        *handler.AuthHandler
	Dashboard   *handler.DashboardHandler
	Sesi        *handler.SesiHandler
	Aktivasi    *handler.AktivasiHandler
	Soal        *handler.SoalHandler
	Jawaban     *handler.JawabanHandler
	Penilaian   *handler.PenilaianHandler
	Konfigurasi *handler.KonfigurasiHandler
	Profile     *handler.ProfileHandler
	Jadwal      *handler.JadwalHandler
	User        *handler.UserHandler
	Kelas       *handler.KelasHandler
	Pedoman     *handler.PedomanHandler
	Praktikum   *handler.PraktikumHandler
	Upload      *handler.UploadHandler
	Ampuan       *handler.AmpuanHandler
	Rekap        *handler.RekapHandler
	AIGrading    *handler.AIGradingHandler
	RekapJawaban *handler.RekapJawabanHandler
}

// HealthCheck GET /api/health
// @Summary Health Check
// @Description Mengecek status service backend
// @Tags Info
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "lab-ap"})
}

// Setup membangun engine Gin + seluruh route & middleware.
func Setup(cfg *config.Config, jm *jwt.Manager, reg *online.Registry, h Handlers) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/health", HealthCheck)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	// ---- Auth (publik) ----
	auth := api.Group("/auth")
	{
		auth.POST("/cek-nim", h.Auth.CekNIM)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/register", h.Auth.Register)
	}

	// ---- Info (publik) ----
	info := api.Group("/info")
	{
		info.GET("/asisten", h.User.ListAsisten)
		info.GET("/laporan", h.Pedoman.List)
		info.GET("/modul", h.Konfigurasi.PublicModul)
		info.GET("/jadwal", h.Jadwal.List)
		info.GET("/jadwal/config", h.Konfigurasi.PublicJadwalConfig)
		info.GET("/kelas/:id/mahasiswa", h.Ampuan.PublicKelasMahasiswa)
	}

	authmw := middleware.Auth(jm, reg)

	// ---- Terautentikasi (semua role) ----
	authed := api.Group("")
	authed.Use(authmw)
	{
		authed.POST("/auth/logout", h.Auth.Logout)
		authed.GET("/auth/me", h.Auth.Me)
		authed.GET("/profile", h.Profile.Get)
		authed.PUT("/profile", h.Profile.Update)
	}

	// ---- Praktikum (role user) ----
	prak := api.Group("/praktikum")
	prak.Use(authmw, middleware.RequireRole(string(entity.RoleUser)))
	{
		prak.GET("/dashboard", h.Praktikum.Dashboard)
		prak.GET("/sesi", h.Praktikum.ListSesi)
		prak.GET("/ruang", h.Jawaban.GetRuang)
		prak.POST("/mulai", h.Jawaban.Mulai)
		prak.POST("/autosave", h.Jawaban.AutoSave)
		prak.POST("/submit", h.Jawaban.Submit)
	}

	// ---- Admin (role admin) ----
	admin := api.Group("/admin")
	admin.Use(authmw, middleware.RequireRole(string(entity.RoleAdmin)))
	{
		admin.GET("/dashboard", h.Dashboard.Statistik)
		admin.GET("/dashboard/online", h.Dashboard.Online)

		// Users (mahasiswa)
		admin.GET("/users", h.User.ListMahasiswa)
		admin.POST("/users", h.User.CreateMahasiswa)
		admin.POST("/users/bulk", h.User.BulkUpsertMahasiswa)
		admin.PUT("/users/:id", h.User.UpdateMahasiswa)
		admin.DELETE("/users/:id", h.User.Delete)
		admin.POST("/users/:id/reset-password", h.User.ResetPassword)

		// Kelas
		admin.GET("/kelas", h.Kelas.List)
		admin.POST("/kelas", h.Kelas.Create)
		admin.PUT("/kelas/:id", h.Kelas.Update)
		admin.DELETE("/kelas/:id", h.Kelas.Delete)
		admin.POST("/kelas-register", h.User.SetRegisterOpen)

		// Asisten
		admin.GET("/asisten", h.User.ListAsisten)
		admin.POST("/asisten", h.User.CreateAsisten)
		admin.PUT("/asisten/:id", h.User.UpdateAsisten)

		// Jadwal
		admin.GET("/jadwal", h.Jadwal.List)
		admin.POST("/jadwal", h.Jadwal.Create)
		admin.PUT("/jadwal/:id", h.Jadwal.Update)
		admin.DELETE("/jadwal/:id", h.Jadwal.Delete)

		// Pedoman
		admin.GET("/pedoman", h.Pedoman.List)
		admin.POST("/pedoman", h.Pedoman.Create)
		admin.PUT("/pedoman/:id", h.Pedoman.Update)
		admin.DELETE("/pedoman/:id", h.Pedoman.Delete)

		// Konfigurasi (termasuk modul & gdrive jadwal URL)
		admin.GET("/konfigurasi", h.Konfigurasi.All)
		admin.POST("/konfigurasi", h.Konfigurasi.Set)

		// Sesi & Course
		admin.GET("/sesi", h.Sesi.List)
		admin.GET("/sesi/:id", h.Sesi.Get)
		admin.POST("/sesi", h.Sesi.Create)
		admin.PUT("/sesi/:id", h.Sesi.Update)
		admin.DELETE("/sesi/:id", h.Sesi.Delete)
		admin.GET("/sesi/:id/course", h.Sesi.ListCourse)
		admin.POST("/sesi/:id/course", h.Sesi.CreateCourse)
		admin.PUT("/course/:courseId", h.Sesi.UpdateCourse)
		admin.DELETE("/course/:courseId", h.Sesi.DeleteCourse)

		// Soal
		admin.GET("/soal", h.Soal.ListByCourse)
		admin.POST("/soal", h.Soal.Create)
		admin.PUT("/soal/:id", h.Soal.Update)
		admin.DELETE("/soal/:id", h.Soal.Delete)

		// Aktivasi
		admin.GET("/aktivasi", h.Aktivasi.List)
		admin.POST("/aktivasi", h.Aktivasi.Aktivasi)
		admin.GET("/aktivasi/:id", h.Aktivasi.Get)
		admin.POST("/aktivasi/:id/token", h.Aktivasi.GenerateToken)
		admin.POST("/aktivasi/:id/susulan", h.Aktivasi.AddSusulan)
		admin.GET("/aktivasi/:id/susulan", h.Aktivasi.ListSusulan)
		admin.DELETE("/aktivasi/:id/susulan/:mahasiswaId", h.Aktivasi.RemoveSusulan)
		admin.POST("/aktivasi-course/buka-tutup", h.Aktivasi.BukaTutupCourse)

		// Penilaian
		admin.GET("/penilaian/rekap", h.Penilaian.Rekap)
		admin.POST("/penilaian", h.Penilaian.SetNilai)
		admin.POST("/penilaian/ai-grade/bulk", h.AIGrading.BulkGrade)
		
		// Rekap Jawaban Global & Bulk Actions
		admin.GET("/rekap-jawaban", h.RekapJawaban.GetRekapJawabanGlobal)
		admin.POST("/penilaian/bulk-action", h.RekapJawaban.BulkAction)

		// Background Jobs
		admin.GET("/jobs/:id", h.AIGrading.GetJobStatus)

		// Rekap
		admin.GET("/rekap/kelas/:id_kelas", h.Rekap.GetRekapKelas)

		// Upload (Supabase)
		admin.POST("/upload", h.Upload.Upload)

		admin.GET("/ampuan", h.Ampuan.List)
		admin.POST("/ampuan", h.Ampuan.Create)
		admin.DELETE("/ampuan/:id", h.Ampuan.Delete)
	}

	return r
}
