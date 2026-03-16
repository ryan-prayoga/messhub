package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/handlers"
	"github.com/ryanprayoga/messhub/backend/internal/middleware"
)

func Register(app *fiber.App, cfg config.Config, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler) {
	api := app.Group("/api/v1")

	api.Get("/health", healthHandler.Health)
	api.Post("/auth/login", authHandler.Login)

	authenticated := api.Group("", middleware.RequireAuth(cfg))
	authenticated.Get("/auth/me", authHandler.Me)

	adminOnly := authenticated.Group("", middleware.RequireRoles("admin", "treasurer"))
	adminOnly.Get("/admin/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "admin or treasurer access granted",
		})
	})
}
