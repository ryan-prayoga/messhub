package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
)

var (
	ErrInvalidProposalInput       = errors.New("invalid proposal input")
	ErrInvalidProposalVote        = errors.New("invalid proposal vote")
	ErrProposalNotFound           = errors.New("proposal not found")
	ErrProposalNotActive          = errors.New("proposal is not active")
	ErrProposalVoteExists         = errors.New("proposal vote already exists")
	ErrProposalVotingUnavailable  = errors.New("proposal voting window is unavailable")
	ErrInvalidProposalFinalStatus = errors.New("final proposal status must be approved or rejected")
)

type CreateProposalInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	VotingStart *string `json:"voting_start"`
	VotingEnd   *string `json:"voting_end"`
}

type VoteProposalInput struct {
	VoteType string `json:"vote_type"`
}

type FinalizeProposalInput struct {
	Status            string  `json:"status"`
	FinalDecisionNote *string `json:"final_decision_note"`
}

type ProposalService struct {
	db           *sql.DB
	proposalRepo *repository.ProposalRepository
	auditService *AuditService
}

func NewProposalService(
	db *sql.DB,
	proposalRepo *repository.ProposalRepository,
	auditService *AuditService,
) *ProposalService {
	return &ProposalService{
		db:           db,
		proposalRepo: proposalRepo,
		auditService: auditService,
	}
}

func (s *ProposalService) List(ctx context.Context, viewerID string) ([]models.Proposal, error) {
	return s.proposalRepo.List(ctx, viewerID)
}

func (s *ProposalService) GetDetail(ctx context.Context, proposalID string, viewerID string) (*models.ProposalDetail, error) {
	proposal, err := s.proposalRepo.FindByID(ctx, proposalID, viewerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}

	votes, err := s.proposalRepo.ListVotes(ctx, proposalID)
	if err != nil {
		return nil, err
	}

	return &models.ProposalDetail{
		Proposal: *proposal,
		Votes:    votes,
	}, nil
}

func (s *ProposalService) Create(ctx context.Context, actorID string, input CreateProposalInput) (*models.ProposalDetail, error) {
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)
	if title == "" || description == "" {
		return nil, ErrInvalidProposalInput
	}

	votingStart, err := parseOptionalDateTime(input.VotingStart)
	if err != nil {
		return nil, ErrInvalidProposalInput
	}
	votingEnd, err := parseOptionalDateTime(input.VotingEnd)
	if err != nil {
		return nil, ErrInvalidProposalInput
	}
	if votingStart != nil && votingEnd != nil && !votingEnd.After(*votingStart) {
		return nil, ErrInvalidProposalInput
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	proposal, err := s.proposalRepo.CreateTx(ctx, tx, repository.CreateProposalParams{
		Title:       title,
		Description: description,
		CreatedBy:   actorID,
		VotingStart: votingStart,
		VotingEnd:   votingEnd,
	}, actorID)
	if err != nil {
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "create_proposal",
		EntityType: "proposal",
		EntityID:   stringPtr(proposal.ID),
		NewValue:   proposal,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetDetail(ctx, proposal.ID, actorID)
}

func (s *ProposalService) Vote(ctx context.Context, proposalID string, actorID string, input VoteProposalInput) (*models.ProposalDetail, error) {
	voteType := strings.TrimSpace(input.VoteType)
	if voteType != models.ProposalVoteAgree && voteType != models.ProposalVoteDisagree {
		return nil, ErrInvalidProposalVote
	}

	proposal, err := s.proposalRepo.FindByID(ctx, proposalID, actorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}
	if proposal.Status != models.ProposalStatusActive {
		return nil, ErrProposalNotActive
	}
	if !isProposalVotingOpen(proposal) {
		return nil, ErrProposalVotingUnavailable
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	vote, err := s.proposalRepo.CreateVoteTx(ctx, tx, repository.CreateProposalVoteParams{
		ProposalID: proposalID,
		UserID:     actorID,
		VoteType:   voteType,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrProposalVoteExists
		}
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "vote_proposal",
		EntityType: "proposal",
		EntityID:   stringPtr(proposalID),
		NewValue:   vote,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetDetail(ctx, proposalID, actorID)
}

func (s *ProposalService) Close(ctx context.Context, proposalID string, actorID string) (*models.ProposalDetail, error) {
	proposal, err := s.proposalRepo.FindByID(ctx, proposalID, actorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}
	if proposal.Status != models.ProposalStatusActive {
		return nil, ErrProposalNotActive
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updated, err := s.proposalRepo.CloseTx(ctx, tx, proposalID, actorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "close_proposal",
		EntityType: "proposal",
		EntityID:   stringPtr(updated.ID),
		OldValue: map[string]any{
			"status": proposal.Status,
		},
		NewValue: map[string]any{
			"status": updated.Status,
		},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetDetail(ctx, proposalID, actorID)
}

func (s *ProposalService) Finalize(ctx context.Context, proposalID string, actorID string, input FinalizeProposalInput) (*models.ProposalDetail, error) {
	status := strings.TrimSpace(input.Status)
	if status != models.ProposalStatusApproved && status != models.ProposalStatusRejected {
		return nil, ErrInvalidProposalFinalStatus
	}

	proposal, err := s.proposalRepo.FindByID(ctx, proposalID, actorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}
	if proposal.Status != models.ProposalStatusActive && proposal.Status != models.ProposalStatusClosed {
		return nil, ErrProposalNotActive
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updated, err := s.proposalRepo.FinalizeTx(ctx, tx, repository.FinalizeProposalParams{
		ProposalID:        proposalID,
		Status:            status,
		FinalDecisionBy:   stringPtr(actorID),
		FinalDecisionNote: normalizeOptionalString(input.FinalDecisionNote),
	}, actorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "finalize_proposal",
		EntityType: "proposal",
		EntityID:   stringPtr(updated.ID),
		OldValue: map[string]any{
			"status": proposal.Status,
		},
		NewValue: updated,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetDetail(ctx, proposalID, actorID)
}

func parseOptionalDateTime(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}

	candidate := strings.TrimSpace(*value)
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, candidate); err == nil {
			utc := parsed.UTC()
			return &utc, nil
		}
	}

	return nil, ErrInvalidProposalInput
}

func isProposalVotingOpen(proposal *models.Proposal) bool {
	now := time.Now().UTC()
	if proposal.VotingStart != nil && now.Before(proposal.VotingStart.UTC()) {
		return false
	}
	if proposal.VotingEnd != nil && now.After(proposal.VotingEnd.UTC()) {
		return false
	}

	return true
}
