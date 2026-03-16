package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
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
		return response.Error(c, fiber.StatusBadRequest, "invalid wifi bill payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	bill, err := h.wifiService.CreateBill(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidWifiBillInput), errors.Is(err, services.ErrInvalidWifiStatus), errors.Is(err, services.ErrWifiNoActiveMembers):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_wifi_bill_input")
		case errors.Is(err, services.ErrDuplicateWifiBill):
			return response.Error(c, fiber.StatusConflict, err.Error(), "wifi_bill_already_exists")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to create wifi bill", "wifi_bill_create_failed")
		}
	}

	return response.Success(c, fiber.StatusCreated, "wifi bill created", bill)
}

func (h *WifiHandler) ListBills(c *fiber.Ctx) error {
	bills, err := h.wifiService.ListBills(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to load wifi bills", "wifi_bills_failed")
	}

	return response.Success(c, fiber.StatusOK, "wifi bills loaded", bills)
}

func (h *WifiHandler) GetBillDetail(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	detail, err := h.wifiService.GetBillDetailForViewer(c.UserContext(), c.Params("id"), user.ID, user.Role)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiBillNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "wifi_bill_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to load wifi bill", "wifi_bill_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi bill loaded", detail)
}

func (h *WifiHandler) GetActiveBill(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	detail, err := h.wifiService.GetActiveBill(c.UserContext(), user.ID, user.Role)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to load active wifi bill", "wifi_active_failed")
	}

	return response.Success(c, fiber.StatusOK, "active wifi bill loaded", detail)
}

func (h *WifiHandler) GetMyBills(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	items, err := h.wifiService.ListMyBills(c.UserContext(), user.ID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to load wifi payments", "wifi_my_failed")
	}

	return response.Success(c, fiber.StatusOK, "wifi payments loaded", items)
}

func (h *WifiHandler) SubmitPaymentProof(c *fiber.Ctx) error {
	request := new(services.SubmitWifiPaymentInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid wifi payment payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	member, err := h.wifiService.SubmitPaymentProof(c.UserContext(), c.Params("id"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiProofRequired), errors.Is(err, services.ErrWifiSubmissionNotAllowed):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_wifi_submission")
		case errors.Is(err, services.ErrWifiBillInactive):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "wifi_bill_inactive")
		case errors.Is(err, services.ErrWifiBillNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "wifi_bill_not_found")
		case errors.Is(err, services.ErrWifiMemberNotFound):
			return response.Error(c, fiber.StatusForbidden, err.Error(), "wifi_member_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to submit wifi payment proof", "wifi_submit_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi payment proof submitted", member)
}

func (h *WifiHandler) VerifyPayment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	member, err := h.wifiService.VerifyPayment(c.UserContext(), c.Params("id"), c.Params("memberId"), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiMemberNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "wifi_member_not_found")
		case errors.Is(err, services.ErrWifiReviewNotAllowed):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "wifi_review_not_allowed")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to verify wifi payment", "wifi_verify_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi payment verified", member)
}

func (h *WifiHandler) RejectPayment(c *fiber.Ctx) error {
	request := new(services.RejectWifiPaymentInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid wifi rejection payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	member, err := h.wifiService.RejectPayment(c.UserContext(), c.Params("id"), c.Params("memberId"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrWifiRejectReasonRequired), errors.Is(err, services.ErrWifiReviewNotAllowed):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_wifi_rejection")
		case errors.Is(err, services.ErrWifiMemberNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "wifi_member_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to reject wifi payment", "wifi_reject_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "wifi payment rejected", member)
}
