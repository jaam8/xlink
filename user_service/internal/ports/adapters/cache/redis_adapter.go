package cache

import (
	"errors"
	"github.com/go-redis/redis/v7"
)

const (
	userIdKeyPrefix = "user-id-"
	tokenKeyPrefix  = "token-"
)

type UserCacheRepositoryRedis struct {
	RedisClient *redis.Client
}

func getUserIdKey(userId string) string {
	return userIdKeyPrefix + userId
}

func getTokenKey(token string) string {
	return tokenKeyPrefix + token
}

func NewUserCacheRepositoryRedis(redisClient *redis.Client) *UserCacheRepositoryRedis {
	return &UserCacheRepositoryRedis{RedisClient: redisClient}
}

func (t UserCacheRepositoryRedis) CheckToken(userId string, token string) (bool, error) {
	commandResult := t.RedisClient.Get(getUserIdKey(userId))
	if commandResult.Err() != nil {
		if errors.Is(commandResult.Err(), redis.Nil) {
			return false, nil
		}
		return false, commandResult.Err()
	}
	if commandResult.Val() != token {
		return false, nil
	}
	return true, nil
}

func (t UserCacheRepositoryRedis) GetToken(userId string) (string, error) {
	commandResult := t.RedisClient.Get(getUserIdKey(userId))
	if commandResult.Err() != nil {
		if errors.Is(commandResult.Err(), redis.Nil) {
			return "", nil
		}
		return "", commandResult.Err()
	}
	return commandResult.Val(), nil
}

func (t UserCacheRepositoryRedis) SetToken(userId string, token string) error {
	commandResult := t.RedisClient.Set(getUserIdKey(userId), token, 0)
	if commandResult.Err() != nil {
		return commandResult.Err()
	}

	return nil
}
