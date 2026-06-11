package database

import (
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string, debug bool) (*gorm.DB, error) {
	logLevel := logger.Warn
	if debug {
		logLevel = logger.Info
	}

	// connect_timeout supaya cold-start serverless gagal cepat (bukan menggantung).
	if !strings.Contains(dsn, "connect_timeout") {
		dsn += " connect_timeout=15"
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		// PreferSimpleProtocol mematikan prepared statement — WAJIB untuk
		// Supabase Transaction Pooler (pgbouncer mode transaction) yang dipakai
		// di serverless. Tanpa ini, query akan error "prepared statement exists".
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// Pool kecil: di serverless koneksi mudah menumpuk lewat pooler.
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("✓ Terhubung ke PostgreSQL")
	return db, nil
}
