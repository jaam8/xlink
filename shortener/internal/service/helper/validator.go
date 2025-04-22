package helper

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetValidatedId(request LinkRequestOnlyLinkId) (uuid.UUID, error) {
	var id uuid.UUID
	var err error

	id, err = uuid.Parse(request.GetLinkId())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid id (can't parse uuid): %w", err)
	}

	return id, nil
}

func GetValidatedUserId(request LinkRequestOnlyUserId) (uuid.UUID, error) {
	var userId uuid.UUID
	var err error

	userId, err = uuid.Parse(request.GetUserId())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user id (can't parse uuid): %w", err)
	}

	return userId, nil
}

func GetValidatedExpireAt(request LinkRequestOnlyExpireAt) (*time.Time, error) {
	var expireAt *time.Time = nil

	if request.GetExpireAt() != nil {
		expireAtValue := request.GetExpireAt().AsTime()
		if expireAtValue.Before(time.Now()) {
			return nil, fmt.Errorf("expire at time is out of date")
		}

		expireAt = &expireAtValue
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

func ValidateShortLink(str string) (string, error) {
	shortLink, err := url.Parse(str)
	if err != nil {
		return "", fmt.Errorf("invalid short_link %s: %v", str, err)
	}
	return shortLink.String(), nil
}

func ValidateIPAddress(ipAddress string) (string, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "", fmt.Errorf("invalid ip address: %s", ipAddress)
	}
	return ip.String(), nil
}

func ValidateNotEmptyStr(str string) (string, error) {
	if len(strings.TrimSpace(str)) == 0 {
		return "", fmt.Errorf("string cannot be empty")
	}
	return str, nil
}
