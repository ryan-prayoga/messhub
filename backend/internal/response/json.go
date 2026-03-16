package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func Success(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}

func Error(c *fiber.Ctx, status int, message string, code string) error {
	return ErrorWithDetails(c, status, message, code, nil)
}

func ErrorWithDetails(c *fiber.Ctx, status int, message string, code string, details any) error {
	payload := fiber.Map{
		"error":   normalizeErrorCode(status, code, details != nil),
		"message": message,
	}

	if details != nil {
		payload["details"] = details
	}

	return c.Status(status).JSON(payload)
}

func InvalidRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message, "invalid_request")
}

func ValidationFailed(c *fiber.Ctx, message string, details any) error {
	return ErrorWithDetails(c, fiber.StatusBadRequest, message, "validation_failed", details)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, "unauthorized")
}

func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, "forbidden")
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, "not_found")
}

func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, message, "conflict")
}

func RateLimited(c *fiber.Ctx, message string, details any) error {
	return ErrorWithDetails(c, fiber.StatusTooManyRequests, message, "rate_limited", details)
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message, "internal_server_error")
}

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(c, fiberErr.Code, fiberErr.Message, "")
	}

	return InternalServerError(c, "internal server error")
}

func normalizeErrorCode(status int, code string, hasDetails bool) string {
	switch status {
	case fiber.StatusBadRequest:
		if code == "validation_failed" || hasDetails {
			return "validation_failed"
		}

		return "invalid_request"
	case fiber.StatusUnauthorized:
		return "unauthorized"
	case fiber.StatusForbidden:
		return "forbidden"
	case fiber.StatusNotFound:
		return "not_found"
	case fiber.StatusConflict:
		return "conflict"
	case fiber.StatusTooManyRequests:
		return "rate_limited"
	default:
		return "internal_server_error"
	}
}
