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
		AppName: cfg.AppName,
	})

	web.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	userRepository := repository.NewUserRepository(db)
	authService := services.NewAuthService(cfg, userRepository)
	userService := services.NewUserService(userRepository)
	authMiddleware := middleware.NewAuthMiddleware(cfg, userRepository)
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	routes.Register(web, healthHandler, authHandler, userHandler, authMiddleware)

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
