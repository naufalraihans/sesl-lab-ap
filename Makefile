# Laboratorium Algoritma dan Pemrograman - Makefile
# Semua perintah dijalankan dari root; backend ada di ./backend

BACKEND_DIR := backend

.PHONY: run build migrate-up migrate-down seed tidy fe-install fe-dev fe-build help

help:
	@echo "Target tersedia:"
	@echo "  make run          - Jalankan server backend (go run ./cmd/server)"
	@echo "  make build        - Build binary server"
	@echo "  make migrate-up   - Jalankan semua migrasi (.up.sql)"
	@echo "  make migrate-down - Rollback satu langkah migrasi (.down.sql)"
	@echo "  make seed         - Isi data awal (admin, kelas, jadwal, soal contoh)"
	@echo "  make swag         - Generate OpenAPI documentation (swag init)"
	@echo "  make tidy         - go mod tidy"
	@echo "  make fe-install   - Install dependency frontend"
	@echo "  make fe-dev       - Jalankan frontend dev server"
	@echo "  make fe-build     - Build frontend"
	@echo "  make test         - Run all tests (backend)"
	@echo "  make mock         - Generate mocks untuk testing (backend)"

run:
	cd $(BACKEND_DIR) && go run ./cmd/server

build:
	cd $(BACKEND_DIR) && go build -o bin/server ./cmd/server

migrate-up:
	cd $(BACKEND_DIR) && go run ./cmd/migrate up

migrate-down:
	cd $(BACKEND_DIR) && go run ./cmd/migrate down

migrate-drop:
	cd $(BACKEND_DIR) && go run ./cmd/migrate down all

seed:
	cd $(BACKEND_DIR) && go run ./database/seed

swag:
	cd $(BACKEND_DIR) && go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go

tidy:
	cd $(BACKEND_DIR) && go mod tidy

test:
	cd $(BACKEND_DIR) && go test -v ./...

mock:
	cd $(BACKEND_DIR) && go run github.com/vektra/mockery/v2@latest --all --keeptree --dir=internal/repository --output=internal/repository/mocks

fe-install:
	cd frontend && npm install

fe-dev:
	cd frontend && npm run dev

fe-build:
	cd frontend && npm run build
