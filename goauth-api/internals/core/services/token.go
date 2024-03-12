package services

import (
	"go-chat/internals/core/ports"
	"time"
)

type TokenService struct {
	repo ports.TokenRepository
}

func NewTokenService(repo ports.TokenRepository) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (t *TokenService) SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) error {
	return t.repo.SetRefreshToken(userID, tokenID, expiresIn)
}

func (t *TokenService) DeleteRefreshToken(userID string, prevTokenID string) error {
	return t.repo.DeleteRefreshToken(userID, prevTokenID)
}
