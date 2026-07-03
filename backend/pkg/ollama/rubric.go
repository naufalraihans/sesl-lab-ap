package ollama

import (
	"fmt"
	"html"
	"math"
	"regexp"
	"sort"
	"strings"
)

var htmlTagRe = regexp.MustCompile(`(?s)<[^>]*>`)

// stripHTML membersihkan teks dari tag HTML (editor edra) + unescape entitas,
// supaya soal/kunci yang dikirim ke AI berupa teks/kode bersih (bukan HTML).
func stripHTML(s string) string {
	s = strings.ReplaceAll(s, "</p>", "\n")
	s = strings.ReplaceAll(s, "</li>", "\n")
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")
	s = htmlTagRe.ReplaceAllString(s, "")
	return strings.TrimSpace(html.UnescapeString(s))
}

// Rubrik bobot AI grading — DEFAULT diambil dari project lama (lab-ap-v2 DEFAULT_RUBRIK).
// Esai (pretest/posttest easy-medium, pretest hard): klasifikasi kategori -> poin.
// Coding (keterampilan, posttest hard, ujian praktik): 3 sub-kriteria
//   Sesuai Petunjuk (SP) + Berjalan Tanpa Error (BTE) + Tepat Waktu/Selesai (TW).

type subCrit struct {
	SesuaiPetunjukMax int
	BteMax            int
	TwMin             int
	TwMax             int
}

var (
	rubrikPretest = map[string]map[string]int{
		"easy":   {"benar": 20, "salah": 8, "kosong": 0},
		"medium": {"benar_penjelasan": 15, "benar_singkat": 10, "salah_penjelasan": 7, "salah": 3, "kosong": 0},
		"hard":   {"benar_penjelasan": 25, "benar_singkat": 15, "salah_penjelasan": 10, "salah": 5, "kosong": 0},
	}
	rubrikPosttest = map[string]map[string]int{
		"easy":   {"benar": 20, "salah": 8, "kosong": 0},
		"medium": {"benar_penjelasan": 35, "benar_singkat": 20, "salah_penjelasan": 15, "salah": 10, "kosong": 0},
	}
	// Keterampilan total 85 (SP35 + BTE30 + TW20). Coding posttest-hard & ujian
	// praktik total 45 (SP25 + BTE7 + TW13) — BTE jauh lebih kecil. Sesuai rubrik resmi.
	subKeterampilan = subCrit{SesuaiPetunjukMax: 35, BteMax: 30, TwMin: 10, TwMax: 20}
	subPosttestHard = subCrit{SesuaiPetunjukMax: 25, BteMax: 7, TwMin: 3, TwMax: 13}
	subUjianPraktik = subCrit{SesuaiPetunjukMax: 25, BteMax: 7, TwMin: 3, TwMax: 13}
)

// subForCourse memilih sub-kriteria coding berdasar jenis course.
func subForCourse(jenisCourse string) subCrit {
	switch jenisCourse {
	case "posttest":
		return subPosttestHard
	case "ujian_praktik":
		return subUjianPraktik
	default: // keterampilan & lainnya
		return subKeterampilan
	}
}

// rubrikEssay mengembalikan tabel kategori->poin untuk esai sesuai course+level.
func rubrikEssay(jenisCourse, level string) map[string]int {
	tbl := rubrikPretest
	if jenisCourse == "posttest" {
		tbl = rubrikPosttest
	}
	if m, ok := tbl[level]; ok {
		return m
	}
	// fallback bila level tak dikenal: skala biner sederhana
	return map[string]int{"benar": 20, "salah": 8, "kosong": 0}
}

// promptEssay: klasifikasi jawaban esai ke salah satu kategori rubrik (port dari
// buildPromptPTClassification). poin maksimal soal = poinMax (cap akhir).
func promptEssay(soal, jawabanRef, jawabanMhs, level string, kategori map[string]int) string {
	// urut kategori menurun supaya konsisten & terbaca
	type kv struct {
		k string
		v int
	}
	var arr []kv
	for k, v := range kategori {
		arr = append(arr, kv{k, v})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].v > arr[j].v })
	var b strings.Builder
	for _, e := range arr {
		fmt.Fprintf(&b, "  - \"%s\": %d poin\n", e.k, e.v)
	}
	if jawabanMhs == "" {
		jawabanMhs = "(kosong / tidak dijawab)"
	}
	if jawabanRef == "" {
		jawabanRef = "(tidak ada kunci, nilai berdasarkan kebenaran umum)"
	}
	return fmt.Sprintf(`Kamu adalah pengoreksi ujian Pre-test/Post-test mata kuliah Algoritma dan Pemrograman.

SOAL:
%s

JAWABAN REFERENSI (kunci jawaban yang benar):
%s

JAWABAN MAHASISWA:
%s

RUBRIK PENILAIAN untuk level "%s":
%s
INSTRUKSI:
1. Bandingkan jawaban mahasiswa dengan jawaban referensi.
2. Tentukan kategori yang PALING sesuai untuk jawaban mahasiswa.
3. "kosong" HANYA untuk jawaban yang benar-benar kosong, hanya spasi, atau
   sama sekali tidak relevan dengan soal. Jangan pakai "kosong" untuk jawaban
   yang salah tapi ada isinya.
4. Jika jawaban SALAH namun mahasiswa menuliskan penjelasan/alasan/usaha yang
   masih relevan dengan soal (mis. jawaban panjang tapi keliru), pilih kategori
   "salah_penjelasan" bila tersedia di rubrik. Jika rubrik tidak punya
   "salah_penjelasan" (mis. level easy), pilih "salah" — BUKAN "kosong".
5. "salah" dipakai untuk jawaban salah yang singkat/tanpa penjelasan berarti.

Jawab HANYA dalam format JSON (tanpa teks lain):
{"kategori":"nama kategori dari rubrik","poin":angka_sesuai_rubrik,"feedback":"alasan singkat 1 kalimat (Indonesia)"}`,
		soal, jawabanRef, jawabanMhs, level, b.String())
}

// promptCoding: penilaian kode 3 sub-kriteria (port dari buildPromptKeterampilan/
// buildPromptUprakCoding). kunci/instruksi disertakan sebagai referensi.
func promptCoding(soal, instruksi, kode string, sc subCrit) string {
	bteMax := sc.BteMax
	bteMinor := int(math.Round(float64(bteMax) * 0.75)) // typo/syntax minor
	bteSedang := int(math.Round(float64(bteMax) * 0.5))
	bteBesar := int(math.Round(float64(bteMax) * 0.25)) // error besar / "kode salah"
	twMin, twMax := sc.TwMin, sc.TwMax
	twMinor := int(math.Round(float64(twMax) * 0.75))
	twIncomplete := int(math.Round(float64(twMax) * 0.3))

	if kode == "" {
		kode = "(kosong / tidak dijawab)"
	}
	instrBlock := ""
	if strings.TrimSpace(instruksi) != "" {
		instrBlock = "\nKUNCI JAWABAN / INSTRUKSI (acuan penilaian):\n" + instruksi + "\n"
	}
	return fmt.Sprintf(`Kamu adalah pengoreksi program (coding) mata kuliah Algoritma dan Pemrograman.

SOAL / DESKRIPSI:
%s
%s
KODE MAHASISWA:
%s

NILAI 3 SUB-KRITERIA:

1. "sesuai_petunjuk" (0-%d): KUNCI JAWABAN memuat instruksi BERBOBOT POIN per item, biasanya
   berformat komentar "# N. deskripsi (X poin)" diikuti kode referensi item tsb.
   Periksa kode mahasiswa terhadap SETIAP item instruksi, lalu JUMLAHKAN poin item yang terpenuhi:
   terpenuhi penuh = poin penuh item, sebagian benar = proporsional, tidak ada = 0.
   (Jika kunci tidak memuat bobot poin per item, nilai proporsional ke kesesuaian.) Maksimal %d.

2. "berjalan_tanpa_error" (0-%d) — BACA SKALA INI:
   - %d: sempurna, zero error
   - %d-%d: typo/syntax minor saja (logika BENAR, salah ketik 1-2 char, mis. "flot"->"float", "whiel"->"while")
   - %d-%d: error sedang (logika sebagian benar, mudah diperbaiki)
   - %d-%d: error besar (logika salah / kode salah / output salah / algoritma tidak paham)
   - 0: kosong / tidak ada kode

3. "tepat_waktu_selesai" (%d-%d):
   - %d: Complete (selesai penuh)
   - %d-%d: Minor missing
   - %d-%d: Incomplete (kode tidak selesai)
   - %d-%d: Minimum

ATURAN: jika kode TIDAK kosong, sesuai_petunjuk minimal 1. berjalan_tanpa_error = 0 HANYA jika kode benar-benar kosong.

Jawab HANYA JSON (tanpa teks lain):
{"sesuai_petunjuk":angka,"berjalan_tanpa_error":angka,"tepat_waktu_selesai":angka,"feedback":"catatan singkat (Indonesia)"}`,
		soal, instrBlock, kode,
		sc.SesuaiPetunjukMax, sc.SesuaiPetunjukMax,
		bteMax, bteMax, bteMinor, bteMax-1, bteSedang, bteMinor-1, bteBesar, bteSedang-1,
		twMin, twMax, twMax, twMinor, twMax-1, twIncomplete, twMinor-1, twMin, twIncomplete-1,
	)
}
