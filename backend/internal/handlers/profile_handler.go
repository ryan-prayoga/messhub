package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
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
		return response.Unauthorized(c, "authentication required")
	}

	profile, err := h.userService.GetProfile(c.UserContext(), user.ID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return response.NotFound(c, err.Error())
		}

		return response.InternalServerError(c, "failed to load profile")
	}

	return response.Success(c, fiber.StatusOK, "profile loaded", profile)
}

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	request := new(services.UpdateProfileInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "profile")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	updated := false
	if request.Name != nil {
		updated = true
		details.RequiredMaxLength("name", *request.Name, maxNameLength, "name is required", "name is too long")
	}
	if request.Phone != nil {
		updated = true
		details.OptionalMaxLength("phone", request.Phone, maxPhoneLength, "phone is too long")
	}
	if request.AvatarURL != nil {
		updated = true
		details.OptionalMaxLength("avatar_url", request.AvatarURL, maxAvatarURLLength, "avatar_url is too long")
		details.URL("avatar_url", request.AvatarURL, "avatar_url must be a valid URL")
	}
	if !updated {
		details.Add("request", "at least one field must be provided")
	}
	if details.HasAny() {
		return validationFailed(c, details)
	}

	profile, err := h.userService.UpdateProfile(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidProfileInput):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrUserNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to update profile")
		}
	}

	return response.Success(c, fiber.StatusOK, "profile updated", profile)
}

func (h *ProfileHandler) ChangePassword(c *fiber.Ctx) error {
	request := new(services.ChangePasswordInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "password")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredString("current_password", request.CurrentPassword, "current_password is required")
	details.RequiredString("new_password", request.NewPassword, "new_password is required")
	details.MinLength("new_password", request.NewPassword, 8, "new_password must be at least 8 characters")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	if err := h.userService.ChangePassword(c.UserContext(), user.ID, *request); err != nil {
		switch {
		case errors.Is(err, services.ErrCurrentPasswordRequired),
			errors.Is(err, services.ErrNewPasswordRequired),
			errors.Is(err, services.ErrPasswordTooShort):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrCurrentPasswordInvalid):
			return response.Unauthorized(c, err.Error())
		case errors.Is(err, services.ErrUserNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to change password")
		}
	}

	return response.Success(c, fiber.StatusOK, "password changed", fiber.Map{
		"changed": true,
	})
}
