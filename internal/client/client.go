package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"
)

// Client — HTTP-клиент к API GophKeeper.
type Client struct {
	baseURL string
	http    *http.Client
}

// New создаёт клиент.
func New(baseURL string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Jar: jar,
		},
	}, nil
}

// parseErrorResponse читает тело ответа и возвращает ErrBadRequest (4xx) или ErrServer (5xx).
func (c *Client) parseErrorResponse(resp *http.Response) error {
	var respErr errResponse
	_ = json.NewDecoder(resp.Body).Decode(&respErr)
	msg := respErr.Error
	if msg == "" {
		msg = resp.Status
	}
	if resp.StatusCode >= 500 {
		return &ErrServer{StatusCode: resp.StatusCode, Message: msg}
	}
	return &ErrBadRequest{StatusCode: resp.StatusCode, Message: msg}
}

// Register регистрирует пользователя. При успехе сервер выставляет cookie сессии.
func (c *Client) Register(ctx context.Context, login, password string) error {
	body := RegisterRequest{Login: login, Password: password}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/user/register", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return c.parseErrorResponse(resp)
}

// Login аутентифицирует пользователя.
func (c *Client) Login(ctx context.Context, login, password string) error {
	body := LoginRequest{Login: login, Password: password}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/user/login", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return c.parseErrorResponse(resp)
}

// Logout сбрасывает сессию.
func (c *Client) Logout(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/user/logout", nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return c.parseErrorResponse(resp)
}

// CreateSecret создаёт секрет.
func (c *Client) CreateSecret(ctx context.Context, secretType string, data json.RawMessage) (int64, error) {
	body := CreateSecretRequest{SecretType: secretType, Data: data}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/secret", bytes.NewReader(jsonBody))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		var out CreateSecretResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return 0, fmt.Errorf("decode response: %w", err)
		}
		return out.ID, nil
	}
	return 0, c.parseErrorResponse(resp)
}

// ListSecrets возвращает список секретов пользователя.
func (c *Client) ListSecrets(ctx context.Context) ([]SecretResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/secret", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, c.parseErrorResponse(resp)
	}
	var list []SecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return list, nil
}

// GetSecret возвращает секрет по id.
func (c *Client) GetSecret(ctx context.Context, id int64) (*SecretResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/secret/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var out SecretResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		return &out, nil
	}
	return nil, c.parseErrorResponse(resp)
}

// UpdateSecret обновляет данные секрета.
func (c *Client) UpdateSecret(ctx context.Context, id int64, data json.RawMessage) error {
	body := UpdateSecretRequest{Data: data}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/api/secret/"+strconv.FormatInt(id, 10), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return c.parseErrorResponse(resp)
}

// DeleteSecret удаляет секрет.
func (c *Client) DeleteSecret(ctx context.Context, id int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/api/secret/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return c.parseErrorResponse(resp)
}
