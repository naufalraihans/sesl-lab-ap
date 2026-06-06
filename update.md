# Rencana Pengembangan Web Lab-AP v3

Dokumen ini berisi peta jalan (roadmap) pengembangan fitur-fitur lanjutan yang diadaptasi dari sistem lama (`lab-ap-v2`) untuk diterapkan ke arsitektur sistem baru.

---

## Aturan Utama Pengembangan (UI/UX & Security)
1. **Tidak Boleh Ada Emoji**: Website bersifat formal untuk lingkungan akademik. Penggunaan emoji dilarang keras di seluruh antarmuka pengguna (UI) maupun *source code*.
2. **Gunakan SVG Icon**: Wajib menggunakan ikon berformat SVG (seperti *Lucide* atau *Heroicons*) untuk keperluan visual.
3. **Bahasa Formal**: Gunakan tata bahasa Indonesia yang baku dan profesional.
4. **Desain Interaktif & Modern**: UI harus dirancang kreatif, premium, dan dinamis (animasi *hover*, *glassmorphism*) tanpa merombak palet warna dasar saat ini.
5. **[SECURITY] Middleware RBAC (Role-Based Access Control)**: Seluruh *endpoint* sensitif yang baru dibuat (AI Grading, Import CSV, Rekap) WAJIB dilindungi oleh *middleware* keamanan di Go. Hanya *User* dengan otorisasi `Admin` atau `Asisten` yang berhak mengeksekusi *endpoint* tersebut untuk mencegah eksploitasi oleh Mahasiswa.
6. **[DATABASE] Soft Deletes**: Mengimplementasikan *Soft Delete* (menggunakan fitur `gorm.DeletedAt`) pada tabel-tabel krusial. Alih-alih menghapus data secara permanen dari PostgreSQL (Supabase), data hanya akan disembunyikan (*flagged*). Ini sebagai jaring pengaman (*safety net*) jika ada data penting yang tidak sengaja terhapus.

---

## 1. Integrasi AI Grading (Penilaian Otomatis) & Background Queue
**Tujuan:** Meringankan beban asisten laboratorium dalam menilai puluhan/ratusan soal menggunakan bantuan Large Language Models (LLM).

**Rencana Implementasi (Backend - Go):**
*   **Modifikasi Skema Database**: Menambahkan kolom `kunci_jawaban` (TEXT) pada tabel `soal` sebagai acuan pasti bagi AI.
*   **Background Worker / Queue System [CRITICAL]**: Karena proses *grading* menggunakan LLM memakan waktu lama, endpoint `POST /api/penilaian/ai-grade/bulk` tidak akan merespons langsung dengan hasil akhir (untuk mencegah *HTTP Timeout*). 
    *   Backend akan memasukkan tugas *grading* ke dalam antrean (*Go Channels* atau *Redis Queue*).
    *   Backend mengembalikan status `202 Accepted` ("Tugas Diterima").
*   **Struktur Prompting Berbasis Konteks:** Go merakit prompt: *"Soal: [Teks Soal]. Kunci Jawaban: [Kunci Jawaban]. Jawaban Mahasiswa: [Jawaban]. Evaluasi dan beri skor 0-100."*

**Rencana Implementasi (Frontend - SvelteKit):**
*   Tambahkan tombol **"Bulk AI Grade"** di menu penilaian.
*   SvelteKit akan menerapkan sistem *Polling* atau WebSockets untuk mengecek status antrean. Layar akan menampilkan progres *real-time* (Misal: *"Memproses: 45/100 jawaban"*).

---

## 2. Dashboard Rekap Jawaban, Pencarian, & Pivot Nilai Akhir
**Tujuan:** Memberikan pandangan menyeluruh (bird's-eye view) performa kelas.

**Rencana Implementasi:**
*   **Backend:** Membuat *Query SQL Custom* (Pivot) untuk merangkai data dari `users`, `course`, dan `pengerjaan_jawaban`. Siapkan endpoint `GET /api/rekap/kelas/:id_kelas`.
*   **Frontend:** Rute `/praktikum/admin/rekap-nilai`. Menambahkan kolom pencarian yang reaktif (NIM/Nama), tabel Pivot dinamis, dan tombol **"Export to Excel/CSV"**.

---

## 3. Sistem Token Ujian (PIN Akses)
**Tujuan:** Memastikan kehadiran fisik mahasiswa dengan mewajibkan input Token/PIN yang dibagikan di dalam ruangan lab.

**Rencana Implementasi:**
*   **Backend:** Tambah kolom `Token` pada tabel `AktivasiSesi` atau `Course`. *Endpoint* `Mulai Ujian` akan memvalidasi *payload* token.
*   **Frontend:** Modifikasi *form* aktivasi sesi agar admin bisa membuat token acak. Buat layar *Lock/PIN* berdesain modern bagi mahasiswa yang akan memulai ujian.

---

## 4. Bulk Update & Import CSV Mahasiswa
**Tujuan:** Mempercepat *onboarding* data mahasiswa baru via CSV.

**Rencana Implementasi:**
*   **Backend:** Eksekusi `Bulk Insert` (atau `Upsert`) menggunakan transaksi GORM untuk file CSV yang diunggah.
*   **Frontend:** Modal *Upload* CSV dengan tabel pratinjau (*preview*). 
    *   **Penanganan Error Terpusat:** SvelteKit akan menampilkan daftar peringatan tegas yang menunjuk nomor baris yang bermasalah (contoh: *"Gagal di Baris 14: Kelas tidak valid"*) berdasarkan *response* spesifik dari backend.

---

## 5. WYSIWYG Rich Text Editor (Edra) untuk Pembuatan Soal
**Tujuan:** Menerapkan editor ala *Notion* (*What You See Is What You Get*) untuk pembuatan soal yang mendukung tabel, format teks, dan rumus matematika (KaTeX).

**Rencana Implementasi:**
*   **Backend:** Ubah kolom `teks_soal` menjadi tipe `TEXT` dan terapkan Sanitasi HTML (XSS Protection) menggunakan `bluemonday`.
*   **Frontend:** Integrasi *library* **Edra** di sisi Admin, dan gunakan komponen *renderer* di sisi Mahasiswa.
*   **Manajemen Media:** Gambar pada soal diunggah langsung ke *Bucket Storage* (seperti Supabase) dari frontend, lalu URL publiknya disematkan ke Edra (via *iframe* / `<img>`). *Database* hanya menyimpan URL, bukan `Base64`.

---

## 6. Dokumentasi API (OpenAPI / Swagger)
**Tujuan:** *Single Source of Truth* untuk kontrak API antara Backend (Go) dan Frontend (SvelteKit).

**Rencana Implementasi:**
*   Implementasi `swaggo/swag`. Sistem otomatis men-generate file dokumentasi interaktif yang bisa diakses di rute `GET /swagger/index.html`.

---

## 7. Kontrak API (DTO & Response Examples)

### A. API Rekap Nilai Pivot (`GET /api/rekap/kelas/:id_kelas`)
```json
{
  "data": {
    "mahasiswa": [
      { "nim": "11223344", "nama": "Udin", "scores": {"course_1": 85.5}, "total_score": 85.5 }
    ]
  }
}
```

### B. API AI Grading Bulk / Queue (`POST /api/penilaian/ai-grade/bulk`)
*Menyesuaikan dengan implementasi Background Worker:*
```json
// Request Body dari SvelteKit
{
  "course_id": "uuid-course-xxx",
  "mahasiswa_ids": ["uuid-1", "uuid-2", "uuid-3"]
}

// Response Status: 202 Accepted (Tugas Diterima di Antrean)
{
  "message": "Permintaan AI Grading sedang diproses di latar belakang.",
  "data": {
    "job_id": "job-ai-grading-8899",
    "total_queued": 3
  }
}
```
*(SvelteKit lalu memanggil `GET /api/jobs/job-ai-grading-8899` setiap 3 detik untuk mengecek status proses)*

### C. API Import CSV (`POST /api/admin/users/import`)
```json
// Response Status: 200 OK
{
  "message": "Import diproses",
  "data": {
    "success_count": 148,
    "failed_count": 2,
    "errors": [
      "Baris 45: NIM 11223344 sudah terdaftar"
    ]
  }
}
```

---

## 8. Rencana Struktur Direktori & File Baru (Tree Format)
Berikut adalah pemetaan *file* dan direktori baru menyesuaikan struktur *project* yang sudah ada saat ini:

### ⚙️ Backend (Go)
```text
backend/
├── cmd/
│   └── server/
│       └── main.go                  <-- [MODIFY] Registrasi rute Swagger & Handler baru
├── database/
│   └── migration/
│       └── 017_update_fitur_v3.sql  <-- [NEW] Migrasi DB (Token ujian, Tipe TEXT Edra, Kolom Kunci Jawaban)
├── docs/                            <-- [NEW] Folder auto-generated oleh Swaggo (WAJIB MASUK .GITIGNORE)
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── middleware/
│   │       │   └── auth_middleware.go     <-- [MODIFY] Pengetatan role Admin/Asisten (RBAC)
│   │       └── handler/
│   │           ├── ai_grading_handler.go  <-- [NEW] Handler endpoint LLM
│   │           └── rekap_handler.go       <-- [NEW] Handler endpoint Matrix Rekap
│   └── usecase/
│       ├── ai_grading_usecase.go          <-- [NEW] Logika prompting AI & Job Queue (Worker)
│       ├── rekap_usecase.go               <-- [NEW] Logika pivot/agresi data rekap
│       └── user_usecase.go                <-- [MODIFY] Tambahan logika parsing CSV
└── pkg/
    └── ollama/
        └── client.go                      <-- [NEW] Package khusus tembak API Ollama
```

### 🎨 Frontend (SvelteKit)
```text
frontend/
└── src/
    ├── lib/
    │   └── components/
    │       ├── EdraEditor.svelte          <-- [NEW] Komponen WYSIWYG untuk form soal
    │       └── TokenGate.svelte           <-- [NEW] Komponen UI gembok/PIN token ujian
    └── routes/
        └── praktikum/
            └── admin/
                ├── rekap-nilai/
                │   └── +page.svelte       <-- [NEW] Halaman Dashboard Pivot Nilai
                ├── users/
                │   └── +page.svelte       <-- [MODIFY] Penambahan UI Modal Import CSV
                ├── penilaian/
                │   └── +page.svelte       <-- [MODIFY] Penambahan UI Polling AI Grading Queue
                └── soal/
                    └── +page.svelte       <-- [MODIFY] Penambahan field Kunci Jawaban & EdraEditor
```

---

## 9. Catatan Penting Version Control (.gitignore)
Untuk menjaga kebersihan *repository* Git (agar developer lain tidak pusing melihat *diff* yang aneh-aneh), aturan berikut wajib ditambahkan ke dalam `.gitignore` di masing-masing direktori:

*   **Backend (`backend/.gitignore`)**:
    *   Wajib mengecualikan folder `docs/` hasil *auto-generate* dari library Swaggo (`docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`). File ini berubah terus menerus dan berpotensi bikin konflik Git (*Merge Conflict*).
    *   Wajib mengecualikan folder temporary tempat Backend menyimpan sementara file `.csv` yang diunggah asisten sebelum di-*parsing* (misal: `/tmp/uploads/`).
*   **Frontend (`frontend/.gitignore`)**:
    *   SvelteKit sudah punya bawaan `.gitignore` yang tangguh, tapi pastikan folder `.svelte-kit/` dan file `.env` selalu dikecualikan.
