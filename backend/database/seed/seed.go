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
	var kelasB entity.Kelas
	if err := db.Where("nama_kelas = ?", "TTL B").First(&kelasB).Error; err != nil {
		return
	}

	shift1 := 1
	shift2 := 2
	registered, _ := hash.Password("mahasiswa123")
	mhs := []entity.User{
		{Role: entity.RoleUser, NIM: "2021001", Nama: "Budi Santoso", KelasID: &kelasA.ID, Shift: &shift1, IsRegistered: true, PasswordHash: &registered},
		{Role: entity.RoleUser, NIM: "2021002", Nama: "Siti Aminah", KelasID: &kelasA.ID, Shift: &shift1, IsRegistered: false},
		{Role: entity.RoleUser, NIM: "2021003", Nama: "Andi Wijaya", KelasID: &kelasA.ID, Shift: &shift1, IsRegistered: false},
		{Role: entity.RoleUser, NIM: "2021004", Nama: "Dewi Lestari", KelasID: &kelasB.ID, Shift: &shift2, IsRegistered: false},
		{Role: entity.RoleUser, NIM: "2021005", Nama: "Eko Pratama", KelasID: &kelasB.ID, Shift: &shift2, IsRegistered: false},
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
	var kelasA, kelasB entity.Kelas
	db.Where("nama_kelas = ?", "TTL A").First(&kelasA)
	db.Where("nama_kelas = ?", "TTL B").First(&kelasB)

	var n int64
	db.Model(&entity.Jadwal{}).Count(&n)
	if n == 0 {
		db.Create(&entity.Jadwal{
			KelasID: kelasA.ID, Shift: 1, Hari: "Senin",
			JamMulai: "08:00:00", JamSelesai: "10:00:00", Keterangan: "Minggu 1-4",
		})
		db.Create(&entity.Jadwal{
			KelasID: kelasB.ID, Shift: 2, Hari: "Selasa",
			JamMulai: "13:00:00", JamSelesai: "15:00:00", Keterangan: "Minggu 1-4",
		})
		log.Println("  + jadwal TTL A shift 1 & TTL B shift 2")
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
			{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "<p>Apa itu variabel dalam bahasa C?</p>", Poin: 20},
			{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &medium, TeksSoal: "<p>Apa fungsi <code>printf</code> dan <code>scanf</code>?</p>", Poin: 20},
		}

		posttestSoal := []entity.Soal{
			{CourseID: posttest.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "<p>Sebutkan langkah kompilasi program C.</p>", Poin: 30},
			{CourseID: posttest.ID, JenisSoal: entity.SoalCoding, Difficulty: &hard, TeksSoal: "<p>Tulis program C yang menjumlahkan dua bilangan dari input.</p>", Poin: 40},
		}

		keterampilanSoal := []entity.Soal{
			{CourseID: keterampilan.ID, JenisSoal: entity.SoalCoding, TeksSoal: "<p>Buat program konversi suhu Celcius ke Fahrenheit dengan output presisi 2 desimal.</p>", Poin: 100},
		}

		all := append([]entity.Soal{}, pretestSoal...)
		all = append(all, posttestSoal...)
		all = append(all, keterampilanSoal...)
		for i := range all {
			all[i].CreatedAt = now
		}
		db.Create(&all)
		log.Printf("  + sesi 'Modul 1 - Pengenalan Dasar Bahasa C' + %d soal contoh", len(all))

		// SEED MODUL 2
		sesi2 := entity.SesiPraktikum{
			JudulSesi: "Modul 2 - Percabangan (If-Else & Switch Case)",
			Deskripsi: "Mempelajari logika kontrol alur program percabangan.",
			Urutan:    2,
			CreatedAt: now, UpdatedAt: now,
		}
		db.Create(&sesi2)

		pre2 := entity.Course{SesiPraktikumID: sesi2.ID, Jenis: entity.CoursePretest, Judul: "Pre-test Modul 2", DurasiMenit: 15}
		post2 := entity.Course{SesiPraktikumID: sesi2.ID, Jenis: entity.CoursePosttest, Judul: "Post-test Modul 2", DurasiMenit: 25}
		ket2 := entity.Course{SesiPraktikumID: sesi2.ID, Jenis: entity.CourseKeterampilan, Judul: "Keterampilan Modul 2", DurasiMenit: 45}
		db.Create(&pre2)
		db.Create(&post2)
		db.Create(&ket2)

		soalMod2 := []entity.Soal{
			{CourseID: pre2.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "<p>Apa perbedaan <code>if</code> dan <code>switch</code>?</p>", Poin: 100},
			{CourseID: post2.ID, JenisSoal: entity.SoalEssay, Difficulty: &medium, TeksSoal: "<p>Berikan contoh blok <code>if-else if-else</code>.</p>", Poin: 100},
			{CourseID: ket2.ID, JenisSoal: entity.SoalCoding, TeksSoal: "<p>Buat program kalkulator sederhana (+, -, *, /) menggunakan <code>switch case</code>.</p>", Poin: 100},
		}
		for i := range soalMod2 {
			soalMod2[i].CreatedAt = now
		}
		db.Create(&soalMod2)
		log.Println("  + sesi 'Modul 2 - Percabangan' + soal")
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
			{CourseID: courseUjian.ID, JenisSoal: entity.SoalEssay, KategoriUjian: &kat1, TeksSoal: "<p>Ini soal Nomor 1 dengan contoh LaTeX: <span data-type=\"inlineMath\" data-latex=\"\\frac{1}{2}\">$\\frac{1}{2}$</span></p><p>Buatkan <code>#include &lt;stdio.h&gt;</code> gitu aja sih.</p>", Poin: 15},
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
		var mhs1, mhs2 entity.User
		db.Where("nim = ?", "2021001").First(&mhs1)
		db.Where("nim = ?", "2021003").First(&mhs2)
		
		var sesi1, sesi2 entity.SesiPraktikum
		db.Where("judul_sesi = ?", "Modul 1 - Pengenalan Dasar Bahasa C").First(&sesi1)
		db.Where("judul_sesi = ?", "Modul 2 - Percabangan (If-Else & Switch Case)").First(&sesi2)
		
		if sesi1.ID != 0 {
			aks1 := entity.AktivasiSesi{ SesiPraktikumID: sesi1.ID, KelasID: *mhs1.KelasID, Shift: *mhs1.Shift, IsActive: true }
			db.Create(&aks1)
			
			var courses []entity.Course
			db.Where("sesi_praktikum_id = ?", sesi1.ID).Find(&courses)
			
			for _, c := range courses {
				// Mhs 1 Modul 1
				nilai1 := float64(85 + c.ID%5)
				db.Create(&entity.PengerjaanCourse{ MahasiswaID: mhs1.ID, AktivasiSesiID: aks1.ID, CourseID: c.ID, Status: entity.StatusSelesai, TotalNilai: &nilai1 })
				
				// Mhs 2 Modul 1
				nilai2 := float64(70 + c.ID%5)
				db.Create(&entity.PengerjaanCourse{ MahasiswaID: mhs2.ID, AktivasiSesiID: aks1.ID, CourseID: c.ID, Status: entity.StatusSelesai, TotalNilai: &nilai2 })
			}
		}

		if sesi2.ID != 0 {
			aks2 := entity.AktivasiSesi{ SesiPraktikumID: sesi2.ID, KelasID: *mhs1.KelasID, Shift: *mhs1.Shift, IsActive: true }
			db.Create(&aks2)

			var courses2 []entity.Course
			db.Where("sesi_praktikum_id = ?", sesi2.ID).Find(&courses2)
			for _, c := range courses2 {
				nilai1 := float64(90)
				db.Create(&entity.PengerjaanCourse{ MahasiswaID: mhs1.ID, AktivasiSesiID: aks2.ID, CourseID: c.ID, Status: entity.StatusSelesai, TotalNilai: &nilai1 })
			}
		}
		log.Println("  + dummy nilai untuk mhs 2021001 & 2021003 di Modul 1 dan 2")
	}
}
