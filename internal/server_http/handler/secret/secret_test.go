package secret_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"

	"github.com/MaxRadzey/gog/internal/auth"
	"github.com/MaxRadzey/gog/internal/constant"
	"github.com/MaxRadzey/gog/internal/crypto"
	"github.com/MaxRadzey/gog/internal/repository"
	"github.com/MaxRadzey/gog/internal/repository/mocks"
	"github.com/MaxRadzey/gog/internal/server_http"
	"github.com/MaxRadzey/gog/internal/server_http/handler"
	"github.com/MaxRadzey/gog/internal/service"
)

const testCookieSecret = "test-secret"
const testEncKey = "dev-encryption-key-32bytes-long!"

type testSecretContainer struct {
	user   *mocks.MockUserRepository
	secret *mocks.MockSecretRepository
}

func (c *testSecretContainer) UserRepository() repository.UserRepository     { return c.user }
func (c *testSecretContainer) SecretRepository() repository.SecretRepository { return c.secret }

func init() {
	gin.SetMode(gin.TestMode)
}

// setupSecretRouter создаёт роутер приложения с моками; возвращает роутер и MockSecretRepository для ожиданий.
func setupSecretRouter(t *testing.T, ctrl *gomock.Controller) (*gin.Engine, *mocks.MockSecretRepository) {
	t.Helper()
	userRepo := mocks.NewMockUserRepository(ctrl)
	secretRepo := mocks.NewMockSecretRepository(ctrl)
	container := &testSecretContainer{user: userRepo, secret: secretRepo}
	services := service.NewServices(container, testEncKey)
	h := handler.New(services, testCookieSecret)
	return server_http.SetupRouter(h), secretRepo
}

func TestSecret_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, secretRepo := setupSecretRouter(t, ctrl)
	secretRepo.EXPECT().
		Create(gomock.Any(), int64(1), constant.SecretTypeLoginPassword, gomock.Any()).
		Return(int64(1), nil)

	body := `{"secret_type":"login_password","data":{"login":"u","password":"p"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/secret", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(auth.NewCookie(1, testCookieSecret))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"id":1`)) {
		t.Errorf("body = %q, want id:1", rec.Body.String())
	}
}

func TestSecret_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, _ := setupSecretRouter(t, ctrl)

	body := `{"secret_type":"login_password","data":{}}`
	req := httptest.NewRequest(http.MethodPost, "/api/secret", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(auth.NewCookie(1, testCookieSecret))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestSecret_Create_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, _ := setupSecretRouter(t, ctrl)

	body := `{"secret_type":"login_password","data":{"login":"u","password":"p"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/secret", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestSecret_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, secretRepo := setupSecretRouter(t, ctrl)
	secretRepo.EXPECT().
		ListByUserID(gomock.Any(), int64(1)).
		Return([]*repository.Secret{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
	req.AddCookie(auth.NewCookie(1, testCookieSecret))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "[]" {
		t.Errorf("body = %q, want []", rec.Body.String())
	}
}

func TestSecret_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encrypted, err := crypto.Encrypt([]byte(`{"content":"note"}`), []byte(testEncKey))
	if err != nil {
		t.Fatal(err)
	}
	router, secretRepo := setupSecretRouter(t, ctrl)
	secretRepo.EXPECT().
		GetByID(gomock.Any(), int64(1), int64(1)).
		Return(&repository.Secret{
			ID: 1, UserID: 1, SecretType: constant.SecretTypeText, Data: encrypted,
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/secret/1", nil)
	req.AddCookie(auth.NewCookie(1, testCookieSecret))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"secret_type":"text"`)) {
		t.Errorf("body %q should contain secret_type text", rec.Body.String())
	}
}

func TestSecret_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, secretRepo := setupSecretRouter(t, ctrl)
	secretRepo.EXPECT().
		GetByID(gomock.Any(), int64(999), int64(1)).
		Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/secret/999", nil)
	req.AddCookie(auth.NewCookie(1, testCookieSecret))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestSecret_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, secretRepo := setupSecretRouter(t, ctrl)
	secretRepo.EXPECT().
		Delete(gomock.Any(), int64(1), int64(1)).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/secret/1", nil)
	req.AddCookie(auth.NewCookie(1, testCookieSecret))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
