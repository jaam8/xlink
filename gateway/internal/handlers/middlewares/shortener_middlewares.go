package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"xlink/gateway/internal/handlers/helpers"
)

type ShortenerOwnerChecker interface {
}

func ShortenerOwnerOnlyMiddleware(idParamName string, shortenerService ShortenerOwnerChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := helpers.ParseUUIDField(c, idParamName)
		if err != nil {
			return helpers.BadRequest(c, fmt.Sprintf("invalid %s", idParamName))
		}

		// check link owner via GetLink

		return c.Next()
	}
}
