package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type PushSubscriptionRepository struct {
	db *sql.DB
}

func NewPushSubscriptionRepository(db *sql.DB) *PushSubscriptionRepository {
	return &PushSubscriptionRepository{db: db}
}

type UpsertPushSubscriptionParams struct {
	UserID    string
	Endpoint  string
	P256DHKey string
	AuthKey   string
}

func (r *PushSubscriptionRepository) Upsert(
	ctx context.Context,
	params UpsertPushSubscriptionParams,
) (*models.PushSubscription, error) {
	query := `
		INSERT INTO push_subscriptions (user_id, endpoint, p256dh_key, auth_key)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (endpoint) DO UPDATE
		SET
			user_id = EXCLUDED.user_id,
			p256dh_key = EXCLUDED.p256dh_key,
			auth_key = EXCLUDED.auth_key
		RETURNING id, user_id, endpoint, p256dh_key, auth_key, created_at
	`

	return scanPushSubscription(r.db.QueryRowContext(
		ctx,
		query,
		params.UserID,
		params.Endpoint,
		params.P256DHKey,
		params.AuthKey,
	))
}

func (r *PushSubscriptionRepository) DeleteByUserAndEndpoint(
	ctx context.Context,
	userID string,
	endpoint string,
) (bool, error) {
	query := `
		DELETE FROM push_subscriptions
		WHERE user_id = $1 AND endpoint = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, endpoint)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *PushSubscriptionRepository) DeleteByEndpoint(ctx context.Context, endpoint string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM push_subscriptions WHERE endpoint = $1`, endpoint)
	return err
}

func (r *PushSubscriptionRepository) ListByUserIDs(
	ctx context.Context,
	userIDs []string,
) (map[string][]models.PushSubscription, error) {
	if len(userIDs) == 0 {
		return map[string][]models.PushSubscription{}, nil
	}

	args := make([]any, 0, len(userIDs))
	placeholders := make([]string, 0, len(userIDs))
	for index, userID := range userIDs {
		args = append(args, userID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", index+1))
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, endpoint, p256dh_key, auth_key, created_at
		FROM push_subscriptions
		WHERE user_id IN (%s)
		ORDER BY created_at DESC
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[string][]models.PushSubscription, len(userIDs))
	for rows.Next() {
		item, err := scanPushSubscription(rows)
		if err != nil {
			return nil, err
		}

		items[item.UserID] = append(items[item.UserID], *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func scanPushSubscription(row scanner) (*models.PushSubscription, error) {
	item := &models.PushSubscription{}

	if err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.Endpoint,
		&item.P256DHKey,
		&item.AuthKey,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}
