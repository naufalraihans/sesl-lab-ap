package dto

type AmpuanRequest struct {
	AsistenID int    `json:"asisten_id" binding:"required"`
	KelasID   int    `json:"kelas_id" binding:"required"`
	Kelompok  string `json:"kelompok" binding:"required"`
}
