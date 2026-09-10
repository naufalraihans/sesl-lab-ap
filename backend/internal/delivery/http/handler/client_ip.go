package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// clientIP mengambil IP asli klien di balik proxy Vercel.
// Gin c.ClientIP() jatuh ke RemoteAddr (127.0.0.1) karena request masuk lewat
// proxy internal Vercel; IP user sebenarnya ada di header X-Forwarded-For.
// Format XFF: "client, proxy1, proxy2" -> ambil entri PERTAMA (paling kiri).
func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if xr := strings.TrimSpace(c.GetHeader("X-Real-IP")); xr != "" {
		return xr
	}
	return c.ClientIP() // fallback (dev/local)
}
