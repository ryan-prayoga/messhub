package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	RequestIDHeader = "X-Request-ID"
	requestIDLocal  = "request_id"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := strings.TrimSpace(c.Get(RequestIDHeader))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Locals(requestIDLocal, requestID)
		c.Set(RequestIDHeader, requestID)
		c.SetUserContext(context.WithValue(c.UserContext(), requestIDContextKey, requestID))

		return c.Next()
	}
}

func RequestID(c *fiber.Ctx) string {
	value, _ := c.Locals(requestIDLocal).(string)
	return strings.TrimSpace(value)
}

func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDContextKey).(string)
	return strings.TrimSpace(value)
}
