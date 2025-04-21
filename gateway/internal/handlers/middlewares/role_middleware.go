package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"xlink/common/gen/user_service"
)

type RoleCheckerService interface {
	GetUserIDByToken(request *user_service.GetUserIDByTokenRequest) (*user_service.GetUserIDByTokenResponse, error)
	GetRole(request *user_service.GetRoleRequest) (*user_service.GetRoleResponse, error)
}

func RoleMiddleware(requireIsStaff bool, requireIsAdmin bool, service RoleCheckerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if requireIsStaff || requireIsAdmin {
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

			userId := userIdData.UserId

			var userRoleData *user_service.GetRoleResponse
			userRoleData, err = service.GetRole(&user_service.GetRoleRequest{UserId: userId})
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).
					JSON(fiber.Map{"error": "invalid token (roles info couldn't be found)"})
			}

			if !userRoleData.IsAdmin && requireIsAdmin || !userRoleData.IsStaff && requireIsStaff {
				return c.Status(fiber.StatusForbidden).
					JSON(fiber.Map{"error": "not enough permissions"})
			}
		}

		return c.Next()
	}
}
