package services

import (
	"go-chat/internals/core/domain"
	"go-chat/internals/core/ports"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) GetUserTokenByID(tokenID string) (string, error) {
	return a.repo.GetUserTokenByID(tokenID)
}

func (a *AuthService) GetUserByID(userID string) (*domain.User, error) {
	return a.repo.GetUserByID(userID)
}
