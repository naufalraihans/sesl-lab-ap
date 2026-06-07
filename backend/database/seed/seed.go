package main

import (
	"log"
	"time"

	"lab-ap/config"
	"lab-ap/database"
	"lab-ap/internal/entity"
	"lab-ap/pkg/hash"

	"gorm.io/gorm"
)

// Seed data untuk pengembangan lokal: admin, kelas, jadwal, mahasiswa, sesi + soal contoh.
func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DSN(), false)
	if err != nil {
		log.Fatalf("Koneksi DB gagal: %v", err)
	}

	seedKelas(db)
	seedAdmin(db)
	seedMahasiswa(db)
	seedJadwal(db)
	seedSesiSoal(db)

	log.Println("✓ Seed selesai.")
}

func seedKelas(db *gorm.DB) {
	kelas := []entity.Kelas{
		{NamaKelas: "TTL A", IsRegisterOpen: true},
		{NamaKelas: "TTL B", IsRegisterOpen: true},
	}
	for _, k := range kelas {
		var existing entity.Kelas
		if err := db.Where("nama_kelas = ?", k.NamaKelas).First(&existing).Error; err == gorm.ErrRecordNotFound {
			db.Create(&k)
			log.Printf("  + kelas %s", k.NamaKelas)
		}
	}
}

func seedAdmin(db *gorm.DB) {
	var n int64
	db.Model(&entity.User{}).Where("nim = ?", "admin").Count(&n)
	if n > 0 {
		return
	}
	h, _ := hash.Password("admin123")
	hp := "081234567890"
	medsos := "https://linkedin.com/in/asisten-lab"
	admin := entity.User{
		Role:         entity.RoleAdmin,
		NIM:          "admin",
		Nama:         "Asisten Lab Utama",
		PasswordHash: &h,
		IsRegistered: true,
		NomorHP:      &hp,
		MedsosLink:   &medsos,
	}
	db.Create(&admin)
	log.Println("  + admin (nim: admin / pass: admin123)")
}

func seedMahasiswa(db *gorm.DB) {
	var kelasA entity.Kelas
	if err := db.Where("nama_kelas = ?", "TTL A").First(&kelasA).Error; err != nil {
		return
	}
	shift1 := 1
	// Satu mahasiswa sudah register (untuk uji login langsung), sisanya belum.
	registered, _ := hash.Password("mahasiswa123")
	mhs := []entity.User{
		{Role: entity.RoleUser, NIM: "2021001", Nama: "Budi Santoso", KelasID: &kelasA.ID, Shift: &shift1, IsRegistered: true, PasswordHash: &registered},
		{Role: entity.RoleUser, NIM: "2021002", Nama: "Siti Aminah", KelasID: &kelasA.ID, Shift: &shift1, IsRegistered: false},
		{Role: entity.RoleUser, NIM: "2021003", Nama: "Andi Wijaya", KelasID: &kelasA.ID, Shift: &shift1, IsRegistered: false},
	}
	for _, m := range mhs {
		var existing entity.User
		if err := db.Where("nim = ?", m.NIM).First(&existing).Error; err == gorm.ErrRecordNotFound {
			db.Create(&m)
			log.Printf("  + mahasiswa %s (%s)", m.Nama, m.NIM)
		}
	}
}

func seedJadwal(db *gorm.DB) {
	var kelasA entity.Kelas
	if err := db.Where("nama_kelas = ?", "TTL A").First(&kelasA).Error; err != nil {
		return
	}
	var n int64
	db.Model(&entity.Jadwal{}).Where("kelas_id = ? AND shift = ?", kelasA.ID, 1).Count(&n)
	if n == 0 {
		db.Create(&entity.Jadwal{
			KelasID: kelasA.ID, Shift: 1, Hari: "Senin",
			JamMulai: "08:00:00", JamSelesai: "10:00:00", Keterangan: "Minggu 1-4",
		})
		log.Println("  + jadwal TTL A shift 1")
	}
}

func seedSesiSoal(db *gorm.DB) {
	now := time.Now()

	var n int64
	db.Model(&entity.SesiPraktikum{}).Where("is_ujian_praktik = ?", false).Count(&n)
	if n == 0 {
		sesi := entity.SesiPraktikum{
			JudulSesi: "Modul 1 - Pengenalan Dasar Bahasa C",
			Deskripsi: "Variabel, tipe data, dan operasi dasar.",
			Urutan:    1,
			CreatedAt: now, UpdatedAt: now,
		}
		db.Create(&sesi)

		pretest := entity.Course{SesiPraktikumID: sesi.ID, Jenis: entity.CoursePretest, Judul: "Pre-test Modul 1", DurasiMenit: 20}
		posttest := entity.Course{SesiPraktikumID: sesi.ID, Jenis: entity.CoursePosttest, Judul: "Post-test Modul 1", DurasiMenit: 25}
		keterampilan := entity.Course{SesiPraktikumID: sesi.ID, Jenis: entity.CourseKeterampilan, Judul: "Keterampilan Modul 1", DurasiMenit: 40}
		db.Create(&pretest)
		db.Create(&posttest)
		db.Create(&keterampilan)

		easy := entity.DiffEasy
		medium := entity.DiffMedium
		hard := entity.DiffHard

		pretestSoal := []entity.Soal{
			{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "Apa itu variabel dalam bahasa C?", Poin: 20},
			{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &medium, TeksSoal: "Apa fungsi printf dan scanf?", Poin: 20},
		}

		posttestSoal := []entity.Soal{
			{CourseID: posttest.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "Sebutkan langkah kompilasi program C.", Poin: 30},
			{CourseID: posttest.ID, JenisSoal: entity.SoalCoding, Difficulty: &hard, TeksSoal: "Tulis program C yang menjumlahkan dua bilangan dari input.", Poin: 40},
		}

		keterampilanSoal := []entity.Soal{
			{CourseID: keterampilan.ID, JenisSoal: entity.SoalCoding, TeksSoal: "Buat program konversi suhu Celcius ke Fahrenheit.", Poin: 100},
		}

		all := append([]entity.Soal{}, pretestSoal...)
		all = append(all, posttestSoal...)
		all = append(all, keterampilanSoal...)
		for i := range all {
			all[i].CreatedAt = now
		}
		db.Create(&all)
		log.Printf("  + sesi 'Modul 1 - Pengenalan Dasar Bahasa C' + %d soal contoh", len(all))
	}

	// Seed Ujian Praktik
	var nUjian int64
	db.Model(&entity.SesiPraktikum{}).Where("is_ujian_praktik = ?", true).Count(&nUjian)
	if nUjian == 0 {
		sesiUjian := entity.SesiPraktikum{
			JudulSesi:      "Ujian Tengah Semester Praktikum",
			Deskripsi:      "Ujian praktik menguasai Modul 1-3.",
			Urutan:         99,
			IsUjianPraktik: true,
			CreatedAt:      now, UpdatedAt: now,
		}
		db.Create(&sesiUjian)

		courseUjian := entity.Course{SesiPraktikumID: sesiUjian.ID, Jenis: entity.CourseUjianPraktik, Judul: "Soal Ujian Praktik", DurasiMenit: 60}
		db.Create(&courseUjian)

		kat1 := entity.KatModul1
		kat2 := entity.KatModul2
		kat3 := entity.KatModul3
		kat45 := entity.KatModul45
		kat6 := entity.KatModul6
		katF := entity.KatFlowchart

		ujianSoal := []entity.Soal{
			{CourseID: courseUjian.ID, JenisSoal: entity.SoalEssay, KategoriUjian: &kat1, TeksSoal: "Soal Ujian Praktik - Modul 1", Poin: 15},
			{CourseID: courseUjian.ID, JenisSoal: entity.SoalEssay, KategoriUjian: &kat2, TeksSoal: "Soal Ujian Praktik - Modul 2", Poin: 15},
			{CourseID: courseUjian.ID, JenisSoal: entity.SoalEssay, KategoriUjian: &kat3, TeksSoal: "Soal Ujian Praktik - Modul 3", Poin: 15},
			{CourseID: courseUjian.ID, JenisSoal: entity.SoalCoding, KategoriUjian: &kat45, TeksSoal: "Soal Ujian Praktik - Modul 4 & 5", Poin: 20},
			{CourseID: courseUjian.ID, JenisSoal: entity.SoalCoding, KategoriUjian: &kat6, TeksSoal: "Soal Ujian Praktik - Modul 6", Poin: 20},
			{CourseID: courseUjian.ID, JenisSoal: entity.SoalEssay, KategoriUjian: &katF, TeksSoal: "Soal Ujian Praktik - Flowchart", Poin: 15},
		}
		for i := range ujianSoal {
			ujianSoal[i].CreatedAt = now
		}
		db.Create(&ujianSoal)
		log.Printf("  + sesi '%s' (Ujian Praktik)", sesiUjian.JudulSesi)
	}

	// Seed Dummy Nilai for Pivot
	var nPengerjaan int64
	db.Model(&entity.PengerjaanCourse{}).Count(&nPengerjaan)
	if nPengerjaan == 0 {
		var mhs entity.User
		if err := db.Where("nim = ?", "2021003").First(&mhs).Error; err == nil {
			var sesi entity.SesiPraktikum
			if err := db.Where("judul_sesi = ?", "Modul 1 - Pengenalan Dasar Bahasa C").First(&sesi).Error; err == nil {
				// Create Aktivasi Sesi
				aks := entity.AktivasiSesi{
					SesiPraktikumID: sesi.ID,
					KelasID:         *mhs.KelasID,
					Shift:           *mhs.Shift,
					IsActive:        true,
				}
				db.Create(&aks)

				var courses []entity.Course
				db.Where("sesi_praktikum_id = ?", sesi.ID).Find(&courses)

				for _, c := range courses {
					nilai := float64(80 + c.ID*2) // Dummy score
					pc := entity.PengerjaanCourse{
						MahasiswaID:    mhs.ID,
						AktivasiSesiID: aks.ID,
						CourseID:       c.ID,
						Status:         entity.StatusSelesai,
						TotalNilai:     &nilai,
					}
					db.Create(&pc)
				}
				log.Printf("  + dummy nilai untuk mhs %s", mhs.NIM)
			}
		}
	}
}
