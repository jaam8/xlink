package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"time"
	"xlink/common/logger"
)

type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"redis"`
	Port     uint16 `yaml:"port" env:"PORT" env-default:"6379"`
	Username string `yaml:"user" env:"USER"`
	Password string `yaml:"user_password" env:"USER_PASSWORD"`

	MaxRetries int `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
	PoolSize   int `yaml:"pool_size" env:"POOL_SIZE" env-default:"10"`

	DialTimeoutSeconds  int `yaml:"dial_timeout_seconds" env:"DIAL_TIMEOUT_SECONDS" env-default:"5"`
	ReadTimeoutSeconds  int `yaml:"read_timeout_seconds" env:"READ_TIMEOUT_SECONDS" env-default:"3"`
	WriteTimeoutSeconds int `yaml:"write_timeout_seconds" env:"WRITE_TIMEOUT_SECONDS" env-default:"3"`
}

// NewRedisClient try to connect to Redis and get the client
func NewRedisClient(ctx context.Context, config Config, redisDB int) (*redis.Client, error) {
	option := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username:     config.Username,
		Password:     config.Password,
		DB:           redisDB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  time.Duration(config.DialTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeoutSeconds) * time.Second,
		PoolSize:     config.PoolSize,
	}
	client := redis.NewClient(option)
	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}
	logger.GetOrCreateLoggerFromCtx(ctx).
		Info(ctx, "connected to a redis database",
			zap.String("addr", option.Addr),
			zap.Int("db", option.DB),
		)
	return client, nil
}
