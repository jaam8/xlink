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
	if result.Val() != tgID {
		return "", fmt.Errorf("(cache value tg_id) %s not equal (response tg_id) %s", result.Val(), tgID)
	}
	return result.Val(), nil
}

func (r *RedisAdapter) SetUserToken(tgID, userToken string) error {
	resultGet := r.RedisClient.Get(tgID)
	if resultGet.Err() == nil {
		if resultGet.Val() != userToken {
			return fmt.Errorf("(cache value user_token) %s not equal (response user_token) %s",
				resultGet.Val(), userToken)
		}
	}
	result := r.RedisClient.Set(tgID, userToken, r.Expiration)
	if result.Err() != nil {
		return fmt.Errorf("failed to set user_token %s by tg_id %s: %v", userToken, tgID, result.Err())
	}

	return nil
}
