package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/database"
	"github.com/ryanprayoga/messhub/backend/internal/handlers"
	"github.com/ryanprayoga/messhub/backend/internal/middleware"
	"github.com/ryanprayoga/messhub/backend/internal/observability"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/routes"
	"github.com/ryanprayoga/messhub/backend/internal/services"
)

type App struct {
	config config.Config
	db     *sql.DB
	fiber  *fiber.App
}

func New() (*App, error) {
	cfg := config.Load()
	logger := observability.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	web := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: response.FiberErrorHandler,
	})

	allowedOrigins := parseAllowedOrigins(cfg.CORSOrigin)

	web.Use(middleware.RequestContext())
	web.Use(middleware.Recover(logger))
	web.Use(middleware.SecurityHeaders(cfg.ContentSecurityPolicy))
	web.Use(middleware.RequestLogger(logger))
	web.Use(cors.New(cors.Config{
		AllowOrigins:  strings.Join(allowedOrigins, ","),
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		ExposeHeaders: middleware.RequestIDHeader,
		MaxAge:        300,
	}))

	userRepository := repository.NewUserRepository(db)
	settingsRepository := repository.NewSettingsRepository(db)
	walletRepository := repository.NewWalletRepository(db)
	wifiRepository := repository.NewWifiRepository(db)
	activityRepository := repository.NewActivityRepository(db)
	sharedExpenseRepository := repository.NewSharedExpenseRepository(db)
	proposalRepository := repository.NewProposalRepository(db)
	auditRepository := repository.NewAuditLogRepository(db)
	importJobRepository := repository.NewImportJobRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)
	pushSubscriptionRepository := repository.NewPushSubscriptionRepository(db)
	authService := services.NewAuthService(cfg, userRepository)
	auditService := services.NewAuditService(auditRepository)
	settingsService := services.NewSettingsService(cfg, db, settingsRepository, auditService)
	systemService := services.NewSystemService(cfg, db)
	pushService := services.NewPushService(cfg, pushSubscriptionRepository)
	userService := services.NewUserService(db, userRepository, auditService)
	walletService := services.NewWalletService(db, walletRepository, auditService)
	importService := services.NewImportService(db, userRepository, walletRepository, importJobRepository, auditService)
	notificationService := services.NewNotificationService(db, notificationRepository, userRepository, pushService, auditService)
	wifiService := services.NewWifiService(db, wifiRepository, settingsService, auditService, notificationService)
	activityService := services.NewActivityService(db, activityRepository, notificationService, auditService)
	sharedExpenseService := services.NewSharedExpenseService(db, sharedExpenseRepository, userRepository, auditService)
	proposalService := services.NewProposalService(db, proposalRepository, auditService)
	authMiddleware := middleware.NewAuthMiddleware(cfg, userRepository)
	healthHandler := handlers.NewHealthHandler(systemService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandler(userService)
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	systemHandler := handlers.NewSystemHandler(systemService)
	walletHandler := handlers.NewWalletHandler(walletService)
	importHandler := handlers.NewImportHandler(importService)
	wifiHandler := handlers.NewWifiHandler(wifiService)
	activityHandler := handlers.NewActivityHandler(activityService)
	sharedExpenseHandler := handlers.NewSharedExpenseHandler(sharedExpenseService)
	proposalHandler := handlers.NewProposalHandler(proposalService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	pushHandler := handlers.NewPushHandler(pushService)

	routes.Register(
		web,
		healthHandler,
		authHandler,
		userHandler,
		profileHandler,
		settingsHandler,
		systemHandler,
		walletHandler,
		importHandler,
		wifiHandler,
		activityHandler,
		sharedExpenseHandler,
		proposalHandler,
		notificationHandler,
		pushHandler,
		authMiddleware,
	)

	return &App{
		config: cfg,
		db:     db,
		fiber:  web,
	}, nil
}

func parseAllowedOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		origins = append(origins, trimmed)
	}

	if len(origins) == 0 {
		return []string{"http://127.0.0.1:4101", "http://localhost:4101"}
	}

	return origins
}

func (a *App) Listen() error {
	return a.fiber.Listen(fmt.Sprintf("%s:%s", a.config.BackendHost, a.config.BackendPort))
}

func (a *App) Close() error {
	return a.db.Close()
}
