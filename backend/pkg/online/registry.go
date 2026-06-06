package online

import (
	"sync"
	"time"
)

// Registry adalah registry session in-memory untuk menghitung user online.
// Memanfaatkan sifat backend yang STATEFUL (bukan serverless): setiap aktivitas
// terautentikasi memperbarui entri user; entri kedaluwarsa setelah TTL atau dihapus saat logout.
type Registry struct {
	mu      sync.RWMutex
	entries map[int]entry
	ttl     time.Duration
}

type entry struct {
	role     string
	lastSeen time.Time
}

func NewRegistry(ttl time.Duration) *Registry {
	return &Registry{
		entries: make(map[int]entry),
		ttl:     ttl,
	}
}

// Touch mendaftarkan / menyegarkan entri online untuk seorang user.
func (r *Registry) Touch(userID int, role string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[userID] = entry{role: role, lastSeen: time.Now()}
}

// Remove menghapus entri (dipanggil saat logout).
func (r *Registry) Remove(userID int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, userID)
}

// Count mengembalikan jumlah user online (belum kedaluwarsa), opsional difilter per role.
// role kosong = semua role.
func (r *Registry) Count(role string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cutoff := time.Now().Add(-r.ttl)
	n := 0
	for _, e := range r.entries {
		if e.lastSeen.Before(cutoff) {
			continue
		}
		if role == "" || e.role == role {
			n++
		}
	}
	return n
}

// OnlineUserIDs mengembalikan daftar id user yang online.
func (r *Registry) OnlineUserIDs() []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cutoff := time.Now().Add(-r.ttl)
	ids := make([]int, 0, len(r.entries))
	for id, e := range r.entries {
		if e.lastSeen.After(cutoff) {
			ids = append(ids, id)
		}
	}
	return ids
}

// StartJanitor menjalankan goroutine pembersih entri kedaluwarsa secara periodik.
func (r *Registry) StartJanitor(interval time.Duration, stop <-chan struct{}) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				r.gc()
			}
		}
	}()
}

func (r *Registry) gc() {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := time.Now().Add(-r.ttl)
	for id, e := range r.entries {
		if e.lastSeen.Before(cutoff) {
			delete(r.entries, id)
		}
	}
}
