package cache

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type ShortenerCacheRepositoryRedis struct {
	RedisClient *redis.Client
	Expiration  time.Duration
}

func NewShortenerCacheRepositoryRedis(redisClient *redis.Client, expiration time.Duration) *ShortenerCacheRepositoryRedis {
	return &ShortenerCacheRepositoryRedis{
		RedisClient: redisClient,
		Expiration:  expiration,
	}
}

func (t *ShortenerCacheRepositoryRedis) GetUrl(shortUrl string) (string, error) {
	commandResult := t.RedisClient.Get(shortUrl)
	if commandResult.Err() != nil {
		if errors.Is(commandResult.Err(), redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("couldn't get url %s: %v", shortUrl, commandResult.Err())
	}

	longUrl := commandResult.Val()
	return longUrl, nil
}

func (t *ShortenerCacheRepositoryRedis) SetUrl(shortUrl string, url string) error {
	commandResult := t.RedisClient.Set(shortUrl, url, t.Expiration)
	if commandResult.Err() != nil {
		return fmt.Errorf("couldn't set url %s to be %s: %v", shortUrl, url, commandResult.Err())
	}
	return nil
}

func (t *ShortenerCacheRepositoryRedis) DeleteUrl(shortUrl string) error {
	commandResult := t.RedisClient.Del(shortUrl)
	if commandResult.Err() != nil {
		return fmt.Errorf("couldn't delete url %s: %v", shortUrl, commandResult.Err())
	}
	return nil
}
