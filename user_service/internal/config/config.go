package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"xlink/common/postgres"
	"xlink/common/redis"
)

type UpstreamNamesConfig struct {
	Shortener string `yaml:"shortener" env:"SHORTENER"`
}

type UpstreamPortsConfig struct {
	Shortener string `yaml:"shortener" env:"SHORTENER"`
}

type TimeoutsConfig struct {
	Shortener int `yaml:"shortener" env:"SHORTENER"`
}

type Config struct {
	Redis                  redis.Config        `yaml:"redis" env-prefix:"REDIS_"`
	Postgres               postgres.Config     `yaml:"postgres" env-prefix:"POSTGRES_"`
	UpstreamNames          UpstreamNamesConfig `yaml:"upstream_names" env-prefix:"UPSTREAM_NAME_"`
	UpstreamPorts          UpstreamPortsConfig `yaml:"upstream_ports" env-prefix:"UPSTREAM_PORT_"`
	Timeouts               TimeoutsConfig      `yaml:"timeouts" env-prefix:"TIMEOUT_"`
	GRPCPort               int                 `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50050"`
	TokenLength            int                 `yaml:"token_len" env:"TOKEN_LEN" env-default:"32"`
	CacheExpirationSeconds int                 `yaml:"cache_expiration_seconds" env:"CACHE_EXPIRATION_SECONDS" env-default:"3600"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env vars: %v", err)
	}
	return cfg, nil
}
