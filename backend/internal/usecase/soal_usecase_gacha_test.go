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
