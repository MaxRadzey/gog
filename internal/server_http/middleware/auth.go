package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/auth"
	"github.com/MaxRadzey/gog/internal/server_http/handler"
)

// RequireAuth возвращает Gin middleware: проверяет куку, при валидной куке кладёт userID в контекст и вызывает Next.
// При отсутствии или невалидной куке отвечает 401 и прерывает цепочку.
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
