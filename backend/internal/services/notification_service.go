package services

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

const (
	defaultNotificationLimit = 20
	maxNotificationLimit     = 100
)

type NotificationList struct {
	Items       []models.Notification `json:"items"`
	UnreadCount int                   `json:"unread_count"`
}

type MarkNotificationsReadInput struct {
	IDs []string `json:"ids"`
	All bool     `json:"all"`
}

type NotificationService struct {
	db                     *sql.DB
	notificationRepository *repository.NotificationRepository
	userRepository         *repository.UserRepository
	auditService           *AuditService
}

func NewNotificationService(
	db *sql.DB,
	notificationRepository *repository.NotificationRepository,
	userRepository *repository.UserRepository,
	auditService *AuditService,
) *NotificationService {
	return &NotificationService{
		db:                     db,
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
		auditService:           auditService,
	}
}

func (s *NotificationService) ListForUser(ctx context.Context, userID string, limit int) (*NotificationList, error) {
	if limit <= 0 {
		limit = defaultNotificationLimit
	}
	if limit > maxNotificationLimit {
		limit = maxNotificationLimit
	}

	items, err := s.notificationRepository.ListByUser(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.notificationRepository.CountUnread(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &NotificationList{
		Items:       items,
		UnreadCount: unreadCount,
	}, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, userID string, input MarkNotificationsReadInput) (int, error) {
	trimmedIDs := make([]string, 0, len(input.IDs))
	for _, id := range input.IDs {
		id = strings.TrimSpace(id)
		if id != "" {
			trimmedIDs = append(trimmedIDs, id)
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var updatedIDs []string
	if input.All || len(trimmedIDs) == 0 {
		updatedIDs, err = s.notificationRepository.MarkAllReadTx(ctx, tx, userID)
	} else {
		updatedIDs, err = s.notificationRepository.MarkReadByIDsTx(ctx, tx, userID, trimmedIDs)
	}
	if err != nil {
		return 0, err
	}

	if len(updatedIDs) > 0 {
		if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
			UserID:     stringPtr(userID),
			Action:     "notification_read",
			EntityType: "notification",
			NewValue: map[string]any{
				"ids":   updatedIDs,
				"count": len(updatedIDs),
				"all":   input.All || len(trimmedIDs) == 0,
			},
		}); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(updatedIDs), nil
}

func (s *NotificationService) NotifyAllActiveExceptTx(
	ctx context.Context,
	tx *sql.Tx,
	actorID string,
	title string,
	message string,
	notificationType string,
	entityID *string,
) error {
	users, err := s.userRepository.ListActive(ctx)
	if err != nil {
		return err
	}

	recipientIDs := make([]string, 0, len(users))
	for _, user := range users {
		if user.ID == actorID {
			continue
		}

		recipientIDs = append(recipientIDs, user.ID)
	}

	return s.notifyUsersTx(ctx, tx, recipientIDs, title, message, notificationType, entityID)
}

func (s *NotificationService) NotifyUserTx(
	ctx context.Context,
	tx *sql.Tx,
	userID string,
	title string,
	message string,
	notificationType string,
	entityID *string,
) error {
	if strings.TrimSpace(userID) == "" {
		return nil
	}

	return s.notifyUsersTx(ctx, tx, []string{userID}, title, message, notificationType, entityID)
}

func (s *NotificationService) notifyUsersTx(
	ctx context.Context,
	tx *sql.Tx,
	userIDs []string,
	title string,
	message string,
	notificationType string,
	entityID *string,
) error {
	title = strings.TrimSpace(title)
	message = strings.TrimSpace(message)
	notificationType = strings.TrimSpace(notificationType)

	if len(userIDs) == 0 || title == "" || message == "" || notificationType == "" {
		return nil
	}

	seen := make(map[string]struct{}, len(userIDs))
	for _, userID := range userIDs {
		userID = strings.TrimSpace(userID)
		if userID == "" {
			continue
		}
		if _, found := seen[userID]; found {
			continue
		}

		seen[userID] = struct{}{}
		if _, err := s.notificationRepository.CreateTx(ctx, tx, repository.CreateNotificationParams{
			UserID:   userID,
			Title:    title,
			Message:  message,
			Type:     notificationType,
			EntityID: entityID,
		}); err != nil {
			return err
		}
	}

	return nil
}
