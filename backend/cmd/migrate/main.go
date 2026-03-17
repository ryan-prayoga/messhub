package main

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	migrationsDir := filepath.Join("db", "migrations")
	if err := database.ApplyMigrations(ctx, db, migrationsDir); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	log.Printf("migrations applied successfully from %s", migrationsDir)
}
