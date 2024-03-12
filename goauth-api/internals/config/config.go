package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                string        `envconfig:"DB_HOST"`
	DBUser                string        `envconfig:"DB_USER"`
	DBPassword            string        `envconfig:"DB_PASSWORD"`
	DBName                string        `envconfig:"DB_NAME"`
	DBPort                string        `envconfig:"DB_PORT"`
	ClientOrigin          string        `envconfig:"CLIENT_ORIGIN"`
	JWTAccessTokenSecret  string        `envconfig:"JWT_ACCESS_TOKEN"`
	AccessTokenExpiredIn  time.Duration `envconfig:"ACCESS_TOKEN_EXPIRED_IN"`
	JWTRefreshTokenSecret string        `envconfig:"JWT_REFRESH_TOKEN"`
	RefreshTokenExpiredIn time.Duration `envconfig:"REFRESH_TOKEN_EXPIRED_IN"`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	config := Config{
		DBHost:                os.Getenv("DB_HOST"),
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBName:                os.Getenv("DB_NAME"),
		DBPort:                os.Getenv("DB_PORT"),
		ClientOrigin:          os.Getenv("CLIENT_ORIGIN"),
		JWTAccessTokenSecret:  os.Getenv("JWT_ACCESS_TOKEN"),
		JWTRefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN"),
	}

	accessTokenExpiredIn, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		return Config{}, err
	}
	config.AccessTokenExpiredIn = accessTokenExpiredIn

	refreshTokenExpiredIn, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		return Config{}, err
	}
	config.RefreshTokenExpiredIn = refreshTokenExpiredIn

	return config, nil
}
