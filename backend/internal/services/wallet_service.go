package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

var (
	ErrInvalidWalletInput     = errors.New("invalid wallet input")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

const (
	defaultWalletPageSize = 20
	maxWalletPageSize     = 50
)

type CreateWalletTransactionInput struct {
	Type        string `json:"type"`
	Category    string `json:"category"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type ListWalletTransactionsInput struct {
	Page     int
	PageSize int
}

type WalletTransactionList struct {
	Items      []models.WalletTransaction `json:"items"`
	Pagination WalletPagination           `json:"pagination"`
}

type WalletPagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type WalletService struct {
	walletRepository *repository.WalletRepository
	db               *sql.DB
	auditService     *AuditService
}

func NewWalletService(db *sql.DB, walletRepository *repository.WalletRepository, auditService *AuditService) *WalletService {
	return &WalletService{
		walletRepository: walletRepository,
		db:               db,
		auditService:     auditService,
	}
}

func (s *WalletService) CalculateBalance(ctx context.Context) (*models.WalletSummary, error) {
	return s.walletRepository.GetSummary(ctx)
}

func (s *WalletService) ListTransactions(ctx context.Context, input ListWalletTransactionsInput) (*WalletTransactionList, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}

	pageSize := input.PageSize
	switch {
	case pageSize <= 0:
		pageSize = defaultWalletPageSize
	case pageSize > maxWalletPageSize:
		pageSize = maxWalletPageSize
	}

	items, totalItems, err := s.walletRepository.List(ctx, repository.ListWalletTransactionsParams{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	return &WalletTransactionList{
		Items: items,
		Pagination: WalletPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *WalletService) CreateTransaction(ctx context.Context, createdBy string, input CreateWalletTransactionInput) (*models.WalletTransaction, error) {
	transactionType := strings.TrimSpace(input.Type)
	category := strings.TrimSpace(input.Category)
	description := strings.TrimSpace(input.Description)

	if transactionType != "income" && transactionType != "expense" {
		return nil, ErrInvalidTransactionType
	}

	if category == "" || description == "" || createdBy == "" || input.Amount <= 0 {
		return nil, ErrInvalidWalletInput
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	transaction, err := s.walletRepository.CreateTx(ctx, tx, repository.CreateWalletTransactionParams{
		TransactionDate: time.Now().UTC(),
		Type:            transactionType,
		Category:        category,
		Amount:          input.Amount,
		Description:     description,
		Source:          "manual",
		CreatedBy:       createdBy,
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(createdBy),
		Action:     "wallet_transaction_created",
		EntityType: "wallet_transaction",
		EntityID:   stringPtr(transaction.ID),
		NewValue:   transaction,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}
