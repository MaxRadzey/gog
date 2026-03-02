// Package secret — HTTP-хендлеры для CRUD секретов: создание, список, получение по id, обновление, удаление.
package secret

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/service"
)

// Create godoc
// @Summary Создать секрет
// @Tags secret
// @Accept json
// @Produce json
// @Param input body CreateRequest true "Новый секрет"
// @Success 201 {object} CreateResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/secret [post]
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

// List godoc
// @Summary Список секретов пользователя
// @Tags secret
// @Produce json
// @Success 200 {array} SecretResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/secret [get]
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

// GetByID godoc
// @Summary Получить секрет по ID
// @Tags secret
// @Produce json
// @Param id path int true "ID секрета"
// @Success 200 {object} SecretResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/secret/{id} [get]
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

// Update godoc
// @Summary Обновить секрет
// @Tags secret
// @Accept json
// @Produce json
// @Param id path int true "ID секрета"
// @Param input body UpdateRequest true "Обновлённые данные секрета"
// @Success 200 {string} string "OK"
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/secret/{id} [put]
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

// Delete godoc
// @Summary Удалить секрет
// @Tags secret
// @Produce json
// @Param id path int true "ID секрета"
// @Success 200 {string} string "OK"
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/secret/{id} [delete]
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
