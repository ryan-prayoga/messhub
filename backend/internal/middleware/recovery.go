package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
)

func Recover(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error(
					"panic recovered",
					slog.String("request_id", RequestID(c)),
					slog.String("method", c.Method()),
					slog.String("path", c.OriginalURL()),
					slog.String("panic", fmt.Sprint(recovered)),
					slog.String("stack", string(debug.Stack())),
				)

				err = response.InternalServerError(c, "internal server error")
			}
		}()

		return c.Next()
	}
}
