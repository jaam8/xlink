package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type BotConfig struct {
	BotToken   string `yaml:"bot_token" env:"BOT_TOKEN"`
	Host       string `yaml:"host" env:"HOST" env-default:"tg_bot"`
	Port       string `yaml:"port" env:"PORT" env-default:"50055"`
	BaseAPIURL string `yaml:"base_api_url" env:"BASE_URL" env-default:"http://nginx:80"`
	Timeouts   int    `yaml:"timeouts" env:"TIMEOUTS" env-default:"1000"`
}

type UserServiceConfig struct {
	UpstreamNames string `yaml:"upstream_names" env:"UPSTREAM_NAMES"`
	UpstreamPorts string `yaml:"upstream_ports" env:"UPSTREAM_PORTS"`
	Timeouts      int    `yaml:"timeouts" env:"TIMEOUTS"`
}

type Config struct {
	BotConfig   BotConfig         `yaml:"bot" env-prefix:"BOT_"`
	UserService UserServiceConfig `yaml:"user_service" env-prefix:"USER_SERVICE_"`
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
