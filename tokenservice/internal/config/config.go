package config

import (
	"fmt"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/storage/redis"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Redis       redis.Config `yaml:"redis" env:"REDIS" env-prefix:""`
	GRPCPort    int          `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50050"`
	TokenLength int          `yaml:"token_len" env:"TOKEN_LEN" env-default:"32"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{},
			fmt.Errorf("failed to read env variables after accessing .env: %w", err)
	}
	return cfg, nil
}
