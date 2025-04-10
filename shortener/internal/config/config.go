package config

import (
	"fmt"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/storage/postgres"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/storage/redis"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Redis                        redis.Config    `yaml:"redis" env:"REDIS" env-prefix:""`
	Postgres                     postgres.Config `yaml:"postgres" env:"POSTGRES" env-prefix:""`
	GRPCPort                     int             `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50050"`
	ExpirationSeconds            int             `yaml:"expiration_seconds" env:"EXPIRATION_SECONDS" env-default:"500"`
	DefaultLinkExpirationMinutes int             `yaml:"default_link_expiration_minutes" env:"DEFAULT_LINK_EXPIRATION_MINUTES" env-default:"14400"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{},
			fmt.Errorf("failed to read env variables after accessing .env: %w", err)
	}
	return cfg, nil
}
