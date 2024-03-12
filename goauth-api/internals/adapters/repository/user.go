package repository

import (
	"errors"
	"fmt"
	"go-chat/internals/config"
	"go-chat/internals/core/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *DB) CreateUser(email, username, password string) (*domain.User, error) {
	if err := u.checkExistingUser(email, username); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password not hashed: %v", err)
	}

	user := &domain.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}
	if err := u.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (u *DB) LoginUser(email, password string) (*domain.LoginResponse, error) {
	user, err := u.findUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := u.VerifyPassword(user.Password, password); err != nil {
		return nil, err
	}

	return u.generateAndStoreTokens(user)
}

func (u *DB) LogoutUser(refreshToken string) error {
	claims, err := u.parseRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	userID, err := u.GetUserTokenByID(claims.ID)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	if _, err := u.GetUserByID(userID); err != nil {
		return errors.New("user not found")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return errors.New("refresh token expired")
	}

	if err := u.cache.Delete(claims.ID); err != nil {
		return err
	}

	return nil
}

func (u *DB) RefreshTokens(refreshToken string) (*domain.LoginResponse, error) {
	claims, user, err := u.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if err := u.cache.Delete(claims.ID); err != nil {
		return nil, err
	}

	return u.generateAndStoreTokens(user)
}

func (u *DB) checkExistingUser(email, username string) error {
	user := &domain.User{}
	if err := u.db.First(user, "email = ?", email).Error; err == nil {
		return errors.New("user with this email already exists")
	}
	if err := u.db.First(user, "username = ?", username).Error; err == nil {
		return errors.New("user with this username already exists")
	}
	return nil
}

func (u *DB) findUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	if err := u.db.First(user, "email = ?", email).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *DB) VerifyPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("password not matched")
	}
	return nil
}

func (u *DB) parseRefreshToken(refreshToken string) (*domain.JWTCustomClaims, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(refreshToken, &domain.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTRefreshTokenSecret), nil
	})
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*domain.JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

func (u *DB) validateRefreshToken(refreshToken string) (*domain.JWTCustomClaims, *domain.User, error) {
	claims, err := u.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	userID, err := u.GetUserTokenByID(claims.ID)
	if err != nil {
		return nil, nil, errors.New("invalid refresh token")
	}

	user, err := u.GetUserByID(userID)
	if err != nil {
		return nil, nil, errors.New("user not found")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, nil, errors.New("refresh token expired")
	}

	return claims, user, nil
}

func (u *DB) generateAndStoreTokens(user *domain.User) (*domain.LoginResponse, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	accessTokenDetails, err := u.generateToken(user, config.JWTAccessTokenSecret, config.AccessTokenExpiredIn)
	if err != nil {
		return nil, err
	}

	refreshTokenDetails, err := u.generateToken(user, config.JWTRefreshTokenSecret, config.RefreshTokenExpiredIn)
	if err != nil {
		return nil, err
	}

	if err := u.storeTokensInCache(user.ID.String(), accessTokenDetails.TokenID, refreshTokenDetails.TokenID, time.Duration(accessTokenDetails.ExpiresIn), time.Duration(refreshTokenDetails.ExpiresIn)); err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		CommonModel:  user.CommonModel,
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessTokenDetails.Token,
		RefreshToken: refreshTokenDetails.Token,
	}, nil
}
