package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

type userQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, phone, avatar_url, role, is_active, joined_at, left_at, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, email)

	return scanUser(row)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, phone, avatar_url, role, is_active, joined_at, left_at, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	return scanUser(row)
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, name, email, password_hash, phone, avatar_url, role, is_active, joined_at, left_at, created_at, updated_at
		FROM users
		ORDER BY is_active DESC, joined_at ASC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) ListActive(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, name, email, password_hash, phone, avatar_url, role, is_active, joined_at, left_at, created_at, updated_at
		FROM users
		WHERE is_active = TRUE
		ORDER BY joined_at ASC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
}

func (r *UserRepository) Create(ctx context.Context, params CreateUserParams) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash, role, is_active, joined_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, name, email, password_hash, phone, avatar_url, role, is_active, joined_at, left_at, created_at, updated_at
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		params.Name,
		params.Email,
		params.PasswordHash,
		params.Role,
		params.IsActive,
	)

	return scanUser(row)
}

type UpdateUserParams struct {
	ID       string
	Name     string
	Role     string
	IsActive bool
}

type UpdateProfileParams struct {
	ID        string
	Name      string
	Phone     *string
	AvatarURL *string
}

type UpdatePasswordParams struct {
	ID           string
	PasswordHash string
}

func (r *UserRepository) Update(ctx context.Context, params UpdateUserParams) (*models.User, error) {
	return r.update(ctx, r.db, params)
}

func (r *UserRepository) UpdateTx(ctx context.Context, tx *sql.Tx, params UpdateUserParams) (*models.User, error) {
	return r.update(ctx, tx, params)
}

func (r *UserRepository) update(ctx context.Context, runner userQueryRunner, params UpdateUserParams) (*models.User, error) {
	query := `
		UPDATE users
		SET
			name = $2,
			role = $3,
			is_active = $4,
			left_at = CASE
				WHEN $4 = FALSE AND left_at IS NULL THEN NOW()
				WHEN $4 = TRUE THEN NULL
				ELSE left_at
			END,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, password_hash, phone, avatar_url, role, is_active, joined_at, left_at, created_at, updated_at
	`

	row := runner.QueryRowContext(
		ctx,
		query,
		params.ID,
		params.Name,
		params.Role,
		params.IsActive,
	)

	return scanUser(row)
}

func (r *UserRepository) UpdateProfileTx(ctx context.Context, tx *sql.Tx, params UpdateProfileParams) (*models.User, error) {
	query := `
		UPDATE users
		SET
			name = $2,
			phone = $3,
			avatar_url = $4,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, password_hash, phone, avatar_url, role, is_active, joined_at, left_at, created_at, updated_at
	`

	return scanUser(
		tx.QueryRowContext(
			ctx,
			query,
			params.ID,
			params.Name,
			nullableString(params.Phone),
			nullableString(params.AvatarURL),
		),
	)
}

func (r *UserRepository) UpdatePasswordTx(ctx context.Context, tx *sql.Tx, params UpdatePasswordParams) error {
	query := `
		UPDATE users
		SET
			password_hash = $2,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := tx.ExecContext(ctx, query, params.ID, params.PasswordHash)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner) (*models.User, error) {
	user := &models.User{}
	var phone sql.NullString
	var avatarURL sql.NullString
	var joinedAt sql.NullTime
	var leftAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&phone,
		&avatarURL,
		&user.Role,
		&user.IsActive,
		&joinedAt,
		&leftAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	user.Phone = nullStringPtr(phone)
	user.AvatarURL = nullStringPtr(avatarURL)
	user.JoinedAt = nullTimePtr(joinedAt)
	user.LeftAt = nullTimePtr(leftAt)

	return user, nil
}

func nullableString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: trimmed,
		Valid:  true,
	}
}
