package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func ParseUUIDField(c *fiber.Ctx, fieldName string) (uuid.UUID, error) {
	uuidField, err := uuid.Parse(c.Params(fieldName))
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidField, nil
}

func ParseDateTimeField(c *fiber.Ctx, fieldName string) (time.Time, error) {
	return ParseDateTime(c.Params(fieldName))
}

func ParseDateTime(dateTime string) (time.Time, error) {
	timeField, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid datetime: %v", err)
	}
	return timeField, nil
}
