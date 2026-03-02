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

// Register возвращает обработчик для POST /api/user/register.
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

// Login возвращает обработчик для POST /api/user/login.
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

// Logout возвращает обработчик для POST /api/user/logout.
func Logout(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.ClearAuthCookie(c.Writer)
		c.Status(http.StatusOK)
	}
}
