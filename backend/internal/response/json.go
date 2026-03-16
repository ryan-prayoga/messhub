package response

import (
	"errors"
	"fmt"
	"strings"

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
		"message": message,
		"data":    nil,
		"error": fiber.Map{
			"code": code,
		},
	}

	if details != nil {
		payload["error"].(fiber.Map)["details"] = details
	}

	return c.Status(status).JSON(payload)
}

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(
			c,
			fiberErr.Code,
			fiberErr.Message,
			statusCodeToErrorCode(fiberErr.Code, fiberErr.Message),
		)
	}

	return Error(c, fiber.StatusInternalServerError, "internal server error", "internal_server_error")
}

func statusCodeToErrorCode(status int, message string) string {
	if status == fiber.StatusNotFound {
		return "not_found"
	}

	if status == fiber.StatusMethodNotAllowed {
		return "method_not_allowed"
	}

	if status == fiber.StatusUnauthorized {
		return "unauthorized"
	}

	if status == fiber.StatusForbidden {
		return "forbidden"
	}

	if status == fiber.StatusBadRequest {
		return "bad_request"
	}

	normalized := strings.ToLower(strings.TrimSpace(message))
	normalized = strings.ReplaceAll(normalized, "-", " ")
	normalized = strings.ReplaceAll(normalized, ".", "")
	code := strings.Join(strings.Fields(normalized), "_")
	if code == "" {
		return fmt.Sprintf("http_%d", status)
	}

	if len(code) > 64 {
		return fmt.Sprintf("http_%d", status)
	}

	return code
}
