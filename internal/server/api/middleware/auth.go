package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/auth"
)

// RequireAuth возвращает Gin middleware: проверяет куку, при валидной куке кладёт userID в контекст.
func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(auth.CookieName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateCookie(cookie, secret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(handler.KeyUserID, userID)
		c.Next()
	}
}
