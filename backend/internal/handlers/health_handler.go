package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
)

type HealthHandler struct {
	systemService *services.SystemService
}

func NewHealthHandler(systemService *services.SystemService) *HealthHandler {
	return &HealthHandler{systemService: systemService}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	status, err := h.systemService.GetStatus(c.UserContext())
	if err != nil {
		return response.InternalServerError(c, "failed to evaluate health status")
	}

	statusCode := fiber.StatusOK
	message := "service ready"
	if !status.DatabaseReachable {
		statusCode = fiber.StatusServiceUnavailable
		message = "service not ready"
	}

	return response.Success(c, statusCode, message, status)
}
