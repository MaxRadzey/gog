package api

import (
	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/api/handler/secret"
	"github.com/MaxRadzey/gog/internal/server/api/handler/user"
	"github.com/MaxRadzey/gog/internal/server/api/middleware"
	"github.com/MaxRadzey/gog/internal/server/logger"
)

// SetupRouter создаёт Gin-роутер.
func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(logger.RequestLogger())
	r.Use(logger.ResponseLogger())

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/user/register", user.Register(h))
		apiGroup.POST("/user/login", user.Login(h))
		apiGroup.POST("/user/logout", user.Logout(h))

		secretGroup := apiGroup.Group("/secret", middleware.RequireAuth(h.CookieSecret))
		{
			secretGroup.POST("", secret.Create(h))
			secretGroup.GET("", secret.List(h))
			secretGroup.GET("/:id", secret.GetByID(h))
			secretGroup.PUT("/:id", secret.Update(h))
			secretGroup.DELETE("/:id", secret.Delete(h))
		}
	}

	return r
}
