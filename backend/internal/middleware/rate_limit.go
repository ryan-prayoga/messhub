package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

const (
	loginRateLimitMax   = 5
	loginRateLimitTTL   = 5 * time.Minute
	postRateLimitMax    = 20
	postRateLimitTTL    = 1 * time.Minute
	rateLimitHeaderName = "Retry-After"
)

func LoginRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        loginRateLimitMax,
		Expiration: loginRateLimitTTL,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("login:%s", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			c.Set(rateLimitHeaderName, fmt.Sprintf("%d", int(loginRateLimitTTL.Seconds())))
			return response.RateLimited(c, "too many login attempts", fiber.Map{
				"retry_after_seconds": int(loginRateLimitTTL.Seconds()),
			})
		},
	})
}

func PostRateLimit(scope string) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        postRateLimitMax,
		Expiration: postRateLimitTTL,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("%s:%s", scope, rateLimitActor(c))
		},
		LimitReached: func(c *fiber.Ctx) error {
			c.Set(rateLimitHeaderName, fmt.Sprintf("%d", int(postRateLimitTTL.Seconds())))
			return response.RateLimited(c, "too many write requests", fiber.Map{
				"retry_after_seconds": int(postRateLimitTTL.Seconds()),
			})
		},
	})
}

func rateLimitActor(c *fiber.Ctx) string {
	if user, ok := c.Locals("user").(types.AuthUser); ok && user.ID != "" {
		return user.ID
	}

	return c.IP()
}
