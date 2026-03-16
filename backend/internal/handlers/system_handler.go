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
		return response.Error(c, fiber.StatusInternalServerError, "failed to load system status", "system_status_failed")
	}

	return response.Success(c, fiber.StatusOK, "system status loaded", status)
}
