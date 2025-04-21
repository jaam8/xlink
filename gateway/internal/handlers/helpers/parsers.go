package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ParseUUID(c *fiber.Ctx, fieldName string) (uuid.UUID, error) {
	uuidField, err := uuid.Parse(c.Params(fieldName))
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidField, nil
}
