package ollama

import (
	"strings"
	"testing"
)

func TestRubrikSelection(t *testing.T) {
	if rubrikEssay("pretest", "easy")["benar"] != 20 {
		t.Fatal("pretest easy benar harus 20")
	}
	if rubrikEssay("posttest", "medium")["benar_penjelasan"] != 35 {
		t.Fatal("posttest medium benar_penjelasan harus 35")
	}
	if subForCourse("keterampilan") != (subCrit{35, 30, 3, 20}) {
		t.Fatal("keterampilan sub-kriteria salah")
	}
	if subForCourse("ujian_praktik") != (subCrit{25, 30, 3, 13}) {
		t.Fatal("ujian_praktik sub-kriteria salah")
	}
}

func TestPromptTidakNgawur(t *testing.T) {
	pe := promptEssay("Apa itu variabel?", "Tempat menyimpan nilai", "wadah data", "easy", rubrikEssay("posttest", "easy"))
	for _, must := range []string{"JAWABAN MAHASISWA", "RUBRIK PENILAIAN", "\"benar\": 20", "kategori", "JSON"} {
		if !strings.Contains(pe, must) {
			t.Fatalf("prompt esai tak memuat %q", must)
		}
	}

	pc := promptCoding("Buat program konversi suhu", "1. include stdio.h\n2. scanf", "int main(){}", subForCourse("keterampilan"))
	for _, must := range []string{"berjalan_tanpa_error", "tepat_waktu_selesai", "sesuai_petunjuk", "error besar", "tidak selesai", "KODE MAHASISWA", "INSTRUKSI"} {
		if !strings.Contains(pc, must) {
			t.Fatalf("prompt coding tak memuat %q", must)
		}
	}
	t.Logf("\n--- CONTOH PROMPT ESAI ---\n%s\n", pe)
	t.Logf("\n--- CONTOH PROMPT CODING ---\n%s\n", pc)
}
