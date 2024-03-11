package services

import (
	"go-chat/internals/core/domain"
	"go-chat/internals/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) CreateUser(email, username, password string) (*domain.User, error) {
	return u.repo.CreateUser(email, username, password)
}

func (u *UserService) LoginUser(email, password string) (*domain.LoginResponse, error) {
	return u.repo.LoginUser(email, password)
}

func (u *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return u.repo.GetUserByUsername(username)
}

func (u *UserService) RefreshTokens(refreshToken string) (*domain.Tokens, error) {
	return u.repo.RefreshTokens(refreshToken)
}
