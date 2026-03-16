package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.walletService.CalculateBalance(c.UserContext())
	if err != nil {
		return response.InternalServerError(c, "failed to load wallet summary")
	}

	return response.Success(c, fiber.StatusOK, "wallet summary loaded", summary)
}

func (h *WalletHandler) ListTransactions(c *fiber.Ctx) error {
	transactions, err := h.walletService.ListTransactions(c.UserContext(), services.ListWalletTransactionsInput{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	})
	if err != nil {
		return response.InternalServerError(c, "failed to load wallet transactions")
	}

	return response.Success(c, fiber.StatusOK, "wallet transactions loaded", transactions)
}

func (h *WalletHandler) CreateTransaction(c *fiber.Ctx) error {
	request := new(services.CreateWalletTransactionInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "wallet")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.Enum("type", strings.TrimSpace(request.Type), []string{"income", "expense"}, "type must be income or expense")
	details.RequiredMaxLength("category", request.Category, maxWalletCategoryLength, "category is required", "category is too long")
	details.RequiredMaxLength("description", request.Description, maxWalletDescLength, "description is required", "description is too long")
	details.PositiveInt64("amount", request.Amount, "amount must be positive")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	transaction, err := h.walletService.CreateTransaction(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidWalletInput), errors.Is(err, services.ErrInvalidTransactionType):
			return response.InvalidRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to create wallet transaction")
		}
	}

	return response.Success(c, fiber.StatusCreated, "wallet transaction created", transaction)
}
