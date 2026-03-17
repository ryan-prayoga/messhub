package repository

import (
	"context"
	"database/sql"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type SharedExpenseRepository struct {
	db *sql.DB
}

type sharedExpenseQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type CreateSharedExpenseParams struct {
	ExpenseDate  string
	Category     string
	Description  string
	Amount       int64
	PaidByUserID string
	Status       string
	Notes        *string
	ProofURL     *string
	CreatedBy    string
}

type UpdateSharedExpenseParams struct {
	ID           string
	ExpenseDate  string
	Category     string
	Description  string
	Amount       int64
	PaidByUserID string
	Status       string
	Notes        *string
	ProofURL     *string
}

func NewSharedExpenseRepository(db *sql.DB) *SharedExpenseRepository {
	return &SharedExpenseRepository{db: db}
}

func (r *SharedExpenseRepository) List(ctx context.Context) ([]models.SharedExpense, error) {
	query := `
		SELECT
			se.id,
			se.expense_date,
			se.category,
			se.description,
			se.amount,
			se.paid_by_user_id,
			paid_by.name,
			se.status,
			se.notes,
			se.proof_url,
			se.created_by,
			creator.name,
			se.created_at,
			se.updated_at
		FROM shared_expenses se
		JOIN users paid_by ON paid_by.id = se.paid_by_user_id
		JOIN users creator ON creator.id = se.created_by
		ORDER BY se.expense_date DESC, se.created_at DESC, se.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.SharedExpense, 0)
	for rows.Next() {
		item, err := scanSharedExpense(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	return items, rows.Err()
}

func (r *SharedExpenseRepository) FindByID(ctx context.Context, expenseID string) (*models.SharedExpense, error) {
	query := `
		SELECT
			se.id,
			se.expense_date,
			se.category,
			se.description,
			se.amount,
			se.paid_by_user_id,
			paid_by.name,
			se.status,
			se.notes,
			se.proof_url,
			se.created_by,
			creator.name,
			se.created_at,
			se.updated_at
		FROM shared_expenses se
		JOIN users paid_by ON paid_by.id = se.paid_by_user_id
		JOIN users creator ON creator.id = se.created_by
		WHERE se.id = $1
		LIMIT 1
	`

	return scanSharedExpense(r.db.QueryRowContext(ctx, query, expenseID))
}

func (r *SharedExpenseRepository) GetSummary(ctx context.Context) (*models.SharedExpenseSummary, error) {
	query := `
		SELECT
			COUNT(1)::integer AS total_count,
			COALESCE(SUM(amount), 0) AS total_amount,
			COUNT(1) FILTER (WHERE status IN ('fronted', 'partially_reimbursed'))::integer AS fronted_count,
			COALESCE(SUM(amount) FILTER (WHERE status IN ('fronted', 'partially_reimbursed')), 0) AS outstanding_amount,
			COALESCE(
				SUM(amount) FILTER (
					WHERE DATE_TRUNC('month', expense_date) = DATE_TRUNC('month', CURRENT_DATE)
				),
				0
			) AS this_month_amount
		FROM shared_expenses
	`

	summary := &models.SharedExpenseSummary{}
	if err := r.db.QueryRowContext(ctx, query).Scan(
		&summary.TotalCount,
		&summary.TotalAmount,
		&summary.FrontedCount,
		&summary.OutstandingAmount,
		&summary.ThisMonthAmount,
	); err != nil {
		return nil, err
	}

	return summary, nil
}

func (r *SharedExpenseRepository) CreateTx(ctx context.Context, tx *sql.Tx, params CreateSharedExpenseParams) (*models.SharedExpense, error) {
	query := `
		WITH inserted AS (
			INSERT INTO shared_expenses (
				expense_date,
				category,
				description,
				amount,
				paid_by_user_id,
				status,
				notes,
				proof_url,
				created_by
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING
				id,
				expense_date,
				category,
				description,
				amount,
				paid_by_user_id,
				status,
				notes,
				proof_url,
				created_by,
				created_at,
				updated_at
		)
		SELECT
			inserted.id,
			inserted.expense_date,
			inserted.category,
			inserted.description,
			inserted.amount,
			inserted.paid_by_user_id,
			paid_by.name,
			inserted.status,
			inserted.notes,
			inserted.proof_url,
			inserted.created_by,
			creator.name,
			inserted.created_at,
			inserted.updated_at
		FROM inserted
		JOIN users paid_by ON paid_by.id = inserted.paid_by_user_id
		JOIN users creator ON creator.id = inserted.created_by
	`

	return scanSharedExpense(tx.QueryRowContext(
		ctx,
		query,
		params.ExpenseDate,
		params.Category,
		params.Description,
		params.Amount,
		params.PaidByUserID,
		params.Status,
		nullableString(params.Notes),
		nullableString(params.ProofURL),
		params.CreatedBy,
	))
}

func (r *SharedExpenseRepository) UpdateTx(ctx context.Context, tx *sql.Tx, params UpdateSharedExpenseParams) (*models.SharedExpense, error) {
	query := `
		WITH updated AS (
			UPDATE shared_expenses
			SET
				expense_date = $2,
				category = $3,
				description = $4,
				amount = $5,
				paid_by_user_id = $6,
				status = $7,
				notes = $8,
				proof_url = $9,
				updated_at = NOW()
			WHERE id = $1
			RETURNING
				id,
				expense_date,
				category,
				description,
				amount,
				paid_by_user_id,
				status,
				notes,
				proof_url,
				created_by,
				created_at,
				updated_at
		)
		SELECT
			updated.id,
			updated.expense_date,
			updated.category,
			updated.description,
			updated.amount,
			updated.paid_by_user_id,
			paid_by.name,
			updated.status,
			updated.notes,
			updated.proof_url,
			updated.created_by,
			creator.name,
			updated.created_at,
			updated.updated_at
		FROM updated
		JOIN users paid_by ON paid_by.id = updated.paid_by_user_id
		JOIN users creator ON creator.id = updated.created_by
	`

	return scanSharedExpense(tx.QueryRowContext(
		ctx,
		query,
		params.ID,
		params.ExpenseDate,
		params.Category,
		params.Description,
		params.Amount,
		params.PaidByUserID,
		params.Status,
		nullableString(params.Notes),
		nullableString(params.ProofURL),
	))
}

func scanSharedExpense(row scanner) (*models.SharedExpense, error) {
	item := &models.SharedExpense{}
	var notes sql.NullString
	var proofURL sql.NullString
	if err := row.Scan(
		&item.ID,
		&item.ExpenseDate,
		&item.Category,
		&item.Description,
		&item.Amount,
		&item.PaidByUserID,
		&item.PaidByUserName,
		&item.Status,
		&notes,
		&proofURL,
		&item.CreatedBy,
		&item.CreatedByName,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	item.Notes = nullStringPtr(notes)
	item.ProofURL = nullStringPtr(proofURL)

	return item, nil
}
