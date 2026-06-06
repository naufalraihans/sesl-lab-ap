package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config menampung seluruh konfigurasi aplikasi yang dimuat dari environment.
type Config struct {
	AppPort string
	AppEnv  string

	CORSOrigins []string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBParams   string

	JWTSecret      string
	JWTExpireHours int

	OnlineTTL       time.Duration
	SweeperInterval time.Duration

	SupabaseURL        string
	SupabaseServiceKey string
	SupabaseBucket     string
}

// Load membaca .env (jika ada) lalu environment OS.
func Load() *Config {
	// Coba muat .env dari beberapa lokasi umum (root backend / root project).
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env")

	cfg := &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		CORSOrigins:        splitCSV(getEnv("CORS_ORIGINS", "http://localhost:5173")),
		DBHost:             getEnv("DB_HOST", "127.0.0.1"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "lab_ap"),
		DBParams:           getEnv("DB_PARAMS", "sslmode=disable TimeZone=Asia/Jakarta"),
		JWTSecret:          getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpireHours:     getEnvInt("JWT_EXPIRE_HOURS", 12),
		OnlineTTL:          time.Duration(getEnvInt("ONLINE_TTL_SECONDS", 120)) * time.Second,
		SweeperInterval:    time.Duration(getEnvInt("SWEEPER_INTERVAL_SECONDS", 30)) * time.Second,
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseBucket:     getEnv("SUPABASE_BUCKET", "public-assets"),
	}
	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBParams)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
