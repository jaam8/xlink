package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"xlink/common/gen/shortener"
	"xlink/gateway/internal/handlers"
	"xlink/gateway/internal/handlers/helpers"
)

type ShortenerOwnerChecker interface {
	GetLink(request *shortener.GetLinkRequest) (*shortener.Link, error)
}

func ShortenerOwnerOnlyMiddleware(idParamName string, shortenerService ShortenerOwnerChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIdValue := c.Context().Value(handlers.UserIdKey)
		if userIdValue == nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "unauthorized (Use auth middleware before Owner only middleware!!!)"})
		}
		userId := userIdValue.(string)

		var err error
		var linkId uuid.UUID

		linkId, err = helpers.ParseUUIDField(c, idParamName)
		if err != nil {
			return helpers.BadRequest(c, fmt.Errorf("invalid %s", idParamName))
		}

		var link *shortener.Link

		request := &shortener.GetLinkRequest{LinkId: linkId.String()}

		link, err = shortenerService.GetLink(request)
		if err != nil {
			return helpers.InternalServerError(c, fmt.Errorf("couldn't check link by id='%s': %w",
				request.LinkId, err))
		}

		if link.UserId != userId {
			return helpers.ForbiddenError(c)
		}

		return c.Next()
	}
}
