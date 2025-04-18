package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"xlink/common/redis"
)

type Config struct {
	Redis       redis.Config `yaml:"redis" env:"REDIS" env-prefix:""`
	GRPCPort    int          `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50050"`
	TokenLength int          `yaml:"token_len" env:"TOKEN_LEN" env-default:"32"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env vars: %v", err)
	}
	return cfg, nil
}
