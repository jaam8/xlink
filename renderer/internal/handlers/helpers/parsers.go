package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func ParseNotEmptyStringField(ctx *fiber.Ctx, fieldName string) (string, error) {
	val := ctx.Params(fieldName)
	if len(val) == 0 {
		return "", fmt.Errorf("%s mustn't be empty", fieldName)
	}
	return val, nil
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
