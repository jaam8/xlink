package adapters

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"tokenservice/internal/utils"
)

type TokensRepositoryRedis struct {
	RedisClient *redis.Client
	tokenLength int8
}

func NewTokensRepositoryRedis(redisClient *redis.Client, tokenLength int8) *TokensRepositoryRedis {
	return &TokensRepositoryRedis{RedisClient: redisClient, tokenLength: tokenLength}
}

func (t TokensRepositoryRedis) Check(userId string, token string) (bool, error) {
	commandResult := t.RedisClient.Get(userId)
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

func (t TokensRepositoryRedis) Create(userId string) (string, error) {
	token := utils.GenerateToken(t.tokenLength)

	commandResult := t.RedisClient.Set(userId, token, 0)
	if commandResult.Err() != nil {
		return "", commandResult.Err()
	}
	return token, nil
}

func (t TokensRepositoryRedis) Refresh(userId string, token string) (string, error) {
	// todo added domain errors
	status, err := t.Check(userId, token)
	if err != nil {
		if !status {
			return "", err
		}
		return "", err
	}

	token, err = t.Create(userId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (t TokensRepositoryRedis) Delete(userId string) (bool, error) {
	result := t.RedisClient.Del(userId)
	if result.Err() != nil {
		return false, result.Err()
	}
	return true, nil
}
