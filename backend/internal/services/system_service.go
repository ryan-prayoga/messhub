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
	startedAt  time.Time
}

func NewSystemService(cfg config.Config, db *sql.DB) *SystemService {
	return &SystemService{
		appVersion: cfg.AppVersion,
		db:         db,
		startedAt:  time.Now().UTC(),
	}
}

func (s *SystemService) GetStatus(ctx context.Context) (*models.SystemStatus, error) {
	status := &models.SystemStatus{
		Status:            "ok",
		DatabaseStatus:    "ok",
		DatabaseReachable: true,
		ServerTime:        time.Now().UTC(),
		AppVersion:        s.appVersion,
		UptimeSeconds:     int64(time.Since(s.startedAt).Seconds()),
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.db.PingContext(pingCtx); err != nil {
		status.Status = "degraded"
		status.DatabaseStatus = "error"
		status.DatabaseReachable = false
	}

	return status, nil
}
