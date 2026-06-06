package usecase

import (
	"fmt"
	"math/rand"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

type SoalUsecase struct {
	soal     repository.SoalRepository
	course   repository.CourseRepository
	terpilih repository.SoalTerpilihRepository
}

func NewSoalUsecase(s repository.SoalRepository, c repository.CourseRepository, t repository.SoalTerpilihRepository) *SoalUsecase {
	return &SoalUsecase{soal: s, course: c, terpilih: t}
}

// ---- CRUD pool ----

func (uc *SoalUsecase) Create(req dto.SoalRequest) (*entity.Soal, error) {
	if _, err := uc.course.FindByID(req.CourseID); err != nil {
		return nil, ErrNotFound
	}
	s := &entity.Soal{
		CourseID:     req.CourseID,
		JenisSoal:    entity.JenisSoal(req.JenisSoal),
		TeksSoal:     req.TeksSoal,
		GambarURL:    req.GambarURL,
		Poin:         req.Poin,
		KunciJawaban: req.KunciJawaban,
	}
	if req.Difficulty != nil {
		d := entity.Difficulty(*req.Difficulty)
		s.Difficulty = &d
	}
	if req.KategoriUjian != nil {
		k := entity.KategoriUjian(*req.KategoriUjian)
		s.KategoriUjian = &k
	}
	if err := uc.soal.Create(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *SoalUsecase) Update(id int, req dto.SoalRequest) (*entity.Soal, error) {
	s, err := uc.soal.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	s.JenisSoal = entity.JenisSoal(req.JenisSoal)
	s.TeksSoal = req.TeksSoal
	s.GambarURL = req.GambarURL
	s.Poin = req.Poin
	s.KunciJawaban = req.KunciJawaban
	s.Difficulty = nil
	if req.Difficulty != nil {
		d := entity.Difficulty(*req.Difficulty)
		s.Difficulty = &d
	}
	s.KategoriUjian = nil
	if req.KategoriUjian != nil {
		k := entity.KategoriUjian(*req.KategoriUjian)
		s.KategoriUjian = &k
	}
	if err := uc.soal.Update(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *SoalUsecase) Delete(id int) error    { return uc.soal.Delete(id) }
func (uc *SoalUsecase) ListByCourse(courseID int) ([]entity.Soal, error) {
	return uc.soal.ListByCourse(courseID)
}

// ---- Pengacakan ke soal_terpilih ----

// distribusiDifficulty mengembalikan jumlah soal per difficulty (hardcoded by jenis).
func distribusiDifficulty(jenis entity.JenisCourse) map[entity.Difficulty]int {
	switch jenis {
	case entity.CoursePretest:
		// 1 easy, 2 medium, 2 hard (total 5)
		return map[entity.Difficulty]int{entity.DiffEasy: 1, entity.DiffMedium: 2, entity.DiffHard: 2}
	case entity.CoursePosttest:
		// 1 easy (essay), 1 medium (essay), 1 hard (coding) (total 3)
		return map[entity.Difficulty]int{entity.DiffEasy: 1, entity.DiffMedium: 1, entity.DiffHard: 1}
	default:
		return nil
	}
}

// AcakUntukAktivasiCourse mengacak soal dari pool sesuai jenis course dan
// menyimpan hasilnya ke soal_terpilih (idempoten: dilewati jika sudah ada).
func (uc *SoalUsecase) AcakUntukAktivasiCourse(aktivasiSesiID int, course *entity.Course) error {
	exists, err := uc.terpilih.ExistsForAktivasiCourse(aktivasiSesiID, course.ID)
	if err != nil {
		return err
	}
	if exists {
		return nil // sudah diacak sebelumnya
	}

	var dipilih []entity.Soal

	switch course.Jenis {
	case entity.CoursePretest, entity.CoursePosttest:
		dist := distribusiDifficulty(course.Jenis)
		for diff, n := range dist {
			pool, err := uc.soal.PoolByDifficulty(course.ID, diff)
			if err != nil {
				return err
			}
			picked, err := pickRandom(pool, n, course.Jenis, string(diff))
			if err != nil {
				return err
			}
			dipilih = append(dipilih, picked...)
		}
	case entity.CourseKeterampilan:
		pool, err := uc.soal.PoolAll(course.ID)
		if err != nil {
			return err
		}
		picked, err := pickRandom(pool, 1, course.Jenis, "keterampilan")
		if err != nil {
			return err
		}
		dipilih = append(dipilih, picked...)
	case entity.CourseUjianPraktik:
		for _, kat := range entity.SemuaKategoriUjian {
			pool, err := uc.soal.PoolByKategori(course.ID, kat)
			if err != nil {
				return err
			}
			picked, err := pickRandom(pool, 1, course.Jenis, string(kat))
			if err != nil {
				return err
			}
			dipilih = append(dipilih, picked...)
		}
	}

	items := make([]entity.SoalTerpilih, 0, len(dipilih))
	for i, s := range dipilih {
		items = append(items, entity.SoalTerpilih{
			AktivasiSesiID: aktivasiSesiID,
			CourseID:       course.ID,
			SoalID:         s.ID,
			Urutan:         i + 1,
		})
	}
	return uc.terpilih.BulkCreate(items)
}

// pickRandom memilih n soal acak dari pool; error jika pool kurang.
func pickRandom(pool []entity.Soal, n int, jenis entity.JenisCourse, label string) ([]entity.Soal, error) {
	if len(pool) < n {
		return nil, fmt.Errorf("%w: pool soal '%s' (%s) hanya %d, butuh %d",
			ErrBadRequest, label, jenis, len(pool), n)
	}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	return pool[:n], nil
}
