package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/models"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/services"
	"github.com/ryanprayoga/messhub/backend/internal/types"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
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
		return response.Unauthorized(c, "authentication required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	items, err := h.activityService.ListActivities(c.UserContext(), user.ID, services.ListActivitiesInput{
		Limit: limit,
	})
	if err != nil {
		return response.InternalServerError(c, "failed to load activities")
	}

	return response.Success(c, fiber.StatusOK, "activities loaded", items)
}

func (h *ActivityHandler) CreateActivity(c *fiber.Ctx) error {
	request := new(services.CreateActivityInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "activity")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.Enum("type", request.Type, []string{
		models.ActivityTypeContribution,
		models.ActivityTypeFood,
		models.ActivityTypeRice,
		models.ActivityTypeAnnouncement,
		models.ActivityTypeOther,
	}, "type must be contribution, food, rice, announcement, or other")
	details.RequiredMaxLength("title", request.Title, maxActivityTitleLength, "title is required", "title is too long")
	details.RequiredMaxLength("content", request.Content, maxActivityBodyLength, "content is required", "content is too long")
	if request.Points != nil && *request.Points <= 0 {
		details.Add("points", "points must be positive")
	}
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.activityService.CreateActivity(c.UserContext(), user.ID, user.Name, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidActivityInput), errors.Is(err, services.ErrInvalidActivityType):
			return response.InvalidRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to create activity")
		}
	}

	return response.Success(c, fiber.StatusCreated, "activity created", item)
}

func (h *ActivityHandler) GetContributionLeaderboard(c *fiber.Ctx) error {
	items, err := h.activityService.GetContributionLeaderboard(c.UserContext(), c.Query("period", "month"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidLeaderboardPeriod):
			return response.InvalidRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to load leaderboard")
		}
	}

	return response.Success(c, fiber.StatusOK, "contribution leaderboard loaded", items)
}

func (h *ActivityHandler) ListComments(c *fiber.Ctx) error {
	items, err := h.activityService.ListComments(c.UserContext(), c.Params("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to load comments")
		}
	}

	return response.Success(c, fiber.StatusOK, "activity comments loaded", items)
}

func (h *ActivityHandler) AddComment(c *fiber.Ctx) error {
	request := new(services.CreateActivityCommentInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "comment")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.RequiredMaxLength("comment", request.Comment, maxCommentLength, "comment is required", "comment is too long")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.activityService.AddComment(c.UserContext(), c.Params("id"), user.ID, user.Name, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCommentInput):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrActivityNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to add comment")
		}
	}

	return response.Success(c, fiber.StatusCreated, "activity comment created", item)
}

func (h *ActivityHandler) ToggleReaction(c *fiber.Ctx) error {
	request := new(services.ToggleActivityReactionInput)
	if err := c.BodyParser(request); err != nil {
		return invalidPayload(c, "reaction")
	}

	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	details := validation.NewErrors()
	details.Enum("reaction_type", request.ReactionType, []string{"like"}, "reaction_type must be like")
	if details.HasAny() {
		return validationFailed(c, details)
	}

	item, err := h.activityService.ToggleReaction(c.UserContext(), c.Params("id"), user.ID, *request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidReactionInput):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrActivityNotFound):
			return response.NotFound(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to update reaction")
		}
	}

	return response.Success(c, fiber.StatusOK, "activity reaction updated", item)
}

func (h *ActivityHandler) ClaimFood(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	item, err := h.activityService.ClaimFood(c.UserContext(), c.Params("id"), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.NotFound(c, err.Error())
		case errors.Is(err, services.ErrFoodClaimNotAllowed):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrFoodClaimAlreadyExists):
			return response.Conflict(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to claim food")
		}
	}

	return response.Success(c, fiber.StatusCreated, "food claimed", item)
}

func (h *ActivityHandler) ListFoodClaims(c *fiber.Ctx) error {
	items, err := h.activityService.ListFoodClaims(c.UserContext(), c.Params("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.NotFound(c, err.Error())
		case errors.Is(err, services.ErrFoodClaimNotAllowed):
			return response.InvalidRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to load food claims")
		}
	}

	return response.Success(c, fiber.StatusOK, "food claims loaded", items)
}

func (h *ActivityHandler) RespondRice(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.AuthUser)
	if !ok {
		return response.Unauthorized(c, "authentication required")
	}

	item, err := h.activityService.RespondRice(c.UserContext(), c.Params("id"), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.NotFound(c, err.Error())
		case errors.Is(err, services.ErrRiceResponseNotAllowed):
			return response.InvalidRequest(c, err.Error())
		case errors.Is(err, services.ErrRiceResponseAlreadyExists):
			return response.Conflict(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to save rice response")
		}
	}

	return response.Success(c, fiber.StatusCreated, "rice response saved", item)
}

func (h *ActivityHandler) ListRiceResponses(c *fiber.Ctx) error {
	items, err := h.activityService.ListRiceResponses(c.UserContext(), c.Params("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrActivityNotFound):
			return response.NotFound(c, err.Error())
		case errors.Is(err, services.ErrRiceResponseNotAllowed):
			return response.InvalidRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to load rice responses")
		}
	}

	return response.Success(c, fiber.StatusOK, "rice responses loaded", items)
}
