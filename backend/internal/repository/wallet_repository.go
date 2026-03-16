package repository

import (
	"context"
	"database/sql"
	"math"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type WalletRepository struct {
	db *sql.DB
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
			wt.type,
			wt.category,
			wt.amount,
			wt.description,
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
	Type        string
	Category    string
	Amount      int64
	Description string
	CreatedBy   string
}

func (r *WalletRepository) Create(ctx context.Context, params CreateWalletTransactionParams) (*models.WalletTransaction, error) {
	query := `
		INSERT INTO wallet_transactions (type, category, amount, description, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, type, category, amount, description, created_by, created_at, updated_at
	`

	transaction := &models.WalletTransaction{}
	if err := r.db.QueryRowContext(
		ctx,
		query,
		params.Type,
		params.Category,
		params.Amount,
		params.Description,
		params.CreatedBy,
	).Scan(
		&transaction.ID,
		&transaction.Type,
		&transaction.Category,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedBy,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return transaction, nil
}

func scanWalletTransaction(row scanner) (*models.WalletTransaction, error) {
	transaction := &models.WalletTransaction{}
	if err := row.Scan(
		&transaction.ID,
		&transaction.Type,
		&transaction.Category,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedBy,
		&transaction.CreatedByName,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return transaction, nil
}
