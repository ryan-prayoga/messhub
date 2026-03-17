package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanprayoga/messhub/backend/internal/config"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid identifier or password")
var ErrInactiveUser = errors.New("user is inactive")
var ErrInvalidLoginInput = errors.New("identifier and password are required")

type LoginInput struct {
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token string         `json:"token"`
	User  types.AuthUser `json:"user"`
}

type AuthService struct {
	config         config.Config
	userRepository *repository.UserRepository
}

func NewAuthService(cfg config.Config, userRepository *repository.UserRepository) *AuthService {
	return &AuthService{
		config:         cfg,
		userRepository: userRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	identifier := normalizeLoginIdentifier(input.Identifier)
	if identifier == "" {
		identifier = normalizeLoginIdentifier(input.Email)
	}

	password := strings.TrimSpace(input.Password)
	if identifier == "" || password == "" {
		return nil, ErrInvalidLoginInput
	}

	user, err := s.userRepository.FindByLoginIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !user.IsActive || user.ArchivedAt != nil {
		return nil, ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User: types.AuthUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}

func (s *AuthService) issueToken(user *models.User) (string, error) {
	expiresAt := time.Now().Add(time.Duration(s.config.JWTExpiresInHours) * time.Hour)
	claims := types.JWTClaims{
		UserID:      user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Username:    user.Username,
		Role:        user.Role,
		AuthVersion: user.AuthVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.config.JWTSecret))
}

func normalizeLoginIdentifier(identifier string) string {
	return strings.ToLower(strings.TrimSpace(identifier))
}
