package user

// RegisterRequest — тело запроса регистрации.
type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse — ответ успешной регистрации.
type RegisterResponse struct {
	UserID int64 `json:"user_id"`
}

// LoginRequest — тело запроса входа (логин + пароль).
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
