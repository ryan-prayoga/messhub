package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
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
		return response.InternalServerError(c, "failed to load settings")
	}

	return response.Success(c, fiber.StatusOK, "settings loaded", settings)
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	request := new(services.UpdateMessSettingsInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "settings")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	updated := false
	if request.MessName != nil {
		updated = true
		details.RequiredMaxLength("mess_name", *request.MessName, maxMessNameLength, "mess_name is required", "mess_name is too long")
	}
	if request.WifiPrice != nil {
		updated = true
		details.PositiveInt64("wifi_price", *request.WifiPrice, "wifi_price must be positive")
	}
	if request.WifiDeadlineDay != nil {
		updated = true
		details.IntRange("wifi_deadline_day", *request.WifiDeadlineDay, 1, 31, "wifi_deadline_day must be between 1 and 31")
	}
	if request.BankAccountName != nil {
		updated = true
		details.RequiredMaxLength("bank_account_name", *request.BankAccountName, maxBankFieldLength, "bank_account_name is required", "bank_account_name is too long")
	}
	if request.BankAccountNumber != nil {
		updated = true
		details.RequiredMaxLength("bank_account_number", *request.BankAccountNumber, maxBankFieldLength, "bank_account_number is required", "bank_account_number is too long")
	}
	if !updated {
		details.Add("request", "at least one field must be provided")
	}
	if details.HasAny() {
		return validationFailed(c, details)
	}

	settings, err := h.settingsService.UpdateSettings(c.UserContext(), user.ID, *request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSettingsInput) {
			return response.InvalidRequest(c, err.Error())
		}

		return response.InternalServerError(c, "failed to update settings")
	}

	return response.Success(c, fiber.StatusOK, "settings updated", settings)
}
