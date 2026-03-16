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
	ErrInvalidRole             = errors.New("invalid role")
	ErrInvalidUserInput        = errors.New("invalid user input")
	ErrInvalidProfileInput     = errors.New("invalid profile input")
	ErrPasswordTooShort        = errors.New("password must be at least 8 characters")
	ErrCurrentPasswordRequired = errors.New("current password is required")
	ErrNewPasswordRequired     = errors.New("new password is required")
	ErrCurrentPasswordInvalid  = errors.New("current password is invalid")
	ErrUserAlreadyExists       = errors.New("user email already exists")
	ErrUserNotFound            = errors.New("user not found")
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

type UpdateProfileInput struct {
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	AvatarURL *string `json:"avatar_url"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserService struct {
	userRepository *repository.UserRepository
	db             *sql.DB
	auditService   *AuditService
}

func NewUserService(db *sql.DB, userRepository *repository.UserRepository, auditService *AuditService) *UserService {
	return &UserService{
		userRepository: userRepository,
		db:             db,
		auditService:   auditService,
	}
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepository.List(ctx)
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
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

	username, err := s.userRepository.FindAvailableUsername(ctx, name, email)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.Create(ctx, repository.CreateUserParams{
		Name:         name,
		Email:        email,
		Username:     username,
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

func (s *UserService) UpdateUser(ctx context.Context, actorID string, userID string, input UpdateUserInput) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	previous := *user
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updatedUser, err := s.userRepository.UpdateTx(ctx, tx, repository.UpdateUserParams{
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

	if previous.Role != updatedUser.Role {
		if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
			UserID:     stringPtr(actorID),
			Action:     "update_user_role",
			EntityType: "user",
			EntityID:   stringPtr(updatedUser.ID),
			OldValue: map[string]any{
				"role": previous.Role,
			},
			NewValue: map[string]any{
				"role": updatedUser.Role,
			},
		}); err != nil {
			return nil, err
		}
	}

	if previous.IsActive != updatedUser.IsActive {
		action := "deactivate_user"
		if updatedUser.IsActive {
			action = "activate_user"
		}

		if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
			UserID:     stringPtr(actorID),
			Action:     action,
			EntityType: "user",
			EntityID:   stringPtr(updatedUser.ID),
			OldValue: map[string]any{
				"is_active": previous.IsActive,
				"left_at":   previous.LeftAt,
			},
			NewValue: map[string]any{
				"is_active": updatedUser.IsActive,
				"left_at":   updatedUser.LeftAt,
			},
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	previous := *user
	updated := false

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrInvalidProfileInput
		}

		user.Name = name
		updated = true
	}

	if input.Phone != nil {
		phone := strings.TrimSpace(*input.Phone)
		if phone == "" {
			user.Phone = nil
		} else {
			user.Phone = &phone
		}
		updated = true
	}

	if input.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*input.AvatarURL)
		if avatarURL == "" {
			user.AvatarURL = nil
		} else {
			user.AvatarURL = &avatarURL
		}
		updated = true
	}

	if !updated {
		return nil, ErrInvalidProfileInput
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updatedUser, err := s.userRepository.UpdateProfileTx(ctx, tx, repository.UpdateProfileParams{
		ID:        user.ID,
		Name:      user.Name,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(userID),
		Action:     "update_profile",
		EntityType: "user",
		EntityID:   stringPtr(updatedUser.ID),
		OldValue: map[string]any{
			"name":       previous.Name,
			"phone":      previous.Phone,
			"avatar_url": previous.AvatarURL,
		},
		NewValue: map[string]any{
			"name":       updatedUser.Name,
			"phone":      updatedUser.Phone,
			"avatar_url": updatedUser.AvatarURL,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID string, input ChangePasswordInput) error {
	currentPassword := strings.TrimSpace(input.CurrentPassword)
	newPassword := strings.TrimSpace(input.NewPassword)

	if currentPassword == "" {
		return ErrCurrentPasswordRequired
	}

	if newPassword == "" {
		return ErrNewPasswordRequired
	}

	if len(newPassword) < 8 {
		return ErrPasswordTooShort
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrCurrentPasswordInvalid
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.userRepository.UpdatePasswordTx(ctx, tx, repository.UpdatePasswordParams{
		ID:           user.ID,
		PasswordHash: string(hashedPassword),
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(userID),
		Action:     "change_password",
		EntityType: "user",
		EntityID:   stringPtr(user.ID),
		NewValue: map[string]any{
			"changed": true,
		},
	}); err != nil {
		return err
	}

	return tx.Commit()
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
