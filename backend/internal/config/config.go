package config

import (
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	AppName           string
	AppEnv            string
	BackendHost       string
	BackendPort       string
	DatabaseURL       string
	JWTSecret         string
	JWTExpiresInHours int
	CORSOrigin        string
	SeedAdminName     string
	SeedAdminEmail    string
	SeedAdminPassword string
	SeedAdminRole     string
}

func Load() Config {
	loadEnvFiles()

	return Config{
		AppName:           getEnv("APP_NAME", "MessHub"),
		AppEnv:            getEnv("APP_ENV", "development"),
		BackendHost:       getEnv("BACKEND_HOST", "0.0.0.0"),
		BackendPort:       getEnv("PORT", getEnv("BACKEND_PORT", "4100")),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://messhub:messhub@127.0.0.1:5432/messhub?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "change-this-secret"),
		JWTExpiresInHours: getEnvInt("JWT_EXPIRES_IN_HOURS", 72),
		CORSOrigin:        getEnv("CORS_ORIGIN", "http://127.0.0.1:4101,http://localhost:4101"),
		SeedAdminName:     getEnv("SEED_ADMIN_NAME", "MessHub Admin"),
		SeedAdminEmail:    getEnv("SEED_ADMIN_EMAIL", "admin@messhub.local"),
		SeedAdminPassword: getEnv("SEED_ADMIN_PASSWORD", "ChangeMe123!"),
		SeedAdminRole:     getEnv("SEED_ADMIN_ROLE", "admin"),
	}
}

func (c Config) DatabaseConfig() *pgx.ConnConfig {
	cfg, err := pgx.ParseConfig(c.DatabaseURL)
	if err != nil {
		panic(err)
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return number
}
