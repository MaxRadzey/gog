package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// KeyUserID — ключ в gin.Context, в который RequireAuth кладёт userID.
const KeyUserID = "user_id"

// GetUserID возвращает userID из контекста; если нет — ошибка.
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
