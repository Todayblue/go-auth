package ports

import (
	"go-chat/internals/core/domain"
)

type UserService interface {
	GetUserByUsername(username string) (*domain.User, error)
	CreateUser(email, username, password string) (*domain.User, error)
	LoginUser(email, password string) (*domain.LoginResponse, error)
	RefreshTokens(refreshToken string) (*domain.Tokens, error)
}

type UserRepository interface {
	GetUserByUsername(username string) (*domain.User, error)
	CreateUser(email, username, password string) (*domain.User, error)
	LoginUser(email, password string) (*domain.LoginResponse, error)
	RefreshTokens(refreshToken string) (*domain.Tokens, error)
}

type BookRepository interface {
	GetBooks() ([]*domain.Book, error)
	CreateBook(title string) (*domain.Book, error)
}

type BookService interface {
	GetBooks() ([]*domain.Book, error)
	CreateBook(title string) (*domain.Book, error)
}
