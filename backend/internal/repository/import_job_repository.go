package repository

import (
	"context"
	"database/sql"

	"github.com/ryanprayoga/messhub/backend/internal/models"
)

type ImportJobRepository struct {
	db *sql.DB
}

type importJobQueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type CreateImportJobParams struct {
	ImportType     string
	Status         string
	Source         string
	FileName       string
	FileHash       string
	CreatedBy      string
	TotalRows      int
	ValidRows      int
	InvalidRows    int
	PreviewPayload []byte
}

type MarkImportJobCommittedParams struct {
	ID                string
	DuplicateStrategy *string
	CommittedRows     int
	SkippedRows       int
	FailedRows        int
	CommitSummary     []byte
}

type MarkImportJobFailedParams struct {
	ID            string
	CommitSummary []byte
}

func NewImportJobRepository(db *sql.DB) *ImportJobRepository {
	return &ImportJobRepository{db: db}
}

func (r *ImportJobRepository) Create(ctx context.Context, params CreateImportJobParams) (*models.ImportJob, error) {
	return r.create(ctx, r.db, params)
}

func (r *ImportJobRepository) CreateTx(ctx context.Context, tx *sql.Tx, params CreateImportJobParams) (*models.ImportJob, error) {
	return r.create(ctx, tx, params)
}

func (r *ImportJobRepository) create(ctx context.Context, runner importJobQueryRunner, params CreateImportJobParams) (*models.ImportJob, error) {
	query := `
		INSERT INTO import_jobs (
			import_type,
			status,
			source,
			file_name,
			file_hash,
			created_by,
			total_rows,
			valid_rows,
			invalid_rows,
			preview_payload
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
			id,
			import_type,
			status,
			source,
			file_name,
			file_hash,
			created_by,
			NULL::TEXT AS created_by_name,
			duplicate_strategy,
			total_rows,
			valid_rows,
			invalid_rows,
			committed_rows,
			skipped_rows,
			failed_rows,
			preview_payload,
			commit_summary,
			committed_at,
			created_at,
			updated_at
	`

	return scanImportJob(
		runner.QueryRowContext(
			ctx,
			query,
			params.ImportType,
			params.Status,
			params.Source,
			params.FileName,
			params.FileHash,
			params.CreatedBy,
			params.TotalRows,
			params.ValidRows,
			params.InvalidRows,
			params.PreviewPayload,
		),
	)
}

func (r *ImportJobRepository) FindByID(ctx context.Context, id string) (*models.ImportJob, error) {
	query := `
		SELECT
			ij.id,
			ij.import_type,
			ij.status,
			ij.source,
			ij.file_name,
			ij.file_hash,
			ij.created_by,
			u.name,
			ij.duplicate_strategy,
			ij.total_rows,
			ij.valid_rows,
			ij.invalid_rows,
			ij.committed_rows,
			ij.skipped_rows,
			ij.failed_rows,
			ij.preview_payload,
			ij.commit_summary,
			ij.committed_at,
			ij.created_at,
			ij.updated_at
		FROM import_jobs ij
		JOIN users u ON u.id = ij.created_by
		WHERE ij.id = $1
		LIMIT 1
	`

	return scanImportJob(r.db.QueryRowContext(ctx, query, id))
}

func (r *ImportJobRepository) FindCommittedByHash(ctx context.Context, importType string, fileHash string) (*models.ImportJob, error) {
	query := `
		SELECT
			ij.id,
			ij.import_type,
			ij.status,
			ij.source,
			ij.file_name,
			ij.file_hash,
			ij.created_by,
			u.name,
			ij.duplicate_strategy,
			ij.total_rows,
			ij.valid_rows,
			ij.invalid_rows,
			ij.committed_rows,
			ij.skipped_rows,
			ij.failed_rows,
			ij.preview_payload,
			ij.commit_summary,
			ij.committed_at,
			ij.created_at,
			ij.updated_at
		FROM import_jobs ij
		JOIN users u ON u.id = ij.created_by
		WHERE ij.import_type = $1
		  AND ij.file_hash = $2
		  AND ij.status = 'committed'
		ORDER BY ij.committed_at DESC NULLS LAST, ij.created_at DESC
		LIMIT 1
	`

	return scanImportJob(r.db.QueryRowContext(ctx, query, importType, fileHash))
}

func (r *ImportJobRepository) MarkCommittedTx(ctx context.Context, tx *sql.Tx, params MarkImportJobCommittedParams) (*models.ImportJob, error) {
	query := `
		UPDATE import_jobs
		SET
			status = 'committed',
			duplicate_strategy = $2,
			committed_rows = $3,
			skipped_rows = $4,
			failed_rows = $5,
			commit_summary = $6,
			committed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			import_type,
			status,
			source,
			file_name,
			file_hash,
			created_by,
			NULL::TEXT AS created_by_name,
			duplicate_strategy,
			total_rows,
			valid_rows,
			invalid_rows,
			committed_rows,
			skipped_rows,
			failed_rows,
			preview_payload,
			commit_summary,
			committed_at,
			created_at,
			updated_at
	`

	return scanImportJob(
		tx.QueryRowContext(
			ctx,
			query,
			params.ID,
			params.DuplicateStrategy,
			params.CommittedRows,
			params.SkippedRows,
			params.FailedRows,
			params.CommitSummary,
		),
	)
}

func (r *ImportJobRepository) MarkFailed(ctx context.Context, params MarkImportJobFailedParams) error {
	query := `
		UPDATE import_jobs
		SET
			status = 'failed',
			commit_summary = $2,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, params.ID, params.CommitSummary)
	return err
}

func scanImportJob(row scanner) (*models.ImportJob, error) {
	job := &models.ImportJob{}
	var createdByName sql.NullString
	var duplicateStrategy sql.NullString
	var previewPayload []byte
	var commitSummary []byte
	var committedAt sql.NullTime

	if err := row.Scan(
		&job.ID,
		&job.ImportType,
		&job.Status,
		&job.Source,
		&job.FileName,
		&job.FileHash,
		&job.CreatedBy,
		&createdByName,
		&duplicateStrategy,
		&job.TotalRows,
		&job.ValidRows,
		&job.InvalidRows,
		&job.CommittedRows,
		&job.SkippedRows,
		&job.FailedRows,
		&previewPayload,
		&commitSummary,
		&committedAt,
		&job.CreatedAt,
		&job.UpdatedAt,
	); err != nil {
		return nil, err
	}

	job.CreatedByName = nullStringPtr(createdByName)
	job.DuplicateStrategy = nullStringPtr(duplicateStrategy)
	job.PreviewPayload = previewPayload
	job.CommitSummary = commitSummary
	job.CommittedAt = nullTimePtr(committedAt)

	return job, nil
}
