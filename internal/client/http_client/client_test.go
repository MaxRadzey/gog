package http_client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetSecret_200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SecretResponse{
			ID:         1,
			SecretType: "text",
		})
	}))
	defer server.Close()

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	resp, err := client.GetSecret(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSecret: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.ID != 1 {
		t.Errorf("ID: got %d", resp.ID)
	}
}

func TestClient_GetSecret_404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "secret not found"})
	}))
	defer server.Close()

	client, _ := New(server.URL)
	_, err := client.GetSecret(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
	var bad *ErrBadRequest
	if !errors.As(err, &bad) {
		t.Fatalf("expected ErrBadRequest, got %T", err)
	}
	if bad.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode: got %d", bad.StatusCode)
	}
	if bad.Message != "secret not found" {
		t.Errorf("Message: got %q", bad.Message)
	}
}

func TestClient_GetSecret_500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
	}))
	defer server.Close()

	client, _ := New(server.URL)
	_, err := client.GetSecret(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
	var srv *ErrServer
	if !errors.As(err, &srv) {
		t.Fatalf("expected ErrServer, got %T", err)
	}
	if srv.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode: got %d", srv.StatusCode)
	}
	if srv.Message != "internal error" {
		t.Errorf("Message: got %q", srv.Message)
	}
}

func TestClient_Register_200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	err := client.Register(context.Background(), "user", "pass")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
}

func TestClient_Register_500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
	}))
	defer server.Close()

	client, _ := New(server.URL)
	err := client.Register(context.Background(), "user", "pass")
	if err == nil {
		t.Fatal("expected error")
	}
	var srv *ErrServer
	if !errors.As(err, &srv) {
		t.Fatalf("expected ErrServer, got %T", err)
	}
	if srv.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode: got %d", srv.StatusCode)
	}
}

func TestClient_Login_200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client, _ := New(server.URL)
	if err := client.Login(context.Background(), "u", "p"); err != nil {
		t.Fatalf("Login: %v", err)
	}
}

func TestClient_Logout_200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client, _ := New(server.URL)
	if err := client.Logout(context.Background()); err != nil {
		t.Fatalf("Logout: %v", err)
	}
}

func TestClient_CreateSecret_201(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(CreateSecretResponse{ID: 42})
	}))
	defer server.Close()
	client, _ := New(server.URL)
	id, err := client.CreateSecret(context.Background(), "text", json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("CreateSecret: %v", err)
	}
	if id != 42 {
		t.Errorf("id: got %d, want 42", id)
	}
}

func TestClient_ListSecrets_200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]SecretResponse{{ID: 1, SecretType: "text"}})
	}))
	defer server.Close()
	client, _ := New(server.URL)
	list, err := client.ListSecrets(context.Background())
	if err != nil {
		t.Fatalf("ListSecrets: %v", err)
	}
	if len(list) != 1 || list[0].ID != 1 {
		t.Errorf("list: got %v", list)
	}
}

func TestClient_UpdateSecret_200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client, _ := New(server.URL)
	if err := client.UpdateSecret(context.Background(), 1, json.RawMessage(`{"content":"updated"}`)); err != nil {
		t.Fatalf("UpdateSecret: %v", err)
	}
}

func TestClient_DeleteSecret_200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client, _ := New(server.URL)
	if err := client.DeleteSecret(context.Background(), 1); err != nil {
		t.Fatalf("DeleteSecret: %v", err)
	}
}
