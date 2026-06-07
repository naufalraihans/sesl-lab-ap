# Konsep & Cara Kerja Automated Testing

Automated Testing (Pengujian Otomatis) adalah sebuah teknik di mana kita menulis "kode/skrip khusus" yang bertugas untuk menguji "kode utama" aplikasi kita. Jika diibaratkan, skrip pengujian ini adalah seorang *Quality Assurance* (QA) robot yang mengecek setiap tombol dan fungsi di web kita berulang kali dalam hitungan detik.

Tujuannya adalah: **Mencegah Regresi**. Regresi terjadi ketika kita menambahkan fitur baru, tetapi fitur tersebut tanpa sengaja merusak fitur lama yang sudah jalan. Dengan Automated Testing, robot akan langsung berteriak (memberikan pesan *Error/Failed*) sebelum kode tersebut sampai ke tahap Production.

Berikut adalah penjelasan cara kerjanya di ekosistem web Lab-AP v3 milikmu:

---

## 1. Backend Testing (Golang)

Di Golang, alat pengujian sudah tertanam (*built-in*) langsung di dalam bahasanya, yaitu melalui *package* `testing`.

### A. Unit Testing
Ini adalah pengujian di level paling kecil. Kita mengetes satu fungsi atau satu metode saja.
Karena arsitektur kita adalah **Clean Architecture**, kita sangat mudah melakukan Unit Testing pada layer **Usecase** (Logika Bisnis).

**Cara Kerjanya:**
1. Kita membuat file `auth_usecase_test.go`.
2. Kita membuat skenario: "Coba berikan NIM '123' dan Password salah".
3. Skrip tes memanggil fungsi `Login("123", "salah")`.
4. Kita melakukan *Assertion* (penegasan): "Pastikan fungsi tersebut mengembalikan error 'Password salah'".
5. Jika saat dites fungsi tersebut malah mengembalikan pesan "Login Berhasil", maka tes akan **GAGAL (Failed)** dan terminal akan berwarna merah.

### B. Teknik Mocking
Bagaimana cara mengetes fitur Login tanpa menghubungkannya ke database PostgreSQL asli?
Jawabannya adalah **Mocking**. Kita membuat "Database Bohongan" (*Mock Repository*) di memori.
- Kita perintah *Mock*: "Hei Mock, kalau ada yang cari NIM 123, kembalikan data si A dengan hash password X".
- Usecase kemudian akan memproses data dari *Mock* tersebut seolah-olah itu dari *database* asli. Ini membuat pengetesan berjalan **ribuan kali lebih cepat** karena tidak ada koneksi internet atau baca-tulis *hardisk* (I/O).

---

## 2. Frontend Testing (SvelteKit & TypeScript)

Di sisi Frontend, pengujian dibagi menjadi dua jenis utama karena selain logika, ada wujud visual yang harus dicek.

### A. Unit Testing dengan Vitest
Vitest adalah *framework* pengetesan generasi baru yang super cepat. Kerjanya mirip dengan pengetesan Backend, tapi fokusnya ke fungsi-fungsi TypeScript atau satu buah Komponen Svelte.

**Cara Kerjanya:**
1. Kita ingin mengetes komponen `<Countdown deadline="..." />`.
2. Vitest akan melakukan *render* komponen tersebut di *virtual memory* (tanpa membuka *browser*).
3. Vitest lalu menyimulasikan waktu yang berjalan dan mengecek apakah teks di dalam komponen berubah dari `00:05` menjadi `00:04`.
4. Jika tidak berubah, tes **GAGAL**.

### B. End-to-End (E2E) Testing dengan Playwright
Ini adalah **level tertinggi** dari pengujian. Playwright adalah robot yang benar-benar membuka *browser* Google Chrome/Firefox secara otomatis!

**Cara Kerjanya:**
1. Kita menulis skenario *User Flow* (Alur Pengguna).
2. Playwright akan menyalakan server lokal secara otomatis.
3. Playwright membuka *browser*, mengetik url `http://localhost:5173/praktikum/login`.
4. Robot mencari kotak teks berlabel "NIM" dan mengetik angka.
5. Robot mengeklik tombol "Masuk".
6. Robot mengecek apakah url pindah ke `/praktikum/dashboard` dan apakah ada teks "Selamat Datang" di layar.
7. Jika di salah satu langkah tombolnya hilang atau ada *error* Svelte yang muncul, tes E2E akan berhenti dan melaporkan **GAGAL**.

---

## Kesimpulan: Alur Kerja (Workflow)

Jika sistem *Automated Testing* ini sudah dipasang, maka cara kerja sehari-harimu sebagai pembuat kode (*Developer*) akan berubah menjadi:

1. Kamu mengedit kode atau menambahkan fitur baru.
2. Kamu menjalankan perintah `make test` (menjalankan seluruh *Backend, Vitest, dan Playwright test*).
3. Komputer akan berpikir selama ~5-15 detik.
4. Jika muncul tulisan **"All Tests Passed ✅"**, kamu bisa dengan tenang melakukan `git push` karena kamu yakin 100% tidak ada yang rusak.
5. Jika muncul tulisan **"1 Test Failed ❌"**, kamu segera membatalkan `push`, melihat baris kode mana yang rusak, dan memperbaikinya.

---

## 3. Rancangan Arsitektur (Implementation Plan)

Untuk mengintegrasikan pengujian ini ke Lab-AP v3, berikut adalah rancangan sistem dan *tools* yang akan diinstal:

### A. Konfigurasi Backend (Golang)
- **Library Tambahan**: Menginstal `github.com/stretchr/testify` (untuk fungsi *assert* / memastikan nilai) dan `github.com/vektra/mockery` (alat pintar untuk membuat "Database Bohongan" alias *Mock Repository* secara otomatis).
- **Struktur Direktori**: Berdasarkan Clean Architecture kita, kita akan meletakkan file tes tepat di sebelah file aslinya. Contoh: `backend/internal/usecase/auth_usecase_test.go`.
- **Target Proof of Concept (PoC)**: Menulis pengujian untuk `Login` dan `Penilaian` Usecase untuk membuktikan *Mocking* berfungsi tanpa *database*.
- **Makefile**: Menambahkan perintah `make test` (untuk menjalankan tes) dan `make mock` (untuk memperbarui database bohongan jika skema berubah).

### B. Konfigurasi Frontend (SvelteKit)
- **Unit Test (Vitest)**:
  - **Library**: Menginstal `vitest`, `@testing-library/svelte`, dan `jsdom`.
  - **Konfigurasi**: Mengubah `vite.config.ts` agar mendukung pengetesan komponen Svelte di memori tanpa membuka browser.
  - **Target PoC**: Membuat `Countdown.test.ts` untuk mengetes logika waktu mundur.
- **End-to-End Test (Playwright)**:
  - **Library**: Menginstal `@playwright/test`.
  - **Konfigurasi**: Membuat `playwright.config.ts` untuk mengatur browser Chrome/Firefox/Safari.
  - **Target PoC**: Membuat skenario `login.spec.ts` (Robot membuka web, mengetik NIM, klik submit, dan mengecek apakah berhasil masuk ke Dashboard).
- **Package.json**: Menambahkan perintah `npm run test:unit` dan `npm run test:e2e`.
