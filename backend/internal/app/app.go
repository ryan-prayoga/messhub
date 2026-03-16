package app

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/database"
	"github.com/ryanprayoga/messhub/backend/internal/handlers"
	"github.com/ryanprayoga/messhub/backend/internal/middleware"
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

	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	web := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: response.FiberErrorHandler,
	})

	web.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	userRepository := repository.NewUserRepository(db)
	settingsRepository := repository.NewSettingsRepository(db)
	walletRepository := repository.NewWalletRepository(db)
	wifiRepository := repository.NewWifiRepository(db)
	activityRepository := repository.NewActivityRepository(db)
	auditRepository := repository.NewAuditLogRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)
	authService := services.NewAuthService(cfg, userRepository)
	auditService := services.NewAuditService(auditRepository)
	settingsService := services.NewSettingsService(cfg, db, settingsRepository, auditService)
	systemService := services.NewSystemService(cfg, db)
	userService := services.NewUserService(db, userRepository, auditService)
	walletService := services.NewWalletService(db, walletRepository, auditService)
	notificationService := services.NewNotificationService(db, notificationRepository, userRepository, auditService)
	wifiService := services.NewWifiService(db, wifiRepository, settingsService, auditService, notificationService)
	activityService := services.NewActivityService(db, activityRepository, notificationService, auditService)
	authMiddleware := middleware.NewAuthMiddleware(cfg, userRepository)
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandler(userService)
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	systemHandler := handlers.NewSystemHandler(systemService)
	walletHandler := handlers.NewWalletHandler(walletService)
	wifiHandler := handlers.NewWifiHandler(wifiService)
	activityHandler := handlers.NewActivityHandler(activityService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	routes.Register(
		web,
		healthHandler,
		authHandler,
		userHandler,
		profileHandler,
		settingsHandler,
		systemHandler,
		walletHandler,
		wifiHandler,
		activityHandler,
		notificationHandler,
		authMiddleware,
	)

	return &App{
		config: cfg,
		db:     db,
		fiber:  web,
	}, nil
}

func (a *App) Listen() error {
	return a.fiber.Listen(fmt.Sprintf("%s:%s", a.config.BackendHost, a.config.BackendPort))
}

func (a *App) Close() error {
	return a.db.Close()
}
