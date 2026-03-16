package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanprayoga/messhub/backend/internal/response"
	"github.com/ryanprayoga/messhub/backend/internal/validation"
)

const (
	maxNameLength           = 120
	maxPhoneLength          = 32
	maxAvatarURLLength      = 512
	maxWalletCategoryLength = 80
	maxWalletDescLength     = 500
	maxWifiProofLength      = 512
	maxWifiNoteLength       = 500
	maxActivityTitleLength  = 140
	maxActivityBodyLength   = 2000
	maxCommentLength        = 500
	maxMessNameLength       = 120
	maxBankFieldLength      = 120
	maxRejectReasonLength   = 300
	maxUsernameLength       = 32
)

func invalidPayload(c *fiber.Ctx, entity string) error {
	return response.InvalidRequest(c, "invalid "+strings.TrimSpace(entity)+" payload")
}

func validationFailed(c *fiber.Ctx, details validation.Errors) error {
	return response.ValidationFailed(c, "invalid input", details)
}
