package user

// RegisterRequest — запрос регистрации (логин не короче 3, пароль не короче 6).
type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest — запрос входа (логин и пароль).
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
