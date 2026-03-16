package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/handlers"
	"github.com/ryanprayoga/messhub/backend/internal/middleware"
	"github.com/ryanprayoga/messhub/backend/internal/response"
)

func Register(
	app *fiber.App,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	api := app.Group("/api/v1")

	app.Get("/health", healthHandler.Health)
	api.Get("/health", healthHandler.Health)
	api.Post("/auth/login", authHandler.Login)

	authenticated := api.Group("", authMiddleware.RequireAuth())
	authenticated.Get("/auth/me", authHandler.Me)

	userReaders := authenticated.Group("", middleware.RequireRoles("admin", "treasurer"))
	userReaders.Get("/users", userHandler.List)

	adminOnly := authenticated.Group("", middleware.RequireRoles("admin"))
	adminOnly.Post("/users", userHandler.Create)
	adminOnly.Patch("/users/:id", userHandler.Update)

	userReaders.Get("/admin/ping", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "admin or treasurer access granted", fiber.Map{
			"status": "ok",
		})
	})
}
