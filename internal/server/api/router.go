// Package api — маршруты HTTP API: /api/user (register, login, logout) и /api/secret (CRUD с авторизацией).
package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/MaxRadzey/gog/docs"
	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/api/handler/secret"
	"github.com/MaxRadzey/gog/internal/server/api/handler/user"
	"github.com/MaxRadzey/gog/internal/server/api/middleware"
)

// SetupRouter собирает роутер Gin с recovery, логгером и маршрутами API.
func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.ResponseLogger())

	apiGroup := r.Group("/api")

	// Swagger UI по адресу /api/swagger/index.html.
	apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	return r
}
