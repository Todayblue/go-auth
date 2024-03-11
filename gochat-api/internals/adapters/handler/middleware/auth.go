package middleware

import (
	"fmt"
	"go-chat/internals/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	ID       string `json:"Id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(c *fiber.Ctx) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	accessToken := c.Cookies("access_token")

	if accessToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := jwt.ParseWithClaims(accessToken, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the algorithm is expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for verification
		return []byte(config.JWTAccessTokenSecret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid access token",
		})
	}

	return c.Next()
}
