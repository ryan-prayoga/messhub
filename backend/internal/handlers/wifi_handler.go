package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

type WifiHandler struct {
	wifiService *services.WifiService
}

func NewWifiHandler(wifiService *services.WifiService) *WifiHandler {
	return &WifiHandler{wifiService: wifiService}
}

func (h *WifiHandler) CreateBill(c *fiber.Ctx) error {
	request := new(services.CreateWifiBillInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "wifi bill")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	currentYear := time.Now().Year()
	details.IntRange("month", request.Month, 1, 12, "month must be between 1 and 12")
	details.IntRange("year", request.Year, 2024, currentYear+5, "year is out of range")
	if request.NominalPerPerson != nil {
		details.PositiveInt64("nominal_per_person", *request.NominalPerPerson, "nominal_per_person must be positive")
	}
	if request.Status != "" {
		details.Enum("status", request.Status, []string{
			models.WifiBillStatusDraft,
			models.WifiBillStatusActive,
			models.WifiBillStatusClosed,
		}, "status must be draft, active, or closed")
	}
	if request.DeadlineDate != nil {
		if _, err := time.Parse("2006-01-02", *request.DeadlineDate); err != nil {
			details.Add("deadline_date", "deadline_date must use YYYY-MM-DD format")
		}
	}
	if details.HasAny() {
		return validationFailed(c, details)
	}

	bill, err := h.wifiService.CreateBill(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidWifiBillInput), errors.Is(err, services.ErrInvalidWifiStatus), errors.Is(err, services.ErrWifiNoActiveMembers):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrDuplicateWifiBill):
			return response.Conflict(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to create wifi bill")
		}
	}

	return response.Success(c, fiber.StatusCreated, "wifi bill created", bill)
}

func (h *WifiHandler) ListBills(c *fiber.Ctx) error {
	bills, err := h.wifiService.ListBills(c.UserContext())
	if err != nil {
		return response.InternalServerError(c, "failed to load wifi bills")
	}

	return response.Success(c, fiber.StatusOK, "wifi bills loaded", bills)
}

func (h *WifiHandler) UpdateBillStatus(c *fiber.Ctx) error {
	request := new(services.UpdateWifiBillStatusInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "wifi bill status")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.Enum("status", request.Status, []string{
		models.WifiBillStatusDraft,
		models.WifiBillStatusActive,
		models.WifiBillStatusClosed,
	}, "status must be draft, active, or closed")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	bill, err := h.wifiService.UpdateBillStatus(c.UserContext(), c.Params("id"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidWifiStatus), errors.Is(err, services.ErrWifiReviewNotAllowed):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrAnotherActiveWifiBill):
			return response.Conflict(c, "Masih ada tagihan wifi lain yang aktif. Tutup atau jadikan draft dulu tagihan aktif sebelumnya.")
		case errors.Is(err, services.ErrWifiBillNotFound):
			return response.NotFound(c, "Tagihan wifi tidak ditemukan.")
		default:
			return response.InternalServerError(c, "failed to update wifi bill status")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi bill status updated", bill)
}

func (h *WifiHandler) GetBillDetail(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	detail, err := h.wifiService.GetBillDetailForViewer(c.UserContext(), c.Params("id"), user.ID, user.Role)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiBillNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to load wifi bill")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi bill loaded", detail)
}

func (h *WifiHandler) GetActiveBill(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	detail, err := h.wifiService.GetActiveBill(c.UserContext(), user.ID, user.Role)
	if err != nil {
		return response.InternalServerError(c, "failed to load active wifi bill")
	}

	return response.Success(c, fiber.StatusOK, "active wifi bill loaded", detail)
}

func (h *WifiHandler) GetMyBills(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	items, err := h.wifiService.ListMyBills(c.UserContext(), user.ID)
	if err != nil {
		return response.InternalServerError(c, "failed to load wifi payments")
	}

	return response.Success(c, fiber.StatusOK, "wifi payments loaded", items)
}

func (h *WifiHandler) SubmitPaymentProof(c *fiber.Ctx) error {
	request := new(services.SubmitWifiPaymentInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "wifi payment")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredMaxLength("proof_url", request.ProofURL, maxWifiProofLength, "proof_url is required", "proof_url is too long")
	details.OptionalMaxLength("note", request.Note, maxWifiNoteLength, "note is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	member, err := h.wifiService.SubmitPaymentProof(c.UserContext(), c.Params("id"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiProofRequired), errors.Is(err, services.ErrWifiSubmissionNotAllowed):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrWifiBillInactive):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrWifiBillNotFound):
			return response.NotFound(c, err.Error())
		case errors.Is(err, services.ErrWifiMemberNotFound):
			return response.Forbidden(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to submit wifi payment proof")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi payment proof submitted", member)
}

func (h *WifiHandler) VerifyPayment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	member, err := h.wifiService.VerifyPayment(c.UserContext(), c.Params("id"), c.Params("memberId"), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiBillNotFound):
			return response.NotFound(c, err.Error())
		case errors.Is(err, services.ErrWifiMemberNotFound):
			return response.NotFound(c, err.Error())
		case errors.Is(err, services.ErrWifiReviewNotAllowed):
			return response.InvalidRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to verify wifi payment")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi payment verified", member)
}

func (h *WifiHandler) RejectPayment(c *fiber.Ctx) error {
	request := new(services.RejectWifiPaymentInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "wifi rejection")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredMaxLength("reason", request.Reason, maxRejectReasonLength, "reason is required", "reason is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	member, err := h.wifiService.RejectPayment(c.UserContext(), c.Params("id"), c.Params("memberId"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiRejectReasonRequired), errors.Is(err, services.ErrWifiReviewNotAllowed):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrWifiMemberNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to reject wifi payment")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi payment rejected", member)
}
