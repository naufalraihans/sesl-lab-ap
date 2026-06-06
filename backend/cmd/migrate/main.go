package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lab-ap/config"
	"lab-ap/database"

	"gorm.io/gorm"
)

// Migrator sederhana: membaca file *.up.sql / *.down.sql berurutan,
// mencatat versi yang sudah diterapkan di tabel schema_migrations.
//
// Pemakaian:
//   go run ./cmd/migrate up        -> terapkan semua migrasi yang belum jalan
//   go run ./cmd/migrate down      -> rollback 1 langkah terakhir
//   go run ./cmd/migrate down all  -> rollback semua

const migrationDir = "database/migration"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("pemakaian: migrate [up|down] [all]")
	}
	cmd := os.Args[1]

	cfg := config.Load()
	db, err := database.Connect(cfg.DSN(), false)
	if err != nil {
		log.Fatalf("koneksi DB gagal: %v", err)
	}
	if err := ensureTable(db); err != nil {
		log.Fatalf("gagal menyiapkan schema_migrations: %v", err)
	}

	switch cmd {
	case "up":
		runUp(db)
	case "down":
		all := len(os.Args) > 2 && os.Args[2] == "all"
		runDown(db, all)
	default:
		log.Fatalf("perintah tidak dikenal: %s", cmd)
	}
}

type schemaMigration struct {
	Version string `gorm:"primaryKey;column:version"`
}

func (schemaMigration) TableName() string { return "schema_migrations" }

func ensureTable(db *gorm.DB) error {
	return db.AutoMigrate(&schemaMigration{})
}

func appliedVersions(db *gorm.DB) (map[string]bool, error) {
	var rows []schemaMigration
	if err := db.Order("version asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[string]bool, len(rows))
	for _, r := range rows {
		m[r.Version] = true
	}
	return m, nil
}

// version mengambil prefix "001" dari "001_create_users.up.sql".
func version(file string) string {
	base := filepath.Base(file)
	if i := strings.Index(base, "_"); i > 0 {
		return base[:i]
	}
	return base
}

func listFiles(suffix string) []string {
	pattern := filepath.Join(migrationDir, "*"+suffix)
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("gagal membaca migrasi: %v", err)
	}
	sort.Strings(files)
	return files
}

func runUp(db *gorm.DB) {
	applied, err := appliedVersions(db)
	if err != nil {
		log.Fatalf("gagal baca versi: %v", err)
	}
	files := listFiles(".up.sql")
	count := 0
	for _, f := range files {
		v := version(f)
		if applied[v] {
			continue
		}
		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("baca %s gagal: %v", f, err)
		}
		if err := execSQL(db, string(sqlBytes)); err != nil {
			log.Fatalf("✗ migrasi %s gagal: %v", v, err)
		}
		if err := db.Create(&schemaMigration{Version: v}).Error; err != nil {
			log.Fatalf("gagal catat versi %s: %v", v, err)
		}
		fmt.Printf("✓ up   %s\n", filepath.Base(f))
		count++
	}
	if count == 0 {
		fmt.Println("Tidak ada migrasi baru. Database sudah terkini.")
	} else {
		fmt.Printf("Selesai: %d migrasi diterapkan.\n", count)
	}
}

func runDown(db *gorm.DB, all bool) {
	applied, err := appliedVersions(db)
	if err != nil {
		log.Fatalf("gagal baca versi: %v", err)
	}
	files := listFiles(".down.sql")
	// urutkan menurun agar rollback dari yang terbaru
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	for _, f := range files {
		v := version(f)
		if !applied[v] {
			continue
		}
		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("baca %s gagal: %v", f, err)
		}
		if err := execSQL(db, string(sqlBytes)); err != nil {
			log.Fatalf("✗ rollback %s gagal: %v", v, err)
		}
		if err := db.Delete(&schemaMigration{Version: v}).Error; err != nil {
			log.Fatalf("gagal hapus versi %s: %v", v, err)
		}
		fmt.Printf("✓ down %s\n", filepath.Base(f))
		if !all {
			return
		}
	}
	fmt.Println("Rollback selesai.")
}

// execSQL menjalankan satu file SQL yang mungkin berisi beberapa statement.
func execSQL(db *gorm.DB, raw string) error {
	stmts := strings.Split(raw, ";")
	return db.Transaction(func(tx *gorm.DB) error {
		for _, s := range stmts {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			if err := tx.Exec(s).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
