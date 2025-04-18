package helper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetValidatedId(request LinkBodyRequestOnlyId) (uuid.UUID, error) {
	var id uuid.UUID
	var err error

	id, err = uuid.Parse(request.GetId())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid id (can't parse uuid): %w", err)
	}

	return id, nil
}

func GetValidatedUserId(request LinkBodyRequest) (uuid.UUID, error) {
	var userId uuid.UUID
	var err error

	userId, err = uuid.Parse(request.GetUserId())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user id (can't parse uuid): %w", err)
	}

	return userId, nil
}

func GetValidatedGroupId(request LinkBodyRequest, defaultValue *uuid.UUID) (*uuid.UUID, error) {
	var err error

	if request.GetGroupId() != "" {
		var parsed uuid.UUID
		parsed, err = uuid.Parse(request.GetGroupId())
		if err != nil {
			return nil, fmt.Errorf("invalid group id (can't parse uuid): %w", err)
		}
		return &parsed, nil
	}

	return defaultValue, nil
}

func GetValidatedExpireAt(request LinkBodyRequest, defaultValue time.Time) (time.Time, error) {
	var expireAt = defaultValue
	var err error

	if request.GetExpireAt() != "" {
		expireAt, err = time.Parse(time.RFC3339, request.GetExpireAt())
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid expire at (can't parse RFC 3339): %w", err)
		}
	}

	return expireAt, nil
}

func ValidateStringNotEmpty(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return fmt.Errorf("string cannot be empty")
	}
	return nil
}

func ValidateUrl(str string) error {
	_, err := url.Parse(str)
	if err != nil {
		return fmt.Errorf("invalid url %s: %w", str, err)
	}
	return nil
}
