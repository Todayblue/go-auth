package ports

import (
	"go-chat/internals/core/domain"
	"time"
)

type UserService interface {
	GetUserByUsername(username string) (*domain.User, error)
	CreateUser(email, username, password string) (*domain.User, error)
	LoginUser(email, password string) (*domain.LoginResponse, error)
	LogoutUser(refreshToken string) error
	RefreshTokens(refreshToken string) (*domain.LoginResponse, error)
}

type UserRepository interface {
	GetUserByUsername(username string) (*domain.User, error)
	CreateUser(email, username, password string) (*domain.User, error)
	LoginUser(email, password string) (*domain.LoginResponse, error)
	LogoutUser(refreshToken string) error
	RefreshTokens(refreshToken string) (*domain.LoginResponse, error)
}

type BookRepository interface {
	GetBooks() ([]*domain.Book, error)
	CreateBook(title string) (*domain.Book, error)
}

type BookService interface {
	GetBooks() ([]*domain.Book, error)
	CreateBook(title string) (*domain.Book, error)
}

type TokenRepository interface {
	SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(userID string, prevTokenID string) error
}
type TokenService interface {
	SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(userID string, prevTokenID string) error
}

type AuthRepository interface {
	GetUserTokenByID(tokenID string) (string, error)
	GetUserByID(userID string) (*domain.User, error)
}
type AuthService interface {
	GetUserTokenByID(tokenID string) (string, error)
	GetUserByID(userID string) (*domain.User, error)
}
