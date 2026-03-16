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

	authUser, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredMaxLength("name", request.Name, maxNameLength, "name is required", "name is too long")
	details.RequiredString("email", request.Email, "email is required")
	details.Email("email", request.Email, "email must be valid")
	details.OptionalMaxLength("username", request.Username, maxUsernameLength, "username is too long")
	details.OptionalMaxLength("phone", request.Phone, maxPhoneLength, "phone is too long")
	details.RequiredString("password", request.Password, "password is required")
	details.MinLength("password", request.Password, 8, "password must be at least 8 characters")
	details.Date("joined_at", request.JoinedAt, "2006-01-02", "joined_at must be a valid date")
	details.Enum("role", strings.TrimSpace(request.Role), []string{
		models.RoleAdmin,
		models.RoleTreasurer,
		models.RoleMember,
	}, "role must be admin, treasurer, or member")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	user, err := h.userService.CreateUser(c.UserContext(), authUser.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserInput):
			return response.InvalidRequest(c, "Data anggota belum lengkap atau formatnya belum valid.")
		case errors.Is(err, services.ErrInvalidUsername):
			return response.InvalidRequest(c, "Username hanya boleh berisi huruf, angka, atau pemisah sederhana dan minimal 3 karakter.")
		case errors.Is(err, services.ErrPasswordTooShort):
			return response.InvalidRequest(c, "Password minimal 8 karakter.")
		case errors.Is(err, services.ErrInvalidRole):
			return response.InvalidRequest(c, "Role anggota harus admin, treasurer, atau member.")
		case errors.Is(err, services.ErrUserAlreadyExists):
			return response.Conflict(c, "Email anggota sudah terdaftar.")
		case errors.Is(err, services.ErrUsernameAlreadyExists):
			return response.Conflict(c, "Username sudah dipakai anggota lain.")
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
	if request.Email != nil {
		updated = true
		details.RequiredString("email", *request.Email, "email is required")
		details.Email("email", *request.Email, "email must be valid")
	}
	if request.Username != nil {
		updated = true
		details.OptionalMaxLength("username", request.Username, maxUsernameLength, "username is too long")
	}
	if request.Phone != nil {
		updated = true
		details.OptionalMaxLength("phone", request.Phone, maxPhoneLength, "phone is too long")
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
	if request.JoinedAt != nil {
		updated = true
		details.Date("joined_at", request.JoinedAt, "2006-01-02", "joined_at must be a valid date")
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
		case errors.Is(err, services.ErrInvalidUserInput):
			return response.InvalidRequest(c, "Data anggota belum lengkap atau formatnya belum valid.")
		case errors.Is(err, services.ErrInvalidUsername):
			return response.InvalidRequest(c, "Username hanya boleh berisi huruf, angka, atau pemisah sederhana dan minimal 3 karakter.")
		case errors.Is(err, services.ErrInvalidRole):
			return response.InvalidRequest(c, "Role anggota harus admin, treasurer, atau member.")
		case errors.Is(err, services.ErrUserAlreadyExists):
			return response.Conflict(c, "Email anggota sudah terdaftar.")
		case errors.Is(err, services.ErrUsernameAlreadyExists):
			return response.Conflict(c, "Username sudah dipakai anggota lain.")
		case errors.Is(err, services.ErrSelfDeactivateBlocked):
			return response.InvalidRequest(c, "Admin tidak bisa menonaktifkan akunnya sendiri dari halaman ini.")
		case errors.Is(err, services.ErrSelfDemoteBlocked):
			return response.InvalidRequest(c, "Admin tidak bisa menurunkan role dirinya sendiri dari halaman ini.")
		case errors.Is(err, services.ErrLastAdminRequired):
			return response.InvalidRequest(c, "Minimal satu admin aktif harus tetap tersedia.")
		case errors.Is(err, services.ErrUserNotFound):
			return response.NotFound(c, "Anggota tidak ditemukan.")
		default:
			return response.InternalServerError(c, "failed to update member")
		}
	}

	return response.Success(c, fiber.StatusOK, "member updated", user)
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	request := new(services.AdminResetPasswordInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "password")
	}

	authUser, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredString("new_password", request.NewPassword, "new_password is required")
	details.MinLength("new_password", request.NewPassword, 8, "new_password must be at least 8 characters")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	if err := h.userService.AdminResetPassword(c.UserContext(), authUser.ID, c.Params("id"), *request); err != nil {
		switch {
		case errors.Is(err, services.ErrNewPasswordRequired), errors.Is(err, services.ErrPasswordTooShort):
			return response.InvalidRequest(c, "Password baru minimal 8 karakter.")
		case errors.Is(err, services.ErrUserNotFound):
			return response.NotFound(c, "Anggota tidak ditemukan.")
		default:
			return response.InternalServerError(c, "failed to reset member password")
		}
	}

	return response.Success(c, fiber.StatusOK, "member password reset", fiber.Map{
		"changed": true,
	})
}
