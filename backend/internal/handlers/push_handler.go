package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

const (
	maxPushEndpointLength = 2048
	maxPushKeyLength      = 512
)

type PushHandler struct {
	pushService *services.PushService
}

func NewPushHandler(pushService *services.PushService) *PushHandler {
	return &PushHandler{pushService: pushService}
}

func (h *PushHandler) Subscribe(c *fiber.Ctx) error {
	request := new(services.PushSubscriptionInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "push subscription")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredMaxLength("endpoint", request.Endpoint, maxPushEndpointLength, "endpoint is required", "endpoint is too long")
	details.RequiredMaxLength("keys.p256dh", request.Keys.P256DH, maxPushKeyLength, "keys.p256dh is required", "keys.p256dh is too long")
	details.RequiredMaxLength("keys.auth", request.Keys.Auth, maxPushKeyLength, "keys.auth is required", "keys.auth is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	subscription, err := h.pushService.Subscribe(c.UserContext(), user.ID, *request)
	if err != nil {
		return response.InternalServerError(c, "failed to save push subscription")
	}

	return response.Success(c, fiber.StatusCreated, "push subscription saved", subscription)
}

func (h *PushHandler) Unsubscribe(c *fiber.Ctx) error {
	request := new(struct {
		Endpoint string `json:"endpoint"`
	})
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "push subscription")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	endpoint := strings.TrimSpace(request.Endpoint)
	details := validation.NewErrors()
	details.RequiredMaxLength("endpoint", endpoint, maxPushEndpointLength, "endpoint is required", "endpoint is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	removed, err := h.pushService.Unsubscribe(c.UserContext(), user.ID, endpoint)
	if err != nil {
		return response.InternalServerError(c, "failed to remove push subscription")
	}

	return response.Success(c, fiber.StatusOK, "push subscription removed", fiber.Map{
		"removed": removed,
	})
}
