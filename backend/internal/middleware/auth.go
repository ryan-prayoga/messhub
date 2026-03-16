package middleware

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

type AuthMiddleware struct {
	config         config.Config
	userRepository *repository.UserRepository
}

func NewAuthMiddleware(cfg config.Config, userRepository *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		config:         cfg,
		userRepository: userRepository,
	}
}

func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" {
			return response.Unauthorized(c, "authentication required")
		}

		if !strings.HasPrefix(authorization, "Bearer ") {
			return response.Unauthorized(c, "authentication required")
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(m.config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return response.Unauthorized(c, "authentication required")
		}

		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok {
			return response.Unauthorized(c, "authentication required")
		}

		user, err := m.userRepository.FindByID(c.UserContext(), claims.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return response.Unauthorized(c, "authentication required")
			}

			return response.InternalServerError(c, "failed to resolve current user")
		}

		if !user.IsActive {
			return response.Forbidden(c, "insufficient permissions")
		}

		c.Locals("user", types.AuthUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
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
		user, ok := c.Locals("user").(types.AuthUser)
		if !ok {
			return response.Unauthorized(c, "authentication required")
		}

		if _, found := allowed[user.Role]; !found {
			return response.Forbidden(c, "insufficient permissions")
		}

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return RequireRoles(roles...)
}
