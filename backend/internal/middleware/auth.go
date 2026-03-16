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
			return response.Error(c, fiber.StatusUnauthorized, "missing authorization header", "missing_authorization")
		}

		if !strings.HasPrefix(authorization, "Bearer ") {
			return response.Error(c, fiber.StatusUnauthorized, "invalid bearer token", "invalid_bearer_token")
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(m.config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, "invalid token", "invalid_token")
		}

		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "invalid token claims", "invalid_token_claims")
		}

		user, err := m.userRepository.FindByID(c.UserContext(), claims.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return response.Error(c, fiber.StatusUnauthorized, "user not found", "user_not_found")
			}

			return response.Error(c, fiber.StatusInternalServerError, "failed to resolve current user", "user_lookup_failed")
		}

		if !user.IsActive {
			return response.Error(c, fiber.StatusUnauthorized, "user is inactive", "user_inactive")
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
			return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
		}

		if _, found := allowed[user.Role]; !found {
			return response.Error(c, fiber.StatusForbidden, "insufficient role", "insufficient_role")
		}

		return c.Next()
	}
}
