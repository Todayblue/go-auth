package domain

import (
	"time"
)

type CommonModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

type User struct {
	CommonModel
	Email    string
	Username string
	Password string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	CommonModel
	Email        string `gorm:"uniqueIndex" json:"email"`
	Username     string `json:"username"`
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Book struct {
	CommonModel
	Title string `json:"title"`
}
