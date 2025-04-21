package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type UpstreamNamesConfig struct {
	Shortener     string `yaml:"shortener" env:"SHORTENER"`
	UserService   string `yaml:"user_service" env:"USER_SERVICE"`
	Analytics     string `yaml:"analytics" env:"ANALYTICS"`
	FileGenerator string `yaml:"file_generator" env:"FILE_GENERATOR"`
}

type UpstreamPortsConfig struct {
	Shortener     string `yaml:"shortener" env:"SHORTENER"`
	UserService   string `yaml:"user_service" env:"USER_SERVICE"`
	Analytics     string `yaml:"analytics" env:"ANALYTICS"`
	FileGenerator string `yaml:"file_generator" env:"FILE_GENERATOR"`
}

type TimeoutsConfig struct {
	Shortener     int `yaml:"shortener" env:"SHORTENER"`
	UserService   int `yaml:"user_service" env:"USER_SERVICE"`
	Analytics     int `yaml:"analytics" env:"ANALYTICS"`
	FileGenerator int `yaml:"file_generator" env:"FILE_GENERATOR"`
}

type Config struct {
	UpstreamNames              UpstreamNamesConfig `yaml:"upstream_names" env-prefix:"UPSTREAM_NAME_"`
	UpstreamPorts              UpstreamPortsConfig `yaml:"upstream_ports" env-prefix:"UPSTREAM_PORT_"`
	Timeouts                   TimeoutsConfig      `yaml:"timeouts" env-prefix:"TIMEOUT_"`
	HTTPPort                   int                 `yaml:"http_port" env:"HTTP_PORT" env-default:"8080"`
	MaxConnections             int                 `yaml:"max_connections" env:"MAX_CONNECTIONS" env-default:"10"`
	MinConnections             int                 `yaml:"min_connections" env:"MIN_CONNECTIONS" env-default:"1"`
	MaxRetries                 uint                `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
	BaseRetryDelayMilliseconds uint                `yaml:"base_retry_delay_milliseconds" env:"BASE_RETRY_DELAY_MILLISECONDS" env-default:"200"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env vars: %v", err)
	}
	return cfg, nil
}
