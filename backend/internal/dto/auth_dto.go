package dto

// CekNIMRequest: langkah pertama login (cek status NIM).
type CekNIMRequest struct {
	NIM string `json:"nim" binding:"required"`
}

// CekNIMResponse: status NIM untuk menentukan alur (login / register / ditolak).
type CekNIMResponse struct {
	NIM            string `json:"nim"`
	Ditemukan      bool   `json:"ditemukan"`
	IsRegistered   bool   `json:"is_registered"`
	IsRegisterOpen bool   `json:"is_register_open"`
	Nama           string `json:"nama,omitempty"`
	Pesan          string `json:"pesan"`
}

// LoginRequest: login normal (NIM + password).
type LoginRequest struct {
	NIM      string `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest: first-time register (set password).
type RegisterRequest struct {
	NIM      string `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse: hasil login/register berhasil.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse: representasi user aman untuk dikirim ke klien.
type UserResponse struct {
	ID         int     `json:"id"`
	Role       string  `json:"role"`
	NIM        string  `json:"nim"`
	Nama       string  `json:"nama"`
	KelasID    *int    `json:"kelas_id"`
	NamaKelas  string  `json:"nama_kelas,omitempty"`
	Shift      *int    `json:"shift"`
	FotoURL    *string `json:"foto_url,omitempty"`
	NomorHP    *string `json:"nomor_hp,omitempty"`
	MedsosLink *string `json:"medsos_link,omitempty"`
}
