package middleware

import (
	"net/http"
	"sync"
	"time"

	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

// rateLimiter adalah pembatas laju per-IP berbasis fixed window in-memory.
// Dipakai untuk melindungi endpoint sensitif (login/register/cek-nim) dari brute-force.
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string]*window
	limit  int
	window time.Duration
}

type window struct {
	count int
	reset time.Time
}

// RateLimit membuat middleware yang membatasi maksimum `limit` request per `window`
// untuk tiap IP. Entri kedaluwarsa dibersihkan periodik oleh janitor.
func RateLimit(limit int, win time.Duration) gin.HandlerFunc {
	rl := &rateLimiter{
		hits:   make(map[string]*window),
		limit:  limit,
		window: win,
	}
	rl.startJanitor()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
			response.Fail(c, http.StatusTooManyRequests,
				"Terlalu banyak percobaan. Coba lagi beberapa saat.", nil)
			return
		}
		c.Next()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	w, ok := rl.hits[ip]
	if !ok || now.After(w.reset) {
		rl.hits[ip] = &window{count: 1, reset: now.Add(rl.window)}
		return true
	}
	if w.count >= rl.limit {
		return false
	}
	w.count++
	return true
}

// startJanitor membersihkan entri kedaluwarsa secara berkala agar map tidak tumbuh tanpa batas.
func (rl *rateLimiter) startJanitor() {
	go func() {
		t := time.NewTicker(rl.window)
		defer t.Stop()
		for range t.C {
			rl.mu.Lock()
			now := time.Now()
			for ip, w := range rl.hits {
				if now.After(w.reset) {
					delete(rl.hits, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
}
