package middleware

import (
	"net/http"
	"strings"

	"lab-ap/pkg/jwt"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID = "user_id"
	CtxNIM    = "nim"
	CtxRole   = "role"
)

// Auth memverifikasi JWT dan menyetel identitas user di context.
func Auth(jm *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, "Token tidak ditemukan", nil)
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jm.Verify(tokenStr)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "Token tidak valid / kedaluwarsa", err.Error())
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxNIM, claims.NIM)
		c.Set(CtxRole, claims.Role)

		c.Next()
	}
}

// UserID helper mengambil user id dari context.
func UserID(c *gin.Context) int {
	if v, ok := c.Get(CtxUserID); ok {
		if id, ok := v.(int); ok {
			return id
		}
	}
	return 0
}

// Role helper.
func Role(c *gin.Context) string {
	if v, ok := c.Get(CtxRole); ok {
		if r, ok := v.(string); ok {
			return r
		}
	}
	return ""
}
