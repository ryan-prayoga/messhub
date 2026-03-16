package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRole       = errors.New("invalid role")
	ErrInvalidUserInput  = errors.New("invalid user input")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrUserAlreadyExists = errors.New("user email already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepository.List(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)
	role := strings.TrimSpace(input.Role)

	if name == "" || email == "" || input.Password == "" {
		return nil, ErrInvalidUserInput
	}

	if len(input.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	if !models.IsValidRole(role) {
		return nil, ErrInvalidRole
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	user, err := s.userRepository.Create(ctx, repository.CreateUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     isActive,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUserAlreadyExists
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, input UpdateUserInput) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	updated := false

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrInvalidUserInput
		}

		user.Name = name
		updated = true
	}

	if input.Role != nil {
		role := strings.TrimSpace(*input.Role)
		if !models.IsValidRole(role) {
			return nil, ErrInvalidRole
		}

		user.Role = role
		updated = true
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
		updated = true
	}

	if !updated {
		return nil, ErrInvalidUserInput
	}

	updatedUser, err := s.userRepository.Update(ctx, repository.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Role:     user.Role,
		IsActive: user.IsActive,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return updatedUser, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505"
}
