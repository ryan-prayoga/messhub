package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type ActivityRepository struct {
	db *sql.DB
}

type activityQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

type CreateActivityParams struct {
	Type      string
	Title     string
	Content   string
	Points    int
	UserID    string
	CreatedBy string
	ExpiresAt *time.Time
}

func (r *ActivityRepository) ListActivities(ctx context.Context, limit int) ([]models.Activity, error) {
	query := `
		SELECT
			a.id,
			a.type,
			a.title,
			a.content,
			a.points,
			a.user_id,
			subject.name,
			a.created_by,
			creator.name,
			a.expires_at,
			a.created_at,
			a.updated_at
		FROM activities a
		JOIN users subject ON subject.id = a.user_id
		JOIN users creator ON creator.id = a.created_by
		ORDER BY
			CASE
				WHEN a.expires_at IS NULL THEN 0
				WHEN a.expires_at >= NOW() THEN 0
				ELSE 1
			END,
			a.created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Activity, 0)
	for rows.Next() {
		item, err := scanActivity(rows)
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

func (r *ActivityRepository) FindActivityByID(ctx context.Context, activityID string) (*models.Activity, error) {
	query := `
		SELECT
			a.id,
			a.type,
			a.title,
			a.content,
			a.points,
			a.user_id,
			subject.name,
			a.created_by,
			creator.name,
			a.expires_at,
			a.created_at,
			a.updated_at
		FROM activities a
		JOIN users subject ON subject.id = a.user_id
		JOIN users creator ON creator.id = a.created_by
		WHERE a.id = $1
		LIMIT 1
	`

	return scanActivity(r.db.QueryRowContext(ctx, query, activityID))
}

func (r *ActivityRepository) CreateActivityTx(ctx context.Context, tx *sql.Tx, params CreateActivityParams) (*models.Activity, error) {
	query := `
		INSERT INTO activities (type, title, content, points, user_id, created_by, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, type, title, content, points, user_id, created_by, expires_at, created_at, updated_at
	`

	activity := &models.Activity{}
	var expiresAt sql.NullTime
	if err := tx.QueryRowContext(
		ctx,
		query,
		params.Type,
		params.Title,
		params.Content,
		params.Points,
		params.UserID,
		params.CreatedBy,
		params.ExpiresAt,
	).Scan(
		&activity.ID,
		&activity.Type,
		&activity.Title,
		&activity.Content,
		&activity.Points,
		&activity.UserID,
		&activity.CreatedBy,
		&expiresAt,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	); err != nil {
		return nil, err
	}

	activity.ExpiresAt = nullTimePtr(expiresAt)

	return activity, nil
}

type CreateActivityCommentParams struct {
	ActivityID string
	UserID     string
	Comment    string
}

func (r *ActivityRepository) ListComments(ctx context.Context, activityID string) ([]models.ActivityComment, error) {
	query := `
		SELECT
			c.id,
			c.activity_id,
			c.user_id,
			u.name,
			c.comment,
			c.created_at,
			c.updated_at
		FROM activity_comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.activity_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.ActivityComment, 0)
	for rows.Next() {
		item, err := scanActivityComment(rows)
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

func (r *ActivityRepository) CreateCommentTx(ctx context.Context, tx *sql.Tx, params CreateActivityCommentParams) (*models.ActivityComment, error) {
	query := `
		INSERT INTO activity_comments (activity_id, user_id, comment)
		VALUES ($1, $2, $3)
		RETURNING id, activity_id, user_id, comment, created_at, updated_at
	`

	item := &models.ActivityComment{}
	if err := tx.QueryRowContext(ctx, query, params.ActivityID, params.UserID, params.Comment).Scan(
		&item.ID,
		&item.ActivityID,
		&item.UserID,
		&item.Comment,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

type ToggleActivityReactionParams struct {
	ActivityID   string
	UserID       string
	ReactionType string
}

func (r *ActivityRepository) ToggleReactionTx(ctx context.Context, tx *sql.Tx, params ToggleActivityReactionParams) (bool, error) {
	const findQuery = `
		SELECT id
		FROM activity_reactions
		WHERE activity_id = $1 AND user_id = $2 AND reaction_type = $3
		LIMIT 1
	`

	var existingID string
	err := tx.QueryRowContext(ctx, findQuery, params.ActivityID, params.UserID, params.ReactionType).Scan(&existingID)
	switch {
	case err == nil:
		if _, execErr := tx.ExecContext(ctx, `DELETE FROM activity_reactions WHERE id = $1`, existingID); execErr != nil {
			return false, execErr
		}

		return false, nil
	case err != nil && err != sql.ErrNoRows:
		return false, err
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO activity_reactions (activity_id, user_id, reaction_type) VALUES ($1, $2, $3)`,
		params.ActivityID,
		params.UserID,
		params.ReactionType,
	); err != nil {
		return false, err
	}

	return true, nil
}

func (r *ActivityRepository) ListReactions(ctx context.Context, activityID string, viewerID string) ([]models.ActivityReactionSummary, error) {
	query := `
		SELECT
			reaction_type,
			COUNT(*)::integer AS total,
			BOOL_OR(user_id = $2) AS reacted
		FROM activity_reactions
		WHERE activity_id = $1
		GROUP BY reaction_type
		ORDER BY reaction_type ASC
	`

	rows, err := r.db.QueryContext(ctx, query, activityID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.ActivityReactionSummary, 0)
	for rows.Next() {
		item := models.ActivityReactionSummary{}
		if err := rows.Scan(&item.ReactionType, &item.Count, &item.Reacted); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

type CreateFoodClaimParams struct {
	ActivityID string
	UserID     string
}

func (r *ActivityRepository) CreateFoodClaimTx(ctx context.Context, tx *sql.Tx, params CreateFoodClaimParams) (*models.FoodClaim, error) {
	query := `
		INSERT INTO food_claims (activity_id, user_id)
		VALUES ($1, $2)
		RETURNING id, activity_id, user_id, created_at
	`

	item := &models.FoodClaim{}
	if err := tx.QueryRowContext(ctx, query, params.ActivityID, params.UserID).Scan(
		&item.ID,
		&item.ActivityID,
		&item.UserID,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ActivityRepository) ListFoodClaims(ctx context.Context, activityID string) ([]models.FoodClaim, error) {
	query := `
		SELECT
			fc.id,
			fc.activity_id,
			fc.user_id,
			u.name,
			fc.created_at
		FROM food_claims fc
		JOIN users u ON u.id = fc.user_id
		WHERE fc.activity_id = $1
		ORDER BY fc.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.FoodClaim, 0)
	for rows.Next() {
		item, err := scanFoodClaim(rows)
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

type CreateRiceResponseParams struct {
	ActivityID string
	UserID     string
}

func (r *ActivityRepository) CreateRiceResponseTx(ctx context.Context, tx *sql.Tx, params CreateRiceResponseParams) (*models.RiceResponse, error) {
	query := `
		INSERT INTO rice_responses (activity_id, user_id)
		VALUES ($1, $2)
		RETURNING id, activity_id, user_id, created_at
	`

	item := &models.RiceResponse{}
	if err := tx.QueryRowContext(ctx, query, params.ActivityID, params.UserID).Scan(
		&item.ID,
		&item.ActivityID,
		&item.UserID,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ActivityRepository) ListRiceResponses(ctx context.Context, activityID string) ([]models.RiceResponse, error) {
	query := `
		SELECT
			rr.id,
			rr.activity_id,
			rr.user_id,
			u.name,
			rr.created_at
		FROM rice_responses rr
		JOIN users u ON u.id = rr.user_id
		WHERE rr.activity_id = $1
		ORDER BY rr.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.RiceResponse, 0)
	for rows.Next() {
		item, err := scanRiceResponse(rows)
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

func (r *ActivityRepository) GetContributionLeaderboard(ctx context.Context, start *time.Time, end *time.Time, limit int) ([]models.ContributionLeaderboardEntry, error) {
	var args []any
	conditions := []string{"a.type = 'contribution'"}

	if start != nil {
		args = append(args, *start)
		conditions = append(conditions, fmt.Sprintf("a.created_at >= $%d", len(args)))
	}

	if end != nil {
		args = append(args, *end)
		conditions = append(conditions, fmt.Sprintf("a.created_at < $%d", len(args)))
	}

	args = append(args, limit)
	query := fmt.Sprintf(`
		SELECT
			u.id,
			u.name,
			COALESCE(SUM(a.points), 0)::integer AS total_points,
			COUNT(a.id)::integer AS total_activities
		FROM activities a
		JOIN users u ON u.id = a.user_id
		WHERE %s
		GROUP BY u.id, u.name
		ORDER BY total_points DESC, total_activities DESC, u.name ASC
		LIMIT $%d
	`, strings.Join(conditions, " AND "), len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.ContributionLeaderboardEntry, 0)
	rank := 1
	for rows.Next() {
		item := models.ContributionLeaderboardEntry{}
		if err := rows.Scan(&item.UserID, &item.UserName, &item.TotalPoints, &item.TotalActivities); err != nil {
			return nil, err
		}

		item.Rank = rank
		rank++
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func scanActivity(row scanner) (*models.Activity, error) {
	item := &models.Activity{}
	var expiresAt sql.NullTime

	if err := row.Scan(
		&item.ID,
		&item.Type,
		&item.Title,
		&item.Content,
		&item.Points,
		&item.UserID,
		&item.UserName,
		&item.CreatedBy,
		&item.CreatedByName,
		&expiresAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	item.ExpiresAt = nullTimePtr(expiresAt)

	return item, nil
}

func scanActivityComment(row scanner) (*models.ActivityComment, error) {
	item := &models.ActivityComment{}
	if err := row.Scan(
		&item.ID,
		&item.ActivityID,
		&item.UserID,
		&item.UserName,
		&item.Comment,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func scanFoodClaim(row scanner) (*models.FoodClaim, error) {
	item := &models.FoodClaim{}
	if err := row.Scan(
		&item.ID,
		&item.ActivityID,
		&item.UserID,
		&item.UserName,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func scanRiceResponse(row scanner) (*models.RiceResponse, error) {
	item := &models.RiceResponse{}
	if err := row.Scan(
		&item.ID,
		&item.ActivityID,
		&item.UserID,
		&item.UserName,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}
