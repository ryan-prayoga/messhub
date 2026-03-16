package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
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
		return response.Error(c, fiber.StatusInternalServerError, "failed to load wallet summary", "wallet_summary_failed")
	}

	return response.Success(c, fiber.StatusOK, "wallet summary loaded", summary)
}

func (h *WalletHandler) ListTransactions(c *fiber.Ctx) error {
	transactions, err := h.walletService.ListTransactions(c.UserContext(), services.ListWalletTransactionsInput{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to load wallet transactions", "wallet_transactions_failed")
	}

	return response.Success(c, fiber.StatusOK, "wallet transactions loaded", transactions)
}

func (h *WalletHandler) CreateTransaction(c *fiber.Ctx) error {
	request := new(services.CreateWalletTransactionInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid wallet payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	transaction, err := h.walletService.CreateTransaction(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidWalletInput), errors.Is(err, services.ErrInvalidTransactionType):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_wallet_input")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to create wallet transaction", "wallet_transaction_create_failed")
		}
	}

	return response.Success(c, fiber.StatusCreated, "wallet transaction created", transaction)
}
