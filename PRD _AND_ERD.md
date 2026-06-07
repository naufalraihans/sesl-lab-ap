# Laboratorium Algoritma dan Pemrograman

```
{{BaseUrl}} = https://lab-ap.vercel.app
```

## A. Web Dashboard (halaman informasi)

### 1.1. Halaman utama (Lobby dari website lab)

Fitur: Menampilkan ucapan selamat datang, deskripsi singkat laboratorium, pengumuman terbaru, dan quick links ke bagian penting.

`url: {{baseUrl}}/info`

### 1.2. sub-halaman web dashboard

Sub-halaman web dashboard terbagi menjadi 4 halaman, yaitu:

- Halaman Jadwal
- Halaman Daftar Asisten Lab
- Halaman Pedoman Laporan
- Halaman Modul Praktikum

#### Halaman jadwal

Fitur: Menampilkan informasi jadwal praktikum per shift/kelas dengan dua opsi tampilan yang diatur oleh Admin:

- **Tautan Google Drive**
  Menampilkan tombol atau link (href) yang mengarah langsung ke dokumen Google Drive eksternal.
- **Tabel Internal Web**
  Menampilkan tabel/kalender jadwal secara read-only. Data ini terhubung langsung (sinkron) dengan data Jadwal Kelas yang diatur oleh Admin. Sehingga, jadwal yang tampil di halaman publik ini adalah representasi global dari seluruh jadwal kelas yang sedang aktif.

`url: {{baseUrl}}/info/jadwal`

#### Halaman Daftar Asisten Lab

Fitur: Menampilkan profil asisten aktif (Foto, Nama, NIM, Nomor WhatsApp, Link LinkedIn/Media Sosial). (Catatan: Data profil asisten ini bersifat dinamis (tidak di-hardcode) dan ditarik dari database yang dikelola oleh Admin).

`url: {{baseUrl}}/info/asisten`

#### Halaman Pedoman Laporan

Fitur: Menampilkan tombol download untuk template atau pedoman laporan. Jumlah tombol/dokumen bersifat dinamis (bisa 1, 2, 3, atau lebih) bergantung pada data file yang diunggah dan diatur oleh Admin melalui sistem.

`url: {{baseUrl}}/info/laporan`

#### Halaman Modul Praktikum

Fitur: Menampilkan file modul untuk didownload, ini bersifat dari admin. bukan hardcode ke url download, tetapi dikendalikan oleh admin. Modul hanya ada satu file saja dalam format pdf.

`url: {{baseUrl}}/info/modul`

## B. Halaman Praktikum (harus Login)

Link utama untuk praktikum adalah:
`{{baseUrl}}/praktikum`

Halaman Login/Register:
`{{baseUrl}}/praktikum/login`

Logout:
`{{baseUrl}}/praktikum/logout`

### 2.1 Sub-Halaman Praktikum

Sub-halaman praktikum terbagi menjadi 2 jenis berdasarkan role:

#### 2.1.1 Admin

Note: Admin = Asisten Lab. Data profil asisten (foto, nomor HP, link medsos) disimpan di tabel `users` dengan `role = admin`.

- Dashboard Admin — `{{baseUrl}}/praktikum/admin`
- Manajemen Data User — `{{baseUrl}}/praktikum/admin/users`
- Manajemen Data Asisten Lab — `{{baseUrl}}/praktikum/admin/asisten`
- Manajemen Data Jadwal Praktikum — `{{baseUrl}}/praktikum/admin/jadwal`
- Manajemen Data Pedoman Laporan — `{{baseUrl}}/praktikum/admin/pedoman`
- Manajemen Data Modul Praktikum — `{{baseUrl}}/praktikum/admin/modul`
- Manajemen Sesi Praktikum (aktivasi sesi, pool soal, rekap jawaban) — `{{baseUrl}}/praktikum/admin/sesi`
- Penilaian Mahasiswa — `{{baseUrl}}/praktikum/admin/penilaian`

#### 2.1.2 User

a. **Dashboard User** — `{{baseUrl}}/praktikum/dashboard`

- Menampilkan profil (Nama, NIM, Kelas, Shift)
- Menampilkan nilai pre-test/post-test, keterampilan, serta nilai ujian praktik per sesi
- Jadwal praktikum (diambil dari `jadwal` where `kelas_id = user.kelas_id AND shift = user.shift` → hari, jam_mulai, jam_selesai)
- Notifikasi sesi yang sedang aktif untuk kelas dan shift user

f. **Daftar Sesi** — `{{baseUrl}}/praktikum/sesi`

- Menampilkan daftar semua sesi praktikum (Modul 1, Modul 2, ..., Ujian Praktik)
- Sesi yang diaktifkan untuk kelas + shift user akan ditandai sebagai "Aktif" dan bisa diakses
- Sesi yang belum diaktifkan ditampilkan dalam status terkunci

b. **Pre-test** — `{{baseUrl}}/praktikum/sesi/[sesiId]/pretest`

- Halaman pengerjaan pre-test (5 soal essay/isian)
- Terdapat timer yang dimulai saat mahasiswa membuka pre-test
- Durasi timer dinamis, diatur oleh admin per course
- Ditampilkan **sebelum** keterampilan
- **Auto-save**: Jawaban otomatis tersimpan secara berkala (jika user disconnect/keluar, jawaban tetap aman)
- **Auto-submit**: Jika timer habis **atau akses ditutup admin**, jawaban otomatis ter-submit

c. **Post-test** — `{{baseUrl}}/praktikum/sesi/[sesiId]/posttest`

- Halaman pengerjaan post-test (2 soal essay/isian + 1 soal coding dengan Monaco Editor)
- Terdapat timer yang dimulai saat mahasiswa membuka post-test
- Durasi timer dinamis, diatur oleh admin per course
- Ditampilkan **sesudah** keterampilan
- **Auto-save**: Jawaban otomatis tersimpan secara berkala
- **Auto-submit**: Jika timer habis **atau akses ditutup admin**, jawaban otomatis ter-submit

d. **Keterampilan** — `{{baseUrl}}/praktikum/sesi/[sesiId]/keterampilan`

- Halaman pengerjaan keterampilan (1 soal live coding dengan Monaco Editor)
- Terdapat timer yang dimulai saat mahasiswa membuka keterampilan
- Durasi timer dinamis, diatur oleh admin per course
- **Auto-save**: Jawaban otomatis tersimpan secara berkala
- **Auto-submit**: Jika timer habis **atau akses ditutup admin**, jawaban otomatis ter-submit

e. **Ujian Praktik** — `{{baseUrl}}/praktikum/sesi/[sesiId]/ujian`

- Halaman pengerjaan ujian praktik (**6 soal live coding** dengan Monaco Editor, 1 soal per kategori)
- Hanya muncul di sesi terakhir (setelah seluruh modul selesai)
- Soal ujian dikelompokkan ke dalam **6 kategori**:
  - **Modul 1** — 1 soal live coding
  - **Modul 2** — 1 soal live coding
  - **Modul 3** — 1 soal live coding
  - **Modul 4 dan 5** — 1 soal live coding
  - **Modul 6** — 1 soal live coding
  - **Flowchart** — 1 soal live coding. Mahasiswa menerjemahkan sebuah **flowchart (gambar)** ke dalam kode. Gambar flowchart diunggah admin dan disimpan di Supabase Storage.
- Setiap kategori diacak 1 soal dari **pool per kategori** (semua mahasiswa di sesi yang sama mendapat soal yang sama).
- Terdapat timer yang dimulai saat mahasiswa membuka ujian
- Durasi timer dinamis, diatur oleh admin per course
- **Auto-save**: Jawaban otomatis tersimpan secara berkala
- **Auto-submit**: Jika timer habis **atau akses ditutup admin**, jawaban otomatis ter-submit

### 2.2 Role

Pada halaman praktikum terbagi menjadi 2 role, yaitu:

1. Admin
2. User

### 2.3 Mekanisme Auth (first-time login)

Data dasar mahasiswa (`NIM`, `Nama`, `kelas`, `shift`) sudah dimasukkan ke dalam database sebelumnya oleh Admin, namun `password` dan data lainnya masih kosong. Alur `login/registrasi`-nya adalah sebagai berikut:

1. Pengguna memasukkan `NIM`.
2. Sistem mengecek status `is_registered` di database.
   - Jika `is_registered = true`: Pengguna diminta memasukkan `password` yang sudah pernah dibuat (Proses Login normal).
   - Jika `is_registered = false`:
     - Sistem mengecek apakah Admin sudah **membuka akses register** untuk kelas user tersebut (`is_register_open = true`).
     - Jika **akses dibuka**: Pengguna diizinkan untuk register dengan menginput `password`. Setelah selesai, `is_registered` diubah menjadi `true`.
     - Jika **akses belum dibuka**: Pengguna **tidak diizinkan register**. Tampilkan pesan bahwa akses belum dibuka.
3. Setelah login berhasil, user langsung masuk ke halaman praktikum.

Note: Mahasiswa tidak bisa login jika belum melakukan register.

NOTE : pastikan melarang SQL Injection pada bagian ini, serta menggunakan password hash dan keamanan yang tajam

### 2.4 Fungsionalitas Admin

Admin (= Asisten Lab) memiliki hak akses penuh untuk melakukan operasi CRUD pada data sistem.

**Catatan Keamanan (Security):**

- Seluruh endpoint sensitif yang baru dibuat (AI Grading, Import CSV, Rekap) dilindungi oleh middleware keamanan (Role-Based Access Control) di backend. Hanya User dengan otorisasi Admin yang berhak mengeksekusi endpoint tersebut untuk mencegah eksploitasi.
- Mengimplementasikan Soft Deletes pada tabel-tabel krusial untuk mencegah kehilangan data permanen akibat ketidaksengajaan.

#### 0. Dashboard Admin

Fitur: Menampilkan ringkasan statistik dan monitoring sistem:

- Jumlah user yang **sedang online** saat ini (real-time). Karena backend berjalan sebagai **server stateful** (bukan serverless), jumlah online dihitung dari **registry session in-memory** di server: setiap login/aktivitas terautentikasi mendaftarkan/menyegarkan entri, dan entri dihapus saat logout atau token kedaluwarsa. (Tidak dihitung dari `last_login_at`, karena itu hanya mencatat waktu login terakhir.)
- Jumlah mahasiswa per kelas dan shift
- Jumlah mahasiswa yang sudah register vs belum register
- Status sesi yang sedang aktif (kelas + shift mana yang sedang mengerjakan)
- Ringkasan progress pengerjaan per sesi (berapa yang sudah selesai, sedang mengerjakan, belum mulai)
- Quick actions: aktivasi sesi, buka/tutup register

#### 1. Manajemen data User (Users)

Fitur: Tambah, edit, hapus, dan lihat data mahasiswa (Nama, NIM, Kelas, Shift, Password). NIM digunakan sebagai username. Admin juga dapat melakukan reset password dan membuka/menutup akses register per kelas. Admin dapat mempercepat proses penambahan data mahasiswa baru melalui fitur Bulk Update & Import CSV Mahasiswa. Sistem dilengkapi penanganan error terpusat yang menampilkan daftar peringatan baris yang bermasalah pada file CSV.

#### 2. Manajemen Data Asisten Lab

Fitur: Tambah, edit, hapus, dan lihat data asisten (Foto, Nama, NIM, Nomor HP, Link Medsos). Karena Admin = Asisten, data ini dikelola langsung dari tabel `users` dengan `role = admin`. (Catatan: Data yang ditambahkan atau diperbarui di sini akan secara otomatis mengubah tampilan daftar asisten di Halaman Portal Publik `/info/asisten`).

#### 3. Manajemen jadwal dan sesi

Fitur: Mengelola konfigurasi dan data jadwal praktikum yang akan ditampilkan di web. Terdapat dua opsi yang bisa diatur oleh Admin:

- **Pengaturan Link GDrive**: Admin dapat memasukkan atau memperbarui link URL Google Drive yang mengarah ke file jadwal eksternal. Link ini disimpan di tabel `konfigurasi` di database.
- **Pengaturan Data Internal**: Membuat dan mengubah data jadwal praktikum secara mandiri di dalam sistem web dengan atribut: Hari, Jam/Shift, Kelas, dan Asisten Jaga. Admin juga dapat memetakan mahasiswa (NIM) ke dalam kelas/jadwal tertentu. (Catatan: Data internal ini terpisah dan tidak akan mengubah isi file yang ada di Google Drive).

#### 4. Manajemen Pedoman Laporan

Fitur: Admin dapat mengunggah (upload), mengubah nama, atau menghapus berbagai file template pedoman laporan. File yang diunggah di sini akan otomatis muncul sebagai tombol download di halaman /info/laporan.

#### 4b. Manajemen Modul Praktikum

Fitur: Admin dapat mengunggah atau mengganti **1 file PDF modul** praktikum. File modul ini bersifat global (bukan per sesi), dan akan ditampilkan di halaman `/info/modul` untuk didownload oleh mahasiswa. URL file disimpan di tabel `konfigurasi`.

#### 5. Manajemen Sesi Praktikum

Fitur: Admin dapat membuat sesi praktikum (misal: "Modul 1 - Pengenalan Dasar Bahasa C"), serta mengaktifkan/menonaktifkan sesi tersebut.

Setiap sesi praktikum (modul) memiliki **2 jenis course**:

1. **Pre-test atau Post-test** — Admin menentukan (gacha/spin) apakah sesi ini menggunakan pre-test atau post-test. Hanya salah satu per sesi.
2. **Keterampilan** — 1 soal live coding.

**Ujian Praktik** hanya terjadi di **sesi paling akhir** setelah seluruh modul selesai. Misalnya ada 6 modul, maka pertemuan ke-7 adalah ujian praktik. Ujian praktik = **6 soal live coding**, terdiri dari 1 soal per kategori (Modul 1, Modul 2, Modul 3, Modul 4 dan 5, Modul 6, dan Flowchart). Untuk kategori **Flowchart**, admin mengunggah **gambar flowchart** (disimpan di Supabase Storage) yang harus diterjemahkan mahasiswa ke dalam kode.

##### Detail Soal Per Course:

| Course        | Jumlah Soal Tampil      | Jenis Soal                               | Distribusi Difficulty                             |
| ------------- | ----------------------- | ---------------------------------------- | ------------------------------------------------- |
| Pre-test      | 5 soal                  | Essay/Isian                              | 1 easy, 2 medium, 2 hard                          |
| Post-test     | 3 soal                  | 2 essay/isian + 1 coding (Monaco Editor) | 1 easy (essay), 1 medium (essay), 1 hard (coding) |
| Keterampilan  | 1 soal                  | Live Coding (Monaco Editor)              | -                                                 |
| Ujian Praktik | 6 soal (1 per kategori) | Live Coding (Monaco Editor)              | -                                                 |

##### Detail Kategori Soal Ujian Praktik:

Course **Ujian Praktik** memiliki pool soal yang dikelompokkan per **kategori ujian**. Saat ujian diaktifkan, sistem mengacak **1 soal dari setiap kategori** sehingga total 6 soal:

| Kategori      | Jumlah Soal Tampil | Jenis Soal                  | Keterangan                                                                   |
| ------------- | ------------------ | --------------------------- | ---------------------------------------------------------------------------- |
| Modul 1       | 1 soal             | Live Coding (Monaco Editor) | -                                                                            |
| Modul 2       | 1 soal             | Live Coding (Monaco Editor) | -                                                                            |
| Modul 3       | 1 soal             | Live Coding (Monaco Editor) | -                                                                            |
| Modul 4 dan 5 | 1 soal             | Live Coding (Monaco Editor) | -                                                                            |
| Modul 6       | 1 soal             | Live Coding (Monaco Editor) | -                                                                            |
| Flowchart     | 1 soal             | Live Coding (Monaco Editor) | Soal menampilkan**gambar flowchart** untuk diterjemahkan ke dalam kode |

##### Mekanisme Pool & Pengacakan Soal:

- Setiap course memiliki **pool soal di database** (misal: 15 soal pre-test untuk Modul 1).
- Pool soal bersifat **per modul** (Modul 1 punya pool sendiri, Modul 2 punya pool sendiri).
- Khusus course **Ujian Praktik**, pool soal dikelompokkan **per kategori** (Modul 1, Modul 2, Modul 3, Modul 4 dan 5, Modul 6, Flowchart). Saat diaktifkan, sistem mengacak **1 soal per kategori** (total 6 soal).
- Saat sesi diaktifkan, soal **diacak dari pool** sesuai jumlah dan distribusi difficulty (atau kategori untuk ujian) yang ditentukan.
- Pengacakan dilakukan **per sesi** (semua mahasiswa dalam sesi yang sama mendapat soal yang sama).

> **Catatan implementasi:**
>
> - **Gacha butuh 2 pool siap**: agar admin bisa memilih (gacha) pre-test **atau** post-test saat aktivasi, setiap modul harus menyiapkan **pool pre-test DAN pool post-test** lebih dulu (plus pool keterampilan). Saat aktivasi, hanya pool dari course yang terpilih yang diacak ke `soal_terpilih`.
> - **Distribusi difficulty bersifat hardcoded** oleh jenis course (Pre-test = 1 easy/2 medium/2 hard; Post-test = 1 easy essay/1 medium essay/1 hard coding), bukan dikonfigurasi per course di database. Ini keputusan sadar untuk menyederhanakan; jika butuh fleksibel di masa depan, distribusi bisa dipindah jadi konfigurasi.

##### Aktivasi Sesi:

- Admin **manual mengaktifkan** sesi untuk **kelas + shift** tertentu (1 baris `aktivasi_sesi` per kombinasi sesi + kelas + shift). Admin wajib mengenerate Token/PIN Acak untuk aktivasi tersebut guna memastikan kehadiran fisik mahasiswa.
- **Akses terisolasi per kelas + shift**: aktivasi untuk "TTL A shift 1" hanya bisa diakses mahasiswa TTL A shift 1. Mahasiswa TTL A shift 2 **tidak bisa** mengaksesnya, walau kelasnya sama. Tiap kombinasi punya aktivasi sendiri.
- **Gacha pre-test/post-test bersifat eksplisit**: saat aktivasi, admin menentukan (gacha/spin) apakah pakai pre-test atau post-test, lalu sistem membuat baris `aktivasi_course` untuk course terpilih saja (+ keterampilan). Course yang tidak terpilih tidak akan punya baris dan tidak dapat diakses.
- Admin **membuka/menutup** tiap course **secara independen per aktivasi** lewat `aktivasi_course.is_open`. Karena open/close ada di level aktivasi (bukan global di `course`), membuka pre-test untuk "TTL A shift 1" **tidak** ikut membuka pre-test untuk shift/kelas lain. Urutan buka ditentukan admin (`urutan`).
- **Menutup course = auto-submit massal**: saat admin set `is_open=false` untuk sebuah course, semua mahasiswa di aktivasi itu yang **belum submit** akan otomatis ter-submit (jawaban = hasil auto-save terakhir), dan `pengerjaan_course.status` mereka menjadi `selesai`. Setelah ditutup, course tidak bisa lagi dikerjakan.
- Contoh: Admin mengaktifkan "Modul 1" untuk "TTL A shift 1", gacha → pre-test, lalu set `is_open=true` untuk course pre-test. Setelah selesai, tutup pre-test (`is_open=false`) dan buka keterampilan. Saat pre-test ditutup, mahasiswa yang belum sempat submit otomatis tersimpan & ter-submit.

###### Mekanisme Timer (Server-Authoritative):

- **Server adalah sumber kebenaran timer**, bukan client. `pengerjaan_course.waktu_mulai` diisi **sekali** saat mahasiswa pertama membuka course (refresh/buka ulang tidak me-reset timer).
- Deadline = `waktu_mulai + course.durasi_menit`. Setiap request save/submit divalidasi server terhadap deadline; jika sudah lewat, server menolak input baru dan menandai jawaban sebagai submitted.
- Karena backend stateful, sebuah **background job** dapat menyapu pengerjaan yang sudah lewat deadline untuk meng-auto-submit walau mahasiswa sudah keluar/disconnect. Ini mencegah manipulasi timer dari sisi browser.

###### Mahasiswa Susulan (Error Handling):

- Kasus: mahasiswa tidak bisa ikut di jadwal kelas/shift-nya sendiri (mis. mahasiswa **shift 1 kelas B** perlu menumpang ke aktivasi **shift 1 kelas A**). Secara default ini **ditolak** karena kelas/shift tidak cocok.
- Admin dapat memberi **akses susulan** dengan mendaftarkan mahasiswa tsb ke `peserta_susulan` pada aktivasi tujuan.
- **Aturan akses** mahasiswa terhadap sebuah course menjadi: course `is_open = true` **DAN** ( kelas+shift mahasiswa cocok dengan aktivasi **ATAU** mahasiswa terdaftar di `peserta_susulan` aktivasi tsb ).
- Jawaban & nilai mahasiswa susulan tetap tersimpan normal (lewat `soal_terpilih` aktivasi tujuan), sehingga admin tahu jawaban itu dari sesi susulan mana.

Untuk setiap course, Admin dapat:

- Membuat, mengedit, dan menghapus **soal** ke dalam pool (teks soal, jenis: essay / coding, difficulty: easy / medium / hard). Admin menggunakan WYSIWYG Rich Text Editor (Edra) untuk pembuatan soal yang mendukung penyisipan tabel, format teks, dan rumus matematika (KaTeX). Admin dianjurkan mengisi kolom `kunci_jawaban` sebagai acuan bagi AI Grading. Khusus course **Ujian Praktik**, setiap soal juga diberi **kategori ujian** (Modul 1 / Modul 2 / Modul 3 / Modul 4 dan 5 / Modul 6 / Flowchart), dan untuk kategori Flowchart admin mengunggah **gambar flowchart**.
- Mengatur **waktu pengerjaan** (durasi timer) per course.
- Membuka/menutup course secara independen **per aktivasi** (per kelas + shift) lewat `aktivasi_course.is_open`.
- Mendaftarkan **mahasiswa susulan** ke sebuah aktivasi (`peserta_susulan`) agar bisa menumpang mengerjakan di kelas/shift lain.
- Melihat dan merekap **jawaban mahasiswa** per course, terkelompokkan berdasarkan sesi, kelas, dan shift.

#### 6. Penilaian Mahasiswa & Rekapitulasi

Fitur: Admin dapat memberikan nilai dan feedback pada jawaban mahasiswa per course (pre-test/post-test, keterampilan, ujian praktik). Penilaian dilakukan secara manual oleh Admin untuk semua jenis soal (essay dan coding), atau menggunakan fitur AI Grading.

**Integrasi AI Grading (Penilaian Otomatis):**

- Admin dapat menggunakan tombol "Bulk AI Grade" untuk menilai puluhan hingga ratusan soal essay menggunakan bantuan Large Language Models (LLM) dengan membandingkan teks soal, kunci jawaban, dan jawaban mahasiswa.
- Proses grading menggunakan sistem Background Queue (antrean pekerja latar belakang) untuk mencegah HTTP Timeout. Tampilan antarmuka akan menerapkan sistem polling untuk mengecek status antrean secara real-time.

**Dashboard Rekap Jawaban & Pivot Nilai Akhir:**

- Admin dapat melihat pandangan menyeluruh (bird's-eye view) performa kelas melalui halaman Rekap Jawaban Global dan Rekap Nilai (Pivot).
- Terdapat kolom pencarian reaktif (NIM/Nama), tabel data dinamis, dan tombol "Export to Excel/CSV".
- Terdapat fitur Bulk Actions pada Rekap Jawaban untuk mereset nilai atau menghapus jawaban secara masal.

**Aturan skor (poin & nilai):**

- Setiap soal punya **bobot/poin** sendiri (skor maksimal), mis. Pre-test → soal 1 = 20, soal 2 = 15, soal 3 = 50, dst. Bobot tiap soal bisa berbeda.
- `nilai` yang diberikan admin untuk sebuah soal berada di rentang **0 sampai `poin`** soal tersebut (tidak boleh melebihi bobot).
- `total_nilai` course = **Σ nilai** seluruh soal di course itu. Disarankan Σ poin per course = 100 agar `total_nilai` langsung berskala 0–100 (tidak dipaksakan sistem, tergantung pembobotan admin).

### 2.5 Fungsionalitas User

#### 1. Dashboard User

Fitur: Menampilkan ringkasan profil (Nama, NIM, Kelas, Shift) dan notifikasi status sesi modul yang sedang aktif untuk kelas dan shift user.

#### 2. Lihat Jadwal Pribadi

Fitur: Menampilkan jadwal spesifik praktikum mahasiswa bersangkutan.

- Mekanisme Otomatis: Sistem mendeteksi `kelas_id` dan `shift` milik User, lalu mengambil baris `jadwal` yang cocok dengan **kombinasi kelas + shift** tersebut (Hari, Jam, dan keterangan periode).
- Catatan skema shift: dalam satu kelas, **shift 1 dan shift 2 berjalan pada periode berbeda** (mis. shift 1 = minggu 1–4, shift 2 = minggu 5–8). Karena itu jadwal disimpan terpisah per kelas + shift, bukan per kelas saja.

#### 3. Akses Sesi Praktikum

Fitur: Menampilkan daftar sesi praktikum. Sesi yang diaktifkan Admin untuk kelas dan shift user akan terbuka dan bisa diakses. Mahasiswa wajib memasukkan Token/PIN Ujian yang dibagikan di dalam ruangan lab untuk dapat memulai sesi, guna memastikan kehadiran fisik. Di dalam setiap sesi, mahasiswa melihat course yang tersedia beserta statusnya (belum dikerjakan, sedang dikerjakan, selesai).

Note: Mahasiswa hanya bisa melihat dan mengakses sesi yang diaktifkan untuk **kelas + shift** mereka. Pengecualian: mahasiswa yang didaftarkan **susulan** oleh admin (`peserta_susulan`) dapat mengakses aktivasi kelas/shift lain yang ditunjuk.

#### 4. Ruang Sesi & Pengerjaan Soal (Detail Modul)

Fitur Utama: Saat mahasiswa mengakses sesi yang aktif, mereka dapat mengerjakan soal dari course yang tersedia:

1. **Pre-test** — 5 soal essay/isian (ditampilkan sebelum keterampilan). Mahasiswa menjawab langsung di web.
2. **Post-test** — 2 soal essay/isian + 1 soal coding dengan **Monaco Editor** (ditampilkan sesudah keterampilan).
3. **Keterampilan** — 1 soal live coding dengan **Monaco Editor**.
4. **Ujian Praktik** — 6 soal live coding dengan **Monaco Editor** (hanya di sesi terakhir), terdiri dari 1 soal per kategori: Modul 1, Modul 2, Modul 3, Modul 4 dan 5, Modul 6, dan Flowchart (menerjemahkan gambar flowchart ke kode).

Mekanisme Pengerjaan:

- Soal yang tampil sudah **diacak dari pool** oleh sistem saat sesi diaktifkan (semua mahasiswa di sesi yang sama mendapat soal yang sama).
- Urutan course (pre-test → keterampilan → post-test) **dibuka secara manual oleh admin**, bukan otomatis oleh sistem.
- Setiap course memiliki **timer** yang berjalan saat mahasiswa mulai mengerjakan.
- **Auto-save**: Jawaban mahasiswa **otomatis tersimpan secara berkala** selama pengerjaan. Jika user disconnect atau keluar dari halaman, jawaban yang sudah diketik tetap tersimpan dan bisa dilanjutkan.
- **Auto-submit** terjadi pada **2 kondisi**, dan jawaban yang ter-submit adalah hasil auto-save terakhir:
  1. **Timer habis** — server menghitung deadline dari `pengerjaan_course.waktu_mulai + course.durasi_menit`. Lewat deadline → jawaban dianggap submitted.
  2. **Akses course ditutup admin** (`aktivasi_course.is_open` → `false`) sementara mahasiswa belum menekan submit → seluruh jawaban yang sedang dikerjakan untuk course itu **otomatis di-submit**.
- Mahasiswa juga dapat **submit manual** sebelum timer habis / sebelum akses ditutup.
- Mahasiswa hanya dapat mengerjakan setiap course **satu kali** (tidak bisa mengulang setelah submit).
- Terdapat indikator status pengerjaan per course (Belum Dikerjakan, Sedang Dikerjakan, Selesai).

## Tech Stack :

### Frontend

- Framework: Svelte (dengan SvelteKit)
- Bahasa: TypeScript
- Styling: CSS Native / Tailwind CSS
- Hosting: Vercel

### Backend

1. Bahasa: Go
2. ORM: GORM (dengan implementasi Soft Delete)
3. Database: PostgreSQL
4. Authentication: JWT
5. API Documentation: OpenAPI (Swagger dengan swaggo/swag)
6. Build Tool: Makefile (`make run`, `make migrate-up`, `make migrate-down`, `make seed`)
7. Caching: In-Memory Caching (`sync.Map`) untuk optimasi endpoint konfigurasi statis.

### Database

- Database Relasional: PostgreSQL (dengan penerapan Composite Index untuk mencegah Slow SQL pada Auto-Submit Worker).
- Object Storage (Untuk PDF dan Foto): Supabase Storage.
- File pedoman laporan, PDF modul, dan foto profil asisten lab akan disimpan di bucket Supabase Storage (misal: bucket public-assets) dan URL publiknya disimpan di dalam database PostgreSQL.
- Catatan: Karena praktikum berbasis soal-jawaban di web (bukan upload file), kebutuhan storage utama adalah untuk file statis (modul PDF, pedoman, foto asisten).

## ERD :

Kode dibawah ini adalah struktur untuk ERD-nya dengan menggunakan gaya penulisan dbdiagram.com (bukan kode SQL):

```dbml
Enum role_type {
  user
  admin
}

Enum jenis_course {
  pretest
  posttest
  keterampilan
  ujian_praktik
}

Enum jenis_soal {
  essay
  coding
}

Enum difficulty {
  easy
  medium
  hard
}

// Kategori soal khusus untuk course Ujian Praktik
Enum kategori_ujian {
  modul_1
  modul_2
  modul_3
  modul_4_5
  modul_6
  flowchart
}

Enum status_pengerjaan {
  belum_dikerjakan
  sedang_dikerjakan
  selesai
}

Table users {
  id int [pk, increment]
  role role_type [note: 'admin = asisten lab']
  nim varchar [unique, note: 'Digunakan sebagai username']
  nama varchar
  password_hash varchar [note: 'Null sebelum mahasiswa melakukan set password pertama kali']
  is_registered boolean [default: false, note: 'Penanda untuk flow first-time login']
  kelas_id int [note: 'Null jika role admin']
  shift int [note: '1 atau 2. Null jika role admin']
  foto_url varchar [note: 'Null jika role user. Disimpan di Supabase Storage']
  nomor_hp varchar [note: 'Null jika role user']
  medsos_link varchar [note: 'Null jika role user']
  last_login_at datetime [note: 'Waktu terakhir login (riwayat). Status online real-time TIDAK dihitung dari sini, tapi dari registry in-memory di server (lihat Dashboard Admin)']
  created_at datetime
  updated_at datetime
}

Table kelas {
  id int [pk, increment]
  nama_kelas varchar [note: 'misal: TTL A, TTL B']
  is_register_open boolean [default: false, note: 'Admin bisa buka/tutup akses register per kelas']
}

// Jadwal bersifat per kelas + shift, karena tiap shift punya jadwal/periode berbeda.
// Contoh: TTL A shift 1 (minggu 1-4) & TTL A shift 2 (minggu 5-8).
Table jadwal {
  id int [pk, increment]
  kelas_id int
  shift int [note: '1 atau 2']
  hari varchar
  jam_mulai time
  jam_selesai time
  keterangan varchar [note: 'Opsional, mis. "Minggu 1-4" untuk menandai periode shift']

  Indexes {
    (kelas_id, shift) [unique, note: 'Satu jadwal per kelas per shift']
  }
}

// Tabel untuk file pedoman yang dinamis
Table pedoman_laporan {
  id int [pk, increment]
  nama_dokumen varchar [note: 'misal: Template Laporan Akhir']
  file_url varchar [note: 'Link ke file dokumen di Supabase Storage']
  diunggah_pada datetime
}

// Tabel konfigurasi sistem (key-value) untuk menyimpan setting global
Table konfigurasi {
  id int [pk, increment]
  key varchar [unique, note: 'misal: gdrive_jadwal_url, modul_file_url']
  value text [note: 'Nilai konfigurasi (URL, teks, dll)']
  updated_at datetime
}

Table sesi_praktikum {
  id int [pk, increment]
  judul_sesi varchar [note: 'misal: Modul 1 - Pengenalan Dasar Bahasa C']
  deskripsi text
  urutan int [note: 'Urutan tampil/sekuens sesi (Modul 1, 2, 3, ... lalu ujian). Tidak bergantung pada id']
  is_ujian_praktik boolean [default: false, note: 'True jika sesi ini adalah ujian praktik (sesi terakhir)']
  created_at datetime
  updated_at datetime
}

// Setiap sesi memiliki course (pre/post-test, keterampilan, atau ujian_praktik).
// Sesi normal biasanya punya 3 course: pretest + posttest + keterampilan.
// Pool pretest DAN posttest sama-sama disiapkan agar admin bisa gacha salah satu saat aktivasi.
Table course {
  id int [pk, increment]
  sesi_praktikum_id int
  jenis jenis_course [note: 'Dalam satu sesi, pretest dan posttest tidak boleh ada bersamaan']
  judul varchar
  deskripsi text
  durasi_menit int [note: 'Durasi timer pengerjaan dalam menit']
  // Catatan: status buka/tutup TIDAK di sini, melainkan per kelas+shift di tabel aktivasi_course

  Indexes {
    (sesi_praktikum_id, jenis) [unique, note: 'Maksimal 1 course per jenis per sesi (tidak boleh pretest dobel, dst)']
  }
}

// Pool soal milik sebuah course (per modul)
Table soal {
  id int [pk, increment]
  course_id int
  jenis_soal jenis_soal [note: 'essay atau coding']
  difficulty difficulty [note: 'easy, medium, hard. Null untuk keterampilan/ujian_praktik']
  kategori_ujian kategori_ujian [note: 'Hanya diisi untuk course ujian_praktik (modul_1..modul_6, flowchart). Null untuk course lain']
  teks_soal text
  gambar_url varchar [note: 'URL gambar flowchart di Supabase Storage. Diisi untuk soal kategori flowchart. Null untuk soal lain']
  poin float [note: 'Bobot / skor MAKSIMAL soal ini (mis. 20, 15, 50). Nilai yang diberikan admin tidak boleh melebihi poin']
  kunci_jawaban text [note: 'Referensi jawaban (opsional, untuk panduan penilaian admin)']
  created_at datetime
}

// Admin mengaktifkan sesi untuk kelas + shift tertentu
Table aktivasi_sesi {
  id int [pk, increment]
  sesi_praktikum_id int
  kelas_id int
  shift int [note: '1 atau 2']
  token varchar [note: 'PIN/Token akses ujian untuk memastikan kehadiran fisik']
  is_active boolean [default: true]
  activated_at datetime

  Indexes {
    (sesi_praktikum_id, kelas_id, shift) [unique, note: 'Satu aktivasi per sesi per kelas per shift']
  }
}

// Status buka/tutup tiap course dalam sebuah aktivasi (per kelas + shift).
// Admin membuka/menutup course satu per satu; urutan ditentukan admin.
// Hanya course yang punya baris di sini yang DIPAKAI pada aktivasi tsb.
// Untuk pretest vs posttest: hasil gacha = course mana yang dibuatkan barisnya.
Table aktivasi_course {
  id int [pk, increment]
  aktivasi_sesi_id int
  course_id int
  is_open boolean [default: false, note: 'True = course terbuka & bisa dikerjakan untuk kelas+shift aktivasi ini']
  urutan int [note: 'Urutan buka course (pretest -> keterampilan -> dst), diatur admin']
  opened_at datetime
  closed_at datetime

  Indexes {
    (aktivasi_sesi_id, course_id) [unique, note: 'Satu status per course per aktivasi']
  }
}

// Mahasiswa SUSULAN: diberi akses ke aktivasi yang BUKAN kelas/shift aslinya.
// Contoh: mahasiswa shift 1 kelas B ikut mengerjakan di aktivasi shift 1 kelas A.
Table peserta_susulan {
  id int [pk, increment]
  aktivasi_sesi_id int
  mahasiswa_id int [note: 'User (role user) yang diizinkan ikut susulan']
  alasan varchar [note: 'Alasan susulan, diisi admin (opsional)']
  created_at datetime

  Indexes {
    (aktivasi_sesi_id, mahasiswa_id) [unique, note: 'Satu mahasiswa hanya bisa didaftarkan susulan sekali per aktivasi']
  }
}

// Soal yang terpilih (hasil acak) untuk setiap aktivasi sesi
Table soal_terpilih {
  id int [pk, increment]
  aktivasi_sesi_id int
  course_id int
  soal_id int
  urutan int [note: 'Urutan tampil soal']

  Indexes {
    (aktivasi_sesi_id, course_id, soal_id) [unique, note: 'Soal tidak boleh duplikat dalam satu aktivasi+course']
  }
}

// Jawaban mahasiswa per soal (auto-save berkala; auto-submit saat timer habis / akses ditutup)
Table jawaban_mahasiswa {
  id int [pk, increment]
  mahasiswa_id int [note: 'Mengacu ke user dengan role user']
  soal_terpilih_id int [note: 'Mengacu ke soal yang sudah terpilih untuk sesi ini']
  jawaban_teks text [note: 'Untuk essay: isi jawaban. Untuk coding: source code. Auto-saved secara berkala']
  is_submitted boolean [default: false, note: 'True jika sudah di-submit: manual, auto-submit saat timer habis, atau auto-submit saat akses course ditutup admin']
  nilai float [note: 'Skor yang diberikan admin untuk soal ini, rentang 0..poin (boleh desimal, mis. 17.5)']
  feedback text [note: 'Feedback dari admin']
  waktu_submit datetime
  updated_at datetime [note: 'Terakhir auto-save']

  Indexes {
    (mahasiswa_id, soal_terpilih_id) [unique, note: 'Memastikan mahasiswa hanya bisa menjawab 1x per soal']
  }
}

// Tracking status pengerjaan mahasiswa per course per aktivasi
Table pengerjaan_course {
  id int [pk, increment]
  mahasiswa_id int
  aktivasi_sesi_id int
  course_id int
  status status_pengerjaan [default: 'belum_dikerjakan']
  waktu_mulai datetime [note: 'Waktu mulai mengerjakan, untuk hitung timer']
  waktu_selesai datetime
  total_nilai float [note: 'CACHED/DERIVED: akumulasi SUM(jawaban_mahasiswa.nilai) untuk course ini. WAJIB di-recalculate tiap admin update nilai per soal. Sumber kebenaran tetap jawaban_mahasiswa.nilai']

  Indexes {
    (mahasiswa_id, aktivasi_sesi_id, course_id) [unique, note: 'Satu mahasiswa satu record per course per aktivasi']
    (status, waktu_mulai) [note: 'Composite Index untuk mempercepat query polling background worker']
  }
}

// Relasi (Relationships)
Ref: users.kelas_id > kelas.id
Ref: jadwal.kelas_id > kelas.id
Ref: course.sesi_praktikum_id > sesi_praktikum.id
Ref: soal.course_id > course.id
Ref: aktivasi_sesi.sesi_praktikum_id > sesi_praktikum.id
Ref: aktivasi_sesi.kelas_id > kelas.id
Ref: aktivasi_course.aktivasi_sesi_id > aktivasi_sesi.id
Ref: aktivasi_course.course_id > course.id
Ref: peserta_susulan.aktivasi_sesi_id > aktivasi_sesi.id
Ref: peserta_susulan.mahasiswa_id > users.id
Ref: soal_terpilih.aktivasi_sesi_id > aktivasi_sesi.id
Ref: soal_terpilih.course_id > course.id
Ref: soal_terpilih.soal_id > soal.id
Ref: jawaban_mahasiswa.mahasiswa_id > users.id
Ref: jawaban_mahasiswa.soal_terpilih_id > soal_terpilih.id
Ref: pengerjaan_course.mahasiswa_id > users.id
Ref: pengerjaan_course.aktivasi_sesi_id > aktivasi_sesi.id
Ref: pengerjaan_course.course_id > course.id
```

## UI/UX

Untuk memastikan website memiliki tampilan yang profesional, akademis, dan konsisten, berikut adalah Panduan Penggunaan Warna (Color Usage Guideline) dan gaya komponen yang akan diterapkan di frontend (khususnya untuk konfigurasi Tailwind CSS).

### Panduan Penggunaan Warna (Color Usage Guideline)

#### 1. Primary Color: `#D03153` (Warna utama sistem — Rose/Crimson)

Digunakan untuk: Logo, Ikon utama, Tombol utama (Primary Button), Link aktif, Badge utama, Border fokus (focus state).

#### 2. Navbar: `#B02A47` (Sedikit lebih gelap dari warna utama)

- Alasan: Memberikan kontras yang jelas dengan area konten, terlihat lebih profesional, dan membuat menu navigasi lebih menonjol.
- Warna teks navbar: `#FFFFFF`
- Hover menu: `#E03A60`

#### 3. Sidebar (Jika Ada): `#942240` (Lebih gelap, untuk kontras dengan konten)

- Menu aktif: `#D03153`
- Menu hover: `#E03A60`

#### 4. Background Utama:

- `#FFFFFF` (Warna dasar — bersih dan profesional)

#### 5. Secondary Background:

- `#FDF2F4` atau `#F9FAFB`

Digunakan pada: Card area, Area filter, Statistik dashboard, Form section.

#### 6. Card:

- Background: `#FFFFFF`
- Border: `#E5E7EB`
- Hover: `#FDF2F4` (Gunakan shadow ringan agar tidak terlihat datar)

#### 7. Teks (Typography): (Jangan menggunakan hitam murni `#000000`)

- Heading: `#1F2937`
- Body Text: `#374151`
- Caption: `#6B7280`

#### 8. Primary Button:

- Normal: `#D03153`
- Hover: `#E03A60`
- Active: `#B02A47`
- Text: `#FFFFFF`

#### 9. Tabel:

- Header: `#B02A47` (Teks Header: `#FFFFFF`)
- Baris ganjil: `#FFFFFF`
- Baris genap: `#FDF2F4`
- Hover baris: `#FADADD`

#### 10. Link / Tautan:

- Normal: `#D03153`
- Hover: `#E03A60`
- Visited: `#942240`

#### 11. Disabled State (Komponen Non-aktif):

- Background: `#D1D5DB`
- Text: `#9CA3AF`
- Cursor: `not-allowed`

#### 12. Warna Status (Feedback States)

- **Success State** (Data tersimpan, Submit berhasil):
  - Teks/Ikon: `#10B981` | Background: `#D1FAE5`
- **Warning State** (Timer mendekat, Data belum lengkap):
  - Teks/Ikon: `#F59E0B` | Background: `#FEF3C7`
- **Error State** (Login gagal, Validasi gagal):
  - Teks/Ikon: `#EF4444` | Background: `#FEE2E2`
- **Info State** (Pengumuman, Notifikasi sistem):
  - Teks/Ikon: `#3B82F6` | Background: `#DBEAFE`

#### 13. Penggunaan Ikon (Iconography)

Sistem ini diwajibkan menggunakan ikon berbasis SVG murni (menggunakan pustaka `lucide-svelte`) untuk menjaga konsistensi resolusi tinggi, kebersihan desain, dan kesan antarmuka yang dinamis/premium.
**ATURAN MUTLAK:** Dilarang keras menggunakan karakter emoji dalam bentuk apa pun pada seluruh komponen UI/UX, penamaan, maupun notifikasi, guna mempertahankan standar estetika yang tinggi dan profesional.

B. Arahan Komponen & Gaya Dasar (Styling)

- Bentuk (Shape): Gunakan komponen dengan sudut sedikit membulat (rounded-lg atau rounded-xl pada Tailwind). Hindari desain kotak tajam murni agar gaya modern pendidikan terasa lebih fleksibel dan bersahabat.
- Glassmorphism & Animasi: Penggunaan efek kaca (backdrop-blur) pada navigasi dan komponen melayang dengan animasi transisi yang mulus sangat direkomendasikan.
- Ruang (Whitespace): Berikan jarak (padding/margin) yang lega antar elemen. Jangan menumpuk teks dan tombol terlalu rapat, terutama di halaman detail modul.
- Responsivitas (Mobile-First): Mahasiswa sering mengakses jadwal dan men-download modul lewat smartphone. Pastikan navbar bisa berubah menjadi hamburger menu, tabel jadwal/nilai bisa di-scroll secara horizontal di layar kecil, dan form pengumpulan tugas tetap nyaman digunakan di HP.
- Caching Browser: Halaman publik seperti jadwal akan di-cache selama 5 menit menggunakan pengaturan header SvelteKit untuk meminimalisasi beban request server.

## Struktur Folder

Backend project ini memakai **Clean Architecture** (Go) agar tanggung jawab tiap layer terpisah dan logika bisnis (usecase) bisa diuji tanpa HTTP/DB.

**Aliran dependency (satu arah):**

```
delivery (handler) ──> usecase ──> repository (interface) ──> DB (GORM / PostgreSQL)
        │                  │
        └─ dto             └─ entity (domain model, dipakai semua layer)
```

**Peran tiap layer:**

| Layer                | Tanggung jawab                                                                            | Boleh tahu                        |
| -------------------- | ----------------------------------------------------------------------------------------- | --------------------------------- |
| **entity**     | Domain model murni (struct + tag GORM). Inti sistem.                                      | Tidak depend ke siapa pun         |
| **dto**        | Bentuk request & response API (validasi input, shaping output query).                     | entity                            |
| **repository** | **Satu-satunya** layer yang menyentuh GORM/PostgreSQL. Berisi interface + implementasi.  | entity                            |
| **usecase**    | Seluruh aturan bisnis (gacha, acak soal, timer, auto-submit, recalc nilai, cek akses).    | interface repository, entity, dto |
| **delivery**   | Layer HTTP: handler (parse request → panggil usecase → balas dto), middleware, routing. | usecase, dto                      |

Aturan kunci: **usecase tidak tahu HTTP maupun GORM** — dia hanya memanggil _interface_ repository, sehingga bisa di-unit-test dengan mock. **Hanya repository yang menyentuh database.**

```
/project_lab_ap
│
├── Makefile                            # make run, make migrate-up, make migrate-down, make seed
├── README.md
│
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                 # Entry point: wiring config -> db -> repository -> usecase -> handler -> router
│   │
│   ├── config/
│   │   └── config.go                   # Load env & konfigurasi (DB, JWT secret, Supabase, dll)
│   │
│   ├── internal/
│   │   │
│   │   ├── entity/                      # === DOMAIN MODEL (struct + tag GORM). Inti, tanpa dependency ===
│   │   │   ├── user.go                  # User (role: admin=asisten, user=mahasiswa)
│   │   │   ├── kelas.go
│   │   │   ├── jadwal.go                # Jadwal per kelas + shift
│   │   │   ├── konfigurasi.go           # Key-value config (GDrive link, modul URL)
│   │   │   ├── pedoman_laporan.go
│   │   │   ├── sesi_praktikum.go
│   │   │   ├── course.go                # jenis: pretest/posttest/keterampilan/ujian_praktik
│   │   │   ├── soal.go                  # Pool soal (essay/coding, difficulty, kategori_ujian, gambar flowchart)
│   │   │   ├── aktivasi_sesi.go         # Aktivasi sesi per kelas + shift
│   │   │   ├── aktivasi_course.go       # Status buka/tutup course per aktivasi
│   │   │   ├── peserta_susulan.go       # Mahasiswa susulan per aktivasi
│   │   │   ├── soal_terpilih.go         # Soal hasil acak per aktivasi
│   │   │   ├── jawaban_mahasiswa.go     # Jawaban mahasiswa per soal (auto-save)
│   │   │   └── pengerjaan_course.go     # Status & tracking pengerjaan per course
│   │   │
│   │   ├── dto/                         # === REQUEST & RESPONSE (validasi input, shaping output query) ===
│   │   │   ├── auth_dto.go              # LoginRequest, RegisterRequest, AuthResponse
│   │   │   ├── dashboard_dto.go         # StatistikResponse, OnlineCountResponse
│   │   │   ├── sesi_dto.go              # SesiRequest, SesiResponse, CourseResponse
│   │   │   ├── aktivasi_dto.go          # AktivasiRequest, GachaRequest, BukaTutupCourseRequest, SusulanRequest
│   │   │   ├── soal_dto.go              # SoalRequest (jenis/difficulty/kategori_ujian), SoalResponse
│   │   │   ├── jawaban_dto.go           # AutoSaveRequest, SubmitRequest, JawabanResponse
│   │   │   ├── penilaian_dto.go         # NilaiRequest (per soal, 0..poin), RekapResponse
│   │   │   ├── profile_dto.go           # UpdateProfileRequest (foto, medsos)
│   │   │   └── konfigurasi_dto.go       # KonfigurasiRequest (gdrive/modul URL)
│   │   │
│   │   ├── repository/                  # === AKSES DB (GORM). Interface + implementasi. SATU-SATUNYA penyentuh DB ===
│   │   │   ├── user_repository.go       # interface UserRepository + impl GORM
│   │   │   ├── kelas_repository.go
│   │   │   ├── jadwal_repository.go
│   │   │   ├── konfigurasi_repository.go
│   │   │   ├── pedoman_repository.go
│   │   │   ├── sesi_repository.go
│   │   │   ├── course_repository.go
│   │   │   ├── soal_repository.go
│   │   │   ├── aktivasi_repository.go   # aktivasi_sesi + aktivasi_course + peserta_susulan
│   │   │   ├── soal_terpilih_repository.go
│   │   │   ├── jawaban_repository.go
│   │   │   └── pengerjaan_repository.go
│   │   │
│   │   ├── usecase/                     # === BUSINESS LOGIC. Depend ke interface repo + entity + dto. Tanpa HTTP/GORM ===
│   │   │   ├── auth_usecase.go          # login, first-time register, cek is_register_open, hash password
│   │   │   ├── dashboard_usecase.go     # agregasi statistik + hitung online (pakai pkg/online)
│   │   │   ├── sesi_usecase.go          # CRUD sesi & course
│   │   │   ├── aktivasi_usecase.go      # aktivasi per kelas+shift, GACHA pretest/posttest, buka/tutup course (auto-submit massal), susulan
│   │   │   ├── soal_usecase.go          # CRUD pool + ACAK soal -> soal_terpilih (distribusi difficulty / 1 per kategori ujian)
│   │   │   ├── jawaban_usecase.go       # auto-save, submit manual, AUTO-SUBMIT (timer habis / akses ditutup), cek akses (kelas+shift / susulan)
│   │   │   ├── penilaian_usecase.go     # set nilai per soal (0..poin) + RECALC total_nilai
│   │   │   ├── konfigurasi_usecase.go   # set/get GDrive link & modul URL
│   │   │   └── profile_usecase.go       # update foto (Supabase) & medsos asisten
│   │   │
│   │   └── delivery/
│   │       └── http/
│   │           ├── handler/             # Parse HTTP request -> panggil usecase -> balas dto. TANPA logika bisnis
│   │           │   ├── auth_handler.go
│   │           │   ├── dashboard_handler.go
│   │           │   ├── sesi_handler.go
│   │           │   ├── aktivasi_handler.go
│   │           │   ├── soal_handler.go
│   │           │   ├── jawaban_handler.go
│   │           │   ├── penilaian_handler.go
│   │           │   ├── konfigurasi_handler.go
│   │           │   ├── modul_handler.go
│   │           │   ├── profile_handler.go
│   │           │   └── jadwal_handler.go
│   │           │
│   │           ├── middleware/
│   │           │   ├── auth_middleware.go   # Verifikasi JWT
│   │           │   └── role_middleware.go   # Cek role (user/admin)
│   │           │
│   │           └── route/
│   │               └── router.go            # Registrasi semua route + pasang middleware
│   │
│   ├── pkg/                             # === Helper reusable lintas-layer (tanpa state domain) ===
│   │   ├── jwt/jwt.go                   # generate & verify JWT
│   │   ├── hash/hash.go                 # hash & verify password (bcrypt/argon2)
│   │   ├── supabase/supabase.go         # upload/download Supabase Storage (foto, modul PDF, gambar flowchart)
│   │   ├── online/registry.go           # Registry session in-memory (hitung user online, server stateful)
│   │   └── response/response.go         # Format response JSON konsisten (success/error)
│   │
│   ├── database/
│   │   ├── connection.go                # Koneksi GORM ke PostgreSQL
│   │   ├── migration/
│   │   │   ├── 001_create_users.sql
│   │   │   ├── 002_create_kelas.sql
│   │   │   ├── 003_create_jadwal.sql
│   │   │   ├── 004_create_konfigurasi.sql
│   │   │   ├── 005_create_pedoman_laporan.sql
│   │   │   ├── 006_create_sesi_praktikum.sql
│   │   │   ├── 007_create_course.sql
│   │   │   ├── 008_create_soal.sql
│   │   │   ├── 009_create_aktivasi_sesi.sql
│   │   │   ├── 010_create_aktivasi_course.sql
│   │   │   ├── 011_create_peserta_susulan.sql
│   │   │   ├── 012_create_soal_terpilih.sql
│   │   │   ├── 013_create_jawaban_mahasiswa.sql
│   │   │   └── 014_create_pengerjaan_course.sql
│   │   └── seed/
│   │       └── seed.go                  # Data seed local dev (admin, kelas, jadwal, soal contoh)
│   │
│   └── docs/
│       └── openapi.yaml                # OpenAPI / Swagger spec
│
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api.ts                  # Axios/Fetch config
│   │   │   ├── utils.ts
│   │   │   └── supabase.ts             # Akses storage langsung dari FE
│   │   │
│   │   ├── stores/
│   │   │   ├── auth.ts                 # Auth store
│   │   │   └── praktikum.ts            # State sesi, course, soal
│   │   │
│   │   ├── routes/
│   │   │   ├── info/                                # Halaman publik (tanpa login)
│   │   │   │   ├── +page.svelte                     # Lobby / landing page
│   │   │   │   ├── jadwal/+page.svelte
│   │   │   │   ├── asisten/+page.svelte
│   │   │   │   ├── laporan/+page.svelte
│   │   │   │   └── modul/+page.svelte
│   │   │   │
│   │   │   ├── praktikum/                           # Halaman praktikum (login required)
│   │   │   │   ├── login/+page.svelte
│   │   │   │   ├── dashboard/+page.svelte            # Dashboard user
│   │   │   │   ├── sesi/+page.svelte                 # Daftar sesi praktikum
│   │   │   │   ├── sesi/[sesiId]/+page.svelte        # Detail sesi (daftar course)
│   │   │   │   ├── sesi/[sesiId]/pretest/+page.svelte
│   │   │   │   ├── sesi/[sesiId]/posttest/+page.svelte
│   │   │   │   ├── sesi/[sesiId]/keterampilan/+page.svelte
│   │   │   │   ├── sesi/[sesiId]/ujian/+page.svelte
│   │   │   │   └── admin/                            # Halaman admin
│   │   │   │       ├── +page.svelte                   # Dashboard admin
│   │   │   │       ├── users/+page.svelte
│   │   │   │       ├── asisten/+page.svelte
│   │   │   │       ├── jadwal/+page.svelte
│   │   │   │       ├── pedoman/+page.svelte
│   │   │   │       ├── modul/+page.svelte
│   │   │   │       ├── sesi/+page.svelte
│   │   │   │       └── penilaian/+page.svelte
│   │   │   │
│   │   │   └── +page.svelte                          # Root redirect
│   │   │
│   │   └── app.html
│   │
│   ├── static/
│   │   └── favicon.png
│   │
│   ├── tailwind.config.ts
│   ├── svelte.config.js
│   └── package.json
│
└── .env.example                        # Template environment variables
```
