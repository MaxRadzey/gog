package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetUserID_Ok(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(KeyUserID, int64(42))

	userID, err := GetUserID(c)
	if err != nil {
		t.Fatalf("GetUserID: %v", err)
	}
	if userID != 42 {
		t.Errorf("userID: got %d, want 42", userID)
	}
}

func TestGetUserID_NotInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	_, err := GetUserID(c)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "user_id not in context" {
		t.Errorf("got %q", err.Error())
	}
}

func TestGetUserID_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(KeyUserID, "not int64")

	_, err := GetUserID(c)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "user_id invalid type" {
		t.Errorf("got %q", err.Error())
	}
}
