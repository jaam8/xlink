package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"xlink/common/redis"
)

type BotConfig struct {
	MaxRetries        uint   `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
	BaseRetryDelay    int    `yaml:"base_retry_delay_second" env:"BASE_RETRY_DELAY_SECOND" env-default:"1"`
	GatewayServerUrl  string `yaml:"local_server_url" env:"LOCAL_SERVER_URL"`
	BotToken          string `yaml:"bot_token" env:"BOT_TOKEN"`
	Host              string `yaml:"host" env:"HOST" env-default:"tg_bot"`
	Port              string `yaml:"port" env:"PORT" env-default:"50055"`
	BaseAPIURL        string `yaml:"base_api_url" env:"BASE_URL" env-default:"http://nginx:80"`
	RedisDB           int    `yaml:"redis_db" env:"REDIS_DB" env-default:"3"`
	ExpirationSeconds int    `yaml:"expiration_seconds" env:"EXPIRATION_SECONDS" env-default:"500"`
	Timeouts          int    `yaml:"timeouts" env:"TIMEOUTS" env-default:"1000"`
}

type UserServiceConfig struct {
	UpstreamNames string `yaml:"upstream_names" env:"UPSTREAM_NAMES"`
	UpstreamPorts int    `yaml:"upstream_ports" env:"UPSTREAM_PORTS"`
	Timeouts      int    `yaml:"timeouts" env:"TIMEOUTS"`
}

type Config struct {
	BotConfig   BotConfig         `yaml:"bot" env-prefix:"BOT_"`
	UserService UserServiceConfig `yaml:"user_service" env-prefix:"USER_SERVICE_"`
	RedisConfig redis.Config      `yaml:"redis" env-prefix:"REDIS_"`
}

func New() (Config, error) {
	var cfg Config
	// docker workdir - app/
	// local workdir - xlink/tg_bot
	if err := cleanenv.ReadConfig("../configs/config.yaml", &cfg); err != nil {
		fmt.Println(err.Error())
		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to read env vars: %v", err)
		}
	}

	return cfg, nil
}
