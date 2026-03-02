package handler

// ErrorResponse описывает стандартный формат ошибки API.
type ErrorResponse struct {
	Error string `json:"error"`
}
