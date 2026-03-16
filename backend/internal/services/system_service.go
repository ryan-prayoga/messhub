package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type SystemService struct {
	appVersion string
	db         *sql.DB
}

func NewSystemService(cfg config.Config, db *sql.DB) *SystemService {
	return &SystemService{
		appVersion: cfg.AppVersion,
		db:         db,
	}
}

func (s *SystemService) GetStatus(ctx context.Context) (*models.SystemStatus, error) {
	status := &models.SystemStatus{
		DatabaseStatus: "ok",
		ServerTime:     time.Now().UTC(),
		AppVersion:     s.appVersion,
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.db.PingContext(pingCtx); err != nil {
		status.DatabaseStatus = "error"
	}

	return status, nil
}
