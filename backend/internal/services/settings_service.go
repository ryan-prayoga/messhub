package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

var ErrInvalidSettingsInput = errors.New("invalid settings input")

type UpdateMessSettingsInput struct {
	MessName          *string `json:"mess_name"`
	WifiPrice         *int64  `json:"wifi_price"`
	WifiDeadlineDay   *int    `json:"wifi_deadline_day"`
	BankAccountName   *string `json:"bank_account_name"`
	BankAccountNumber *string `json:"bank_account_number"`
}

type SettingsService struct {
	config             config.Config
	db                 *sql.DB
	settingsRepository *repository.SettingsRepository
	auditService       *AuditService
}

func NewSettingsService(
	cfg config.Config,
	db *sql.DB,
	settingsRepository *repository.SettingsRepository,
	auditService *AuditService,
) *SettingsService {
	return &SettingsService{
		config:             cfg,
		db:                 db,
		settingsRepository: settingsRepository,
		auditService:       auditService,
	}
}

func (s *SettingsService) GetSettings(ctx context.Context) (*models.MessSettings, error) {
	settings, err := s.settingsRepository.Get(ctx)
	if err == nil {
		return settings, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return s.settingsRepository.Upsert(ctx, repository.UpsertMessSettingsParams{
		MessName:          s.config.AppName,
		WifiPrice:         DefaultWifiNominalPerPerson,
		WifiDeadlineDay:   defaultWifiDeadlineDay,
		BankAccountName:   defaultBankAccountName,
		BankAccountNumber: defaultBankAccountNumber,
	})
}

func (s *SettingsService) UpdateSettings(ctx context.Context, actorID string, input UpdateMessSettingsInput) (*models.MessSettings, error) {
	current, err := s.GetSettings(ctx)
	if err != nil {
		return nil, err
	}

	next := *current
	updated := false

	if input.MessName != nil {
		next.MessName = strings.TrimSpace(*input.MessName)
		updated = true
	}

	if input.WifiPrice != nil {
		next.WifiPrice = *input.WifiPrice
		updated = true
	}

	if input.WifiDeadlineDay != nil {
		next.WifiDeadlineDay = *input.WifiDeadlineDay
		updated = true
	}

	if input.BankAccountName != nil {
		next.BankAccountName = strings.TrimSpace(*input.BankAccountName)
		updated = true
	}

	if input.BankAccountNumber != nil {
		next.BankAccountNumber = strings.TrimSpace(*input.BankAccountNumber)
		updated = true
	}

	if !updated || !isValidMessSettings(next) {
		return nil, ErrInvalidSettingsInput
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	saved, err := s.settingsRepository.UpsertTx(ctx, tx, repository.UpsertMessSettingsParams{
		MessName:          next.MessName,
		WifiPrice:         next.WifiPrice,
		WifiDeadlineDay:   next.WifiDeadlineDay,
		BankAccountName:   next.BankAccountName,
		BankAccountNumber: next.BankAccountNumber,
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "update_settings",
		EntityType: "mess_settings",
		EntityID:   stringPtr(saved.ID),
		OldValue:   current,
		NewValue:   saved,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return saved, nil
}

func isValidMessSettings(settings models.MessSettings) bool {
	return strings.TrimSpace(settings.MessName) != "" &&
		settings.WifiPrice > 0 &&
		settings.WifiDeadlineDay >= 1 &&
		settings.WifiDeadlineDay <= 31 &&
		strings.TrimSpace(settings.BankAccountName) != "" &&
		strings.TrimSpace(settings.BankAccountNumber) != ""
}
