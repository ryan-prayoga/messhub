package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
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
		return response.Error(c, fiber.StatusInternalServerError, "failed to load members", "users_list_failed")
	}

	return response.Success(c, fiber.StatusOK, "members loaded", users)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	request := new(services.CreateUserInput)

	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user payload", "invalid_payload")
	}

	user, err := h.userService.CreateUser(c.UserContext(), *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserInput), errors.Is(err, services.ErrPasswordTooShort), errors.Is(err, services.ErrInvalidRole):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_user_input")
		case errors.Is(err, services.ErrUserAlreadyExists):
			return response.Error(c, fiber.StatusConflict, err.Error(), "user_already_exists")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to create member", "user_create_failed")
		}
	}

	return response.Success(c, fiber.StatusCreated, "member created", user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	request := new(services.UpdateUserInput)

	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user payload", "invalid_payload")
	}

	user, err := h.userService.UpdateUser(c.UserContext(), c.Params("id"), *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserInput), errors.Is(err, services.ErrInvalidRole):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_user_input")
		case errors.Is(err, services.ErrUserNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "user_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to update member", "user_update_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "member updated", user)
}
