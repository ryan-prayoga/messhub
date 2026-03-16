package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	request := new(services.LoginInput)

	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid login payload", "invalid_payload")
	}

	loginResponse, err := h.authService.Login(c.UserContext(), *request)
	if err != nil {
		status := fiber.StatusInternalServerError
		code := "login_failed"
		switch err {
		case services.ErrInvalidLoginInput:
			status = fiber.StatusBadRequest
			code = "invalid_login_input"
		case services.ErrInvalidCredentials:
			status = fiber.StatusUnauthorized
			code = "invalid_credentials"
		case services.ErrInactiveUser:
			status = fiber.StatusForbidden
			code = "user_inactive"
		}

		return response.Error(c, status, err.Error(), code)
	}

	return response.Success(c, fiber.StatusOK, "login success", loginResponse)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	return response.Success(c, fiber.StatusOK, "current user", user)
}
