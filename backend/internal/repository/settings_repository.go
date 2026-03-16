package repository

import (
	"context"
	"database/sql"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type SettingsRepository struct {
	db *sql.DB
}

type settingsQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type UpsertMessSettingsParams struct {
	MessName          string
	WifiPrice         int64
	WifiDeadlineDay   int
	BankAccountName   string
	BankAccountNumber string
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context) (*models.MessSettings, error) {
	query := `
		SELECT
			id,
			mess_name,
			wifi_price,
			wifi_deadline_day,
			bank_account_name,
			bank_account_number,
			created_at,
			updated_at
		FROM mess_settings
		ORDER BY created_at ASC
		LIMIT 1
	`

	return scanMessSettings(r.db.QueryRowContext(ctx, query))
}

func (r *SettingsRepository) Upsert(ctx context.Context, params UpsertMessSettingsParams) (*models.MessSettings, error) {
	return r.upsert(ctx, r.db, params)
}

func (r *SettingsRepository) UpsertTx(ctx context.Context, tx *sql.Tx, params UpsertMessSettingsParams) (*models.MessSettings, error) {
	return r.upsert(ctx, tx, params)
}

func (r *SettingsRepository) upsert(ctx context.Context, runner settingsQueryRunner, params UpsertMessSettingsParams) (*models.MessSettings, error) {
	query := `
		INSERT INTO mess_settings (
			singleton,
			mess_name,
			wifi_price,
			wifi_deadline_day,
			bank_account_name,
			bank_account_number
		)
		VALUES (TRUE, $1, $2, $3, $4, $5)
		ON CONFLICT (singleton)
		DO UPDATE SET
			mess_name = EXCLUDED.mess_name,
			wifi_price = EXCLUDED.wifi_price,
			wifi_deadline_day = EXCLUDED.wifi_deadline_day,
			bank_account_name = EXCLUDED.bank_account_name,
			bank_account_number = EXCLUDED.bank_account_number,
			updated_at = NOW()
		RETURNING
			id,
			mess_name,
			wifi_price,
			wifi_deadline_day,
			bank_account_name,
			bank_account_number,
			created_at,
			updated_at
	`

	return scanMessSettings(
		runner.QueryRowContext(
			ctx,
			query,
			params.MessName,
			params.WifiPrice,
			params.WifiDeadlineDay,
			params.BankAccountName,
			params.BankAccountNumber,
		),
	)
}

func scanMessSettings(row scanner) (*models.MessSettings, error) {
	settings := &models.MessSettings{}
	if err := row.Scan(
		&settings.ID,
		&settings.MessName,
		&settings.WifiPrice,
		&settings.WifiDeadlineDay,
		&settings.BankAccountName,
		&settings.BankAccountNumber,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return settings, nil
}
