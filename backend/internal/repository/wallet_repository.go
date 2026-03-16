package repository

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type WalletRepository struct {
	db *sql.DB
}

type walletQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetSummary(ctx context.Context) (*models.WalletSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE type = 'income'), 0) AS total_income,
			COALESCE(SUM(amount) FILTER (WHERE type = 'expense'), 0) AS total_expense
		FROM wallet_transactions
	`

	summary := &models.WalletSummary{}
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.TotalIncome, &summary.TotalExpense); err != nil {
		return nil, err
	}

	summary.Balance = summary.TotalIncome - summary.TotalExpense

	return summary, nil
}

type ListWalletTransactionsParams struct {
	Page     int
	PageSize int
}

func (r *WalletRepository) List(ctx context.Context, params ListWalletTransactionsParams) ([]models.WalletTransaction, int, error) {
	const countQuery = `SELECT COUNT(*) FROM wallet_transactions`

	var totalItems int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			wt.id,
			wt.transaction_date,
			wt.type,
			wt.category,
			wt.amount,
			wt.description,
			wt.proof_url,
			wt.source,
			wt.import_job_id,
			wt.created_by,
			u.name,
			wt.created_at,
			wt.updated_at
		FROM wallet_transactions wt
		JOIN users u ON u.id = wt.created_by
		ORDER BY wt.created_at DESC, wt.id DESC
		LIMIT $1 OFFSET $2
	`

	offset := int(math.Max(float64((params.Page-1)*params.PageSize), 0))
	rows, err := r.db.QueryContext(ctx, query, params.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	transactions := make([]models.WalletTransaction, 0)
	for rows.Next() {
		transaction, err := scanWalletTransaction(rows)
		if err != nil {
			return nil, 0, err
		}

		transactions = append(transactions, *transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return transactions, totalItems, nil
}

type CreateWalletTransactionParams struct {
	TransactionDate time.Time
	Type            string
	Category        string
	Amount          int64
	Description     string
	ProofURL        *string
	Source          string
	ImportJobID     *string
	CreatedBy       string
}

func (r *WalletRepository) Create(ctx context.Context, params CreateWalletTransactionParams) (*models.WalletTransaction, error) {
	return r.create(ctx, r.db, params)
}

func (r *WalletRepository) CreateTx(ctx context.Context, tx *sql.Tx, params CreateWalletTransactionParams) (*models.WalletTransaction, error) {
	return r.create(ctx, tx, params)
}

func (r *WalletRepository) create(ctx context.Context, runner walletQueryRunner, params CreateWalletTransactionParams) (*models.WalletTransaction, error) {
	query := `
		WITH inserted AS (
			INSERT INTO wallet_transactions (
				transaction_date,
				type,
				category,
				amount,
				description,
				proof_url,
				source,
				import_job_id,
				created_by
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING
				id,
				transaction_date,
				type,
				category,
				amount,
				description,
				proof_url,
				source,
				import_job_id,
				created_by,
				created_at,
				updated_at
		)
		SELECT
			inserted.id,
			inserted.transaction_date,
			inserted.type,
			inserted.category,
			inserted.amount,
			inserted.description,
			inserted.proof_url,
			inserted.source,
			inserted.import_job_id,
			inserted.created_by,
			u.name,
			inserted.created_at,
			inserted.updated_at
		FROM inserted
		JOIN users u ON u.id = inserted.created_by
	`

	return scanWalletTransaction(
		runner.QueryRowContext(
			ctx,
			query,
			params.TransactionDate,
			params.Type,
			params.Category,
			params.Amount,
			params.Description,
			nullableString(params.ProofURL),
			params.Source,
			params.ImportJobID,
			params.CreatedBy,
		),
	)
}

func scanWalletTransaction(row scanner) (*models.WalletTransaction, error) {
	transaction := &models.WalletTransaction{}
	var proofURL sql.NullString
	var importJobID sql.NullString
	if err := row.Scan(
		&transaction.ID,
		&transaction.TransactionDate,
		&transaction.Type,
		&transaction.Category,
		&transaction.Amount,
		&transaction.Description,
		&proofURL,
		&transaction.Source,
		&importJobID,
		&transaction.CreatedBy,
		&transaction.CreatedByName,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	); err != nil {
		return nil, err
	}

	transaction.ProofURL = nullStringPtr(proofURL)
	transaction.ImportJobID = nullStringPtr(importJobID)

	return transaction, nil
}
