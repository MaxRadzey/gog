package server_http

import (
	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/logger"
	"github.com/MaxRadzey/gog/internal/server_http/handler"
	"github.com/MaxRadzey/gog/internal/server_http/handler/user"
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
	}

	return r
}
