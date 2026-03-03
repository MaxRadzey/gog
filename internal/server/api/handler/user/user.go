// Package user — HTTP-хендлеры пользователей: регистрация, вход, выход.
package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/auth"
	"github.com/MaxRadzey/gog/internal/server/service"
)

// Register godoc
// @Summary Регистрация пользователя
// @Tags user
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные для регистрации"
// @Success 200 {string} string "OK"
// @Failure 400 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/user/register [post]
func Register(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
				e := errs[0]
				field := strings.ToLower(e.Field())
				switch e.Tag() {
				case "required":
					c.JSON(http.StatusBadRequest, gin.H{"error": field + " is required"})
				case "min":
					c.JSON(http.StatusBadRequest, gin.H{"error": field + " must be at least " + e.Param() + " characters"})
				default:
					c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed: " + field})
				}
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		userID, err := h.Services.User.Register(c.Request.Context(), req.Login, req.Password)
		if err != nil {
			var dup *service.ErrDuplicateLogin
			if errors.As(err, &dup) {
				c.JSON(http.StatusConflict, gin.H{"error": "login already taken"})
				return
			}
			var val *service.ErrValidation
			if errors.As(err, &val) {
				c.JSON(http.StatusBadRequest, gin.H{"error": val.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		auth.SetAuthCookie(c.Writer, userID, h.CookieSecret)
		c.Status(http.StatusOK)
	}
}

// Login godoc
// @Summary Вход пользователя
// @Tags user
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Учетные данные"
// @Success 200 {string} string "OK"
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/user/login [post]
func Login(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		u, err := h.Services.User.Authenticate(c.Request.Context(), req.Login, req.Password)
		if err != nil {
			var invalid *service.ErrInvalidCredentials
			if errors.As(err, &invalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		auth.SetAuthCookie(c.Writer, u.ID, h.CookieSecret)
		c.Status(http.StatusOK)
	}
}

// Logout godoc
// @Summary Выход пользователя (очистка сессии)
// @Tags user
// @Produce json
// @Success 200 {string} string "OK"
// @Router /api/user/logout [post]
func Logout(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.ClearAuthCookie(c.Writer)
		c.Status(http.StatusOK)
	}
}
