package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

const (
	DefaultWifiNominalPerPerson int64 = 20000
	defaultWifiBillStatus             = models.WifiBillStatusActive
)

var (
	ErrInvalidWifiBillInput     = errors.New("invalid wifi bill input")
	ErrInvalidWifiStatus        = errors.New("invalid wifi bill status")
	ErrDuplicateWifiBill        = errors.New("wifi bill for the selected month already exists")
	ErrWifiBillNotFound         = errors.New("wifi bill not found")
	ErrWifiBillInactive         = errors.New("wifi bill is not active")
	ErrWifiMemberNotFound       = errors.New("wifi member record not found")
	ErrWifiProofRequired        = errors.New("payment proof is required")
	ErrWifiSubmissionNotAllowed = errors.New("payment proof cannot be submitted for this bill")
	ErrWifiReviewNotAllowed     = errors.New("payment cannot be reviewed in its current status")
	ErrWifiRejectReasonRequired = errors.New("rejection reason is required")
	ErrWifiNoActiveMembers      = errors.New("wifi bill requires at least one active member")
)

type CreateWifiBillInput struct {
	Month            int     `json:"month"`
	Year             int     `json:"year"`
	NominalPerPerson *int64  `json:"nominal_per_person"`
	DeadlineDate     *string `json:"deadline_date"`
	Status           string  `json:"status"`
}

type SubmitWifiPaymentInput struct {
	ProofURL string  `json:"proof_url"`
	Note     *string `json:"note"`
}

type RejectWifiPaymentInput struct {
	Reason string `json:"reason"`
}

type WifiService struct {
	db             *sql.DB
	wifiRepository *repository.WifiRepository
	auditService   *AuditService
}

func NewWifiService(db *sql.DB, wifiRepository *repository.WifiRepository, auditService *AuditService) *WifiService {
	return &WifiService{
		db:             db,
		wifiRepository: wifiRepository,
		auditService:   auditService,
	}
}

func (s *WifiService) CreateBill(ctx context.Context, createdBy string, input CreateWifiBillInput) (*models.WifiBillDetail, error) {
	nominal := DefaultWifiNominalPerPerson
	if input.NominalPerPerson != nil {
		nominal = *input.NominalPerPerson
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = defaultWifiBillStatus
	}

	if !isValidWifiBillStatus(status) {
		return nil, ErrInvalidWifiStatus
	}

	if !isValidWifiBillDate(input.Month, input.Year) || nominal <= 0 || strings.TrimSpace(createdBy) == "" {
		return nil, ErrInvalidWifiBillInput
	}

	deadlineDate, err := resolveWifiDeadline(input.Year, input.Month, input.DeadlineDate)
	if err != nil {
		return nil, ErrInvalidWifiBillInput
	}

	if _, err := s.wifiRepository.FindByMonthYear(ctx, input.Month, input.Year); err == nil {
		return nil, ErrDuplicateWifiBill
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bill, err := s.wifiRepository.CreateBillTx(ctx, tx, repository.CreateWifiBillParams{
		Month:            input.Month,
		Year:             input.Year,
		NominalPerPerson: nominal,
		DeadlineDate:     deadlineDate.Format("2006-01-02"),
		Status:           status,
		CreatedBy:        createdBy,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrDuplicateWifiBill
		}

		return nil, err
	}

	createdMembers, err := s.wifiRepository.GenerateMembersTx(ctx, tx, bill.ID, nominal)
	if err != nil {
		return nil, err
	}

	if createdMembers == 0 {
		return nil, ErrWifiNoActiveMembers
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(createdBy),
		Action:     "wifi_bill_created",
		EntityType: "wifi_bill",
		EntityID:   stringPtr(bill.ID),
		NewValue: map[string]any{
			"id":                 bill.ID,
			"month":              bill.Month,
			"year":               bill.Year,
			"nominal_per_person": bill.NominalPerPerson,
			"deadline_date":      bill.DeadlineDate,
			"status":             bill.Status,
			"generated_members":  createdMembers,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetBillDetail(ctx, bill.ID)
}

func (s *WifiService) ListBills(ctx context.Context) ([]models.WifiBillWithSummary, error) {
	return s.wifiRepository.ListBills(ctx)
}

func (s *WifiService) GetActiveBill(ctx context.Context, viewerID string, viewerRole string) (*models.WifiBillDetail, error) {
	bill, err := s.wifiRepository.FindActiveBill(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return s.getBillDetail(ctx, bill, viewerID, viewerRole)
}

func (s *WifiService) GetBillDetail(ctx context.Context, billID string) (*models.WifiBillDetail, error) {
	bill, err := s.wifiRepository.FindByID(ctx, billID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiBillNotFound
		}

		return nil, err
	}

	return s.getBillDetail(ctx, bill, "", models.RoleAdmin)
}

func (s *WifiService) GetBillDetailForViewer(ctx context.Context, billID string, viewerID string, viewerRole string) (*models.WifiBillDetail, error) {
	bill, err := s.wifiRepository.FindByID(ctx, billID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiBillNotFound
		}

		return nil, err
	}

	return s.getBillDetail(ctx, bill, viewerID, viewerRole)
}

func (s *WifiService) ListMyBills(ctx context.Context, userID string) ([]models.WifiMyBill, error) {
	return s.wifiRepository.ListMyBills(ctx, userID)
}

func (s *WifiService) SubmitPaymentProof(ctx context.Context, billID string, userID string, input SubmitWifiPaymentInput) (*models.WifiBillMember, error) {
	proofURL := strings.TrimSpace(input.ProofURL)
	if proofURL == "" {
		return nil, ErrWifiProofRequired
	}

	bill, err := s.wifiRepository.FindByID(ctx, billID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiBillNotFound
		}

		return nil, err
	}

	if bill.Status != models.WifiBillStatusActive {
		return nil, ErrWifiBillInactive
	}

	member, err := s.wifiRepository.FindBillMemberByBillAndUser(ctx, billID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiMemberNotFound
		}

		return nil, err
	}

	if member.PaymentStatus == models.WifiPaymentStatusVerified {
		return nil, ErrWifiSubmissionNotAllowed
	}

	var note *string
	if input.Note != nil {
		trimmedNote := strings.TrimSpace(*input.Note)
		if trimmedNote != "" {
			note = &trimmedNote
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updated, err := s.wifiRepository.UpdateMemberSubmissionTx(ctx, tx, repository.UpdateWifiBillMemberSubmissionParams{
		BillID:   billID,
		UserID:   userID,
		ProofURL: proofURL,
		Note:     note,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiMemberNotFound
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(userID),
		Action:     "wifi_payment_submitted",
		EntityType: "wifi_bill_member",
		EntityID:   stringPtr(updated.ID),
		OldValue:   member,
		NewValue: map[string]any{
			"id":             updated.ID,
			"wifi_bill_id":   updated.WifiBillID,
			"user_id":        updated.UserID,
			"payment_status": updated.PaymentStatus,
			"proof_url":      updated.ProofURL,
			"note":           updated.Note,
			"submitted_at":   updated.SubmittedAt,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *WifiService) VerifyPayment(ctx context.Context, billID string, memberID string, reviewerID string) (*models.WifiBillMember, error) {
	member, err := s.wifiRepository.FindBillMemberByID(ctx, billID, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiMemberNotFound
		}

		return nil, err
	}

	if member.PaymentStatus != models.WifiPaymentStatusPendingVerification {
		return nil, ErrWifiReviewNotAllowed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updated, err := s.wifiRepository.ReviewMemberTx(ctx, tx, repository.ReviewWifiBillMemberParams{
		BillID:        billID,
		MemberID:      memberID,
		PaymentStatus: models.WifiPaymentStatusVerified,
		VerifiedBy:    stringPtr(reviewerID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiMemberNotFound
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(reviewerID),
		Action:     "wifi_payment_verified",
		EntityType: "wifi_bill_member",
		EntityID:   stringPtr(updated.ID),
		OldValue:   member,
		NewValue:   updated,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *WifiService) RejectPayment(ctx context.Context, billID string, memberID string, reviewerID string, input RejectWifiPaymentInput) (*models.WifiBillMember, error) {
	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return nil, ErrWifiRejectReasonRequired
	}

	member, err := s.wifiRepository.FindBillMemberByID(ctx, billID, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiMemberNotFound
		}

		return nil, err
	}

	if member.PaymentStatus != models.WifiPaymentStatusPendingVerification {
		return nil, ErrWifiReviewNotAllowed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updated, err := s.wifiRepository.ReviewMemberTx(ctx, tx, repository.ReviewWifiBillMemberParams{
		BillID:          billID,
		MemberID:        memberID,
		PaymentStatus:   models.WifiPaymentStatusRejected,
		RejectionReason: &reason,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWifiMemberNotFound
		}

		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(reviewerID),
		Action:     "wifi_payment_rejected",
		EntityType: "wifi_bill_member",
		EntityID:   stringPtr(updated.ID),
		OldValue:   member,
		NewValue:   updated,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *WifiService) getBillDetail(ctx context.Context, bill *models.WifiBill, viewerID string, viewerRole string) (*models.WifiBillDetail, error) {
	summary, err := s.wifiRepository.GetBillSummary(ctx, bill.ID)
	if err != nil {
		return nil, err
	}

	members, err := s.wifiRepository.ListBillMembers(ctx, bill.ID)
	if err != nil {
		return nil, err
	}

	if viewerRole == models.RoleMember {
		filtered := make([]models.WifiBillMemberDetail, 0, 1)
		for _, member := range members {
			if member.UserID == viewerID {
				filtered = append(filtered, member)
				break
			}
		}
		members = filtered
	}

	return &models.WifiBillDetail{
		Bill:    *bill,
		Summary: *summary,
		Members: members,
	}, nil
}

func isValidWifiBillStatus(status string) bool {
	switch status {
	case models.WifiBillStatusDraft, models.WifiBillStatusActive, models.WifiBillStatusClosed:
		return true
	default:
		return false
	}
}

func isValidWifiBillDate(month int, year int) bool {
	currentYear := time.Now().Year()

	return month >= 1 && month <= 12 && year >= 2024 && year <= currentYear+5
}

func resolveWifiDeadline(year int, month int, value *string) (time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return time.Date(year, time.Month(month), 10, 0, 0, 0, 0, time.UTC), nil
	}

	deadline, err := time.Parse("2006-01-02", strings.TrimSpace(*value))
	if err != nil {
		return time.Time{}, err
	}

	if deadline.Year() != year || int(deadline.Month()) != month {
		return time.Time{}, fmt.Errorf("deadline must match selected month and year")
	}

	return deadline, nil
}

func stringPtr(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return &value
}
