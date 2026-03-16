package repository

import (
	"context"
	"database/sql"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type AuditLogRepository struct {
	db *sql.DB
}

type CreateAuditLogParams struct {
	UserID     *string
	Action     string
	EntityType string
	EntityID   *string
	OldValue   []byte
	NewValue   []byte
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, params CreateAuditLogParams) (*models.AuditLog, error) {
	return r.create(ctx, r.db, params)
}

func (r *AuditLogRepository) CreateTx(ctx context.Context, tx *sql.Tx, params CreateAuditLogParams) (*models.AuditLog, error) {
	return r.create(ctx, tx, params)
}

func (r *AuditLogRepository) create(ctx context.Context, runner auditQueryRunner, params CreateAuditLogParams) (*models.AuditLog, error) {
	query := `
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_value, new_value)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, action, entity_type, entity_id, old_value, new_value, created_at
	`

	log := &models.AuditLog{}
	if err := runner.QueryRowContext(
		ctx,
		query,
		params.UserID,
		params.Action,
		params.EntityType,
		params.EntityID,
		params.OldValue,
		params.NewValue,
	).Scan(
		&log.ID,
		&log.UserID,
		&log.Action,
		&log.EntityType,
		&log.EntityID,
		&log.OldValue,
		&log.NewValue,
		&log.CreatedAt,
	); err != nil {
		return nil, err
	}

	return log, nil
}

type auditQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
