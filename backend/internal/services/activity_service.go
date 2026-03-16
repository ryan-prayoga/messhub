package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

const (
	defaultActivityLimit     = 20
	maxActivityLimit         = 50
	defaultContributionLimit = 10
)

var (
	ErrInvalidActivityInput      = errors.New("invalid activity input")
	ErrInvalidActivityType       = errors.New("invalid activity type")
	ErrActivityNotFound          = errors.New("activity not found")
	ErrInvalidCommentInput       = errors.New("comment is required")
	ErrInvalidReactionInput      = errors.New("reaction type is required")
	ErrInvalidLeaderboardPeriod  = errors.New("period must be month or all")
	ErrFoodClaimAlreadyExists    = errors.New("user already claimed this food")
	ErrRiceResponseAlreadyExists = errors.New("user already responded to this rice plan")
	ErrFoodClaimNotAllowed       = errors.New("food claim is only available for food activities")
	ErrRiceResponseNotAllowed    = errors.New("rice response is only available for rice activities")
)

type ListActivitiesInput struct {
	Limit int
}

type CreateActivityInput struct {
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Points    *int    `json:"points"`
	ExpiresAt *string `json:"expires_at"`
}

type CreateActivityCommentInput struct {
	Comment string `json:"comment"`
}

type ToggleActivityReactionInput struct {
	ReactionType string `json:"reaction_type"`
}

type ActivityService struct {
	db                  *sql.DB
	activityRepository  *repository.ActivityRepository
	notificationService *NotificationService
	auditService        *AuditService
}

func NewActivityService(
	db *sql.DB,
	activityRepository *repository.ActivityRepository,
	notificationService *NotificationService,
	auditService *AuditService,
) *ActivityService {
	return &ActivityService{
		db:                  db,
		activityRepository:  activityRepository,
		notificationService: notificationService,
		auditService:        auditService,
	}
}

func (s *ActivityService) ListActivities(ctx context.Context, viewerID string, input ListActivitiesInput) ([]models.ActivityFeedItem, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultActivityLimit
	case limit > maxActivityLimit:
		limit = maxActivityLimit
	}

	activities, err := s.activityRepository.ListActivities(ctx, limit)
	if err != nil {
		return nil, err
	}

	items := make([]models.ActivityFeedItem, 0, len(activities))
	for _, activity := range activities {
		item, err := s.buildFeedItem(ctx, activity, viewerID)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	return items, nil
}

func (s *ActivityService) CreateActivity(ctx context.Context, actorID string, actorName string, input CreateActivityInput) (*models.ActivityFeedItem, error) {
	activityType := strings.TrimSpace(input.Type)
	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)

	if !isValidActivityType(activityType) {
		return nil, ErrInvalidActivityType
	}

	if title == "" || content == "" || strings.TrimSpace(actorID) == "" {
		return nil, ErrInvalidActivityInput
	}

	points := 0
	if activityType == models.ActivityTypeContribution {
		points = 1
		if input.Points != nil {
			points = *input.Points
		}
		if points <= 0 {
			return nil, ErrInvalidActivityInput
		}
	}

	expiresAt, err := resolveActivityExpiry(activityType, input.ExpiresAt)
	if err != nil {
		return nil, ErrInvalidActivityInput
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	activity, err := s.activityRepository.CreateActivityTx(ctx, tx, repository.CreateActivityParams{
		Type:      activityType,
		Title:     title,
		Content:   content,
		Points:    points,
		UserID:    actorID,
		CreatedBy: actorID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "activity_created",
		EntityType: "activity",
		EntityID:   stringPtr(activity.ID),
		NewValue: map[string]any{
			"type":       activityType,
			"title":      title,
			"points":     points,
			"expires_at": expiresAt,
		},
	}); err != nil {
		return nil, err
	}

	if err := s.notificationService.NotifyAllActiveExceptTx(
		ctx,
		tx,
		actorID,
		"Activity baru",
		fmt.Sprintf("%s membuat %s baru: %s", displayActorName(actorName), activityTypeLabel(activityType), title),
		"activity_created",
		stringPtr(activity.ID),
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetActivityItem(ctx, activity.ID, actorID)
}

func (s *ActivityService) GetActivityItem(ctx context.Context, activityID string, viewerID string) (*models.ActivityFeedItem, error) {
	activity, err := s.activityRepository.FindActivityByID(ctx, activityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrActivityNotFound
		}

		return nil, err
	}

	return s.buildFeedItem(ctx, *activity, viewerID)
}

func (s *ActivityService) ListComments(ctx context.Context, activityID string) ([]models.ActivityComment, error) {
	if _, err := s.ensureActivityExists(ctx, activityID); err != nil {
		return nil, err
	}

	return s.activityRepository.ListComments(ctx, activityID)
}

func (s *ActivityService) AddComment(ctx context.Context, activityID string, actorID string, actorName string, input CreateActivityCommentInput) (*models.ActivityFeedItem, error) {
	comment := strings.TrimSpace(input.Comment)
	if comment == "" {
		return nil, ErrInvalidCommentInput
	}

	activity, err := s.ensureActivityExists(ctx, activityID)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := s.activityRepository.CreateCommentTx(ctx, tx, repository.CreateActivityCommentParams{
		ActivityID: activityID,
		UserID:     actorID,
		Comment:    comment,
	}); err != nil {
		return nil, err
	}

	if err := s.notificationService.NotifyAllActiveExceptTx(
		ctx,
		tx,
		actorID,
		"Komentar baru",
		fmt.Sprintf("%s mengomentari %s", displayActorName(actorName), activity.Title),
		"comment_created",
		stringPtr(activity.ID),
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetActivityItem(ctx, activityID, actorID)
}

func (s *ActivityService) ToggleReaction(ctx context.Context, activityID string, actorID string, input ToggleActivityReactionInput) (*models.ActivityFeedItem, error) {
	reactionType := strings.TrimSpace(input.ReactionType)
	if reactionType == "" {
		return nil, ErrInvalidReactionInput
	}

	if _, err := s.ensureActivityExists(ctx, activityID); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := s.activityRepository.ToggleReactionTx(ctx, tx, repository.ToggleActivityReactionParams{
		ActivityID:   activityID,
		UserID:       actorID,
		ReactionType: reactionType,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetActivityItem(ctx, activityID, actorID)
}

func (s *ActivityService) ClaimFood(ctx context.Context, activityID string, actorID string) (*models.FoodClaim, error) {
	activity, err := s.ensureActivityExists(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if activity.Type != models.ActivityTypeFood {
		return nil, ErrFoodClaimNotAllowed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	claim, err := s.activityRepository.CreateFoodClaimTx(ctx, tx, repository.CreateFoodClaimParams{
		ActivityID: activityID,
		UserID:     actorID,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrFoodClaimAlreadyExists
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "food_claim",
		EntityType: "activity",
		EntityID:   stringPtr(activityID),
		NewValue: map[string]any{
			"activity_id": activityID,
			"user_id":     actorID,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.getFoodClaimByUser(ctx, activityID, actorID, claim)
}

func (s *ActivityService) ListFoodClaims(ctx context.Context, activityID string) ([]models.FoodClaim, error) {
	activity, err := s.ensureActivityExists(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if activity.Type != models.ActivityTypeFood {
		return nil, ErrFoodClaimNotAllowed
	}

	return s.activityRepository.ListFoodClaims(ctx, activityID)
}

func (s *ActivityService) RespondRice(ctx context.Context, activityID string, actorID string) (*models.RiceResponse, error) {
	activity, err := s.ensureActivityExists(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if activity.Type != models.ActivityTypeRice {
		return nil, ErrRiceResponseNotAllowed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	responseItem, err := s.activityRepository.CreateRiceResponseTx(ctx, tx, repository.CreateRiceResponseParams{
		ActivityID: activityID,
		UserID:     actorID,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrRiceResponseAlreadyExists
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "rice_response",
		EntityType: "activity",
		EntityID:   stringPtr(activityID),
		NewValue: map[string]any{
			"activity_id": activityID,
			"user_id":     actorID,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.getRiceResponseByUser(ctx, activityID, actorID, responseItem)
}

func (s *ActivityService) ListRiceResponses(ctx context.Context, activityID string) ([]models.RiceResponse, error) {
	activity, err := s.ensureActivityExists(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if activity.Type != models.ActivityTypeRice {
		return nil, ErrRiceResponseNotAllowed
	}

	return s.activityRepository.ListRiceResponses(ctx, activityID)
}

func (s *ActivityService) GetContributionLeaderboard(ctx context.Context, period string) ([]models.ContributionLeaderboardEntry, error) {
	period = strings.TrimSpace(strings.ToLower(period))
	if period == "" {
		period = "month"
	}
	if period != "month" && period != "all" {
		return nil, ErrInvalidLeaderboardPeriod
	}

	var start *time.Time
	var end *time.Time
	if period == "month" {
		yearMonthStart, yearMonthEnd := currentMonthRange()
		start = &yearMonthStart
		end = &yearMonthEnd
	}

	return s.activityRepository.GetContributionLeaderboard(ctx, start, end, defaultContributionLimit)
}

func (s *ActivityService) ensureActivityExists(ctx context.Context, activityID string) (*models.Activity, error) {
	activity, err := s.activityRepository.FindActivityByID(ctx, activityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrActivityNotFound
		}

		return nil, err
	}

	return activity, nil
}

func (s *ActivityService) buildFeedItem(ctx context.Context, activity models.Activity, viewerID string) (*models.ActivityFeedItem, error) {
	comments, err := s.activityRepository.ListComments(ctx, activity.ID)
	if err != nil {
		return nil, err
	}

	reactions, err := s.activityRepository.ListReactions(ctx, activity.ID, viewerID)
	if err != nil {
		return nil, err
	}

	claims, err := s.activityRepository.ListFoodClaims(ctx, activity.ID)
	if err != nil {
		return nil, err
	}

	riceResponses, err := s.activityRepository.ListRiceResponses(ctx, activity.ID)
	if err != nil {
		return nil, err
	}

	return &models.ActivityFeedItem{
		Activity:      activity,
		Comments:      comments,
		Reactions:     reactions,
		Claims:        claims,
		RiceResponses: riceResponses,
	}, nil
}

func (s *ActivityService) getFoodClaimByUser(ctx context.Context, activityID string, actorID string, fallback *models.FoodClaim) (*models.FoodClaim, error) {
	items, err := s.activityRepository.ListFoodClaims(ctx, activityID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.UserID == actorID {
			return &item, nil
		}
	}

	return fallback, nil
}

func (s *ActivityService) getRiceResponseByUser(ctx context.Context, activityID string, actorID string, fallback *models.RiceResponse) (*models.RiceResponse, error) {
	items, err := s.activityRepository.ListRiceResponses(ctx, activityID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.UserID == actorID {
			return &item, nil
		}
	}

	return fallback, nil
}

func isValidActivityType(activityType string) bool {
	switch activityType {
	case models.ActivityTypeContribution, models.ActivityTypeFood, models.ActivityTypeRice, models.ActivityTypeAnnouncement, models.ActivityTypeOther:
		return true
	default:
		return false
	}
}

func resolveActivityExpiry(activityType string, raw *string) (*time.Time, error) {
	if raw != nil && strings.TrimSpace(*raw) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*raw))
		if err != nil {
			return nil, err
		}

		return &parsed, nil
	}

	if activityType == models.ActivityTypeContribution {
		return nil, nil
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	return &expiresAt, nil
}

func currentMonthRange() (time.Time, time.Time) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
	}

	now := time.Now().In(location)
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 1, 0)

	return start.UTC(), end.UTC()
}

func activityTypeLabel(activityType string) string {
	switch activityType {
	case models.ActivityTypeContribution:
		return "kontribusi"
	case models.ActivityTypeFood:
		return "info makanan"
	case models.ActivityTypeRice:
		return "rencana nasi"
	case models.ActivityTypeAnnouncement:
		return "pengumuman"
	default:
		return "aktivitas"
	}
}

func displayActorName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Seseorang"
	}

	return name
}
