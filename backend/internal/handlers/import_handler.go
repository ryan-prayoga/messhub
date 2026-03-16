package handlers

import (
	"errors"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

type ImportHandler struct {
	importService *services.ImportService
}

func NewImportHandler(importService *services.ImportService) *ImportHandler {
	return &ImportHandler{importService: importService}
}

func (h *ImportHandler) PreviewMembers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.InvalidRequest(c, "File CSV wajib dipilih terlebih dahulu.")
	}

	content, err := readImportFile(fileHeader)
	if err != nil {
		return response.InvalidRequest(c, "File CSV tidak dapat dibaca.")
	}

	preview, err := h.importService.PreviewMembers(c.UserContext(), user.ID, fileHeader.Filename, content)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrImportFileEmpty):
			return response.InvalidRequest(c, "File CSV kosong atau belum berisi data.")
		case errors.Is(err, services.ErrImportMissingHeaders):
			return response.InvalidRequest(c, "Header CSV anggota tidak sesuai template. Gunakan kolom name,email,role,is_active.")
		default:
			return response.InternalServerError(c, "failed to preview member import")
		}
	}

	return response.Success(c, fiber.StatusOK, "member import preview ready", preview)
}

func (h *ImportHandler) CommitMembers(c *fiber.Ctx) error {
	request := new(services.CommitMemberImportInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "member import")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredString("job_id", request.JobID, "job_id is required")
	details.RequiredString("duplicate_strategy", request.DuplicateStrategy, "duplicate_strategy is required")
	details.RequiredString("temporary_password", request.TemporaryPassword, "temporary_password is required")
	details.MinLength("temporary_password", request.TemporaryPassword, 8, "temporary_password must be at least 8 characters")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	result, err := h.importService.CommitMembers(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrImportInvalidJob):
			return response.NotFound(c, "Preview impor anggota tidak ditemukan.")
		case errors.Is(err, services.ErrImportJobTypeMismatch), errors.Is(err, services.ErrImportJobAlreadyCommitted):
			return response.Conflict(c, "Preview impor anggota ini sudah tidak bisa dipakai.")
		case errors.Is(err, services.ErrImportJobNotOwned):
			return response.Forbidden(c, "Preview impor hanya bisa dikonfirmasi oleh admin yang membuat preview.")
		case errors.Is(err, services.ErrImportTemporaryPasswordRequired):
			return response.InvalidRequest(c, "Password sementara minimal 8 karakter.")
		case errors.Is(err, services.ErrImportInvalidDuplicateStrategy):
			return response.InvalidRequest(c, "Pilih strategi duplikasi email yang valid.")
		case errors.Is(err, services.ErrImportNoValidRows):
			return response.Conflict(c, "Tidak ada baris anggota valid yang siap diimpor.")
		case errors.Is(err, services.ErrImportDuplicateRowsPresent):
			return response.Conflict(c, "Masih ada email duplikat. Gunakan opsi skip duplicate atau rapikan CSV terlebih dahulu.")
		default:
			return response.InternalServerError(c, "failed to commit member import")
		}
	}

	return response.Success(c, fiber.StatusOK, "member import committed", result)
}

func (h *ImportHandler) PreviewWallet(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.InvalidRequest(c, "File CSV wajib dipilih terlebih dahulu.")
	}

	content, err := readImportFile(fileHeader)
	if err != nil {
		return response.InvalidRequest(c, "File CSV tidak dapat dibaca.")
	}

	preview, err := h.importService.PreviewWallet(c.UserContext(), user.ID, fileHeader.Filename, content)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrImportFileEmpty):
			return response.InvalidRequest(c, "File CSV kosong atau belum berisi data.")
		case errors.Is(err, services.ErrImportMissingHeaders):
			return response.InvalidRequest(c, "Header CSV transaksi tidak sesuai template. Gunakan kolom transaction_date,description,income,expense,proof.")
		default:
			return response.InternalServerError(c, "failed to preview wallet import")
		}
	}

	return response.Success(c, fiber.StatusOK, "wallet import preview ready", preview)
}

func (h *ImportHandler) CommitWallet(c *fiber.Ctx) error {
	request := new(services.CommitWalletImportInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "wallet import")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredString("job_id", request.JobID, "job_id is required")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	result, err := h.importService.CommitWallet(c.UserContext(), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrImportInvalidJob):
			return response.NotFound(c, "Preview impor transaksi tidak ditemukan.")
		case errors.Is(err, services.ErrImportJobTypeMismatch), errors.Is(err, services.ErrImportJobAlreadyCommitted):
			return response.Conflict(c, "Preview impor transaksi ini sudah tidak bisa dipakai.")
		case errors.Is(err, services.ErrImportJobNotOwned):
			return response.Forbidden(c, "Preview impor hanya bisa dikonfirmasi oleh admin yang membuat preview.")
		case errors.Is(err, services.ErrImportNoValidRows):
			return response.Conflict(c, "Tidak ada transaksi valid yang siap diimpor.")
		case errors.Is(err, services.ErrImportAlreadyCommittedFile):
			return response.Conflict(c, "File CSV ini sudah pernah diimpor sebelumnya.")
		default:
			return response.InternalServerError(c, "failed to commit wallet import")
		}
	}

	return response.Success(c, fiber.StatusOK, "wallet import committed", result)
}

func readImportFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
