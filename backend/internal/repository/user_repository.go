package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

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
		SELECT id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, email)

	return scanUser(row)
}

func (r *UserRepository) FindByLoginIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1) OR LOWER(username) = LOWER($1)
		ORDER BY CASE WHEN LOWER(email) = LOWER($1) THEN 0 ELSE 1 END
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, identifier)

	return scanUser(row)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	return scanUser(row)
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
		FROM users
		ORDER BY
			CASE
				WHEN archived_at IS NULL AND is_active = TRUE THEN 0
				WHEN archived_at IS NULL AND is_active = FALSE THEN 1
				ELSE 2
			END,
			joined_at ASC NULLS LAST,
			created_at ASC
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
		SELECT id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
		FROM users
		WHERE is_active = TRUE AND archived_at IS NULL
		ORDER BY joined_at ASC NULLS LAST, created_at ASC
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
	Username     string
	Phone        *string
	PasswordHash string
	Role         string
	IsActive     bool
	JoinedAt     time.Time
}

func (r *UserRepository) Create(ctx context.Context, params CreateUserParams) (*models.User, error) {
	return r.create(ctx, r.db, params)
}

func (r *UserRepository) CreateTx(ctx context.Context, tx *sql.Tx, params CreateUserParams) (*models.User, error) {
	return r.create(ctx, tx, params)
}

func (r *UserRepository) create(ctx context.Context, runner userQueryRunner, params CreateUserParams) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, username, phone, password_hash, role, is_active, joined_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
	`

	row := runner.QueryRowContext(
		ctx,
		query,
		params.Name,
		params.Email,
		params.Username,
		nullableString(params.Phone),
		params.PasswordHash,
		params.Role,
		params.IsActive,
		params.JoinedAt,
	)

	return scanUser(row)
}

type UpdateUserParams struct {
	ID       string
	Name     string
	Email    string
	Username string
	Phone    *string
	Role     string
	IsActive bool
	JoinedAt *time.Time
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

type UserRelationCounts struct {
	WalletTransactions int `json:"wallet_transactions"`
	WifiBills          int `json:"wifi_bills"`
	WifiBillMembers    int `json:"wifi_bill_members"`
	Activities         int `json:"activities"`
	ActivityComments   int `json:"activity_comments"`
	ActivityReactions  int `json:"activity_reactions"`
	FoodClaims         int `json:"food_claims"`
	RiceResponses      int `json:"rice_responses"`
	AuditLogs          int `json:"audit_logs"`
	ImportJobs         int `json:"import_jobs"`
	Notifications      int `json:"notifications"`
	PushSubscriptions  int `json:"push_subscriptions"`
	SharedExpensesPaid int `json:"shared_expenses_paid"`
	SharedExpensesMade int `json:"shared_expenses_created"`
	ProposalsCreated   int `json:"proposals_created"`
	ProposalVotes      int `json:"proposal_votes"`
}

func (c UserRelationCounts) HasAny() bool {
	return c.WalletTransactions > 0 ||
		c.WifiBills > 0 ||
		c.WifiBillMembers > 0 ||
		c.Activities > 0 ||
		c.ActivityComments > 0 ||
		c.ActivityReactions > 0 ||
		c.FoodClaims > 0 ||
		c.RiceResponses > 0 ||
		c.AuditLogs > 0 ||
		c.ImportJobs > 0 ||
		c.Notifications > 0 ||
		c.PushSubscriptions > 0 ||
		c.SharedExpensesPaid > 0 ||
		c.SharedExpensesMade > 0 ||
		c.ProposalsCreated > 0 ||
		c.ProposalVotes > 0
}

func (r *UserRepository) Update(ctx context.Context, params UpdateUserParams) (*models.User, error) {
	return r.update(ctx, r.db, params)
}

func (r *UserRepository) UpdateTx(ctx context.Context, tx *sql.Tx, params UpdateUserParams) (*models.User, error) {
	return r.update(ctx, tx, params)
}

func (r *UserRepository) ArchiveTx(ctx context.Context, tx *sql.Tx, userID string) (*models.User, error) {
	query := `
		UPDATE users
		SET
			is_active = FALSE,
			archived_at = COALESCE(archived_at, NOW()),
			left_at = COALESCE(left_at, NOW()),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
	`

	return scanUser(tx.QueryRowContext(ctx, query, userID))
}

func (r *UserRepository) ReactivateTx(ctx context.Context, tx *sql.Tx, userID string) (*models.User, error) {
	query := `
		UPDATE users
		SET
			is_active = TRUE,
			archived_at = NULL,
			left_at = NULL,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
	`

	return scanUser(tx.QueryRowContext(ctx, query, userID))
}

func (r *UserRepository) DeleteTx(ctx context.Context, tx *sql.Tx, userID string) error {
	result, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
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

func (r *UserRepository) CountRelations(ctx context.Context, userID string) (*UserRelationCounts, error) {
	query := `
		SELECT
			(SELECT COUNT(1)::integer FROM wallet_transactions WHERE created_by = $1),
			(SELECT COUNT(1)::integer FROM wifi_bills WHERE created_by = $1),
			(SELECT COUNT(1)::integer FROM wifi_bill_members WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM activities WHERE user_id = $1 OR created_by = $1),
			(SELECT COUNT(1)::integer FROM activity_comments WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM activity_reactions WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM food_claims WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM rice_responses WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM audit_logs WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM import_jobs WHERE created_by = $1),
			(SELECT COUNT(1)::integer FROM notifications WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM push_subscriptions WHERE user_id = $1),
			(SELECT COUNT(1)::integer FROM shared_expenses WHERE paid_by_user_id = $1),
			(SELECT COUNT(1)::integer FROM shared_expenses WHERE created_by = $1),
			(SELECT COUNT(1)::integer FROM proposals WHERE created_by = $1),
			(SELECT COUNT(1)::integer FROM proposal_votes WHERE user_id = $1)
	`

	counts := &UserRelationCounts{}
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&counts.WalletTransactions,
		&counts.WifiBills,
		&counts.WifiBillMembers,
		&counts.Activities,
		&counts.ActivityComments,
		&counts.ActivityReactions,
		&counts.FoodClaims,
		&counts.RiceResponses,
		&counts.AuditLogs,
		&counts.ImportJobs,
		&counts.Notifications,
		&counts.PushSubscriptions,
		&counts.SharedExpensesPaid,
		&counts.SharedExpensesMade,
		&counts.ProposalsCreated,
		&counts.ProposalVotes,
	); err != nil {
		return nil, err
	}

	return counts, nil
}

func (r *UserRepository) update(ctx context.Context, runner userQueryRunner, params UpdateUserParams) (*models.User, error) {
	query := `
		UPDATE users
		SET
			name = $2,
			email = $3,
			username = $4,
			phone = $5,
			role = $6,
			is_active = $7,
			left_at = CASE
				WHEN archived_at IS NOT NULL THEN COALESCE(left_at, archived_at)
				WHEN $7 = FALSE AND left_at IS NULL THEN NOW()
				WHEN $7 = TRUE THEN NULL
				ELSE left_at
			END,
			joined_at = $8,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
	`

	row := runner.QueryRowContext(
		ctx,
		query,
		params.ID,
		params.Name,
		params.Email,
		params.Username,
		nullableString(params.Phone),
		params.Role,
		params.IsActive,
		nullableTime(params.JoinedAt),
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
		RETURNING id, name, email, username, password_hash, phone, avatar_url, role, is_active, auth_version, joined_at, left_at, archived_at, created_at, updated_at
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
			auth_version = auth_version + 1,
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

func (r *UserRepository) FindAvailableUsername(ctx context.Context, name string, email string) (string, error) {
	return r.findAvailableUsername(ctx, r.db, name, email)
}

func (r *UserRepository) CountActiveAdmins(ctx context.Context) (int, error) {
	const query = `
		SELECT COUNT(1)
		FROM users
		WHERE role = 'admin' AND is_active = TRUE AND archived_at IS NULL
	`

	var count int
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) FindAvailableUsernameTx(ctx context.Context, tx *sql.Tx, name string, email string) (string, error) {
	return r.findAvailableUsername(ctx, tx, name, email)
}

func (r *UserRepository) findAvailableUsername(ctx context.Context, runner userQueryRunner, name string, email string) (string, error) {
	base := buildUsernameBase(name, email)
	for suffix := 1; suffix <= 9999; suffix++ {
		candidate := usernameCandidate(base, suffix)
		exists, err := usernameExists(ctx, runner, candidate)
		if err != nil {
			return "", err
		}

		if !exists {
			return candidate, nil
		}
	}

	return "", errors.New("unable to allocate unique username")
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
	var archivedAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&phone,
		&avatarURL,
		&user.Role,
		&user.IsActive,
		&user.AuthVersion,
		&joinedAt,
		&leftAt,
		&archivedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	user.Phone = nullStringPtr(phone)
	user.AvatarURL = nullStringPtr(avatarURL)
	user.JoinedAt = nullTimePtr(joinedAt)
	user.LeftAt = nullTimePtr(leftAt)
	user.ArchivedAt = nullTimePtr(archivedAt)

	return user, nil
}

func usernameExists(ctx context.Context, runner userQueryRunner, username string) (bool, error) {
	var exists bool
	err := runner.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`,
		username,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func buildUsernameBase(name string, email string) string {
	candidates := []string{
		sanitizeUsernamePart(name),
		sanitizeUsernamePart(strings.Split(strings.TrimSpace(email), "@")[0]),
	}

	for _, candidate := range candidates {
		if candidate != "" {
			return candidate
		}
	}

	return "user"
}

func usernameCandidate(base string, suffix int) string {
	if suffix <= 1 {
		return truncateUsernameBase(base, 0)
	}

	suffixValue := "-" + strconv.Itoa(suffix)
	return fmt.Sprintf("%s%s", truncateUsernameBase(base, len(suffixValue)), suffixValue)
}

func truncateUsernameBase(base string, reserved int) string {
	const maxUsernameLength = 32

	trimmed := strings.Trim(strings.ToLower(strings.TrimSpace(base)), "-")
	if trimmed == "" {
		trimmed = "user"
	}

	limit := maxUsernameLength - reserved
	if limit < 1 {
		limit = 1
	}

	if len(trimmed) <= limit {
		return trimmed
	}

	return strings.Trim(trimmed[:limit], "-")
}

func sanitizeUsernamePart(value string) string {
	var builder strings.Builder
	lastWasHyphen := false

	for _, character := range strings.ToLower(strings.TrimSpace(value)) {
		switch {
		case character >= 'a' && character <= 'z':
			builder.WriteRune(character)
			lastWasHyphen = false
		case character >= '0' && character <= '9':
			builder.WriteRune(character)
			lastWasHyphen = false
		case unicode.IsSpace(character) || character == '-' || character == '_' || character == '.':
			if builder.Len() == 0 || lastWasHyphen {
				continue
			}

			builder.WriteRune('-')
			lastWasHyphen = true
		default:
			if builder.Len() == 0 || lastWasHyphen {
				continue
			}

			builder.WriteRune('-')
			lastWasHyphen = true
		}
	}

	return strings.Trim(builder.String(), "-")
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

func nullableTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  *value,
		Valid: true,
	}
}
