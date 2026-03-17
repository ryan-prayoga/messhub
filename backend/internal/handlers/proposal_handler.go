package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

type ProposalHandler struct {
	service *services.ProposalService
}

func NewProposalHandler(service *services.ProposalService) *ProposalHandler {
	return &ProposalHandler{service: service}
}

func (h *ProposalHandler) List(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	items, err := h.service.List(c.UserContext(), user.ID)
	if err != nil {
		return response.InternalServerError(c, "failed to load proposals")
	}

	return response.Success(c, fiber.StatusOK, "proposals loaded", items)
}

func (h *ProposalHandler) Get(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	item, err := h.service.GetDetail(c.UserContext(), c.Params("id"), user.ID)
	if err != nil {
		if errors.Is(err, services.ErrProposalNotFound) {
			return response.NotFound(c, "Usulan tidak ditemukan.")
		}

		return response.InternalServerError(c, "failed to load proposal")
	}

	return response.Success(c, fiber.StatusOK, "proposal loaded", item)
}

func (h *ProposalHandler) Create(c *fiber.Ctx) error {
	request := new(services.CreateProposalInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "proposal")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredMaxLength("title", request.Title, maxActivityTitleLength, "title is required", "title is too long")
	details.RequiredMaxLength("description", request.Description, maxActivityBodyLength, "description is required", "description is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.service.Create(c.UserContext(), user.ID, *request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidProposalInput) {
			return response.InvalidRequest(c, "Data usulan belum lengkap atau format tanggal voting belum valid.")
		}

		return response.InternalServerError(c, "failed to create proposal")
	}

	return response.Success(c, fiber.StatusCreated, "proposal created", item)
}

func (h *ProposalHandler) Vote(c *fiber.Ctx) error {
	request := new(services.VoteProposalInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "proposal vote")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.Enum("vote_type", request.VoteType, []string{
		models.ProposalVoteAgree,
		models.ProposalVoteDisagree,
	}, "vote_type must be agree or disagree")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.service.Vote(c.UserContext(), c.Params("id"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidProposalVote):
			return response.InvalidRequest(c, "Jenis vote harus agree atau disagree.")
		case errors.Is(err, services.ErrProposalNotFound):
			return response.NotFound(c, "Usulan tidak ditemukan.")
		case errors.Is(err, services.ErrProposalVoteExists):
			return response.Conflict(c, "Kamu sudah memberikan vote untuk usulan ini.")
		case errors.Is(err, services.ErrProposalNotActive), errors.Is(err, services.ErrProposalVotingUnavailable):
			return response.InvalidRequest(c, "Voting untuk usulan ini sudah tidak tersedia.")
		default:
			return response.InternalServerError(c, "failed to vote proposal")
		}
	}

	return response.Success(c, fiber.StatusOK, "proposal voted", item)
}

func (h *ProposalHandler) Close(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	item, err := h.service.Close(c.UserContext(), c.Params("id"), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProposalNotFound):
			return response.NotFound(c, "Usulan tidak ditemukan.")
		case errors.Is(err, services.ErrProposalNotActive):
			return response.InvalidRequest(c, "Usulan ini sudah tidak aktif.")
		default:
			return response.InternalServerError(c, "failed to close proposal")
		}
	}

	return response.Success(c, fiber.StatusOK, "proposal closed", item)
}

func (h *ProposalHandler) Finalize(c *fiber.Ctx) error {
	request := new(services.FinalizeProposalInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "proposal finalization")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.Enum("status", request.Status, []string{
		models.ProposalStatusApproved,
		models.ProposalStatusRejected,
	}, "status must be approved or rejected")
	details.OptionalMaxLength("final_decision_note", request.FinalDecisionNote, maxActivityBodyLength, "final_decision_note is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.service.Finalize(c.UserContext(), c.Params("id"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidProposalFinalStatus):
			return response.InvalidRequest(c, "Status final usulan harus approved atau rejected.")
		case errors.Is(err, services.ErrProposalNotFound):
			return response.NotFound(c, "Usulan tidak ditemukan.")
		case errors.Is(err, services.ErrProposalNotActive):
			return response.InvalidRequest(c, "Usulan ini tidak bisa difinalisasi lagi.")
		default:
			return response.InternalServerError(c, "failed to finalize proposal")
		}
	}

	return response.Success(c, fiber.StatusOK, "proposal finalized", item)
}
