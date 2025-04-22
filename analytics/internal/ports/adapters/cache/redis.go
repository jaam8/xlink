package cache

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisAdapter struct {
	RedisClient *redis.Client
	Timezone    string
}

func NewRedisAdapter(redisClient *redis.Client, timezone string) *RedisAdapter {
	return &RedisAdapter{RedisClient: redisClient, Timezone: timezone}
}

func (r *RedisAdapter) CheckVisitorToken(visitorToken, shortLink string) (bool, error) {
	result := r.RedisClient.Get(visitorToken)
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get visitor_token %s: %v", visitorToken, result.Err())
	}

	resultShortLink := result.Val()
	if shortLink == resultShortLink {
		return true, nil
	}
	return false, nil
}

func (r *RedisAdapter) SetVisitorToken(visitorToken, shortLink string) error {
	loc, err := time.LoadLocation(r.Timezone)
	if err != nil {
		return fmt.Errorf("failed to load timezone: %v", err)
	}
	now := time.Now().In(loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day()+1,
		0, 0, 0, 0, loc).Add(-time.Nanosecond)

	expiration := endOfDay.Sub(now)
	if expiration <= 0 {
		return fmt.Errorf("calculated expiration is non-positive: %v", expiration)
	}

	result := r.RedisClient.Set(visitorToken, shortLink, expiration)
	if result.Err() != nil {
		return fmt.Errorf("failed to set visitor_token %s: %v", visitorToken, result.Err())
	}

	return nil
}
