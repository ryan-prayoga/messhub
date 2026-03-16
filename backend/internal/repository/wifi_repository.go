package repository

import (
	"context"
	"database/sql"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type WifiRepository struct {
	db *sql.DB
}

func NewWifiRepository(db *sql.DB) *WifiRepository {
	return &WifiRepository{db: db}
}

type CreateWifiBillParams struct {
	Month            int
	Year             int
	NominalPerPerson int64
	DeadlineDate     string
	Status           string
	CreatedBy        string
}

func (r *WifiRepository) FindByID(ctx context.Context, billID string) (*models.WifiBill, error) {
	query := `
		SELECT id, month, year, nominal_per_person, deadline_date, status, created_by, created_at, updated_at
		FROM wifi_bills
		WHERE id = $1
		LIMIT 1
	`

	return scanWifiBill(r.db.QueryRowContext(ctx, query, billID))
}

func (r *WifiRepository) FindByMonthYear(ctx context.Context, month int, year int) (*models.WifiBill, error) {
	query := `
		SELECT id, month, year, nominal_per_person, deadline_date, status, created_by, created_at, updated_at
		FROM wifi_bills
		WHERE month = $1 AND year = $2
		LIMIT 1
	`

	return scanWifiBill(r.db.QueryRowContext(ctx, query, month, year))
}

func (r *WifiRepository) FindActiveBill(ctx context.Context) (*models.WifiBill, error) {
	query := `
		SELECT id, month, year, nominal_per_person, deadline_date, status, created_by, created_at, updated_at
		FROM wifi_bills
		WHERE status = 'active'
		ORDER BY year DESC, month DESC, created_at DESC
		LIMIT 1
	`

	return scanWifiBill(r.db.QueryRowContext(ctx, query))
}

func (r *WifiRepository) ListBills(ctx context.Context) ([]models.WifiBillWithSummary, error) {
	query := `
		SELECT
			wb.id,
			wb.month,
			wb.year,
			wb.nominal_per_person,
			wb.deadline_date,
			wb.status,
			wb.created_by,
			wb.created_at,
			wb.updated_at,
			COUNT(wbm.id)::integer AS total_members,
			COUNT(wbm.id) FILTER (WHERE wbm.payment_status = 'verified')::integer AS verified_count,
			COUNT(wbm.id) FILTER (WHERE wbm.payment_status = 'pending_verification')::integer AS pending_count,
			COUNT(wbm.id) FILTER (WHERE wbm.payment_status = 'unpaid')::integer AS unpaid_count,
			COUNT(wbm.id) FILTER (WHERE wbm.payment_status = 'rejected')::integer AS rejected_count,
			COALESCE(SUM(wbm.amount) FILTER (WHERE wbm.payment_status = 'verified'), 0) AS total_collected,
			COALESCE(SUM(wbm.amount), 0) AS total_target
		FROM wifi_bills wb
		LEFT JOIN wifi_bill_members wbm ON wbm.wifi_bill_id = wb.id
		GROUP BY wb.id
		ORDER BY wb.year DESC, wb.month DESC, wb.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.WifiBillWithSummary, 0)
	for rows.Next() {
		item, err := scanWifiBillWithSummary(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *WifiRepository) GetBillSummary(ctx context.Context, billID string) (*models.WifiBillSummary, error) {
	query := `
		SELECT
			COUNT(id)::integer AS total_members,
			COUNT(id) FILTER (WHERE payment_status = 'verified')::integer AS verified_count,
			COUNT(id) FILTER (WHERE payment_status = 'pending_verification')::integer AS pending_count,
			COUNT(id) FILTER (WHERE payment_status = 'unpaid')::integer AS unpaid_count,
			COUNT(id) FILTER (WHERE payment_status = 'rejected')::integer AS rejected_count,
			COALESCE(SUM(amount) FILTER (WHERE payment_status = 'verified'), 0) AS total_collected,
			COALESCE(SUM(amount), 0) AS total_target
		FROM wifi_bill_members
		WHERE wifi_bill_id = $1
	`

	summary := &models.WifiBillSummary{}
	if err := r.db.QueryRowContext(ctx, query, billID).Scan(
		&summary.TotalMembers,
		&summary.VerifiedCount,
		&summary.PendingCount,
		&summary.UnpaidCount,
		&summary.RejectedCount,
		&summary.TotalCollected,
		&summary.TotalTarget,
	); err != nil {
		return nil, err
	}

	return summary, nil
}

func (r *WifiRepository) ListBillMembers(ctx context.Context, billID string) ([]models.WifiBillMemberDetail, error) {
	query := `
		SELECT
			wbm.id,
			wbm.wifi_bill_id,
			wbm.user_id,
			wbm.amount,
			wbm.payment_status,
			wbm.proof_url,
			wbm.note,
			wbm.submitted_at,
			wbm.verified_at,
			wbm.verified_by,
			wbm.rejection_reason,
			wbm.created_at,
			wbm.updated_at,
			u.name,
			u.email,
			verifier.name
		FROM wifi_bill_members wbm
		JOIN users u ON u.id = wbm.user_id
		LEFT JOIN users verifier ON verifier.id = wbm.verified_by
		WHERE wbm.wifi_bill_id = $1
		ORDER BY
			CASE wbm.payment_status
				WHEN 'pending_verification' THEN 0
				WHEN 'rejected' THEN 1
				WHEN 'unpaid' THEN 2
				WHEN 'verified' THEN 3
				ELSE 4
			END,
			u.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.WifiBillMemberDetail, 0)
	for rows.Next() {
		item, err := scanWifiBillMemberDetail(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *WifiRepository) ListMyBills(ctx context.Context, userID string) ([]models.WifiMyBill, error) {
	query := `
		SELECT
			wbm.id,
			wbm.wifi_bill_id,
			wb.month,
			wb.year,
			wb.nominal_per_person,
			wb.deadline_date,
			wb.status,
			wbm.amount,
			wbm.payment_status,
			wbm.proof_url,
			wbm.note,
			wbm.submitted_at,
			wbm.verified_at,
			wbm.rejection_reason,
			wbm.verified_by,
			verifier.name,
			wbm.created_at,
			wbm.updated_at
		FROM wifi_bill_members wbm
		JOIN wifi_bills wb ON wb.id = wbm.wifi_bill_id
		LEFT JOIN users verifier ON verifier.id = wbm.verified_by
		WHERE wbm.user_id = $1
		ORDER BY wb.year DESC, wb.month DESC, wbm.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.WifiMyBill, 0)
	for rows.Next() {
		item, err := scanWifiMyBill(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *WifiRepository) FindBillMemberByBillAndUser(ctx context.Context, billID string, userID string) (*models.WifiBillMember, error) {
	query := `
		SELECT
			id,
			wifi_bill_id,
			user_id,
			amount,
			payment_status,
			proof_url,
			note,
			submitted_at,
			verified_at,
			verified_by,
			rejection_reason,
			created_at,
			updated_at
		FROM wifi_bill_members
		WHERE wifi_bill_id = $1 AND user_id = $2
		LIMIT 1
	`

	return scanWifiBillMember(r.db.QueryRowContext(ctx, query, billID, userID))
}

func (r *WifiRepository) FindBillMemberByID(ctx context.Context, billID string, memberID string) (*models.WifiBillMember, error) {
	query := `
		SELECT
			id,
			wifi_bill_id,
			user_id,
			amount,
			payment_status,
			proof_url,
			note,
			submitted_at,
			verified_at,
			verified_by,
			rejection_reason,
			created_at,
			updated_at
		FROM wifi_bill_members
		WHERE wifi_bill_id = $1 AND id = $2
		LIMIT 1
	`

	return scanWifiBillMember(r.db.QueryRowContext(ctx, query, billID, memberID))
}

func (r *WifiRepository) CreateBillTx(ctx context.Context, tx *sql.Tx, params CreateWifiBillParams) (*models.WifiBill, error) {
	query := `
		INSERT INTO wifi_bills (month, year, nominal_per_person, deadline_date, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, month, year, nominal_per_person, deadline_date, status, created_by, created_at, updated_at
	`

	return scanWifiBill(tx.QueryRowContext(
		ctx,
		query,
		params.Month,
		params.Year,
		params.NominalPerPerson,
		params.DeadlineDate,
		params.Status,
		params.CreatedBy,
	))
}

func (r *WifiRepository) GenerateMembersTx(ctx context.Context, tx *sql.Tx, billID string, amount int64) (int, error) {
	query := `
		INSERT INTO wifi_bill_members (wifi_bill_id, user_id, amount, payment_status)
		SELECT $1, id, $2, 'unpaid'
		FROM users
		WHERE is_active = TRUE
	`

	result, err := tx.ExecContext(ctx, query, billID, amount)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(affected), nil
}

type UpdateWifiBillMemberSubmissionParams struct {
	BillID   string
	UserID   string
	ProofURL string
	Note     *string
}

func (r *WifiRepository) UpdateMemberSubmissionTx(ctx context.Context, tx *sql.Tx, params UpdateWifiBillMemberSubmissionParams) (*models.WifiBillMember, error) {
	query := `
		UPDATE wifi_bill_members
		SET
			proof_url = $3,
			note = $4,
			payment_status = 'pending_verification',
			submitted_at = NOW(),
			verified_at = NULL,
			verified_by = NULL,
			rejection_reason = NULL,
			updated_at = NOW()
		WHERE wifi_bill_id = $1 AND user_id = $2
		RETURNING
			id,
			wifi_bill_id,
			user_id,
			amount,
			payment_status,
			proof_url,
			note,
			submitted_at,
			verified_at,
			verified_by,
			rejection_reason,
			created_at,
			updated_at
	`

	return scanWifiBillMember(tx.QueryRowContext(ctx, query, params.BillID, params.UserID, params.ProofURL, params.Note))
}

type ReviewWifiBillMemberParams struct {
	BillID          string
	MemberID        string
	PaymentStatus   string
	VerifiedBy      *string
	RejectionReason *string
}

func (r *WifiRepository) ReviewMemberTx(ctx context.Context, tx *sql.Tx, params ReviewWifiBillMemberParams) (*models.WifiBillMember, error) {
	query := `
		UPDATE wifi_bill_members
		SET
			payment_status = $3,
			verified_at = CASE WHEN $3 = 'verified' THEN NOW() ELSE NULL END,
			verified_by = $4,
			rejection_reason = $5,
			updated_at = NOW()
		WHERE wifi_bill_id = $1 AND id = $2
		RETURNING
			id,
			wifi_bill_id,
			user_id,
			amount,
			payment_status,
			proof_url,
			note,
			submitted_at,
			verified_at,
			verified_by,
			rejection_reason,
			created_at,
			updated_at
	`

	return scanWifiBillMember(tx.QueryRowContext(
		ctx,
		query,
		params.BillID,
		params.MemberID,
		params.PaymentStatus,
		params.VerifiedBy,
		params.RejectionReason,
	))
}

func scanWifiBill(row scanner) (*models.WifiBill, error) {
	bill := &models.WifiBill{}
	if err := row.Scan(
		&bill.ID,
		&bill.Month,
		&bill.Year,
		&bill.NominalPerPerson,
		&bill.DeadlineDate,
		&bill.Status,
		&bill.CreatedBy,
		&bill.CreatedAt,
		&bill.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return bill, nil
}

func scanWifiBillMember(row scanner) (*models.WifiBillMember, error) {
	member := &models.WifiBillMember{}
	var proofURL sql.NullString
	var note sql.NullString
	var submittedAt sql.NullTime
	var verifiedAt sql.NullTime
	var verifiedBy sql.NullString
	var rejectionReason sql.NullString

	if err := row.Scan(
		&member.ID,
		&member.WifiBillID,
		&member.UserID,
		&member.Amount,
		&member.PaymentStatus,
		&proofURL,
		&note,
		&submittedAt,
		&verifiedAt,
		&verifiedBy,
		&rejectionReason,
		&member.CreatedAt,
		&member.UpdatedAt,
	); err != nil {
		return nil, err
	}

	member.ProofURL = nullStringPtr(proofURL)
	member.Note = nullStringPtr(note)
	member.SubmittedAt = nullTimePtr(submittedAt)
	member.VerifiedAt = nullTimePtr(verifiedAt)
	member.VerifiedBy = nullStringPtr(verifiedBy)
	member.RejectionReason = nullStringPtr(rejectionReason)

	return member, nil
}

func scanWifiBillWithSummary(row scanner) (*models.WifiBillWithSummary, error) {
	item := &models.WifiBillWithSummary{}
	if err := row.Scan(
		&item.ID,
		&item.Month,
		&item.Year,
		&item.NominalPerPerson,
		&item.DeadlineDate,
		&item.Status,
		&item.CreatedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.Summary.TotalMembers,
		&item.Summary.VerifiedCount,
		&item.Summary.PendingCount,
		&item.Summary.UnpaidCount,
		&item.Summary.RejectedCount,
		&item.Summary.TotalCollected,
		&item.Summary.TotalTarget,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func scanWifiBillMemberDetail(row scanner) (*models.WifiBillMemberDetail, error) {
	item := &models.WifiBillMemberDetail{}
	var proofURL sql.NullString
	var note sql.NullString
	var submittedAt sql.NullTime
	var verifiedAt sql.NullTime
	var verifiedBy sql.NullString
	var rejectionReason sql.NullString
	var verifiedByName sql.NullString

	if err := row.Scan(
		&item.ID,
		&item.WifiBillID,
		&item.UserID,
		&item.Amount,
		&item.PaymentStatus,
		&proofURL,
		&note,
		&submittedAt,
		&verifiedAt,
		&verifiedBy,
		&rejectionReason,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.UserName,
		&item.UserEmail,
		&verifiedByName,
	); err != nil {
		return nil, err
	}

	item.ProofURL = nullStringPtr(proofURL)
	item.Note = nullStringPtr(note)
	item.SubmittedAt = nullTimePtr(submittedAt)
	item.VerifiedAt = nullTimePtr(verifiedAt)
	item.VerifiedBy = nullStringPtr(verifiedBy)
	item.RejectionReason = nullStringPtr(rejectionReason)
	item.VerifiedByName = nullStringPtr(verifiedByName)

	return item, nil
}

func scanWifiMyBill(row scanner) (*models.WifiMyBill, error) {
	item := &models.WifiMyBill{}
	var proofURL sql.NullString
	var note sql.NullString
	var submittedAt sql.NullTime
	var verifiedAt sql.NullTime
	var rejectionReason sql.NullString
	var verifiedBy sql.NullString
	var verifiedByName sql.NullString

	if err := row.Scan(
		&item.MemberID,
		&item.WifiBillID,
		&item.Month,
		&item.Year,
		&item.NominalPerPerson,
		&item.DeadlineDate,
		&item.BillStatus,
		&item.Amount,
		&item.PaymentStatus,
		&proofURL,
		&note,
		&submittedAt,
		&verifiedAt,
		&rejectionReason,
		&verifiedBy,
		&verifiedByName,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	item.ProofURL = nullStringPtr(proofURL)
	item.Note = nullStringPtr(note)
	item.SubmittedAt = nullTimePtr(submittedAt)
	item.VerifiedAt = nullTimePtr(verifiedAt)
	item.RejectionReason = nullStringPtr(rejectionReason)
	item.VerifiedBy = nullStringPtr(verifiedBy)
	item.VerifiedByName = nullStringPtr(verifiedByName)

	return item, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}
