package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func InvalidDateBadRequest(c *fiber.Ctx, fieldName string) error {
	return BadRequest(c, fmt.Errorf("invalid %s: must follow %s pattern", fieldName, time.DateOnly))
}

func BadRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func InternalServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func NotAuthenticatedError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func NotFoundError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": message,
	})
}

func ForbiddenError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": "you don't have access to this resource",
	})
}
