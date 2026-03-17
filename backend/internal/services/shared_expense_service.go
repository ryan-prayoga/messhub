package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

var (
	ErrInvalidSharedExpenseInput  = errors.New("invalid shared expense input")
	ErrInvalidSharedExpenseStatus = errors.New("invalid shared expense status")
	ErrSharedExpenseNotFound      = errors.New("shared expense not found")
)

type CreateSharedExpenseInput struct {
	ExpenseDate  string  `json:"expense_date"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
	Amount       int64   `json:"amount"`
	PaidByUserID string  `json:"paid_by_user_id"`
	Status       string  `json:"status"`
	Notes        *string `json:"notes"`
	ProofURL     *string `json:"proof_url"`
}

type UpdateSharedExpenseInput struct {
	ExpenseDate  *string `json:"expense_date"`
	Category     *string `json:"category"`
	Description  *string `json:"description"`
	Amount       *int64  `json:"amount"`
	PaidByUserID *string `json:"paid_by_user_id"`
	Status       *string `json:"status"`
	Notes        *string `json:"notes"`
	ProofURL     *string `json:"proof_url"`
}

type SharedExpenseService struct {
	db                *sql.DB
	sharedExpenseRepo *repository.SharedExpenseRepository
	userRepository    *repository.UserRepository
	auditService      *AuditService
}

func NewSharedExpenseService(
	db *sql.DB,
	sharedExpenseRepo *repository.SharedExpenseRepository,
	userRepository *repository.UserRepository,
	auditService *AuditService,
) *SharedExpenseService {
	return &SharedExpenseService{
		db:                db,
		sharedExpenseRepo: sharedExpenseRepo,
		userRepository:    userRepository,
		auditService:      auditService,
	}
}

func (s *SharedExpenseService) List(ctx context.Context) (*models.SharedExpenseList, error) {
	items, err := s.sharedExpenseRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	summary, err := s.sharedExpenseRepo.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &models.SharedExpenseList{
		Items:   items,
		Summary: *summary,
	}, nil
}

func (s *SharedExpenseService) GetByID(ctx context.Context, expenseID string) (*models.SharedExpense, error) {
	item, err := s.sharedExpenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSharedExpenseNotFound
		}

		return nil, err
	}

	return item, nil
}

func (s *SharedExpenseService) Create(ctx context.Context, actorID string, input CreateSharedExpenseInput) (*models.SharedExpense, error) {
	category := strings.TrimSpace(input.Category)
	description := strings.TrimSpace(input.Description)
	paidByUserID := strings.TrimSpace(input.PaidByUserID)
	status := strings.TrimSpace(input.Status)
	if category == "" || description == "" || paidByUserID == "" || input.Amount <= 0 || strings.TrimSpace(input.ExpenseDate) == "" {
		return nil, ErrInvalidSharedExpenseInput
	}
	if !isValidSharedExpenseStatus(status) {
		return nil, ErrInvalidSharedExpenseStatus
	}

	expenseDate, err := parseDateOnly(input.ExpenseDate)
	if err != nil {
		return nil, ErrInvalidSharedExpenseInput
	}

	if _, err := s.userRepository.FindByID(ctx, paidByUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidSharedExpenseInput
		}
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	item, err := s.sharedExpenseRepo.CreateTx(ctx, tx, repository.CreateSharedExpenseParams{
		ExpenseDate:  expenseDate.Format("2006-01-02"),
		Category:     category,
		Description:  description,
		Amount:       input.Amount,
		PaidByUserID: paidByUserID,
		Status:       status,
		Notes:        normalizeOptionalString(input.Notes),
		ProofURL:     normalizeOptionalString(input.ProofURL),
		CreatedBy:    actorID,
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "create_shared_expense",
		EntityType: "shared_expense",
		EntityID:   stringPtr(item.ID),
		NewValue:   item,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *SharedExpenseService) Update(ctx context.Context, actorID string, expenseID string, input UpdateSharedExpenseInput) (*models.SharedExpense, error) {
	current, err := s.GetByID(ctx, expenseID)
	if err != nil {
		return nil, err
	}

	next := *current
	updated := false

	if input.ExpenseDate != nil {
		date, err := parseDateOnly(strings.TrimSpace(*input.ExpenseDate))
		if err != nil {
			return nil, ErrInvalidSharedExpenseInput
		}
		next.ExpenseDate = date
		updated = true
	}

	if input.Category != nil {
		next.Category = strings.TrimSpace(*input.Category)
		updated = true
	}

	if input.Description != nil {
		next.Description = strings.TrimSpace(*input.Description)
		updated = true
	}

	if input.Amount != nil {
		next.Amount = *input.Amount
		updated = true
	}

	if input.PaidByUserID != nil {
		next.PaidByUserID = strings.TrimSpace(*input.PaidByUserID)
		updated = true
	}

	if input.Status != nil {
		next.Status = strings.TrimSpace(*input.Status)
		updated = true
	}

	if input.Notes != nil {
		next.Notes = normalizeOptionalString(input.Notes)
		updated = true
	}

	if input.ProofURL != nil {
		next.ProofURL = normalizeOptionalString(input.ProofURL)
		updated = true
	}

	if !updated || next.Category == "" || next.Description == "" || next.PaidByUserID == "" || next.Amount <= 0 {
		return nil, ErrInvalidSharedExpenseInput
	}
	if !isValidSharedExpenseStatus(next.Status) {
		return nil, ErrInvalidSharedExpenseStatus
	}

	if _, err := s.userRepository.FindByID(ctx, next.PaidByUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidSharedExpenseInput
		}
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	item, err := s.sharedExpenseRepo.UpdateTx(ctx, tx, repository.UpdateSharedExpenseParams{
		ID:           current.ID,
		ExpenseDate:  next.ExpenseDate.Format("2006-01-02"),
		Category:     next.Category,
		Description:  next.Description,
		Amount:       next.Amount,
		PaidByUserID: next.PaidByUserID,
		Status:       next.Status,
		Notes:        next.Notes,
		ProofURL:     next.ProofURL,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSharedExpenseNotFound
		}
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "update_shared_expense",
		EntityType: "shared_expense",
		EntityID:   stringPtr(item.ID),
		OldValue:   current,
		NewValue:   item,
	}); err != nil {
		return nil, err
	}

	if current.Status != item.Status {
		if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
			UserID:     stringPtr(actorID),
			Action:     "reimburse_status_change",
			EntityType: "shared_expense",
			EntityID:   stringPtr(item.ID),
			OldValue: map[string]any{
				"status": current.Status,
			},
			NewValue: map[string]any{
				"status": item.Status,
			},
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return item, nil
}

func isValidSharedExpenseStatus(status string) bool {
	switch status {
	case models.SharedExpenseStatusPersonal,
		models.SharedExpenseStatusFronted,
		models.SharedExpenseStatusPartiallyReimbursed,
		models.SharedExpenseStatusReimbursed:
		return true
	default:
		return false
	}
}
