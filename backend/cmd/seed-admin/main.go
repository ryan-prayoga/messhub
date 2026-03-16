package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanprayoga/messhub/backend/internal/config"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	query := `
		INSERT INTO users (name, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, TRUE)
		ON CONFLICT (email)
		DO UPDATE SET
			name = EXCLUDED.name,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			is_active = TRUE,
			updated_at = NOW()
	`

	if _, err := db.ExecContext(
		ctx,
		query,
		cfg.SeedAdminName,
		cfg.SeedAdminEmail,
		string(hashedPassword),
		cfg.SeedAdminRole,
	); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	log.Printf("admin seeded: %s", cfg.SeedAdminEmail)
}
