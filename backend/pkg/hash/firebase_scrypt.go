package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

// FbScryptConfig memuat parameter hash project Firebase Authentication
// (Firebase Console -> Auth -> Password hash parameters).
type FbScryptConfig struct {
	SignerKey     string // base64
	SaltSeparator string // base64
	Rounds        int
	MemCost       int
}

// Enabled menandakan konfigurasi scrypt Firebase tersedia.
func (c FbScryptConfig) Enabled() bool { return c.SignerKey != "" }

// VerifyFirebaseScrypt memverifikasi password plaintext terhadap hash+salt hasil
// export Firebase Auth, memakai algoritma scrypt termodifikasi Firebase:
//
//	dk   = scrypt(password, base64(salt)+base64(saltSeparator), N=2^memCost, r=rounds, p=1, dkLen=64)
//	out  = AES-256-CTR(key=dk[:32], iv=0) atas base64(signerKey)
//	cocok bila base64(out) == passwordHash
func VerifyFirebaseScrypt(password, saltB64, hashB64 string, cfg FbScryptConfig) (bool, error) {
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}
	sep, err := base64.StdEncoding.DecodeString(cfg.SaltSeparator)
	if err != nil {
		return false, err
	}
	signer, err := base64.StdEncoding.DecodeString(cfg.SignerKey)
	if err != nil {
		return false, err
	}

	dk, err := scrypt.Key([]byte(password), append(salt, sep...), 1<<cfg.MemCost, cfg.Rounds, 1, 64)
	if err != nil {
		return false, err
	}

	block, err := aes.NewCipher(dk[:32])
	if err != nil {
		return false, err
	}
	out := make([]byte, len(signer))
	cipher.NewCTR(block, make([]byte, aes.BlockSize)).XORKeyStream(out, signer)

	got := base64.StdEncoding.EncodeToString(out)
	return subtle.ConstantTimeCompare([]byte(got), []byte(hashB64)) == 1, nil
}
