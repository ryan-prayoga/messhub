package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/services"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid login payload",
		})
	}

	response, err := h.authService.Login(c.UserContext(), *request)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err == services.ErrInvalidCredentials {
			status = fiber.StatusUnauthorized
		}

		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "login success",
		"data":    response,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user")

	return c.JSON(fiber.Map{
		"message": "current user",
		"data":    user,
	})
}
