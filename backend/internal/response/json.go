package response

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}

func Error(c *fiber.Ctx, status int, message string, code string) error {
	payload := fiber.Map{
		"message": message,
	}

	if code != "" {
		payload["error"] = fiber.Map{
			"code": code,
		}
	}

	return c.Status(status).JSON(payload)
}
