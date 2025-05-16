package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Account  string `json:"account" binding:"required,min=6,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account" binding:"required,min=6,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}
