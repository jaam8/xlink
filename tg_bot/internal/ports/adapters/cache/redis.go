package cache

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisAdapter struct {
	RedisClient *redis.Client
	Expiration  time.Duration
}

func NewRedisAdapter(redisClient *redis.Client, expiration time.Duration) *RedisAdapter {
	return &RedisAdapter{RedisClient: redisClient, Expiration: expiration}
}

func (r *RedisAdapter) GetUserToken(tgID string) (string, error) {
	result := r.RedisClient.Get(tgID)
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get user_token %s: %v", tgID, result.Err())
	}
	return result.Val(), nil
}

func (r *RedisAdapter) SetUserToken(tgID, userToken string) error {
	result := r.RedisClient.Set(tgID, userToken, r.Expiration)
	if result.Err() != nil {
		return fmt.Errorf("failed to set user_token %s by tg_id %s: %v", userToken, tgID, result.Err())
	}
	return nil
}
