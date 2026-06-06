package handler

import (
	"net/http"

	"lab-ap/internal/delivery/http/middleware"
	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth    *usecase.AuthUsecase
	profile *usecase.ProfileUsecase
}

func NewAuthHandler(a *usecase.AuthUsecase, p *usecase.ProfileUsecase) *AuthHandler {
	return &AuthHandler{auth: a, profile: p}
}

// CekNIM POST /api/auth/cek-nim
func (h *AuthHandler) CekNIM(c *gin.Context) {
	var req dto.CekNIMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "NIM wajib diisi", err.Error())
		return
	}
	res, err := h.auth.CekNIM(req.NIM)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Status NIM", res)
}

// Login POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.auth.Login(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Login berhasil", res)
}

// Register POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Password minimal 6 karakter", err.Error())
		return
	}
	res, err := h.auth.Register(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Registrasi berhasil", res)
}

// Logout POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	h.auth.Logout(middleware.UserID(c))
	response.OK(c, http.StatusOK, "Logout berhasil", nil)
}

// Me GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	u, err := h.profile.Get(middleware.UserID(c))
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Profil", u)
}
