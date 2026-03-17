package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

type SharedExpenseHandler struct {
	service *services.SharedExpenseService
}

func NewSharedExpenseHandler(service *services.SharedExpenseService) *SharedExpenseHandler {
	return &SharedExpenseHandler{service: service}
}

func (h *SharedExpenseHandler) List(c *fiber.Ctx) error {
	items, err := h.service.List(c.UserContext())
	if err != nil {
		return response.InternalServerError(c, "failed to load shared expenses")
	}

	return response.Success(c, fiber.StatusOK, "shared expenses loaded", items)
}

func (h *SharedExpenseHandler) Get(c *fiber.Ctx) error {
	item, err := h.service.GetByID(c.UserContext(), c.Params("id"))
	if err != nil {
		if errors.Is(err, services.ErrSharedExpenseNotFound) {
			return response.NotFound(c, "Pengeluaran bersama tidak ditemukan.")
		}

		return response.InternalServerError(c, "failed to load shared expense")
	}

	return response.Success(c, fiber.StatusOK, "shared expense loaded", item)
}

func (h *SharedExpenseHandler) Create(c *fiber.Ctx) error {
	request := new(services.CreateSharedExpenseInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "shared expense")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredString("expense_date", request.ExpenseDate, "expense_date is required")
	details.Date("expense_date", &request.ExpenseDate, "2006-01-02", "expense_date must use YYYY-MM-DD format")
	details.RequiredMaxLength("category", request.Category, maxWalletCategoryLength, "category is required", "category is too long")
	details.RequiredMaxLength("description", request.Description, maxWalletDescLength, "description is required", "description is too long")
	details.PositiveInt64("amount", request.Amount, "amount must be positive")
	details.RequiredString("paid_by_user_id", request.PaidByUserID, "paid_by_user_id is required")
	details.Enum("status", request.Status, []string{
		models.SharedExpenseStatusPersonal,
		models.SharedExpenseStatusFronted,
		models.SharedExpenseStatusPartiallyReimbursed,
		models.SharedExpenseStatusReimbursed,
	}, "status must be personal, fronted, partially_reimbursed, or reimbursed")
	details.OptionalMaxLength("notes", request.Notes, maxWalletDescLength, "notes is too long")
	details.OptionalMaxLength("proof_url", request.ProofURL, maxWifiProofLength, "proof_url is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.service.Create(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidSharedExpenseInput), errors.Is(err, services.ErrInvalidSharedExpenseStatus):
			return response.InvalidRequest(c, "Data pengeluaran bersama belum valid.")
		default:
			return response.InternalServerError(c, "failed to create shared expense")
		}
	}

	return response.Success(c, fiber.StatusCreated, "shared expense created", item)
}

func (h *SharedExpenseHandler) Update(c *fiber.Ctx) error {
	request := new(services.UpdateSharedExpenseInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "shared expense")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	updated := false
	if request.ExpenseDate != nil {
		updated = true
		details.Date("expense_date", request.ExpenseDate, "2006-01-02", "expense_date must use YYYY-MM-DD format")
	}
	if request.Category != nil {
		updated = true
		details.RequiredMaxLength("category", *request.Category, maxWalletCategoryLength, "category is required", "category is too long")
	}
	if request.Description != nil {
		updated = true
		details.RequiredMaxLength("description", *request.Description, maxWalletDescLength, "description is required", "description is too long")
	}
	if request.Amount != nil {
		updated = true
		details.PositiveInt64("amount", *request.Amount, "amount must be positive")
	}
	if request.PaidByUserID != nil {
		updated = true
		details.RequiredString("paid_by_user_id", *request.PaidByUserID, "paid_by_user_id is required")
	}
	if request.Status != nil {
		updated = true
		details.Enum("status", *request.Status, []string{
			models.SharedExpenseStatusPersonal,
			models.SharedExpenseStatusFronted,
			models.SharedExpenseStatusPartiallyReimbursed,
			models.SharedExpenseStatusReimbursed,
		}, "status must be personal, fronted, partially_reimbursed, or reimbursed")
	}
	if request.Notes != nil {
		updated = true
		details.OptionalMaxLength("notes", request.Notes, maxWalletDescLength, "notes is too long")
	}
	if request.ProofURL != nil {
		updated = true
		details.OptionalMaxLength("proof_url", request.ProofURL, maxWifiProofLength, "proof_url is too long")
	}
	if !updated {
		details.Add("request", "at least one field must be provided")
	}
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.service.Update(c.UserContext(), user.ID, c.Params("id"), *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidSharedExpenseInput), errors.Is(err, services.ErrInvalidSharedExpenseStatus):
			return response.InvalidRequest(c, "Data pengeluaran bersama belum valid.")
		case errors.Is(err, services.ErrSharedExpenseNotFound):
			return response.NotFound(c, "Pengeluaran bersama tidak ditemukan.")
		default:
			return response.InternalServerError(c, "failed to update shared expense")
		}
	}

	return response.Success(c, fiber.StatusOK, "shared expense updated", item)
}
