package http_client

import "fmt"

// ErrBadRequest — ошибка запроса (4xx): неверные данные от клиента.
type ErrBadRequest struct {
	StatusCode int
	Message    string
}

func (e *ErrBadRequest) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("bad request: %d", e.StatusCode)
}

// ErrServer — ошибка на стороне сервера (5xx).
// Проверка: errors.As(err, &http_client.ErrServer{}).
type ErrServer struct {
	StatusCode int    // HTTP-код ответа
	Message    string // сообщение из поля "error" в JSON
}

func (e *ErrServer) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("server error: %d", e.StatusCode)
}
