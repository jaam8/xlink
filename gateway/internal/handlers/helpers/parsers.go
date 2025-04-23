package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func ParseUUID(value string) (uuid.UUID, error) {
	uuidField, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidField, nil
}

func ParseUUIDField(c *fiber.Ctx, fieldName string) (uuid.UUID, error) {
	return ParseUUID(c.Params(fieldName))
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

func ParseDateField(c *fiber.Ctx, fieldName string) (time.Time, error) {
	return ParseDate(c.Params(fieldName))
}

func ParseDate(date string) (time.Time, error) {
	dateField, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid '%s format' date: %v", time.DateOnly, err)
	}
	return dateField, nil
}
