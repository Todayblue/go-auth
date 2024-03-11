package handler

import (
	"go-chat/internals/config"
	"go-chat/internals/core/domain"
	"go-chat/internals/core/ports"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandlers(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	var req domain.LoginRequest
	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid request body",
		})
	}

	user, err := h.userService.LoginUser(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	setTokenCookie(c, "access_token", user.AccessToken, int(config.AccessTokenExpiredIn.Minutes())*60)
	setTokenCookie(c, "refresh_token", user.RefreshToken, int(config.RefreshTokenExpiredIn.Minutes())*60)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid request body",
		})
	}

	user, err := h.userService.CreateUser(req.Email, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}

func setTokenCookie(c *fiber.Ctx, name, value string, maxAge int) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		HTTPOnly: true,
		Secure:   true,
	})
}

func (h *UserHandler) RefreshTokens(c *fiber.Ctx) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "refresh token not provided",
		})
	}

	// Call your service method to refresh the token
	token, err := h.userService.RefreshTokens(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	setTokenCookie(c, "access_token", token.AccessToken, int(config.AccessTokenExpiredIn.Minutes())*60)
	setTokenCookie(c, "refresh_token", token.RefreshToken, int(config.RefreshTokenExpiredIn.Minutes())*60)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":        "success",
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}
