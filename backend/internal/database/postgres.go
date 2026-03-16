package database

import (
	"database/sql"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanprayoga/messhub/backend/internal/config"
)

func NewPostgres(cfg config.Config) (*sql.DB, error) {
	db := stdlib.OpenDB(*cfg.DatabaseConfig())

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
