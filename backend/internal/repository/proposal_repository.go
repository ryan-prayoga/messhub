package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type ProposalRepository struct {
	db *sql.DB
}

type proposalQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type CreateProposalParams struct {
	Title       string
	Description string
	CreatedBy   string
	VotingStart *time.Time
	VotingEnd   *time.Time
}

type FinalizeProposalParams struct {
	ProposalID        string
	Status            string
	FinalDecisionBy   *string
	FinalDecisionNote *string
}

type CreateProposalVoteParams struct {
	ProposalID string
	UserID     string
	VoteType   string
}

func NewProposalRepository(db *sql.DB) *ProposalRepository {
	return &ProposalRepository{db: db}
}

func (r *ProposalRepository) List(ctx context.Context, viewerID string) ([]models.Proposal, error) {
	query := proposalSelectQuery() + `
		GROUP BY
			p.id,
			creator.name,
			finalizer.name
		ORDER BY
			CASE p.status
				WHEN 'active' THEN 0
				WHEN 'approved' THEN 1
				WHEN 'rejected' THEN 2
				ELSE 3
			END,
			p.updated_at DESC,
			p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Proposal, 0)
	for rows.Next() {
		item, err := scanProposal(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	return items, rows.Err()
}

func (r *ProposalRepository) FindByID(ctx context.Context, proposalID string, viewerID string) (*models.Proposal, error) {
	query := proposalSelectQuery() + `
		WHERE p.id = $2
		GROUP BY
			p.id,
			creator.name,
			finalizer.name
		LIMIT 1
	`

	return scanProposal(r.db.QueryRowContext(ctx, query, viewerID, proposalID))
}

func (r *ProposalRepository) ListVotes(ctx context.Context, proposalID string) ([]models.ProposalVote, error) {
	query := `
		SELECT
			pv.id,
			pv.proposal_id,
			pv.user_id,
			u.name,
			pv.vote_type,
			pv.created_at
		FROM proposal_votes pv
		JOIN users u ON u.id = pv.user_id
		WHERE pv.proposal_id = $1
		ORDER BY pv.created_at ASC, u.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, proposalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.ProposalVote, 0)
	for rows.Next() {
		item, err := scanProposalVote(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	return items, rows.Err()
}

func (r *ProposalRepository) CreateTx(ctx context.Context, tx *sql.Tx, params CreateProposalParams, viewerID string) (*models.Proposal, error) {
	query := `
		INSERT INTO proposals (title, description, created_by, voting_start, voting_end, status)
		VALUES ($1, $2, $3, $4, $5, 'active')
		RETURNING id
	`

	var proposalID string
	if err := tx.QueryRowContext(
		ctx,
		query,
		params.Title,
		params.Description,
		params.CreatedBy,
		nullableTime(params.VotingStart),
		nullableTime(params.VotingEnd),
	).Scan(&proposalID); err != nil {
		return nil, err
	}

	return r.findByID(ctx, tx, proposalID, viewerID)
}

func (r *ProposalRepository) CloseTx(ctx context.Context, tx *sql.Tx, proposalID string, viewerID string) (*models.Proposal, error) {
	query := `
		UPDATE proposals
		SET
			status = 'closed',
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := tx.ExecContext(ctx, query, proposalID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.findByID(ctx, tx, proposalID, viewerID)
}

func (r *ProposalRepository) FinalizeTx(ctx context.Context, tx *sql.Tx, params FinalizeProposalParams, viewerID string) (*models.Proposal, error) {
	query := `
		UPDATE proposals
		SET
			status = $2,
			final_decision_by = $3,
			final_decision_note = $4,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		params.ProposalID,
		params.Status,
		params.FinalDecisionBy,
		nullableString(params.FinalDecisionNote),
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.findByID(ctx, tx, params.ProposalID, viewerID)
}

func (r *ProposalRepository) CreateVoteTx(ctx context.Context, tx *sql.Tx, params CreateProposalVoteParams) (*models.ProposalVote, error) {
	query := `
		INSERT INTO proposal_votes (proposal_id, user_id, vote_type)
		VALUES ($1, $2, $3)
		RETURNING id, proposal_id, user_id, vote_type, created_at
	`

	item := &models.ProposalVote{}
	if err := tx.QueryRowContext(ctx, query, params.ProposalID, params.UserID, params.VoteType).Scan(
		&item.ID,
		&item.ProposalID,
		&item.UserID,
		&item.VoteType,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func proposalSelectQuery() string {
	return proposalSelectFrom("proposals")
}

func proposalSelectFrom(source string) string {
	return `
		SELECT
			p.id,
			p.title,
			p.description,
			p.created_by,
			creator.name,
			p.voting_start,
			p.voting_end,
			p.status,
			p.final_decision_by,
			finalizer.name,
			p.final_decision_note,
			COALESCE(COUNT(pv.id) FILTER (WHERE pv.vote_type = 'agree'), 0)::integer AS agree_count,
			COALESCE(COUNT(pv.id) FILTER (WHERE pv.vote_type = 'disagree'), 0)::integer AS disagree_count,
			COALESCE(COUNT(pv.id), 0)::integer AS total_votes,
			MAX(CASE WHEN pv.user_id = $1 THEN pv.vote_type END) AS current_user_vote,
			p.created_at,
			p.updated_at
		FROM ` + source + ` p
		JOIN users creator ON creator.id = p.created_by
		LEFT JOIN users finalizer ON finalizer.id = p.final_decision_by
		LEFT JOIN proposal_votes pv ON pv.proposal_id = p.id
	`
}

func (r *ProposalRepository) findByID(ctx context.Context, runner proposalQueryRunner, proposalID string, viewerID string) (*models.Proposal, error) {
	query := proposalSelectQuery() + `
		WHERE p.id = $2
		GROUP BY
			p.id,
			creator.name,
			finalizer.name
		LIMIT 1
	`

	return scanProposal(runner.QueryRowContext(ctx, query, viewerID, proposalID))
}

func scanProposal(row scanner) (*models.Proposal, error) {
	item := &models.Proposal{}
	var votingStart sql.NullTime
	var votingEnd sql.NullTime
	var finalDecisionBy sql.NullString
	var finalDecisionByName sql.NullString
	var finalDecisionNote sql.NullString
	var currentUserVote sql.NullString
	if err := row.Scan(
		&item.ID,
		&item.Title,
		&item.Description,
		&item.CreatedBy,
		&item.CreatedByName,
		&votingStart,
		&votingEnd,
		&item.Status,
		&finalDecisionBy,
		&finalDecisionByName,
		&finalDecisionNote,
		&item.AgreeCount,
		&item.DisagreeCount,
		&item.TotalVotes,
		&currentUserVote,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	item.VotingStart = nullTimePtr(votingStart)
	item.VotingEnd = nullTimePtr(votingEnd)
	item.FinalDecisionBy = nullStringPtr(finalDecisionBy)
	item.FinalDecisionByName = nullStringPtr(finalDecisionByName)
	item.FinalDecisionNote = nullStringPtr(finalDecisionNote)
	item.CurrentUserVote = nullStringPtr(currentUserVote)

	return item, nil
}

func scanProposalVote(row scanner) (*models.ProposalVote, error) {
	item := &models.ProposalVote{}
	if err := row.Scan(
		&item.ID,
		&item.ProposalID,
		&item.UserID,
		&item.UserName,
		&item.VoteType,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}
