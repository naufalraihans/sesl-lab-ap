# Rencana Bug Fix & Optimasi Performa (Tahap 2)

Dokumen ini berisi catatan mengenai isu performa dan *bug* kritis yang ditemukan pada aplikasi, serta strategi penyelesaiannya untuk mempersiapkan aplikasi agar siap masuk tahap produksi (Production).

## 1. Syntax Error pada Konfigurasi (PostgreSQL)
**Masalah:** 
Pada log terminal saat menjalankan server, ditemukan *error*: `ERROR: syntax error at or near "=" (SQLSTATE 42601)`.
Penyebabnya ada di `konfigurasi_repository.go` baris 26: `Where("\`key\` = ?")`. Penggunaan karakter *backtick* (\`) adalah sintaks bawaan MySQL. Karena database sudah dimigrasikan ke PostgreSQL, PostgreSQL tidak mengenali sintaks tersebut.

**Cara Memperbaikinya:**
- Mengganti sintaks MySQL *backtick* menjadi sintaks *double quotes* standar PostgreSQL atau GORM di dalam file `konfigurasi_repository.go`.
- Perbaikan kode: Mengubah `r.db.Where("\`key\` = ?", key)` menjadi `r.db.Where("\"key\" = ?", key)`.

## 2. Slow SQL (Database Bottleneck) pada Worker Pengerjaan Course
**Masalah:**
*Log* terminal mendeteksi peringatan `SLOW SQL >= 200ms` pada *query* milik `pengerjaan_repository.go` yang digunakan oleh *background worker* (Auto-Submit / timer habis). 
*Query*: `SELECT pc.mahasiswa_id, pc.aktivasi_sesi_id, pc.course_id FROM pengerjaan_course AS pc JOIN course c ON c.id = pc.course_id WHERE pc.status = 'sedang_dikerjakan' AND pc.waktu_mulai IS NOT NULL AND NOW() > pc.waktu_mulai + (c.durasi_menit * interval '1 minute')`
Waktu eksekusi mencapai 257ms karena PostgreSQL harus memindai seluruh isi tabel (*Full Table Scan*). 

**Cara Memperbaikinya:**
- Menambahkan **Composite Index** (*Database Indexing*) pada kolom `status` dan `waktu_mulai` di dalam file migrasi `014_create_pengerjaan_course.sql`.
- Index ini akan memungkinkan PostgreSQL untuk langsung melompat ke baris yang memiliki status `'sedang_dikerjakan'`, memangkas waktu *query* dari 250ms+ menjadi di bawah 5ms, sehingga server tidak akan tercekik ketika data mahasiswa membengkak.

## 3. Implementasi Server Caching (Backend In-Memory)
**Masalah:**
Pengguna mengeluhkan waktu muat (loading) yang lama pada beberapa skenario. Konfigurasi global (seperti `jadwal_mode` dan `gdrive_link`) dipanggil dari *database* secara repetitif setiap kali halaman publik atau login dimuat.

**Cara Memperbaikinya:**
- Menerapkan mekanisme *In-Memory Caching* di Go (menggunakan `sync.Map` atau `go-cache`) pada layer Usecase atau Repository.
- Menyimpan nilai-nilai konfigurasi di dalam memori server dengan masa aktif (TTL - Time To Live) sekitar 5 menit. Jika data sudah ada di *cache*, server tidak perlu lagi mengirim *query* ke database PostgreSQL.
- Teknik yang sama juga dapat diterapkan untuk *endpoint* Dashboard Admin yang melakukan agregasi berat (menghitung total pengguna, total asisten, dll).

## 4. Frontend Caching & SvelteKit Optimization
**Masalah:**
Halaman statis di frontend (seperti `/info`) masih selalu meminta data baru ke server setiap kali berpindah navigasi, dan *preloading* SvelteKit yang terlalu agresif menyebabkan peringatan jaringan.

**Cara Memperbaikinya:**
- Menambahkan `setHeaders({ 'cache-control': 'public, max-age=300' })` pada blok `load` di `+page.server.ts` atau `+page.ts` untuk halaman publik. Ini menginstruksikan *browser* pengguna untuk menggunakan versi *cache* selama 5 menit.
- Mengubah opsi agresif SvelteKit `data-sveltekit-preload-data="hover"` menjadi `"tap"` pada tautan-tautan tertentu agar *preloading* hanya dilakukan saat tombol benar-benar disentuh/diklik.
