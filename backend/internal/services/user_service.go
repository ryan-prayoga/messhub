package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRole                    = errors.New("invalid role")
	ErrInvalidUserInput               = errors.New("invalid user input")
	ErrInvalidUsername                = errors.New("invalid username")
	ErrInvalidProfileInput            = errors.New("invalid profile input")
	ErrPasswordTooShort               = errors.New("password must be at least 8 characters")
	ErrCurrentPasswordRequired        = errors.New("current password is required")
	ErrNewPasswordRequired            = errors.New("new password is required")
	ErrCurrentPasswordInvalid         = errors.New("current password is invalid")
	ErrUserAlreadyExists              = errors.New("user email already exists")
	ErrUsernameAlreadyExists          = errors.New("username already exists")
	ErrUserNotFound                   = errors.New("user not found")
	ErrSelfDeactivateBlocked          = errors.New("cannot deactivate your own account")
	ErrSelfDemoteBlocked              = errors.New("cannot remove your own admin role")
	ErrSelfArchiveBlocked             = errors.New("cannot archive your own account")
	ErrSelfDeleteBlocked              = errors.New("cannot delete your own account")
	ErrUserAlreadyArchived            = errors.New("user is already archived")
	ErrUserNotArchived                = errors.New("user is not archived")
	ErrArchivedUserReactivate         = errors.New("archived user must be reactivated through lifecycle action")
	ErrPermanentDeleteBlocked         = errors.New("permanent delete is not allowed while relations exist")
	ErrPermanentDeleteRequiresArchive = errors.New("user must be archived before permanent delete")
	ErrLastAdminRequired              = errors.New("at least one active admin must remain")
)

type CreateUserInput struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Username *string `json:"username"`
	Phone    *string `json:"phone"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	IsActive *bool   `json:"is_active"`
	JoinedAt *string `json:"joined_at"`
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
	JoinedAt *string `json:"joined_at"`
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

type AdminResetPasswordInput struct {
	NewPassword string `json:"new_password"`
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

func (s *UserService) CreateUser(ctx context.Context, actorID string, input CreateUserInput) (*models.User, error) {
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

	username, err := resolveRequestedUsername(ctx, s.userRepository, input.Username, name, email)
	if err != nil {
		return nil, err
	}

	phone := normalizeOptionalString(input.Phone)
	joinedAt, err := resolveJoinedAt(input.JoinedAt)
	if err != nil {
		return nil, ErrInvalidUserInput
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := s.userRepository.CreateTx(ctx, tx, repository.CreateUserParams{
		Name:         name,
		Email:        email,
		Username:     username,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     isActive,
		JoinedAt:     joinedAt,
	})
	if err != nil {
		switch uniqueViolationField(err) {
		case "email":
			return nil, ErrUserAlreadyExists
		case "username":
			return nil, ErrUsernameAlreadyExists
		default:
			return nil, err
		}
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "create_user",
		EntityType: "user",
		EntityID:   stringPtr(user.ID),
		NewValue: map[string]any{
			"name":      user.Name,
			"email":     user.Email,
			"username":  user.Username,
			"phone":     user.Phone,
			"role":      user.Role,
			"is_active": user.IsActive,
			"joined_at": user.JoinedAt,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
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

	if input.Email != nil {
		email := normalizeEmail(*input.Email)
		if email == "" {
			return nil, ErrInvalidUserInput
		}

		user.Email = email
		updated = true
	}

	if input.Username != nil {
		username, err := normalizeUsername(*input.Username)
		if err != nil {
			return nil, err
		}

		user.Username = username
		updated = true
	}

	if input.Phone != nil {
		user.Phone = normalizeOptionalString(input.Phone)
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
		if user.ArchivedAt != nil && *input.IsActive {
			return nil, ErrArchivedUserReactivate
		}
		user.IsActive = *input.IsActive
		updated = true
	}

	if input.JoinedAt != nil {
		trimmed := strings.TrimSpace(*input.JoinedAt)
		if trimmed != "" {
			joinedAt, err := parseDateOnly(trimmed)
			if err != nil {
				return nil, ErrInvalidUserInput
			}

			user.JoinedAt = &joinedAt
			updated = true
		}
	}

	if !updated {
		return nil, ErrInvalidUserInput
	}

	if actorID == user.ID {
		if previous.Role == models.RoleAdmin && user.Role != models.RoleAdmin {
			return nil, ErrSelfDemoteBlocked
		}

		if previous.IsActive && !user.IsActive {
			return nil, ErrSelfDeactivateBlocked
		}
	}

	adminRoleRemoved := previous.Role == models.RoleAdmin && user.Role != models.RoleAdmin
	adminDeactivated := previous.Role == models.RoleAdmin && previous.IsActive && !user.IsActive
	if adminRoleRemoved || adminDeactivated {
		activeAdmins, err := s.userRepository.CountActiveAdmins(ctx)
		if err != nil {
			return nil, err
		}

		if activeAdmins <= 1 {
			return nil, ErrLastAdminRequired
		}
	}

	joinedAt := user.JoinedAt

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updatedUser, err := s.userRepository.UpdateTx(ctx, tx, repository.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
		Role:     user.Role,
		IsActive: user.IsActive,
		JoinedAt: joinedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		case uniqueViolationField(err) == "email":
			return nil, ErrUserAlreadyExists
		case uniqueViolationField(err) == "username":
			return nil, ErrUsernameAlreadyExists
		default:
			return nil, err
		}
	}

	if basicUserFieldsChanged(previous, *updatedUser) {
		if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
			UserID:     stringPtr(actorID),
			Action:     "update_user",
			EntityType: "user",
			EntityID:   stringPtr(updatedUser.ID),
			OldValue: map[string]any{
				"name":      previous.Name,
				"email":     previous.Email,
				"username":  previous.Username,
				"phone":     previous.Phone,
				"joined_at": previous.JoinedAt,
			},
			NewValue: map[string]any{
				"name":      updatedUser.Name,
				"email":     updatedUser.Email,
				"username":  updatedUser.Username,
				"phone":     updatedUser.Phone,
				"joined_at": updatedUser.JoinedAt,
			},
		}); err != nil {
			return nil, err
		}
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

func (s *UserService) ArchiveUser(ctx context.Context, actorID string, userID string) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if actorID == user.ID {
		return nil, ErrSelfArchiveBlocked
	}

	if user.ArchivedAt != nil {
		return nil, ErrUserAlreadyArchived
	}

	if user.Role == models.RoleAdmin && user.IsActive {
		activeAdmins, err := s.userRepository.CountActiveAdmins(ctx)
		if err != nil {
			return nil, err
		}

		if activeAdmins <= 1 {
			return nil, ErrLastAdminRequired
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	archivedUser, err := s.userRepository.ArchiveTx(ctx, tx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "archive_user",
		EntityType: "user",
		EntityID:   stringPtr(archivedUser.ID),
		OldValue: map[string]any{
			"is_active":   user.IsActive,
			"left_at":     user.LeftAt,
			"archived_at": user.ArchivedAt,
		},
		NewValue: map[string]any{
			"is_active":   archivedUser.IsActive,
			"left_at":     archivedUser.LeftAt,
			"archived_at": archivedUser.ArchivedAt,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return archivedUser, nil
}

func (s *UserService) ReactivateUser(ctx context.Context, actorID string, userID string) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if user.IsActive && user.ArchivedAt == nil {
		return user, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var updatedUser *models.User
	if user.ArchivedAt != nil {
		updatedUser, err = s.userRepository.ReactivateTx(ctx, tx, user.ID)
	} else {
		updatedUser, err = s.userRepository.UpdateTx(ctx, tx, repository.UpdateUserParams{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Phone:    user.Phone,
			Role:     user.Role,
			IsActive: true,
			JoinedAt: user.JoinedAt,
		})
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "reactivate_user",
		EntityType: "user",
		EntityID:   stringPtr(updatedUser.ID),
		OldValue: map[string]any{
			"is_active":   user.IsActive,
			"left_at":     user.LeftAt,
			"archived_at": user.ArchivedAt,
		},
		NewValue: map[string]any{
			"is_active":   updatedUser.IsActive,
			"left_at":     updatedUser.LeftAt,
			"archived_at": updatedUser.ArchivedAt,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUserPermanent(ctx context.Context, actorID string, userID string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	if actorID == user.ID {
		return ErrSelfDeleteBlocked
	}

	if user.ArchivedAt == nil {
		return ErrPermanentDeleteRequiresArchive
	}

	relationCounts, err := s.userRepository.CountRelations(ctx, user.ID)
	if err != nil {
		return err
	}

	if relationCounts.HasAny() {
		_ = s.auditService.Log(ctx, AuditLogInput{
			UserID:     stringPtr(actorID),
			Action:     "failed_delete_user_due_to_relations",
			EntityType: "user",
			EntityID:   stringPtr(user.ID),
			NewValue: map[string]any{
				"user_id":         user.ID,
				"relation_counts": relationCounts,
			},
		})

		return ErrPermanentDeleteBlocked
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "delete_user_permanent",
		EntityType: "user",
		EntityID:   stringPtr(user.ID),
		OldValue: map[string]any{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"username":    user.Username,
			"role":        user.Role,
			"is_active":   user.IsActive,
			"left_at":     user.LeftAt,
			"archived_at": user.ArchivedAt,
		},
		NewValue: map[string]any{
			"deleted": true,
		},
	}); err != nil {
		return err
	}

	if err := s.userRepository.DeleteTx(ctx, tx, user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	return tx.Commit()
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

func (s *UserService) AdminResetPassword(ctx context.Context, actorID string, userID string, input AdminResetPasswordInput) error {
	newPassword := strings.TrimSpace(input.NewPassword)
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
		UserID:     stringPtr(actorID),
		Action:     "admin_reset_password",
		EntityType: "user",
		EntityID:   stringPtr(user.ID),
		NewValue: map[string]any{
			"changed":  true,
			"admin_id": actorID,
		},
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func resolveRequestedUsername(
	ctx context.Context,
	userRepository *repository.UserRepository,
	value *string,
	name string,
	email string,
) (string, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return userRepository.FindAvailableUsername(ctx, name, email)
	}

	return normalizeUsername(*value)
}

func normalizeUsername(value string) (string, error) {
	candidate := sanitizeUsernameValue(value)
	if candidate == "" || len(candidate) < 3 || len(candidate) > 32 {
		return "", ErrInvalidUsername
	}

	return candidate, nil
}

func resolveJoinedAt(value *string) (time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return time.Now().UTC(), nil
	}

	return parseDateOnly(strings.TrimSpace(*value))
}

func parseDateOnly(value string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}

	return date.UTC(), nil
}

func sanitizeUsernameValue(value string) string {
	var builder strings.Builder
	lastWasHyphen := false

	for _, character := range strings.ToLower(strings.TrimSpace(value)) {
		switch {
		case character >= 'a' && character <= 'z':
			builder.WriteRune(character)
			lastWasHyphen = false
		case character >= '0' && character <= '9':
			builder.WriteRune(character)
			lastWasHyphen = false
		case unicode.IsSpace(character) || character == '-' || character == '_' || character == '.':
			if builder.Len() == 0 || lastWasHyphen {
				continue
			}

			builder.WriteRune('-')
			lastWasHyphen = true
		default:
			if builder.Len() == 0 || lastWasHyphen {
				continue
			}

			builder.WriteRune('-')
			lastWasHyphen = true
		}
	}

	return strings.Trim(builder.String(), "-")
}

func basicUserFieldsChanged(previous models.User, next models.User) bool {
	phoneChanged := (previous.Phone == nil) != (next.Phone == nil)
	if !phoneChanged && previous.Phone != nil && next.Phone != nil {
		phoneChanged = *previous.Phone != *next.Phone
	}

	joinedAtChanged := (previous.JoinedAt == nil) != (next.JoinedAt == nil)
	if !joinedAtChanged && previous.JoinedAt != nil && next.JoinedAt != nil {
		joinedAtChanged = !previous.JoinedAt.Equal(*next.JoinedAt)
	}

	return previous.Name != next.Name ||
		previous.Email != next.Email ||
		previous.Username != next.Username ||
		phoneChanged ||
		joinedAtChanged
}

func uniqueViolationField(err error) string {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return ""
	}

	if pgErr.Code != "23505" {
		return ""
	}

	switch pgErr.ConstraintName {
	case "users_email_key":
		return "email"
	case "idx_users_username_lower":
		return "username"
	}

	detail := strings.ToLower(strings.Join([]string{
		pgErr.ConstraintName,
		pgErr.Detail,
		pgErr.Message,
	}, " "))

	switch {
	case strings.Contains(detail, "username"):
		return "username"
	case strings.Contains(detail, "email"):
		return "email"
	default:
		return ""
	}
}

func isUniqueViolation(err error) bool {
	return uniqueViolationField(err) != ""
}
