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
	"golang.org/x/crypto/bcrypt"

	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/auth"
	"github.com/MaxRadzey/gog/internal/server/repository"
	"github.com/MaxRadzey/gog/internal/server/repository/mocks"
	"github.com/MaxRadzey/gog/internal/server/service"
)

// testRepoContainer реализует repository.RepositoryContainer для тестов.
type testRepoContainer struct {
	repo repository.UserRepository
}

func (c *testRepoContainer) UserRepository() repository.UserRepository {
	return c.repo
}

func (c *testRepoContainer) SecretRepository() repository.SecretRepository {
	return nil
}

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		Create(gomock.Any(), "alice", gomock.Any()).
		Return(int64(1), nil)

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
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
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
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
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
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
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
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
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
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
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
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
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
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

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	password := "secret12"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		GetByLogin(gomock.Any(), "alice").
		Return(&repository.User{ID: 1, Login: "alice", PasswordHash: string(hash)}, nil)

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
	h := handler.New(services, "test-secret")

	body := `{"login":"alice","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router := setupTestRouterLogin(Login(h))
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

func TestLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		GetByLogin(gomock.Any(), "alice").
		Return(nil, nil)

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
	h := handler.New(services, "test-secret")

	body := `{"login":"alice","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router := setupTestRouterLogin(Login(h))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(rec.Body.String(), "invalid login or password") {
		t.Errorf("body %q does not contain %q", rec.Body.String(), "invalid login or password")
	}
}

func TestLogin_InvalidCredentialsWrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		GetByLogin(gomock.Any(), "alice").
		Return(&repository.User{ID: 1, Login: "alice", PasswordHash: string(hash)}, nil)

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
	h := handler.New(services, "test-secret")

	body := `{"login":"alice","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router := setupTestRouterLogin(Login(h))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	container := &testRepoContainer{repo: mocks.NewMockUserRepository(gomock.NewController(t))}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
	h := handler.New(services, "test-secret")

	body := `{`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router := setupTestRouterLogin(Login(h))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid request") {
		t.Errorf("body %q missing invalid request", rec.Body.String())
	}
}

func TestLogin_EmptyLogin(t *testing.T) {
	container := &testRepoContainer{repo: mocks.NewMockUserRepository(gomock.NewController(t))}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
	h := handler.New(services, "test-secret")

	body := `{"login":"","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router := setupTestRouterLogin(Login(h))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLogin_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		GetByLogin(gomock.Any(), "alice").
		Return(nil, errors.New("db error"))

	container := &testRepoContainer{repo: mockRepo}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
	h := handler.New(services, "test-secret")

	body := `{"login":"alice","password":"secret12"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router := setupTestRouterLogin(Login(h))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(rec.Body.String(), "internal error") {
		t.Errorf("body %q missing internal error", rec.Body.String())
	}
}

func TestLogout_Success(t *testing.T) {
	container := &testRepoContainer{repo: mocks.NewMockUserRepository(gomock.NewController(t))}
	services := service.NewServices(container, "dev-encryption-key-32bytes-long!")
	h := handler.New(services, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/api/user/logout", nil)
	rec := httptest.NewRecorder()

	router := setupTestRouterLogout(Logout(h))
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
	if cookies[0].MaxAge != -1 {
		t.Errorf("cookie MaxAge = %d, want -1 (deleted)", cookies[0].MaxAge)
	}
}

// setupTestRouter создаёт Gin-роутер с одной ручкой для тестов.
func setupTestRouter(registerHandler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/user/register", registerHandler)
	return r
}

func setupTestRouterLogin(loginHandler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/user/login", loginHandler)
	return r
}

func setupTestRouterLogout(logoutHandler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/user/logout", logoutHandler)
	return r
}
