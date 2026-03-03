package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const testSecret = "test-secret"

func TestNewCookie_ValidateCookie_RoundTrip(t *testing.T) {
	cookie := NewCookie(52, testSecret)
	userID, err := ValidateCookie(cookie, testSecret)
	if err != nil {
		t.Fatalf("ValidateCookie: %v", err)
	}
	if userID != 52 {
		t.Errorf("userID: got %d, want 52", userID)
	}
}

func TestValidateCookie_WrongSecret(t *testing.T) {
	cookie := NewCookie(1, "secret-a")
	_, err := ValidateCookie(cookie, "secret-b")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid signature" {
		t.Errorf("got %q", err.Error())
	}
}

func TestValidateCookie_NilCookie(t *testing.T) {
	_, err := ValidateCookie(nil, testSecret)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "empty cookie value" {
		t.Errorf("got %q", err.Error())
	}
}

func TestValidateCookie_EmptyValue(t *testing.T) {
	_, err := ValidateCookie(&http.Cookie{Name: CookieName, Value: ""}, testSecret)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "empty cookie value" {
		t.Errorf("got %q", err.Error())
	}
}

func TestValidateCookie_InvalidFormat(t *testing.T) {
	_, err := ValidateCookie(&http.Cookie{Name: CookieName, Value: "no-dot"}, testSecret)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid cookie format" {
		t.Errorf("got %q", err.Error())
	}
}

func TestValidateCookie_InvalidUserID(t *testing.T) {
	_, err := ValidateCookie(&http.Cookie{Name: CookieName, Value: "abc.signature"}, testSecret)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid user id in cookie" {
		t.Errorf("got %q", err.Error())
	}
}

func TestSetAuthCookie_ValidateFromResponse(t *testing.T) {
	w := httptest.NewRecorder()
	SetAuthCookie(w, 100, testSecret)
	resp := w.Result()
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	userID, err := ValidateCookie(cookies[0], testSecret)
	if err != nil {
		t.Fatalf("ValidateCookie: %v", err)
	}
	if userID != 100 {
		t.Errorf("userID: got %d, want 100", userID)
	}
}

func TestClearAuthCookie_SetsEmptyCookie(t *testing.T) {
	w := httptest.NewRecorder()
	ClearAuthCookie(w)
	resp := w.Result()
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Value != "" || cookies[0].MaxAge != -1 {
		t.Errorf("expected empty value and MaxAge -1, got Value=%q MaxAge=%d", cookies[0].Value, cookies[0].MaxAge)
	}
}
