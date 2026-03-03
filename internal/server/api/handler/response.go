package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondUnauthorized отвечает 401 в едином формате API (JSON).
func RespondUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
}
