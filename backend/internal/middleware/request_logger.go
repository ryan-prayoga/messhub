package middleware

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startedAt := time.Now()
		err := c.Next()

		statusCode := c.Response().StatusCode()
		if statusCode == 0 {
			statusCode = statusCodeFromError(err)
		}

		requestID := RequestID(c)
		userID := ""
		if user, ok := c.Locals("user").(types.AuthUser); ok {
			userID = user.ID
		}

		attrs := []any{
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.OriginalURL()),
			slog.Int("status_code", statusCode),
			slog.Int64("duration_ms", time.Since(startedAt).Milliseconds()),
		}

		if userID != "" {
			attrs = append(attrs, slog.String("user_id", userID))
		}

		if err != nil {
			attrs = append(attrs, slog.String("error", sanitizeLogValue(err.Error())))
		}

		if err != nil || statusCode >= fiber.StatusInternalServerError {
			logger.Error("request completed", attrs...)
		} else {
			logger.Info("request completed", attrs...)
		}

		return err
	}
}

func statusCodeFromError(err error) int {
	if err == nil {
		return fiber.StatusOK
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	return fiber.StatusInternalServerError
}

func sanitizeLogValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	return value
}
