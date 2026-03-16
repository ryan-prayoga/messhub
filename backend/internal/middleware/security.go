package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const defaultContentSecurityPolicy = "default-src 'self'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'; object-src 'none'"

func SecurityHeaders(contentSecurityPolicy string) fiber.Handler {
	policy := strings.TrimSpace(contentSecurityPolicy)
	if policy == "" {
		policy = defaultContentSecurityPolicy
	}

	return func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Content-Security-Policy", policy)

		return c.Next()
	}
}
