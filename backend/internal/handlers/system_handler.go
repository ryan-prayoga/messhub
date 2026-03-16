package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
)

type SystemHandler struct {
	systemService *services.SystemService
}

func NewSystemHandler(systemService *services.SystemService) *SystemHandler {
	return &SystemHandler{systemService: systemService}
}

func (h *SystemHandler) GetStatus(c *fiber.Ctx) error {
	status, err := h.systemService.GetStatus(c.UserContext())
	if err != nil {
		return response.InternalServerError(c, "failed to load system status")
	}

	statusCode := fiber.StatusOK
	message := "system status loaded"
	if !status.DatabaseReachable {
		statusCode = fiber.StatusServiceUnavailable
		message = "system status degraded"
	}

	return response.Success(c, statusCode, message, status)
}
