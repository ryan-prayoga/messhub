package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrImportFileEmpty                 = errors.New("import file is empty")
	ErrImportMissingHeaders            = errors.New("csv headers do not match the required template")
	ErrImportInvalidJob                = errors.New("import job not found")
	ErrImportJobTypeMismatch           = errors.New("import job type does not match the requested import")
	ErrImportJobAlreadyCommitted       = errors.New("import job has already been committed")
	ErrImportJobNotOwned               = errors.New("import job must be committed by the same admin")
	ErrImportNoValidRows               = errors.New("no valid rows are ready to import")
	ErrImportInvalidDuplicateStrategy  = errors.New("duplicate strategy must be skip or fail")
	ErrImportTemporaryPasswordRequired = errors.New("temporary password must be at least 8 characters")
	ErrImportDuplicateRowsPresent      = errors.New("duplicate rows must be skipped or cleaned before import")
	ErrImportAlreadyCommittedFile      = errors.New("this CSV file has already been imported")
)

const (
	importTypeMembers     = "members"
	importTypeWallet      = "wallet"
	importStatusPreviewed = "previewed"
	importStatusCommitted = "committed"
	importStatusFailed    = "failed"
	importSourceSheet     = "spreadsheet_import"
	importStatusValid     = "valid"
	importStatusInvalid   = "invalid"
	importStatusDuplicate = "duplicate"
	duplicateStrategySkip = "skip"
	duplicateStrategyFail = "fail"
	defaultWalletCategory = "lainnya"
)

type ImportWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MemberImportPreviewSummary struct {
	TotalRows      int `json:"total_rows"`
	ValidRows      int `json:"valid_rows"`
	InvalidRows    int `json:"invalid_rows"`
	DuplicateRows  int `json:"duplicate_rows"`
	ImportableRows int `json:"importable_rows"`
}

type MemberImportPreviewRow struct {
	RowNumber          int      `json:"row_number"`
	Status             string   `json:"status"`
	Name               string   `json:"name"`
	Email              string   `json:"email"`
	Role               string   `json:"role"`
	NormalizedRole     string   `json:"normalized_role"`
	IsActive           string   `json:"is_active"`
	NormalizedIsActive *bool    `json:"normalized_is_active,omitempty"`
	Errors             []string `json:"errors"`
	Warnings           []string `json:"warnings"`
}

type MemberImportPreview struct {
	JobID                     string                     `json:"job_id"`
	FileName                  string                     `json:"file_name"`
	Summary                   MemberImportPreviewSummary `json:"summary"`
	Rows                      []MemberImportPreviewRow   `json:"rows"`
	Warnings                  []ImportWarning            `json:"warnings"`
	CanCommit                 bool                       `json:"can_commit"`
	RequiresTemporaryPassword bool                       `json:"requires_temporary_password"`
}

type WalletImportPreviewSummary struct {
	TotalRows      int   `json:"total_rows"`
	ValidRows      int   `json:"valid_rows"`
	InvalidRows    int   `json:"invalid_rows"`
	ImportableRows int   `json:"importable_rows"`
	TotalIncome    int64 `json:"total_income"`
	TotalExpense   int64 `json:"total_expense"`
}

type WalletImportPreviewRow struct {
	RowNumber                 int      `json:"row_number"`
	Status                    string   `json:"status"`
	TransactionDate           string   `json:"transaction_date"`
	NormalizedTransactionDate *string  `json:"normalized_transaction_date,omitempty"`
	Description               string   `json:"description"`
	Income                    string   `json:"income"`
	Expense                   string   `json:"expense"`
	Type                      string   `json:"type"`
	Amount                    *int64   `json:"amount,omitempty"`
	Category                  string   `json:"category"`
	Proof                     string   `json:"proof"`
	Errors                    []string `json:"errors"`
	Warnings                  []string `json:"warnings"`
}

type WalletImportPreview struct {
	JobID     string                     `json:"job_id"`
	FileName  string                     `json:"file_name"`
	Summary   WalletImportPreviewSummary `json:"summary"`
	Rows      []WalletImportPreviewRow   `json:"rows"`
	Warnings  []ImportWarning            `json:"warnings"`
	CanCommit bool                       `json:"can_commit"`
}

type CommitMemberImportInput struct {
	JobID             string `json:"job_id"`
	DuplicateStrategy string `json:"duplicate_strategy"`
	TemporaryPassword string `json:"temporary_password"`
}

type CommitWalletImportInput struct {
	JobID string `json:"job_id"`
}

type ImportCommitResult struct {
	JobID             string  `json:"job_id"`
	ImportType        string  `json:"import_type"`
	ImportedRows      int     `json:"imported_rows"`
	SkippedRows       int     `json:"skipped_rows"`
	FailedRows        int     `json:"failed_rows"`
	TotalRows         int     `json:"total_rows"`
	DuplicateStrategy *string `json:"duplicate_strategy,omitempty"`
	TotalIncome       int64   `json:"total_income,omitempty"`
	TotalExpense      int64   `json:"total_expense,omitempty"`
}

type memberImportJobPayload struct {
	Summary  MemberImportPreviewSummary `json:"summary"`
	Rows     []MemberImportPreviewRow   `json:"rows"`
	Warnings []ImportWarning            `json:"warnings"`
}

type walletImportJobPayload struct {
	Summary  WalletImportPreviewSummary `json:"summary"`
	Rows     []WalletImportPreviewRow   `json:"rows"`
	Warnings []ImportWarning            `json:"warnings"`
}

type ImportService struct {
	db                  *sql.DB
	userRepository      *repository.UserRepository
	walletRepository    *repository.WalletRepository
	importJobRepository *repository.ImportJobRepository
	auditService        *AuditService
}

func NewImportService(
	db *sql.DB,
	userRepository *repository.UserRepository,
	walletRepository *repository.WalletRepository,
	importJobRepository *repository.ImportJobRepository,
	auditService *AuditService,
) *ImportService {
	return &ImportService{
		db:                  db,
		userRepository:      userRepository,
		walletRepository:    walletRepository,
		importJobRepository: importJobRepository,
		auditService:        auditService,
	}
}

func (s *ImportService) PreviewMembers(ctx context.Context, actorID string, fileName string, content []byte) (*MemberImportPreview, error) {
	headers, rows, fileHash, err := parseCSVFile(content)
	if err != nil {
		return nil, err
	}

	indexes, err := resolveHeaders(headers, []string{"name", "email", "role", "is_active"}, memberHeaderAliases())
	if err != nil {
		return nil, err
	}

	existingUsers, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	existingEmails := make(map[string]struct{}, len(existingUsers))
	for _, user := range existingUsers {
		existingEmails[normalizeEmail(user.Email)] = struct{}{}
	}

	previewRows := make([]MemberImportPreviewRow, 0, len(rows))
	seenEmails := map[string]int{}
	summary := MemberImportPreviewSummary{}

	for rawIndex, rawRow := range rows {
		if isBlankCSVRow(rawRow) {
			continue
		}

		summary.TotalRows++
		rowNumber := rawIndex + 2
		name := csvValue(rawRow, indexes["name"])
		email := normalizeEmail(csvValue(rawRow, indexes["email"]))
		roleValue := csvValue(rawRow, indexes["role"])
		isActiveValue := csvValue(rawRow, indexes["is_active"])
		normalizedRole := normalizeMemberRole(roleValue)
		normalizedIsActive, hasIsActive := normalizeImportBool(isActiveValue)

		row := MemberImportPreviewRow{
			RowNumber:      rowNumber,
			Name:           name,
			Email:          email,
			Role:           roleValue,
			NormalizedRole: normalizedRole,
			IsActive:       isActiveValue,
		}

		if normalizedIsActive != nil {
			value := *normalizedIsActive
			row.NormalizedIsActive = &value
		}

		errorsFound := make([]string, 0)
		warningsFound := make([]string, 0)

		if strings.TrimSpace(name) == "" {
			errorsFound = append(errorsFound, "Nama wajib diisi.")
		}

		if email == "" {
			errorsFound = append(errorsFound, "Email wajib diisi.")
		} else if !isValidImportEmail(email) {
			errorsFound = append(errorsFound, "Format email tidak valid.")
		}

		if normalizedRole == "" {
			errorsFound = append(errorsFound, "Role harus admin, treasurer, atau member.")
		}

		if !hasIsActive {
			errorsFound = append(errorsFound, "Kolom is_active harus berisi true/false, ya/tidak, atau aktif/nonaktif.")
		}

		if len(errorsFound) > 0 {
			row.Status = importStatusInvalid
			row.Errors = errorsFound
			summary.InvalidRows++
			previewRows = append(previewRows, row)
			continue
		}

		if _, exists := existingEmails[email]; exists {
			row.Status = importStatusDuplicate
			warningsFound = append(warningsFound, "Email ini sudah terdaftar di MessHub.")
			summary.DuplicateRows++
		} else if firstRow, exists := seenEmails[email]; exists {
			row.Status = importStatusDuplicate
			warningsFound = append(warningsFound, fmt.Sprintf("Email yang sama sudah muncul di baris %d.", firstRow))
			summary.DuplicateRows++
		} else {
			row.Status = importStatusValid
			seenEmails[email] = row.RowNumber
			summary.ValidRows++
			summary.ImportableRows++
		}

		row.Warnings = warningsFound
		previewRows = append(previewRows, row)
	}

	if summary.TotalRows == 0 {
		return nil, ErrImportFileEmpty
	}

	warnings := make([]ImportWarning, 0)
	payload := memberImportJobPayload{
		Summary:  summary,
		Rows:     previewRows,
		Warnings: warnings,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	job, err := s.importJobRepository.Create(ctx, repository.CreateImportJobParams{
		ImportType:     importTypeMembers,
		Status:         importStatusPreviewed,
		Source:         importSourceSheet,
		FileName:       sanitizeImportFileName(fileName),
		FileHash:       fileHash,
		CreatedBy:      actorID,
		TotalRows:      summary.TotalRows,
		ValidRows:      summary.ValidRows,
		InvalidRows:    summary.InvalidRows,
		PreviewPayload: payloadBytes,
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditService.Log(ctx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "member_import_preview",
		EntityType: "import_job",
		EntityID:   stringPtr(job.ID),
		NewValue: map[string]any{
			"file_name": job.FileName,
			"summary":   summary,
		},
	}); err != nil {
		return nil, err
	}

	return &MemberImportPreview{
		JobID:                     job.ID,
		FileName:                  job.FileName,
		Summary:                   summary,
		Rows:                      previewRows,
		Warnings:                  warnings,
		CanCommit:                 summary.ImportableRows > 0,
		RequiresTemporaryPassword: true,
	}, nil
}

func (s *ImportService) CommitMembers(ctx context.Context, actorID string, input CommitMemberImportInput) (*ImportCommitResult, error) {
	job, payload, err := s.loadMemberImportJob(ctx, actorID, input.JobID)
	if err != nil {
		return nil, err
	}

	duplicateStrategy := normalizeDuplicateStrategy(input.DuplicateStrategy)
	if duplicateStrategy == "" {
		duplicateStrategy = duplicateStrategySkip
	}
	if duplicateStrategy != duplicateStrategySkip && duplicateStrategy != duplicateStrategyFail {
		return nil, ErrImportInvalidDuplicateStrategy
	}

	if len(strings.TrimSpace(input.TemporaryPassword)) < 8 {
		return nil, ErrImportTemporaryPasswordRequired
	}

	if payload.Summary.ImportableRows == 0 {
		return nil, ErrImportNoValidRows
	}

	if duplicateStrategy == duplicateStrategyFail && payload.Summary.DuplicateRows > 0 {
		return nil, ErrImportDuplicateRowsPresent
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(input.TemporaryPassword)), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	importedRows := 0
	skippedRows := 0
	failedRows := payload.Summary.InvalidRows

	for _, row := range payload.Rows {
		switch row.Status {
		case importStatusInvalid:
			continue
		case importStatusDuplicate:
			skippedRows++
			continue
		case importStatusValid:
			isActive := true
			if row.NormalizedIsActive != nil {
				isActive = *row.NormalizedIsActive
			}

			username, err := s.userRepository.FindAvailableUsernameTx(
				ctx,
				tx,
				strings.TrimSpace(row.Name),
				normalizeEmail(row.Email),
			)
			if err != nil {
				tx.Rollback()
				s.markImportJobFailed(ctx, job.ID, err)
				return nil, err
			}

			_, err = s.userRepository.CreateTx(ctx, tx, repository.CreateUserParams{
				Name:         strings.TrimSpace(row.Name),
				Email:        normalizeEmail(row.Email),
				Username:     username,
				Phone:        nil,
				PasswordHash: string(passwordHash),
				Role:         row.NormalizedRole,
				IsActive:     isActive,
				JoinedAt:     time.Now().UTC(),
			})
			if err != nil {
				if isUniqueViolation(err) {
					if duplicateStrategy == duplicateStrategyFail {
						tx.Rollback()
						return nil, ErrImportDuplicateRowsPresent
					}

					skippedRows++
					continue
				}

				tx.Rollback()
				s.markImportJobFailed(ctx, job.ID, err)
				return nil, err
			}

			importedRows++
		}
	}

	commitSummary := ImportCommitResult{
		JobID:        job.ID,
		ImportType:   importTypeMembers,
		ImportedRows: importedRows,
		SkippedRows:  skippedRows,
		FailedRows:   failedRows,
		TotalRows:    payload.Summary.TotalRows,
	}
	commitSummary.DuplicateStrategy = stringPtr(duplicateStrategy)

	commitSummaryBytes, err := json.Marshal(commitSummary)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := s.importJobRepository.MarkCommittedTx(ctx, tx, repository.MarkImportJobCommittedParams{
		ID:                job.ID,
		DuplicateStrategy: stringPtr(duplicateStrategy),
		CommittedRows:     importedRows,
		SkippedRows:       skippedRows,
		FailedRows:        failedRows,
		CommitSummary:     commitSummaryBytes,
	}); err != nil {
		tx.Rollback()
		s.markImportJobFailed(ctx, job.ID, err)
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "member_import_commit",
		EntityType: "import_job",
		EntityID:   stringPtr(job.ID),
		NewValue:   commitSummary,
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.markImportJobFailed(ctx, job.ID, err)
		return nil, err
	}

	return &commitSummary, nil
}

func (s *ImportService) PreviewWallet(ctx context.Context, actorID string, fileName string, content []byte) (*WalletImportPreview, error) {
	headers, rows, fileHash, err := parseCSVFile(content)
	if err != nil {
		return nil, err
	}

	indexes, err := resolveHeaders(headers, []string{"transaction_date", "description", "income", "expense"}, walletHeaderAliases())
	if err != nil {
		return nil, err
	}

	previewRows := make([]WalletImportPreviewRow, 0, len(rows))
	summary := WalletImportPreviewSummary{}
	warnings := make([]ImportWarning, 0)

	if committedJob, err := s.importJobRepository.FindCommittedByHash(ctx, importTypeWallet, fileHash); err == nil {
		warnings = append(warnings, ImportWarning{
			Code:    "already_imported",
			Message: fmt.Sprintf("File yang sama sudah pernah diimpor pada %s.", formatImportJobTime(committedJob.CommittedAt, committedJob.CreatedAt)),
		})
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	for rawIndex, rawRow := range rows {
		if isBlankCSVRow(rawRow) {
			continue
		}

		summary.TotalRows++
		rowNumber := rawIndex + 2
		dateValue := csvValue(rawRow, indexes["transaction_date"])
		description := csvValue(rawRow, indexes["description"])
		incomeValue := csvValue(rawRow, indexes["income"])
		expenseValue := csvValue(rawRow, indexes["expense"])
		proofValue := csvValue(rawRow, indexes["proof"])

		row := WalletImportPreviewRow{
			RowNumber:       rowNumber,
			TransactionDate: dateValue,
			Description:     description,
			Income:          incomeValue,
			Expense:         expenseValue,
			Proof:           proofValue,
			Category:        defaultWalletCategory,
		}

		errorsFound := make([]string, 0)
		warningsFound := make([]string, 0)

		parsedDate, hasDate := parseImportDate(dateValue)
		if !hasDate {
			errorsFound = append(errorsFound, "Tanggal transaksi wajib diisi dengan format yang valid.")
		} else {
			normalizedDate := parsedDate.Format("2006-01-02")
			row.NormalizedTransactionDate = &normalizedDate
		}

		if strings.TrimSpace(description) == "" {
			errorsFound = append(errorsFound, "Deskripsi wajib diisi.")
		}

		incomeAmount, incomeOk := parseCurrencyAmount(incomeValue)
		if !incomeOk {
			errorsFound = append(errorsFound, "Nilai pemasukan tidak valid.")
		}

		expenseAmount, expenseOk := parseCurrencyAmount(expenseValue)
		if !expenseOk {
			errorsFound = append(errorsFound, "Nilai pengeluaran tidak valid.")
		}

		if incomeOk && expenseOk {
			switch {
			case incomeAmount > 0 && expenseAmount > 0:
				errorsFound = append(errorsFound, "Isi salah satu kolom pemasukan atau pengeluaran, bukan keduanya.")
			case incomeAmount <= 0 && expenseAmount <= 0:
				errorsFound = append(errorsFound, "Nominal transaksi harus diisi di salah satu kolom pemasukan atau pengeluaran.")
			case incomeAmount > 0:
				row.Type = "income"
				row.Amount = &incomeAmount
				row.Category = inferWalletImportCategory(description)
			case expenseAmount > 0:
				row.Type = "expense"
				row.Amount = &expenseAmount
				row.Category = inferWalletImportCategory(description)
			}
		}

		if len(errorsFound) > 0 {
			row.Status = importStatusInvalid
			row.Errors = errorsFound
			row.Warnings = warningsFound
			summary.InvalidRows++
			previewRows = append(previewRows, row)
			continue
		}

		row.Status = importStatusValid
		row.Warnings = warningsFound
		summary.ValidRows++
		summary.ImportableRows++
		if row.Type == "income" && row.Amount != nil {
			summary.TotalIncome += *row.Amount
		}
		if row.Type == "expense" && row.Amount != nil {
			summary.TotalExpense += *row.Amount
		}
		previewRows = append(previewRows, row)
	}

	if summary.TotalRows == 0 {
		return nil, ErrImportFileEmpty
	}

	payload := walletImportJobPayload{
		Summary:  summary,
		Rows:     previewRows,
		Warnings: warnings,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	job, err := s.importJobRepository.Create(ctx, repository.CreateImportJobParams{
		ImportType:     importTypeWallet,
		Status:         importStatusPreviewed,
		Source:         importSourceSheet,
		FileName:       sanitizeImportFileName(fileName),
		FileHash:       fileHash,
		CreatedBy:      actorID,
		TotalRows:      summary.TotalRows,
		ValidRows:      summary.ValidRows,
		InvalidRows:    summary.InvalidRows,
		PreviewPayload: payloadBytes,
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditService.Log(ctx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "wallet_import_preview",
		EntityType: "import_job",
		EntityID:   stringPtr(job.ID),
		NewValue: map[string]any{
			"file_name": job.FileName,
			"summary":   summary,
			"warnings":  warnings,
		},
	}); err != nil {
		return nil, err
	}

	return &WalletImportPreview{
		JobID:     job.ID,
		FileName:  job.FileName,
		Summary:   summary,
		Rows:      previewRows,
		Warnings:  warnings,
		CanCommit: summary.ImportableRows > 0 && !hasImportWarning(warnings, "already_imported"),
	}, nil
}

func (s *ImportService) CommitWallet(ctx context.Context, actorID string, input CommitWalletImportInput) (*ImportCommitResult, error) {
	job, payload, err := s.loadWalletImportJob(ctx, actorID, input.JobID)
	if err != nil {
		return nil, err
	}

	if payload.Summary.ImportableRows == 0 {
		return nil, ErrImportNoValidRows
	}

	if committedJob, err := s.importJobRepository.FindCommittedByHash(ctx, importTypeWallet, job.FileHash); err == nil && committedJob.ID != job.ID {
		return nil, ErrImportAlreadyCommittedFile
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	importedRows := 0
	failedRows := payload.Summary.InvalidRows
	var totalIncome int64
	var totalExpense int64

	for _, row := range payload.Rows {
		if row.Status != importStatusValid || row.Amount == nil || row.NormalizedTransactionDate == nil {
			continue
		}

		transactionDate, err := time.Parse("2006-01-02", *row.NormalizedTransactionDate)
		if err != nil {
			tx.Rollback()
			s.markImportJobFailed(ctx, job.ID, err)
			return nil, err
		}

		_, err = s.walletRepository.CreateTx(ctx, tx, repository.CreateWalletTransactionParams{
			TransactionDate: transactionDate,
			Type:            row.Type,
			Category:        row.Category,
			Amount:          *row.Amount,
			Description:     strings.TrimSpace(row.Description),
			ProofURL:        stringPtr(row.Proof),
			Source:          importSourceSheet,
			ImportJobID:     stringPtr(job.ID),
			CreatedBy:       actorID,
		})
		if err != nil {
			tx.Rollback()
			s.markImportJobFailed(ctx, job.ID, err)
			return nil, err
		}

		importedRows++
		if row.Type == "income" {
			totalIncome += *row.Amount
		} else if row.Type == "expense" {
			totalExpense += *row.Amount
		}
	}

	commitSummary := ImportCommitResult{
		JobID:        job.ID,
		ImportType:   importTypeWallet,
		ImportedRows: importedRows,
		SkippedRows:  0,
		FailedRows:   failedRows,
		TotalRows:    payload.Summary.TotalRows,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
	}

	commitSummaryBytes, err := json.Marshal(commitSummary)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := s.importJobRepository.MarkCommittedTx(ctx, tx, repository.MarkImportJobCommittedParams{
		ID:            job.ID,
		CommittedRows: importedRows,
		SkippedRows:   0,
		FailedRows:    failedRows,
		CommitSummary: commitSummaryBytes,
	}); err != nil {
		tx.Rollback()
		if isUniqueViolation(err) {
			return nil, ErrImportAlreadyCommittedFile
		}
		s.markImportJobFailed(ctx, job.ID, err)
		return nil, err
	}

	if err := s.auditService.LogTx(ctx, tx, AuditLogInput{
		UserID:     stringPtr(actorID),
		Action:     "wallet_import_commit",
		EntityType: "import_job",
		EntityID:   stringPtr(job.ID),
		NewValue:   commitSummary,
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.markImportJobFailed(ctx, job.ID, err)
		return nil, err
	}

	return &commitSummary, nil
}

func (s *ImportService) loadMemberImportJob(ctx context.Context, actorID string, jobID string) (*models.ImportJob, *memberImportJobPayload, error) {
	job, err := s.importJobRepository.FindByID(ctx, strings.TrimSpace(jobID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrImportInvalidJob
		}
		return nil, nil, err
	}

	if job.ImportType != importTypeMembers {
		return nil, nil, ErrImportJobTypeMismatch
	}

	if job.CreatedBy != actorID {
		return nil, nil, ErrImportJobNotOwned
	}

	if job.Status != importStatusPreviewed {
		return nil, nil, ErrImportJobAlreadyCommitted
	}

	payload := new(memberImportJobPayload)
	if err := json.Unmarshal(job.PreviewPayload, payload); err != nil {
		return nil, nil, err
	}

	return job, payload, nil
}

func (s *ImportService) loadWalletImportJob(ctx context.Context, actorID string, jobID string) (*models.ImportJob, *walletImportJobPayload, error) {
	job, err := s.importJobRepository.FindByID(ctx, strings.TrimSpace(jobID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrImportInvalidJob
		}
		return nil, nil, err
	}

	if job.ImportType != importTypeWallet {
		return nil, nil, ErrImportJobTypeMismatch
	}

	if job.CreatedBy != actorID {
		return nil, nil, ErrImportJobNotOwned
	}

	if job.Status != importStatusPreviewed {
		return nil, nil, ErrImportJobAlreadyCommitted
	}

	payload := new(walletImportJobPayload)
	if err := json.Unmarshal(job.PreviewPayload, payload); err != nil {
		return nil, nil, err
	}

	return job, payload, nil
}

func (s *ImportService) markImportJobFailed(ctx context.Context, jobID string, err error) {
	_ = s.importJobRepository.MarkFailed(ctx, repository.MarkImportJobFailedParams{
		ID:            jobID,
		CommitSummary: []byte(fmt.Sprintf(`{"error":%q}`, err.Error())),
	})
}

func parseCSVFile(content []byte) ([]string, [][]string, string, error) {
	if len(content) == 0 {
		return nil, nil, "", ErrImportFileEmpty
	}

	hash := sha256.Sum256(content)
	reader := csv.NewReader(strings.NewReader(trimUTF8BOM(string(content))))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil, "", ErrImportFileEmpty
		}
		return nil, nil, "", err
	}

	if len(records) == 0 {
		return nil, nil, "", ErrImportFileEmpty
	}

	headers := make([]string, len(records[0]))
	copy(headers, records[0])

	return headers, records[1:], hex.EncodeToString(hash[:]), nil
}

func resolveHeaders(headers []string, required []string, aliases map[string]string) (map[string]int, error) {
	indexes := make(map[string]int, len(required))
	for index, header := range headers {
		canonical, ok := aliases[normalizeImportHeader(header)]
		if !ok {
			continue
		}
		if _, exists := indexes[canonical]; !exists {
			indexes[canonical] = index
		}
	}

	missing := make([]string, 0)
	for _, requiredField := range required {
		if _, exists := indexes[requiredField]; !exists {
			missing = append(missing, requiredField)
		}
	}

	if len(missing) > 0 {
		return nil, ErrImportMissingHeaders
	}

	if _, exists := indexes["proof"]; !exists {
		indexes["proof"] = -1
	}

	return indexes, nil
}

func memberHeaderAliases() map[string]string {
	return map[string]string{
		"name":     "name",
		"nama":     "name",
		"email":    "email",
		"role":     "role",
		"peran":    "role",
		"isactive": "is_active",
		"aktif":    "is_active",
		"active":   "is_active",
	}
}

func walletHeaderAliases() map[string]string {
	return map[string]string{
		"transactiondate": "transaction_date",
		"tanggal":         "transaction_date",
		"description":     "description",
		"deskripsi":       "description",
		"income":          "income",
		"pemasukan":       "income",
		"pemasukanrp":     "income",
		"expense":         "expense",
		"pengeluaran":     "expense",
		"pengeluaranrp":   "expense",
		"proof":           "proof",
		"bukti":           "proof",
	}
}

func normalizeImportHeader(value string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			return unicode.ToLower(r)
		default:
			return -1
		}
	}, strings.TrimSpace(value))
}

func csvValue(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[index])
}

func isBlankCSVRow(row []string) bool {
	for _, value := range row {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}

	return true
}

func normalizeMemberRole(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "admin":
		return models.RoleAdmin
	case "treasurer", "bendahara":
		return models.RoleTreasurer
	case "member", "anggota":
		return models.RoleMember
	default:
		return ""
	}
}

func normalizeImportBool(value string) (*bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "ya", "aktif", "active":
		return boolPtr(true), true
	case "false", "0", "no", "tidak", "nonaktif", "inactive":
		return boolPtr(false), true
	default:
		return nil, false
	}
}

func boolPtr(value bool) *bool {
	return &value
}

func isValidImportEmail(value string) bool {
	parsed, err := mail.ParseAddress(value)
	if err != nil {
		return false
	}

	return strings.EqualFold(parsed.Address, value)
}

func parseCurrencyAmount(value string) (int64, bool) {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return 0, true
	}

	negative := strings.Contains(trimmed, "-")
	digits := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, trimmed)

	if digits == "" {
		return 0, false
	}

	amount, err := strconv.ParseInt(digits, 10, 64)
	if err != nil {
		return 0, false
	}

	if negative {
		amount *= -1
	}

	return amount, true
}

func parseImportDate(value string) (time.Time, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, false
	}

	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
		"2/1/2006",
		"02-01-2006",
		"2-1-2006",
		"2 Jan 2006",
		"02 Jan 2006",
		"2 January 2006",
		"02 January 2006",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, trimmed, time.UTC)
		if err == nil {
			return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), true
		}
	}

	return time.Time{}, false
}

func inferWalletImportCategory(description string) string {
	normalized := strings.ToLower(strings.TrimSpace(description))
	switch {
	case strings.Contains(normalized, "wifi"):
		return "wifi"
	case strings.Contains(normalized, "hibah"), strings.Contains(normalized, "donasi"), strings.Contains(normalized, "sumbangan"):
		return "hibah"
	case strings.Contains(normalized, "galon"):
		return "galon"
	case strings.Contains(normalized, "plastik"), strings.Contains(normalized, "sabun"), strings.Contains(normalized, "kebersihan"), strings.Contains(normalized, "pel"), strings.Contains(normalized, "sapu"):
		return "kebersihan"
	default:
		return defaultWalletCategory
	}
}

func normalizeDuplicateStrategy(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case duplicateStrategySkip:
		return duplicateStrategySkip
	case duplicateStrategyFail:
		return duplicateStrategyFail
	default:
		return ""
	}
}

func trimUTF8BOM(value string) string {
	return strings.TrimPrefix(value, "\uFEFF")
}

func sanitizeImportFileName(fileName string) string {
	trimmed := strings.TrimSpace(fileName)
	if trimmed == "" {
		return "import.csv"
	}

	return trimmed
}

func formatImportJobTime(committedAt *time.Time, fallback time.Time) string {
	if committedAt != nil {
		return committedAt.In(time.Local).Format("02 Jan 2006 15:04")
	}

	return fallback.In(time.Local).Format("02 Jan 2006 15:04")
}

func hasImportWarning(warnings []ImportWarning, code string) bool {
	for _, warning := range warnings {
		if warning.Code == code {
			return true
		}
	}

	return false
}
