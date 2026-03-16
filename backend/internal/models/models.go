package models

import "time"

const (
	RoleAdmin     = "admin"
	RoleTreasurer = "treasurer"
	RoleMember    = "member"
)

func IsValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleTreasurer, RoleMember:
		return true
	default:
		return false
	}
}

type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	JoinedAt     *time.Time `json:"joined_at"`
	LeftAt       *time.Time `json:"left_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type WalletTransaction struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Category      string    `json:"category"`
	Amount        int64     `json:"amount"`
	Description   string    `json:"description"`
	CreatedBy     string    `json:"created_by"`
	CreatedByName string    `json:"created_by_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WalletSummary struct {
	Balance      int64 `json:"balance"`
	TotalIncome  int64 `json:"total_income"`
	TotalExpense int64 `json:"total_expense"`
}

type WifiBill struct {
	ID               string    `json:"id"`
	Month            int       `json:"month"`
	Year             int       `json:"year"`
	NominalPerPerson int64     `json:"nominal_per_person"`
	DeadlineDate     time.Time `json:"deadline_date"`
	Status           string    `json:"status"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
}

type WifiBillMember struct {
	ID              string     `json:"id"`
	WifiBillID      string     `json:"wifi_bill_id"`
	UserID          string     `json:"user_id"`
	Amount          int64      `json:"amount"`
	PaymentStatus   string     `json:"payment_status"`
	ProofURL        *string    `json:"proof_url"`
	SubmittedAt     *time.Time `json:"submitted_at"`
	VerifiedAt      *time.Time `json:"verified_at"`
	VerifiedBy      *string    `json:"verified_by"`
	RejectionReason *string    `json:"rejection_reason"`
}

type SharedExpense struct {
	ID           string    `json:"id"`
	ExpenseDate  time.Time `json:"expense_date"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	Amount       int64     `json:"amount"`
	PaidByUserID string    `json:"paid_by_user_id"`
	Status       string    `json:"status"`
	Notes        *string   `json:"notes"`
	ProofURL     *string   `json:"proof_url"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type Contribution struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Category      string    `json:"category"`
	Description   string    `json:"description"`
	Points        int       `json:"points"`
	ContributedAt time.Time `json:"contributed_at"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type Proposal struct {
	ID                string     `json:"id"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	CreatedBy         string     `json:"created_by"`
	VotingStart       *time.Time `json:"voting_start"`
	VotingEnd         *time.Time `json:"voting_end"`
	Status            string     `json:"status"`
	FinalDecisionBy   *string    `json:"final_decision_by"`
	FinalDecisionNote *string    `json:"final_decision_note"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ProposalVote struct {
	ID         string    `json:"id"`
	ProposalID string    `json:"proposal_id"`
	UserID     string    `json:"user_id"`
	VoteType   string    `json:"vote_type"`
	CreatedAt  time.Time `json:"created_at"`
}

type Post struct {
	ID           string     `json:"id"`
	Type         string     `json:"type"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	LocationNote *string    `json:"location_note"`
	CreatedBy    string     `json:"created_by"`
	ExpiresAt    *time.Time `json:"expires_at"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PostComment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostReaction struct {
	ID           string    `json:"id"`
	PostID       string    `json:"post_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	CreatedAt    time.Time `json:"created_at"`
}

type Notification struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	Type          string    `json:"type"`
	ReferenceType *string   `json:"reference_type"`
	ReferenceID   *string   `json:"reference_id"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
}

type AuditLog struct {
	ID         string    `json:"id"`
	UserID     *string   `json:"user_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   *string   `json:"entity_id"`
	OldValue   []byte    `json:"old_value"`
	NewValue   []byte    `json:"new_value"`
	CreatedAt  time.Time `json:"created_at"`
}
