package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, joined_at, left_at, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &models.User{}
	var joinedAt sql.NullTime
	var leftAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&joinedAt,
		&leftAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	user.JoinedAt = nullTimePtr(joinedAt)
	user.LeftAt = nullTimePtr(leftAt)

	return user, nil
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}
