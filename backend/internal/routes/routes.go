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
	profileHandler *handlers.ProfileHandler,
	settingsHandler *handlers.SettingsHandler,
	systemHandler *handlers.SystemHandler,
	walletHandler *handlers.WalletHandler,
	importHandler *handlers.ImportHandler,
	wifiHandler *handlers.WifiHandler,
	activityHandler *handlers.ActivityHandler,
	notificationHandler *handlers.NotificationHandler,
	pushHandler *handlers.PushHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	api := app.Group("/api/v1")

	app.Get("/health", healthHandler.Health)
	api.Get("/health", healthHandler.Health)
	api.Post("/auth/login", middleware.LoginRateLimit(), authHandler.Login)

	authenticated := api.Group("", authMiddleware.RequireAuth())
	authenticated.Get("/auth/me", authHandler.Me)
	authenticated.Get("/profile", profileHandler.Get)
	authenticated.Patch("/profile", profileHandler.Update)
	authenticated.Patch("/profile/password", profileHandler.ChangePassword)
	authenticated.Get("/settings", settingsHandler.Get)
	authenticated.Get("/activities", activityHandler.ListActivities)
	authenticated.Post("/activities", middleware.PostRateLimit("activities"), activityHandler.CreateActivity)
	authenticated.Get("/activities/:id/comments", activityHandler.ListComments)
	authenticated.Post("/activities/:id/comments", middleware.PostRateLimit("comments"), activityHandler.AddComment)
	authenticated.Post("/activities/:id/reactions", middleware.PostRateLimit("reactions"), activityHandler.ToggleReaction)
	authenticated.Post("/activities/:id/claim", activityHandler.ClaimFood)
	authenticated.Get("/activities/:id/claims", activityHandler.ListFoodClaims)
	authenticated.Post("/activities/:id/rice-response", activityHandler.RespondRice)
	authenticated.Get("/activities/:id/rice-responses", activityHandler.ListRiceResponses)
	authenticated.Get("/contributions/leaderboard", activityHandler.GetContributionLeaderboard)
	authenticated.Get("/notifications", notificationHandler.List)
	authenticated.Post("/notifications/read", notificationHandler.MarkRead)
	authenticated.Post("/push/subscribe", pushHandler.Subscribe)
	authenticated.Delete("/push/unsubscribe", pushHandler.Unsubscribe)
	authenticated.Get("/wallet", walletHandler.GetSummary)
	authenticated.Get("/wallet/transactions", walletHandler.ListTransactions)
	authenticated.Get("/wifi/active", wifiHandler.GetActiveBill)
	authenticated.Get("/wifi/my", wifiHandler.GetMyBills)
	authenticated.Post("/wifi/bills/:id/submit", wifiHandler.SubmitPaymentProof)

	userReaders := authenticated.Group("", middleware.RequireRole("admin", "treasurer"))
	userReaders.Get("/users", userHandler.List)
	userReaders.Post("/wallet/transactions", walletHandler.CreateTransaction)
	userReaders.Post("/wifi/bills", wifiHandler.CreateBill)
	userReaders.Get("/wifi/bills", wifiHandler.ListBills)
	userReaders.Get("/wifi/bills/:id", wifiHandler.GetBillDetail)
	userReaders.Patch("/wifi/bills/:id/verify/:memberId", wifiHandler.VerifyPayment)
	userReaders.Patch("/wifi/bills/:id/reject/:memberId", wifiHandler.RejectPayment)

	adminOnly := authenticated.Group("", middleware.RequireRole("admin"))
	adminOnly.Post("/users", userHandler.Create)
	adminOnly.Patch("/users/:id", userHandler.Update)
	adminOnly.Patch("/users/:id/password", userHandler.ResetPassword)
	adminOnly.Patch("/settings", settingsHandler.Update)
	adminOnly.Get("/system/status", systemHandler.GetStatus)
	adminOnly.Post("/import/members/preview", importHandler.PreviewMembers)
	adminOnly.Post("/import/members/commit", importHandler.CommitMembers)
	adminOnly.Post("/import/wallet/preview", importHandler.PreviewWallet)
	adminOnly.Post("/import/wallet/commit", importHandler.CommitWallet)

	adminOnly.Get("/admin/ping", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "admin access granted", fiber.Map{
			"status": "ok",
		})
	})
}
