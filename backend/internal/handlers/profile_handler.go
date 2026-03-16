package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

type ProfileHandler struct {
	userService *services.UserService
}

func NewProfileHandler(userService *services.UserService) *ProfileHandler {
	return &ProfileHandler{userService: userService}
}

func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	profile, err := h.userService.GetProfile(c.UserContext(), user.ID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), "user_not_found")
		}

		return response.Error(c, fiber.StatusInternalServerError, "failed to load profile", "profile_load_failed")
	}

	return response.Success(c, fiber.StatusOK, "profile loaded", profile)
}

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	request := new(services.UpdateProfileInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid profile payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	profile, err := h.userService.UpdateProfile(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidProfileInput):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_profile_input")
		case errors.Is(err, services.ErrUserNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "user_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to update profile", "profile_update_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "profile updated", profile)
}

func (h *ProfileHandler) ChangePassword(c *fiber.Ctx) error {
	request := new(services.ChangePasswordInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid password payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	if err := h.userService.ChangePassword(c.UserContext(), user.ID, *request); err != nil {
		switch {
		case errors.Is(err, services.ErrCurrentPasswordRequired),
			errors.Is(err, services.ErrNewPasswordRequired),
			errors.Is(err, services.ErrPasswordTooShort):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_password_input")
		case errors.Is(err, services.ErrCurrentPasswordInvalid):
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), "invalid_current_password")
		case errors.Is(err, services.ErrUserNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "user_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to change password", "change_password_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "password changed", fiber.Map{
		"changed": true,
	})
}
