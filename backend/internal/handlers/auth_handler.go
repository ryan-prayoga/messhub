package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
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
		return invalidPayload(c, "login")
	}

	details := validation.NewErrors()
	details.RequiredString("email", request.Email, "email is required")
	details.Email("email", request.Email, "email must be valid")
	details.RequiredString("password", request.Password, "password is required")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	loginResponse, err := h.authService.Login(c.UserContext(), *request)
	if err != nil {
		status := fiber.StatusInternalServerError
		switch err {
		case services.ErrInvalidLoginInput:
			status = fiber.StatusBadRequest
		case services.ErrInvalidCredentials:
			status = fiber.StatusUnauthorized
		case services.ErrInactiveUser:
			status = fiber.StatusForbidden
		}

		switch status {
		case fiber.StatusBadRequest:
			return response.InvalidRequest(c, err.Error())
		case fiber.StatusUnauthorized:
			return response.Unauthorized(c, "authentication required")
		case fiber.StatusForbidden:
			return response.Forbidden(c, "insufficient permissions")
		default:
			return response.InternalServerError(c, "failed to sign in")
		}
	}

	return response.Success(c, fiber.StatusOK, "login success", loginResponse)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	return response.Success(c, fiber.StatusOK, "current user", user)
}
