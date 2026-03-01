package secret

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/server_http/handler"
	"github.com/MaxRadzey/gog/internal/service"
)

// Create возвращает обработчик POST /api/secret — создание секрета.
func Create(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := handler.GetUserID(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		id, err := h.Services.Secret.Create(c.Request.Context(), userID, req.SecretType, []byte(req.Data))
		if err != nil {
			var val *service.ErrValidation
			if errors.As(err, &val) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusCreated, CreateResponse{ID: id})
	}
}

// List возвращает обработчик GET /api/secret — список секретов пользователя.
func List(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := handler.GetUserID(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		list, err := h.Services.Secret.ListByUserID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		resp := make([]SecretResponse, 0, len(list))
		for _, s := range list {
			resp = append(resp, secretDTOToResponse(s))
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetByID возвращает обработчик GET /api/secret/:id — один секрет по id.
func GetByID(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := handler.GetUserID(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		s, err := h.Services.Secret.GetByID(c.Request.Context(), id, userID)
		if err != nil {
			if errors.Is(err, service.ErrSecretNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, secretDTOToResponse(s))
	}
}

// Update возвращает обработчик PUT /api/secret/:id — обновление данных секрета.
func Update(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := handler.GetUserID(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		err = h.Services.Secret.Update(c.Request.Context(), id, userID, []byte(req.Data))
		if err != nil {
			if errors.Is(err, service.ErrSecretNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
				return
			}
			var val *service.ErrValidation
			if errors.As(err, &val) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.Status(http.StatusOK)
	}
}

// Delete возвращает обработчик DELETE /api/secret/:id — удаление секрета.
func Delete(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := handler.GetUserID(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		err = h.Services.Secret.Delete(c.Request.Context(), id, userID)
		if err != nil {
			if errors.Is(err, service.ErrSecretNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.Status(http.StatusOK)
	}
}

func secretDTOToResponse(s *service.SecretDTO) SecretResponse {
	return SecretResponse{
		ID:         s.ID,
		UserID:     s.UserID,
		SecretType: s.SecretType,
		Data:       json.RawMessage(s.Data),
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}
