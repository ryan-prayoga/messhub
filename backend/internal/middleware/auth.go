package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

func RequireAuth(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authorization, "Bearer ")
		if tokenString == authorization {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid bearer token",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token claims",
			})
		}

		c.Locals("user", fiber.Map{
			"id":    claims.UserID,
			"email": claims.Email,
			"name":  claims.Name,
			"role":  claims.Role,
		})

		return c.Next()
	}
}

func RequireRoles(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(fiber.Map)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing authenticated user",
			})
		}

		role, _ := user["role"].(string)
		if _, found := allowed[role]; !found {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "insufficient role",
			})
		}

		return c.Next()
	}
}
