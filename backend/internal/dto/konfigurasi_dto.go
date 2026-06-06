package dto

// KonfigurasiRequest: set nilai konfigurasi (GDrive jadwal / modul URL / dll).
type KonfigurasiRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}
