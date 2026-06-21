package usecase

import (
	"testing"

	"lab-ap/internal/entity"
)

func intp(n int) *int { return &n }

func TestDistribusiDifficulty(t *testing.T) {
	// Default per jenis bila kuota kosong.
	pre := distribusiDifficulty(&entity.Course{Jenis: entity.CoursePretest})
	if pre[entity.DiffEasy] != 1 || pre[entity.DiffMedium] != 2 || pre[entity.DiffHard] != 2 {
		t.Fatalf("default pretest salah: %v", pre)
	}
	post := distribusiDifficulty(&entity.Course{Jenis: entity.CoursePosttest})
	if post[entity.DiffEasy] != 1 || post[entity.DiffMedium] != 1 || post[entity.DiffHard] != 1 {
		t.Fatalf("default posttest salah: %v", post)
	}

	// Kuota per course menimpa default (termasuk yang nil dianggap 0).
	c := &entity.Course{Jenis: entity.CoursePretest, KuotaEasy: intp(1), KuotaMedium: intp(1), KuotaHard: intp(1)}
	got := distribusiDifficulty(c)
	if got[entity.DiffEasy] != 1 || got[entity.DiffMedium] != 1 || got[entity.DiffHard] != 1 {
		t.Fatalf("kuota override salah: %v", got)
	}

	// Sebagian kuota diisi → yang nil = 0, bukan fallback default.
	c2 := &entity.Course{Jenis: entity.CoursePretest, KuotaEasy: intp(3)}
	got2 := distribusiDifficulty(c2)
	if got2[entity.DiffEasy] != 3 || got2[entity.DiffMedium] != 0 || got2[entity.DiffHard] != 0 {
		t.Fatalf("kuota parsial salah: %v", got2)
	}
}

func TestPickRandom(t *testing.T) {
	pool := []entity.Soal{{ID: 1}, {ID: 2}, {ID: 3}}

	got, err := pickRandom(pool, 2, entity.CoursePretest, "easy")
	if err != nil || len(got) != 2 {
		t.Fatalf("pool cukup: harap 2 tanpa error, dapat %d err=%v", len(got), err)
	}

	if _, err := pickRandom(pool, 5, entity.CoursePretest, "easy"); err == nil {
		t.Fatal("pool kurang harus error")
	}

	got0, err := pickRandom(pool, 0, entity.CoursePretest, "easy")
	if err != nil || len(got0) != 0 {
		t.Fatalf("n=0 harap 0 tanpa error, dapat %d err=%v", len(got0), err)
	}
}
