package users

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct{}

type LoginResponse struct {
	Token    string
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RegisterRequest struct{}

type RegisterResponse struct{}
