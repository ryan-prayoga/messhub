package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type migrationFile struct {
	Version int
	Name    string
	Path    string
}

func ApplyMigrations(ctx context.Context, db *sql.DB, migrationsDir string) error {
	if err := ensureSchemaMigrationsTable(ctx, db); err != nil {
		return err
	}

	files, err := loadMigrationFiles(migrationsDir)
	if err != nil {
		return err
	}

	applied, err := loadAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	if len(applied) == 0 {
		if err := bootstrapExistingMigrations(ctx, db, files); err != nil {
			return err
		}

		applied, err = loadAppliedMigrations(ctx, db)
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		if _, ok := applied[file.Name]; ok {
			continue
		}

		if err := applyMigrationFile(ctx, db, file); err != nil {
			return err
		}
	}

	return nil
}

func ensureSchemaMigrationsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)

	return err
}

func loadMigrationFiles(migrationsDir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	files := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version, err := parseMigrationVersion(entry.Name())
		if err != nil {
			return nil, err
		}

		files = append(files, migrationFile{
			Version: version,
			Name:    entry.Name(),
			Path:    filepath.Join(migrationsDir, entry.Name()),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	return files, nil
}

func parseMigrationVersion(name string) (int, error) {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid migration name: %s", name)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid migration version %q: %w", name, err)
	}

	return version, nil
}

func loadAppliedMigrations(ctx context.Context, db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.QueryContext(ctx, `SELECT name FROM schema_migrations ORDER BY version ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := map[string]struct{}{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		applied[name] = struct{}{}
	}

	return applied, rows.Err()
}

func bootstrapExistingMigrations(ctx context.Context, db *sql.DB, files []migrationFile) error {
	existingUsersTable, err := relationExists(ctx, db, "users")
	if err != nil {
		return err
	}
	if !existingUsersTable {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, file := range files {
		applied, err := migrationAlreadyReflected(ctx, tx, file.Version)
		if err != nil {
			return err
		}
		if !applied {
			continue
		}

		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO schema_migrations (version, name) VALUES ($1, $2) ON CONFLICT (version) DO NOTHING`,
			file.Version,
			file.Name,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func migrationAlreadyReflected(ctx context.Context, db queryable, version int) (bool, error) {
	switch version {
	case 1:
		return relationExistsQueryable(ctx, db, "wallet_transactions")
	case 2:
		return constraintExists(ctx, db, "users_left_at_after_joined_at")
	case 3:
		return indexExists(ctx, db, "idx_wallet_transactions_created_at")
	case 4:
		return columnExists(ctx, db, "wifi_bills", "updated_at")
	case 5:
		return relationExistsQueryable(ctx, db, "activities")
	case 6:
		return relationExistsQueryable(ctx, db, "mess_settings")
	case 7:
		return relationExistsQueryable(ctx, db, "push_subscriptions")
	case 8:
		return relationExistsQueryable(ctx, db, "import_jobs")
	case 9:
		return indexExists(ctx, db, "idx_users_username_lower")
	default:
		return false, nil
	}
}

func applyMigrationFile(ctx context.Context, db *sql.DB, file migrationFile) error {
	payload, err := os.ReadFile(file.Path)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, string(payload)); err != nil {
		return fmt.Errorf("apply migration %s: %w", file.Name, err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`,
		file.Version,
		file.Name,
	); err != nil {
		return fmt.Errorf("record migration %s: %w", file.Name, err)
	}

	return tx.Commit()
}

type queryable interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func relationExists(ctx context.Context, db *sql.DB, name string) (bool, error) {
	return relationExistsQueryable(ctx, db, name)
}

func relationExistsQueryable(ctx context.Context, db queryable, name string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)`,
		name,
	).Scan(&exists)

	return exists, err
}

func columnExists(ctx context.Context, db queryable, table string, column string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = $1 AND column_name = $2
		)`,
		table,
		column,
	).Scan(&exists)

	return exists, err
}

func constraintExists(ctx context.Context, db queryable, name string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM pg_constraint
			WHERE conname = $1
		)`,
		name,
	).Scan(&exists)

	return exists, err
}

func indexExists(ctx context.Context, db queryable, name string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE schemaname = 'public' AND indexname = $1
		)`,
		name,
	).Scan(&exists)

	return exists, err
}
