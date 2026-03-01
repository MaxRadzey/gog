package user

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/mock/gomock"

	"github.com/MaxRadzey/gog/internal/auth"
	"github.com/MaxRadzey/gog/internal/repository"
	"github.com/MaxRadzey/gog/internal/repository/mocks"
	"github.com/MaxRadzey/gog/internal/server_http/handler"
	"github.com/MaxRadzey/gog/internal/service"
)

// testRepoContainer реализует repository.RepositoryContainer для тестов.
type testRepoContainer struct {
	repo repository.UserRepository
}

func (c *testRepoContainer) UserRepository() repository.UserRepository {
	return c.repo
}

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		Create(gomock.Any(), "alice", gomock.Any()).
		Return(int64(1), nil)

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container)
	h := handler.New(services, "test-secret")

	body := `{"login":"alice","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ginHandler := Register(h)
	router := setupTestRouter(ginHandler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies count = %d, want 1", len(cookies))
	}
	if cookies[0].Name != auth.CookieName {
		t.Errorf("cookie name = %q, want %q", cookies[0].Name, auth.CookieName)
	}
	if cookies[0].Value == "" {
		t.Error("cookie value is empty")
	}
}

func TestRegister_DuplicateLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		Create(gomock.Any(), "taken", gomock.Any()).
		Return(int64(0), &pgconn.PgError{Code: "23505"})

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container)
	h := handler.New(services, "test-secret")

	body := `{"login":"taken","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ginHandler := Register(h)
	router := setupTestRouter(ginHandler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
	if !strings.Contains(rec.Body.String(), "login already taken") {
		t.Errorf("body %q does not contain %q", rec.Body.String(), "login already taken")
	}
}

func TestRegister_ValidationLoginTooShort(t *testing.T) {
	container := &testRepoContainer{repo: mocks.NewMockUserRepository(gomock.NewController(t))}
	services := service.NewServices(container)
	h := handler.New(services, "test-secret")

	body := `{"login":"ab","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ginHandler := Register(h)
	router := setupTestRouter(ginHandler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	respBody := rec.Body.String()
	if !strings.Contains(respBody, "login") || !strings.Contains(respBody, "at least 3") {
		t.Errorf("body %q missing login/at least 3", respBody)
	}
}

func TestRegister_ValidationPasswordTooShort(t *testing.T) {
	container := &testRepoContainer{repo: mocks.NewMockUserRepository(gomock.NewController(t))}
	services := service.NewServices(container)
	h := handler.New(services, "test-secret")

	body := `{"login":"alice","password":"12345"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ginHandler := Register(h)
	router := setupTestRouter(ginHandler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	respBody := rec.Body.String()
	if !strings.Contains(respBody, "password") || !strings.Contains(respBody, "at least 6") {
		t.Errorf("body %q missing password/at least 6", respBody)
	}
}

func TestRegister_ValidationEmptyLogin(t *testing.T) {
	container := &testRepoContainer{repo: mocks.NewMockUserRepository(gomock.NewController(t))}
	services := service.NewServices(container)
	h := handler.New(services, "test-secret")

	body := `{"login":"","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ginHandler := Register(h)
	router := setupTestRouter(ginHandler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(strings.ToLower(rec.Body.String()), "required") {
		t.Errorf("body %q missing required", rec.Body.String())
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	container := &testRepoContainer{repo: mocks.NewMockUserRepository(gomock.NewController(t))}
	services := service.NewServices(container)
	h := handler.New(services, "test-secret")

	body := `{`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ginHandler := Register(h)
	router := setupTestRouter(ginHandler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid request") {
		t.Errorf("body %q missing invalid request", rec.Body.String())
	}
}

func TestRegister_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		Create(gomock.Any(), "alice", gomock.Any()).
		Return(int64(0), errors.New("db error"))

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container)
	h := handler.New(services, "test-secret")

	body := `{"login":"alice","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ginHandler := Register(h)
	router := setupTestRouter(ginHandler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(rec.Body.String(), "internal error") {
		t.Errorf("body %q missing internal error", rec.Body.String())
	}
}

// setupTestRouter создаёт Gin-роутер с одной ручкой для тестов.
func setupTestRouter(registerHandler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/user/register", registerHandler)
	return r
}
