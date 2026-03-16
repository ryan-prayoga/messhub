package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	db := stdlib.OpenDB(*cfg.DatabaseConfig())
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	if !models.IsValidRole(cfg.SeedAdminRole) {
		log.Fatalf("invalid SEED_ADMIN_ROLE: %s", cfg.SeedAdminRole)
	}

	var existingID string
	err := db.QueryRowContext(
		ctx,
		`SELECT id FROM users WHERE LOWER(email) = LOWER($1) LIMIT 1`,
		cfg.SeedAdminEmail,
	).Scan(&existingID)
	if err == nil {
		log.Printf("admin already exists: %s", cfg.SeedAdminEmail)
		return
	}

	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("check seed admin: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	userRepository := repository.NewUserRepository(db)
	username, err := userRepository.FindAvailableUsername(ctx, cfg.SeedAdminName, cfg.SeedAdminEmail)
	if err != nil {
		log.Fatalf("generate username: %v", err)
	}

	query := `
		INSERT INTO users (name, email, username, password_hash, role, is_active, joined_at)
		VALUES ($1, LOWER($2), $3, $4, $5, TRUE, NOW())
	`

	if _, err := db.ExecContext(
		ctx,
		query,
		cfg.SeedAdminName,
		cfg.SeedAdminEmail,
		username,
		string(hashedPassword),
		cfg.SeedAdminRole,
	); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	log.Printf("admin seeded: %s", cfg.SeedAdminEmail)
}
