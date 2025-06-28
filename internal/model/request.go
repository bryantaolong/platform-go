package model

type RegisterRequest struct {
	Username string `json:"username" validate:"required" binding:"required"`
	Password string `json:"password" validate:"required,min=6" binding:"required,min=6"`
	Email    string `json:"email" validate:"email" binding:"email"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required" binding:"required"`
	Password string `json:"password" validate:"required" binding:"required"`
}
