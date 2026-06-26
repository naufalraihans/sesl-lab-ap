package hash

import "testing"

// Vector dibuat oleh library referensi `firebase-scrypt` memakai parameter hash
// project nyata. Jika impl Go cocok dgn referensi, password Firebase asli verifiable.
func TestVerifyFirebaseScrypt(t *testing.T) {
	cfg := FbScryptConfig{
		SignerKey:     "WFit0t6dexC2zq+oswaJr+eOhM+0WS6ZWW30iMK9jDhiQQyVbNlYyZLMGKUfC+eH3TwOirSQhZspvWJupfillg==",
		SaltSeparator: "Bw==",
		Rounds:        8,
		MemCost:       14,
	}
	const salt = "bWhzTGFiQVAtc2FsdC0wMQ=="
	// hash dalam base64URL (- dan _), seperti yang dikembalikan Firebase Admin SDK.
	// Memastikan decoder menangani base64url, bukan hanya base64 standar.
	const hash = "n1Q5DOrTUsB993pdPHX9aCYUtz0gvu4p6VGSdmJvLpaOk4LFC88k1LxhIJCqmvUy-iQfVIIeawh_pA3JX_UlVg=="

	ok, err := VerifyFirebaseScrypt("RahasiaMhs123!", salt, hash, cfg)
	if err != nil || !ok {
		t.Fatalf("password benar harus cocok; ok=%v err=%v", ok, err)
	}

	bad, err := VerifyFirebaseScrypt("SalahPassword", salt, hash, cfg)
	if err != nil || bad {
		t.Fatalf("password salah harus ditolak; bad=%v err=%v", bad, err)
	}
}
