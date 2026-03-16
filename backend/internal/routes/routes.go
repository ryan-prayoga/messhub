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
	walletHandler *handlers.WalletHandler,
	wifiHandler *handlers.WifiHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	api := app.Group("/api/v1")

	app.Get("/health", healthHandler.Health)
	api.Get("/health", healthHandler.Health)
	api.Post("/auth/login", authHandler.Login)

	authenticated := api.Group("", authMiddleware.RequireAuth())
	authenticated.Get("/auth/me", authHandler.Me)
	authenticated.Get("/wallet", walletHandler.GetSummary)
	authenticated.Get("/wallet/transactions", walletHandler.ListTransactions)
	authenticated.Get("/wifi/active", wifiHandler.GetActiveBill)
	authenticated.Get("/wifi/my", wifiHandler.GetMyBills)
	authenticated.Post("/wifi/bills/:id/submit", wifiHandler.SubmitPaymentProof)

	userReaders := authenticated.Group("", middleware.RequireRoles("admin", "treasurer"))
	userReaders.Get("/users", userHandler.List)
	userReaders.Post("/wallet/transactions", walletHandler.CreateTransaction)
	userReaders.Post("/wifi/bills", wifiHandler.CreateBill)
	userReaders.Get("/wifi/bills", wifiHandler.ListBills)
	userReaders.Get("/wifi/bills/:id", wifiHandler.GetBillDetail)
	userReaders.Patch("/wifi/bills/:id/verify/:memberId", wifiHandler.VerifyPayment)
	userReaders.Patch("/wifi/bills/:id/reject/:memberId", wifiHandler.RejectPayment)

	adminOnly := authenticated.Group("", middleware.RequireRoles("admin"))
	adminOnly.Post("/users", userHandler.Create)
	adminOnly.Patch("/users/:id", userHandler.Update)

	userReaders.Get("/admin/ping", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "admin or treasurer access granted", fiber.Map{
			"status": "ok",
		})
	})
}
