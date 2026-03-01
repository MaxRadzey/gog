package server_http

import (
	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/logger"
	"github.com/MaxRadzey/gog/internal/server_http/handler"
	"github.com/MaxRadzey/gog/internal/server_http/handler/secret"
	"github.com/MaxRadzey/gog/internal/server_http/handler/user"
	"github.com/MaxRadzey/gog/internal/server_http/middleware"
)

// SetupRouter создаёт Gin-роутер.
func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(logger.RequestLogger())
	r.Use(logger.ResponseLogger())

	api := r.Group("/api")
	{
		api.POST("/user/register", user.Register(h))
		api.POST("/user/login", user.Login(h))
		api.POST("/user/logout", user.Logout(h))

		secretGroup := api.Group("/secret", middleware.RequireAuth(h.CookieSecret))
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
