package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "MessHub API is healthy", fiber.Map{
		"status": "ok",
	})
}
