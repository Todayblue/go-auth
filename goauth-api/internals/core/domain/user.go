package domain

import (
	"time"

	"github.com/google/uuid"
)

type CommonModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
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
