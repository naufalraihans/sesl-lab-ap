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
	var n int64
	db.Model(&entity.SesiPraktikum{}).Count(&n)
	if n > 0 {
		return
	}
	now := time.Now()
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

	// Pool pretest: butuh minimal 1 easy, 2 medium, 2 hard. Sediakan lebih agar acak bervariasi.
	pretestSoal := []entity.Soal{
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "Apa itu variabel dalam bahasa C?", Poin: 20},
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "Sebutkan tipe data dasar di C.", Poin: 20},
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &medium, TeksSoal: "Jelaskan perbedaan int dan float.", Poin: 20},
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &medium, TeksSoal: "Apa fungsi printf dan scanf?", Poin: 20},
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &medium, TeksSoal: "Jelaskan operator aritmatika di C.", Poin: 20},
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &hard, TeksSoal: "Jelaskan konsep pointer secara singkat.", Poin: 20},
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &hard, TeksSoal: "Apa itu type casting dan kapan dipakai?", Poin: 20},
		{CourseID: pretest.ID, JenisSoal: entity.SoalEssay, Difficulty: &hard, TeksSoal: "Jelaskan perbedaan ++i dan i++.", Poin: 20},
	}

	// Pool posttest: 1 easy essay, 1 medium essay, 1 hard coding.
	posttestSoal := []entity.Soal{
		{CourseID: posttest.ID, JenisSoal: entity.SoalEssay, Difficulty: &easy, TeksSoal: "Sebutkan langkah kompilasi program C.", Poin: 30},
		{CourseID: posttest.ID, JenisSoal: entity.SoalEssay, Difficulty: &medium, TeksSoal: "Jelaskan alur eksekusi fungsi main().", Poin: 30},
		{CourseID: posttest.ID, JenisSoal: entity.SoalCoding, Difficulty: &hard, TeksSoal: "Tulis program C yang menjumlahkan dua bilangan dari input.", Poin: 40},
		{CourseID: posttest.ID, JenisSoal: entity.SoalCoding, Difficulty: &hard, TeksSoal: "Tulis program C yang mencetak bilangan genap 1..N.", Poin: 40},
	}

	// Pool keterampilan: live coding.
	keterampilanSoal := []entity.Soal{
		{CourseID: keterampilan.ID, JenisSoal: entity.SoalCoding, TeksSoal: "Buat program konversi suhu Celcius ke Fahrenheit.", Poin: 100},
		{CourseID: keterampilan.ID, JenisSoal: entity.SoalCoding, TeksSoal: "Buat program menghitung luas & keliling lingkaran.", Poin: 100},
	}

	all := append([]entity.Soal{}, pretestSoal...)
	all = append(all, posttestSoal...)
	all = append(all, keterampilanSoal...)
	for i := range all {
		all[i].CreatedAt = now
	}
	db.Create(&all)

	log.Printf("  + sesi '%s' + %d soal contoh", sesi.JudulSesi, len(all))
}
