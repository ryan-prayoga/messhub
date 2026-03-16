package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
}

func NewSettingsHandler(settingsService *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: settingsService}
}

func (h *SettingsHandler) Get(c *fiber.Ctx) error {
	settings, err := h.settingsService.GetSettings(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to load settings", "settings_load_failed")
	}

	return response.Success(c, fiber.StatusOK, "settings loaded", settings)
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	request := new(services.UpdateMessSettingsInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid settings payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	settings, err := h.settingsService.UpdateSettings(c.UserContext(), user.ID, *request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSettingsInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_settings_input")
		}

		return response.Error(c, fiber.StatusInternalServerError, "failed to update settings", "settings_update_failed")
	}

	return response.Success(c, fiber.StatusOK, "settings updated", settings)
}
