package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"xlink/common/gen/user_service"
	"xlink/gateway/internal/handlers"
)

type IdCheckerService interface {
	GetUserIDByToken(request *user_service.GetUserIDByTokenRequest) (*user_service.GetUserIDByTokenResponse, error)
}

func AuthMiddleware(service IdCheckerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "missing Authorization header"})
		}

		userIdData, err := service.GetUserIDByToken(&user_service.GetUserIDByTokenRequest{Token: token})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid token (couldn't be found)"})
		}

		c.Context().SetUserValue(handlers.UserIdKey, userIdData.UserId)

		return c.Next()
	}
}

func AuthMiddlewareTokenParam(service IdCheckerService, paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query(paramName)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": fmt.Errorf("missing '%s' parameter", paramName)})
		}

		userIdData, err := service.GetUserIDByToken(&user_service.GetUserIDByTokenRequest{Token: token})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid token (couldn't be found)"})
		}

		c.Context().SetUserValue(handlers.UserIdKey, userIdData.UserId)

		return c.Next()
	}
}
