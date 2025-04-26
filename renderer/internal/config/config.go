package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type UpstreamNamesConfig struct {
	Analytics string `yaml:"analytics" env:"ANALYTICS"`
}

type UpstreamPortsConfig struct {
	Analytics string `yaml:"analytics" env:"ANALYTICS"`
}

type TimeoutsConfig struct {
	Analytics int `yaml:"analytics" env:"ANALYTICS"`
}

type GrpcPoolConfig struct {
	MaxConnections             int  `yaml:"max_connections" env:"MAX_CONNECTIONS" env-default:"10"`
	MinConnections             int  `yaml:"min_connections" env:"MIN_CONNECTIONS" env-default:"1"`
	MaxRetries                 uint `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
	BaseRetryDelayMilliseconds uint `yaml:"base_retry_delay_milliseconds" env:"BASE_RETRY_DELAY_MILLISECONDS" env-default:"200"`
}

type Config struct {
	UpstreamNames UpstreamNamesConfig `yaml:"upstream_names" env-prefix:"UPSTREAM_NAMES_"`
	UpstreamPorts UpstreamPortsConfig `yaml:"upstream_ports" env-prefix:"UPSTREAM_PORTS_"`
	GrpcPool      GrpcPoolConfig      `yaml:"grpc_pool" env-prefix:"GRPC_POOL_"`
	Timeouts      TimeoutsConfig      `yaml:"timeouts" env-prefix:"TIMEOUT_"`
	HTTPPort      int                 `yaml:"http_port_renderer" env:"HTTP_PORT_RENDERER" env-default:"8085"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env vars: %v", err)
	}
	return cfg, nil
}
