package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
)

type ActivityHandler struct {
	activityService *services.ActivityService
}

func NewActivityHandler(activityService *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

func (h *ActivityHandler) ListActivities(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	items, err := h.activityService.ListActivities(c.UserContext(), user.ID, services.ListActivitiesInput{
		Limit: limit,
	})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to load activities", "activities_failed")
	}

	return response.Success(c, fiber.StatusOK, "activities loaded", items)
}

func (h *ActivityHandler) CreateActivity(c *fiber.Ctx) error {
	request := new(services.CreateActivityInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid activity payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	item, err := h.activityService.CreateActivity(c.UserContext(), user.ID, user.Name, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidActivityInput), errors.Is(err, services.ErrInvalidActivityType):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_activity_input")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to create activity", "activity_create_failed")
		}
	}

	return response.Success(c, fiber.StatusCreated, "activity created", item)
}

func (h *ActivityHandler) GetContributionLeaderboard(c *fiber.Ctx) error {
	items, err := h.activityService.GetContributionLeaderboard(c.UserContext(), c.Query("period", "month"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidLeaderboardPeriod):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_leaderboard_period")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to load leaderboard", "leaderboard_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "contribution leaderboard loaded", items)
}

func (h *ActivityHandler) ListComments(c *fiber.Ctx) error {
	items, err := h.activityService.ListComments(c.UserContext(), c.Params("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "activity_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to load comments", "activity_comments_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "activity comments loaded", items)
}

func (h *ActivityHandler) AddComment(c *fiber.Ctx) error {
	request := new(services.CreateActivityCommentInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid comment payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	item, err := h.activityService.AddComment(c.UserContext(), c.Params("id"), user.ID, user.Name, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCommentInput):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_comment_input")
		case errors.Is(err, services.ErrActivityNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "activity_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to add comment", "activity_comment_failed")
		}
	}

	return response.Success(c, fiber.StatusCreated, "activity comment created", item)
}

func (h *ActivityHandler) ToggleReaction(c *fiber.Ctx) error {
	request := new(services.ToggleActivityReactionInput)
	if err := c.BodyParser(request); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid reaction payload", "invalid_payload")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	item, err := h.activityService.ToggleReaction(c.UserContext(), c.Params("id"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidReactionInput):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "invalid_reaction_input")
		case errors.Is(err, services.ErrActivityNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "activity_not_found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to update reaction", "activity_reaction_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "activity reaction updated", item)
}

func (h *ActivityHandler) ClaimFood(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	item, err := h.activityService.ClaimFood(c.UserContext(), c.Params("id"), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "activity_not_found")
		case errors.Is(err, services.ErrFoodClaimNotAllowed):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "food_claim_not_allowed")
		case errors.Is(err, services.ErrFoodClaimAlreadyExists):
			return response.Error(c, fiber.StatusConflict, err.Error(), "food_claim_exists")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to claim food", "food_claim_failed")
		}
	}

	return response.Success(c, fiber.StatusCreated, "food claimed", item)
}

func (h *ActivityHandler) ListFoodClaims(c *fiber.Ctx) error {
	items, err := h.activityService.ListFoodClaims(c.UserContext(), c.Params("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "activity_not_found")
		case errors.Is(err, services.ErrFoodClaimNotAllowed):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "food_claim_not_allowed")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to load food claims", "food_claims_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "food claims loaded", items)
}

func (h *ActivityHandler) RespondRice(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "missing authenticated user", "missing_authenticated_user")
	}

	item, err := h.activityService.RespondRice(c.UserContext(), c.Params("id"), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "activity_not_found")
		case errors.Is(err, services.ErrRiceResponseNotAllowed):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "rice_response_not_allowed")
		case errors.Is(err, services.ErrRiceResponseAlreadyExists):
			return response.Error(c, fiber.StatusConflict, err.Error(), "rice_response_exists")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to save rice response", "rice_response_failed")
		}
	}

	return response.Success(c, fiber.StatusCreated, "rice response saved", item)
}

func (h *ActivityHandler) ListRiceResponses(c *fiber.Ctx) error {
	items, err := h.activityService.ListRiceResponses(c.UserContext(), c.Params("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), "activity_not_found")
		case errors.Is(err, services.ErrRiceResponseNotAllowed):
			return response.Error(c, fiber.StatusBadRequest, err.Error(), "rice_response_not_allowed")
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to load rice responses", "rice_responses_failed")
		}
	}

	return response.Success(c, fiber.StatusOK, "rice responses loaded", items)
}
