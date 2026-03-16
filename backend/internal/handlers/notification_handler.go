package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) List(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	items, err := h.notificationService.ListForUser(c.UserContext(), user.ID, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to load notifications", "notifications_failed")
	}

	return response.Success(c, fiber.StatusOK, "notifications loaded", items)
}

func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	request := new(services.MarkNotificationsReadInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid notification payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	count, err := h.notificationService.MarkRead(c.UserContext(), user.ID, *request)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to update notifications", "notifications_read_failed")
	}

	return response.Success(c, fiber.StatusOK, "notifications updated", fiber.Map{
		"updated_count": count,
	})
}
