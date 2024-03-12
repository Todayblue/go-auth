package handler

import (
	"fmt"
	"go-chat/internals/config"
	"go-chat/internals/core/domain"
	"go-chat/internals/core/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandlers(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Middleware(c *fiber.Ctx) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := jwt.ParseWithClaims(accessToken, &domain.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the algorithm is expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for verification
		return []byte(config.JWTAccessTokenSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	claims := token.Claims.(*domain.JWTCustomClaims)
	userID, err := h.authService.GetUserTokenByID(claims.ID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "Token is invalid or session has expired"})
	}

	_, err = h.authService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	if !token.Valid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid access token",
		})
	}

	return c.Next()
}
