package repository

import (
	"errors"
	"go-chat/internals/core/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (a *DB) GetUserTokenByID(tokenID string) (string, error) {
	var userID string
	err := a.cache.Get(tokenID, &userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (a *DB) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	result := a.db.First(&user)
	if result.RowsAffected != 0 {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (a *DB) generateToken(user *domain.User, jwtSecret string, duration time.Duration) (*domain.TokenDetails, error) {
	expirationTime := time.Now().UTC().Add(duration)
	tokenID := uuid.New().String()

	claims := domain.JWTCustomClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	tokenDetails := &domain.TokenDetails{
		Token:     signedToken,
		TokenID:   tokenID,
		UserID:    user.ID.String(),
		ExpiresIn: duration.Nanoseconds(),
	}

	return tokenDetails, nil
}

func (a *DB) storeTokensInCache(userID string, accessTokenID, refreshTokenID string, accessTokenExp, refreshTokenExp time.Duration) error {
	err := a.cache.Set(accessTokenID, userID, accessTokenExp)
	if err != nil {
		return err
	}

	err = a.cache.Set(refreshTokenID, userID, refreshTokenExp)
	if err != nil {
		// If storing refresh token fails, delete the previously stored access token as well
		a.cache.Delete(refreshTokenID)
		return err
	}

	return nil
}

func (a *DB) GetUserByID(userID string) (*domain.User, error) {
	user := &domain.User{}
	result := a.db.First(&user, "id = ?", userID)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}
