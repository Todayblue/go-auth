package repository

import (
	"errors"
	"fmt"
	"go-chat/internals/adapters/handler/middleware"
	"go-chat/internals/config"
	"go-chat/internals/core/domain"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *DB) findUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	req := u.db.First(&user, "email = ?", email)
	if req.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *DB) VerifyPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password not matched")
	}
	return nil
}

func (u *DB) LoginUser(email, password string) (*domain.LoginResponse, error) {
	user, err := u.findUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = u.VerifyPassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	accessToken, err := u.generateToken(user, config.JWTAccessTokenSecret, config.AccessTokenExpiredIn)
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.generateToken(user, config.JWTRefreshTokenSecret, config.RefreshTokenExpiredIn)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		CommonModel:  user.CommonModel,
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *DB) CreateUser(email, username, password string) (*domain.User, error) {
	user := &domain.User{}
	existingUserByEmail := u.db.First(&user, "email = ?", email)

	if existingUserByEmail.RowsAffected != 0 {
		return nil, errors.New("user with this email already exists")
	}
	existingUserByUserName := u.db.First(&user, "username = ?", username)

	if existingUserByUserName.RowsAffected != 0 {
		return nil, errors.New("user with this username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password not hashed: %v", err)
	}

	user = &domain.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}
	createUser := u.db.Create(&user)
	if createUser.RowsAffected == 0 {
		return nil, fmt.Errorf("failed to create user: %v", createUser.Error)
	}
	return user, nil
}

func (u *DB) generateToken(user *domain.User, jwtSecret string, duration time.Duration) (string, error) {
	expirationTime := time.Now().UTC().Add(duration)

	claims := middleware.JWTCustomClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    strconv.Itoa(int(user.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (u *DB) RefreshTokens(refreshToken string) (*domain.Tokens, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(refreshToken, &middleware.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWTRefreshTokenSecret), nil
	})

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*middleware.JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	user := &domain.User{}
	result := u.db.First(&user, claims.ID)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	newAccessToken, err := u.generateToken(user, config.JWTAccessTokenSecret, config.AccessTokenExpiredIn)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := u.generateToken(user, config.JWTRefreshTokenSecret, config.RefreshTokenExpiredIn)
	if err != nil {
		return nil, err
	}

	return &domain.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (u *DB) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	result := u.db.First(&user)
	if result.RowsAffected != 0 {
		return nil, errors.New("user not found")
	}

	return user, nil
}
