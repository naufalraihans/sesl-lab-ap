package handler

import (
	"errors"
	"net/http"
	"strconv"

	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

// mapError menerjemahkan error domain usecase ke status HTTP + pesan.
func mapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, usecase.ErrUnauthorized):
		response.Fail(c, http.StatusUnauthorized, "NIM atau password salah", nil)
	case errors.Is(err, usecase.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error(), nil)
	case errors.Is(err, usecase.ErrConflict):
		response.Fail(c, http.StatusConflict, "Data sudah ada / konflik", nil)
	case errors.Is(err, usecase.ErrRegisterClosed):
		response.Fail(c, http.StatusForbidden, err.Error(), nil)
	case errors.Is(err, usecase.ErrAlreadyDone):
		response.Fail(c, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, usecase.ErrTimeUp):
		response.Fail(c, http.StatusForbidden, err.Error(), nil)
	case errors.Is(err, usecase.ErrBadRequest):
		response.Fail(c, http.StatusBadRequest, err.Error(), nil)
	default:
		response.Fail(c, http.StatusInternalServerError, "Terjadi kesalahan server", err.Error())
	}
}

// idParam mengambil parameter path numerik.
func idParam(c *gin.Context, name string) (int, bool) {
	v, err := strconv.Atoi(c.Param(name))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "Parameter "+name+" tidak valid", nil)
		return 0, false
	}
	return v, true
}

// queryIntPtr mengambil query int opsional sebagai *int.
func queryIntPtr(c *gin.Context, name string) *int {
	s := c.Query(name)
	if s == "" {
		return nil
	}
	if v, err := strconv.Atoi(s); err == nil {
		return &v
	}
	return nil
}
