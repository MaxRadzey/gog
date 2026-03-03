// Package middleware — Gin middleware: проверка авторизации по куке и логирование запросов/ответов.
package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/auth"
)

// RequireAuth проверяет сессионную куку; при успехе кладёт userID в контекст, иначе 401 (JSON).
func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(auth.CookieName)
		if err != nil {
			handler.RespondUnauthorized(c)
			return
		}
		userID, err := auth.ValidateCookie(cookie, secret)
		if err != nil {
			handler.RespondUnauthorized(c)
			return
		}
		c.Set(handler.KeyUserID, userID)
		c.Next()
	}
}
