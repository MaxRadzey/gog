package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// KeyUserID — ключ в gin.Context, под которым сохраняется userID (int64) после успешной проверки куки.
const KeyUserID = "user_id"

// GetUserID достаёт userID из контекста (должен быть установлен RequireAuth).
func GetUserID(c *gin.Context) (int64, error) {
	v, ok := c.Get(KeyUserID)
	if !ok {
		return 0, errors.New("user_id not in context")
	}
	userID, ok := v.(int64)
	if !ok {
		return 0, errors.New("user_id invalid type")
	}
	return userID, nil
}
