package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

type PushSubscriptionInput struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256DH string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

type PushService struct {
	config                     config.Config
	pushSubscriptionRepository *repository.PushSubscriptionRepository
}

type webPushNotificationPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon"`
	Badge string `json:"badge"`
	Tag   string `json:"tag"`
	Data  struct {
		EntityID *string `json:"entity_id,omitempty"`
		Type     string  `json:"type"`
		URL      string  `json:"url"`
	} `json:"data"`
}

func NewPushService(
	cfg config.Config,
	pushSubscriptionRepository *repository.PushSubscriptionRepository,
) *PushService {
	return &PushService{
		config:                     cfg,
		pushSubscriptionRepository: pushSubscriptionRepository,
	}
}

func (s *PushService) Enabled() bool {
	return strings.TrimSpace(s.config.VAPIDPublicKey) != "" &&
		strings.TrimSpace(s.config.VAPIDPrivateKey) != ""
}

func (s *PushService) Subscribe(
	ctx context.Context,
	userID string,
	input PushSubscriptionInput,
) (*models.PushSubscription, error) {
	return s.pushSubscriptionRepository.Upsert(ctx, repository.UpsertPushSubscriptionParams{
		UserID:    strings.TrimSpace(userID),
		Endpoint:  strings.TrimSpace(input.Endpoint),
		P256DHKey: strings.TrimSpace(input.Keys.P256DH),
		AuthKey:   strings.TrimSpace(input.Keys.Auth),
	})
}

func (s *PushService) Unsubscribe(ctx context.Context, userID string, endpoint string) (bool, error) {
	return s.pushSubscriptionRepository.DeleteByUserAndEndpoint(
		ctx,
		strings.TrimSpace(userID),
		strings.TrimSpace(endpoint),
	)
}

func (s *PushService) DispatchNotifications(ctx context.Context, notifications []models.Notification) {
	if !s.Enabled() || len(notifications) == 0 {
		return
	}

	userIDs := uniqueNotificationUserIDs(notifications)
	subscriptionsByUser, err := s.pushSubscriptionRepository.ListByUserIDs(ctx, userIDs)
	if err != nil {
		slog.Error("push dispatch failed to load subscriptions", "error", err, "user_count", len(userIDs))
		return
	}

	for _, notification := range notifications {
		subscriptions := subscriptionsByUser[notification.UserID]
		if len(subscriptions) == 0 {
			continue
		}

		payloadBytes, err := json.Marshal(buildWebPushPayload(notification))
		if err != nil {
			slog.Error("push dispatch failed to marshal payload", "error", err, "notification_id", notification.ID)
			continue
		}

		for _, subscription := range subscriptions {
			if err := s.sendNotification(ctx, subscription, payloadBytes); err != nil {
				slog.Error(
					"push dispatch failed",
					"error",
					err,
					"notification_id",
					notification.ID,
					"user_id",
					notification.UserID,
				)
			}
		}
	}
}

func (s *PushService) sendNotification(
	ctx context.Context,
	subscription models.PushSubscription,
	payload []byte,
) error {
	resp, err := webpush.SendNotification(payload, &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			Auth:   subscription.AuthKey,
			P256dh: subscription.P256DHKey,
		},
	}, &webpush.Options{
		Subscriber:      s.config.VAPIDSubject,
		VAPIDPublicKey:  s.config.VAPIDPublicKey,
		VAPIDPrivateKey: s.config.VAPIDPrivateKey,
		TTL:             60,
		Urgency:         webpush.UrgencyNormal,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound {
		if deleteErr := s.pushSubscriptionRepository.DeleteByEndpoint(ctx, subscription.Endpoint); deleteErr != nil {
			slog.Error("failed to delete stale push subscription", "error", deleteErr, "endpoint", subscription.Endpoint)
		}
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return &PushSendError{
			StatusCode: resp.StatusCode,
			Endpoint:   subscription.Endpoint,
		}
	}

	return nil
}

type PushSendError struct {
	StatusCode int
	Endpoint   string
}

func (e *PushSendError) Error() string {
	return "push notification rejected by remote endpoint"
}

func uniqueNotificationUserIDs(notifications []models.Notification) []string {
	seen := make(map[string]struct{}, len(notifications))
	userIDs := make([]string, 0, len(notifications))

	for _, notification := range notifications {
		if _, found := seen[notification.UserID]; found {
			continue
		}

		seen[notification.UserID] = struct{}{}
		userIDs = append(userIDs, notification.UserID)
	}

	return userIDs
}

func buildWebPushPayload(notification models.Notification) webPushNotificationPayload {
	payload := webPushNotificationPayload{
		Title: notification.Title,
		Body:  notification.Message,
		Icon:  "/icons/icon-192.png",
		Badge: "/icons/icon-192.png",
		Tag:   notification.Type,
	}

	payload.Data.EntityID = notification.EntityID
	payload.Data.Type = notification.Type
	payload.Data.URL = notificationTargetURL(notification)

	return payload
}

func notificationTargetURL(notification models.Notification) string {
	switch notification.Type {
	case "wifi_bill_created", "wifi_payment_verified":
		return "/wifi"
	case "activity_created", "comment_created":
		return "/feed"
	default:
		return "/notifications"
	}
}
