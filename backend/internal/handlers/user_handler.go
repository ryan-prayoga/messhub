package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.userService.ListUsers(c.UserContext())
	if err != nil {
		return response.InternalServerError(c, "failed to load members")
	}

	return response.Success(c, fiber.StatusOK, "members loaded", users)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	request := new(services.CreateUserInput)

	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "user")
	}

	details := validation.NewErrors()
	details.RequiredMaxLength("name", request.Name, maxNameLength, "name is required", "name is too long")
	details.RequiredString("email", request.Email, "email is required")
	details.Email("email", request.Email, "email must be valid")
	details.RequiredString("password", request.Password, "password is required")
	details.MinLength("password", request.Password, 8, "password must be at least 8 characters")
	details.Enum("role", strings.TrimSpace(request.Role), []string{
		models.RoleAdmin,
		models.RoleTreasurer,
		models.RoleMember,
	}, "role must be admin, treasurer, or member")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	user, err := h.userService.CreateUser(c.UserContext(), *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserInput), errors.Is(err, services.ErrPasswordTooShort), errors.Is(err, services.ErrInvalidRole):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrUserAlreadyExists):
			return response.Conflict(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to create member")
		}
	}

	return response.Success(c, fiber.StatusCreated, "member created", user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	request := new(services.UpdateUserInput)

	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "user")
	}

	authUser, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	updated := false
	if request.Name != nil {
		updated = true
		details.RequiredMaxLength("name", *request.Name, maxNameLength, "name is required", "name is too long")
	}
	if request.Role != nil {
		updated = true
		details.Enum("role", strings.TrimSpace(*request.Role), []string{
			models.RoleAdmin,
			models.RoleTreasurer,
			models.RoleMember,
		}, "role must be admin, treasurer, or member")
	}
	if request.IsActive != nil {
		updated = true
	}
	if !updated {
		details.Add("request", "at least one field must be provided")
	}
	if details.HasAny() {
		return validationFailed(c, details)
	}

	user, err := h.userService.UpdateUser(c.UserContext(), authUser.ID, c.Params("id"), *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserInput), errors.Is(err, services.ErrInvalidRole):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrUserNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to update member")
		}
	}

	return response.Success(c, fiber.StatusOK, "member updated", user)
}
