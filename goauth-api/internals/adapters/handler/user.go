package handler

import (
	"go-chat/internals/config"
	"go-chat/internals/core/domain"
	"go-chat/internals/core/ports"
	"time"

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

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.userService.CreateUser(req.Email, req.Username, req.Password)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}

func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.userService.LoginUser(req.Email, req.Password)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	setTokenCookies(c, user.AccessToken, user.RefreshToken, config.AccessTokenExpiredIn, config.RefreshTokenExpiredIn)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}

func (h *UserHandler) LogoutUser(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "refresh token not found")
	}

	if err := h.userService.LogoutUser(refreshToken); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	cookieNames := []string{"refresh_token", "access_token"}
	clearTokenCookies(cookieNames, c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "logout successfully"})
}

func (h *UserHandler) RefreshTokens(c *fiber.Ctx) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "refresh token not provided")
	}

	result, err := h.userService.RefreshTokens(refreshToken)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	setTokenCookies(c, result.AccessToken, result.RefreshToken, config.AccessTokenExpiredIn, config.RefreshTokenExpiredIn)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"userID":        result.ID,
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
		},
	})
}

func sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  "fail",
		"message": message,
	})
}

func setTokenCookies(c *fiber.Ctx, accessToken, refreshToken string, accessTokenExp, refreshTokenExp time.Duration) {
	setTokenCookie(c, "access_token", accessToken, int(accessTokenExp.Minutes())*60)
	setTokenCookie(c, "refresh_token", refreshToken, int(refreshTokenExp.Minutes())*60)
}

func clearTokenCookies(cookieNames []string, c *fiber.Ctx) {
	expiredTime := time.Now().Add(-1 * time.Hour)
	for _, name := range cookieNames {
		c.Cookie(&fiber.Cookie{
			Name:    name,
			Value:   "",
			Expires: expiredTime,
		})
	}
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
