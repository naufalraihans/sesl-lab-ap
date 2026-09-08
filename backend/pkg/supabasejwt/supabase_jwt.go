package supabasejwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SupabaseClaims hanya identitas (sub/email). Role TIDAK dari token — resolve DB.
type SupabaseClaims struct {
	Subject string
	Email   string
}

type jwksCache struct {
	mu      sync.RWMutex
	keys    map[string]crypto.PublicKey
	fetched time.Time
	url     string
}

var supaJWKS = &jwksCache{keys: map[string]crypto.PublicKey{}}

const jwksTTL = 10 * time.Minute

type jwkKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type jwksDoc struct {
	Keys []jwkKey `json:"keys"`
}

func (c *jwksCache) refresh(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal ambil JWKS: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS status %d", resp.StatusCode)
	}
	var doc jwksDoc
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return fmt.Errorf("JWKS tidak terbaca: %w", err)
	}
	keys := map[string]crypto.PublicKey{}
	for _, k := range doc.Keys {
		pub, err := jwkToPublicKey(k)
		if err != nil || pub == nil {
			continue
		}
		keys[k.Kid] = pub
	}
	if len(keys) == 0 {
		return errors.New("JWKS tidak berisi kunci yang didukung (EC/RSA)")
	}
	c.mu.Lock()
	c.keys = keys
	c.fetched = time.Now()
	c.url = url
	c.mu.Unlock()
	return nil
}

func (c *jwksCache) keyFor(url, kid string) (crypto.PublicKey, error) {
	c.mu.RLock()
	k, ok := c.keys[kid]
	stale := time.Since(c.fetched) > jwksTTL || c.url != url
	c.mu.RUnlock()
	if ok && !stale {
		return k, nil
	}
	if err := c.refresh(url); err != nil {
		if ok {
			return k, nil
		}
		return nil, err
	}
	c.mu.RLock()
	k, ok = c.keys[kid]
	c.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("kid %q tidak ditemukan di JWKS", kid)
	}
	return k, nil
}

func jwkToPublicKey(k jwkKey) (crypto.PublicKey, error) {
	switch k.Kty {
	case "EC":
		return jwkToECDSA(k)
	case "RSA":
		return jwkToRSA(k)
	default:
		return nil, fmt.Errorf("kty tidak didukung: %q", k.Kty)
	}
}

func jwkToECDSA(k jwkKey) (*ecdsa.PublicKey, error) {
	var curve elliptic.Curve
	switch k.Crv {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("kurva EC tidak didukung: %q", k.Crv)
	}
	xB, err := base64.RawURLEncoding.DecodeString(k.X)
	if err != nil {
		return nil, err
	}
	yB, err := base64.RawURLEncoding.DecodeString(k.Y)
	if err != nil {
		return nil, err
	}
	return &ecdsa.PublicKey{Curve: curve, X: new(big.Int).SetBytes(xB), Y: new(big.Int).SetBytes(yB)}, nil
}

func jwkToRSA(k jwkKey) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nBytes)
	e := 0
	for _, b := range eBytes {
		e = e<<8 | int(b)
	}
	if e == 0 {
		return nil, errors.New("eksponen JWK tidak valid")
	}
	return &rsa.PublicKey{N: n, E: e}, nil
}

// VerifySupabaseToken validasi JWT Supabase ketat (fail-closed).
func VerifySupabaseToken(tokenString, jwksURL, issuer, audience string) (*SupabaseClaims, error) {
	if jwksURL == "" || issuer == "" || audience == "" {
		return nil, errors.New("konfigurasi Supabase Auth (JWKS/issuer/aud) belum lengkap")
	}
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		switch t.Method.(type) {
		case *jwt.SigningMethodECDSA, *jwt.SigningMethodRSA:
		default:
			return nil, fmt.Errorf("metode tanda tangan tak terduga: %v", t.Header["alg"])
		}
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("token tanpa kid")
		}
		return supaJWKS.keyFor(jwksURL, kid)
	}
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, keyfunc,
		jwt.WithValidMethods([]string{"ES256", "ES384", "ES512", "RS256"}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("token Supabase tidak valid: %w", err)
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, errors.New("token tanpa sub")
	}
	email, _ := claims["email"].(string)
	return &SupabaseClaims{Subject: sub, Email: email}, nil
}
