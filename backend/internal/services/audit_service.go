package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

type AuditLogInput struct {
	UserID     *string
	Action     string
	EntityType string
	EntityID   *string
	OldValue   any
	NewValue   any
}

type AuditService struct {
	auditRepository *repository.AuditLogRepository
}

func NewAuditService(auditRepository *repository.AuditLogRepository) *AuditService {
	return &AuditService{
		auditRepository: auditRepository,
	}
}

func (s *AuditService) Log(ctx context.Context, input AuditLogInput) error {
	return s.log(ctx, nil, input)
}

func (s *AuditService) LogTx(ctx context.Context, tx *sql.Tx, input AuditLogInput) error {
	return s.log(ctx, tx, input)
}

func (s *AuditService) log(ctx context.Context, tx *sql.Tx, input AuditLogInput) error {
	oldValue, err := marshalAuditValue(input.OldValue)
	if err != nil {
		return err
	}

	newValue, err := marshalAuditValue(input.NewValue)
	if err != nil {
		return err
	}

	params := repository.CreateAuditLogParams{
		UserID:     input.UserID,
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		OldValue:   oldValue,
		NewValue:   newValue,
	}

	if tx != nil {
		_, err = s.auditRepository.CreateTx(ctx, tx, params)
	} else {
		_, err = s.auditRepository.Create(ctx, params)
	}

	return err
}

func marshalAuditValue(value any) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
