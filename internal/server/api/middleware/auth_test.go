package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/auth"
	"github.com/gin-gonic/gin"
)

const testSecret = "test-secret"

func TestRequireAuth_NoCookie_401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequireAuth(testSecret))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want 401", w.Code)
	}
}

func TestRequireAuth_InvalidCookie_401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequireAuth(testSecret))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(auth.NewCookie(1, "wrong-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want 401", w.Code)
	}
}

func TestRequireAuth_ValidCookie_SetsUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequireAuth(testSecret))
	r.GET("/", func(c *gin.Context) {
		userID, err := handler.GetUserID(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(auth.NewCookie(99, testSecret))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want 200", w.Code)
	}
	// проверяем, что хендлер получил userID из контекста (через ответ)
	body := w.Body.String()
	if body == "" || body == "{}" {
		t.Errorf("expected user_id in response, got %q", body)
	}
}
