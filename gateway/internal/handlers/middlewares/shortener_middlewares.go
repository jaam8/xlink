package middlewares

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"xlink/common/gen/shortener"
	"xlink/gateway/internal/handlers"
	"xlink/gateway/internal/handlers/helpers"
)

type ShortenerOwnerChecker interface {
	GetLinkIdByShortLink(request *shortener.GetLinkIdByShortLinkRequest) (*shortener.GetLinkIdByShortLinkResponse, error)
	GetLink(request *shortener.GetLinkRequest) (*shortener.Link, error)
}

func ShortenerOwnerOnlyMiddleware(shortLinkParamName string, shortenerService ShortenerOwnerChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIdValue := c.Context().Value(handlers.UserIdKey)
		if userIdValue == nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "unauthorized (Use auth middleware before Owner only middleware!!!)"})
		}
		userId := userIdValue.(string)

		shortLink := c.Params(shortLinkParamName)
		if shortLink == "" {
			return helpers.BadRequest(c, errors.New("short link must not be empty"))
		}

		var err error
		var linkId string
		var linkIdResponse *shortener.GetLinkIdByShortLinkResponse

		linkIdResponse, err = shortenerService.GetLinkIdByShortLink(&shortener.GetLinkIdByShortLinkRequest{ShortLink: shortLink})
		if err != nil {
			return helpers.InternalServerError(c,
				fmt.Errorf("error while trying to get link id from short link: %w", err))
		}

		linkId = linkIdResponse.LinkId

		var link *shortener.Link

		request := &shortener.GetLinkRequest{LinkId: linkId}

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
