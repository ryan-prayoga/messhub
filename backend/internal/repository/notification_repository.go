package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

type CreateNotificationParams struct {
	UserID   string
	Title    string
	Message  string
	Type     string
	EntityID *string
}

func (r *NotificationRepository) CreateTx(ctx context.Context, tx *sql.Tx, params CreateNotificationParams) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, title, message, type, entity_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, message, type, entity_id, is_read, created_at
	`

	return scanNotification(tx.QueryRowContext(
		ctx,
		query,
		params.UserID,
		params.Title,
		params.Message,
		params.Type,
		params.EntityID,
	))
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, limit int) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, title, message, type, entity_id, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Notification, 0)
	for rows.Next() {
		item, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *NotificationRepository) CountUnread(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)::integer
		FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
	`

	var count int
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *NotificationRepository) MarkAllReadTx(ctx context.Context, tx *sql.Tx, userID string) ([]string, error) {
	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = $1 AND is_read = FALSE
		RETURNING id
	`

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *NotificationRepository) MarkReadByIDsTx(ctx context.Context, tx *sql.Tx, userID string, ids []string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	args := make([]any, 0, len(ids)+1)
	args = append(args, userID)

	placeholders := make([]string, 0, len(ids))
	for index, id := range ids {
		args = append(args, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d", index+2))
	}

	query := fmt.Sprintf(`
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = $1 AND is_read = FALSE AND id IN (%s)
		RETURNING id
	`, strings.Join(placeholders, ", "))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updated := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		updated = append(updated, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return updated, nil
}

func scanNotification(row scanner) (*models.Notification, error) {
	item := &models.Notification{}
	var entityID sql.NullString

	if err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.Title,
		&item.Message,
		&item.Type,
		&entityID,
		&item.IsRead,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	item.EntityID = nullStringPtr(entityID)

	return item, nil
}
