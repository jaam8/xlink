package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"xlink/common/postgres"
	"xlink/common/redis"
)

type UserServiceConfig struct {
	GRPCPort               int    `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50050"`
	RedisDB                int    `yaml:"redis_db" env:"REDIS_DB"`
	TokenLength            int    `yaml:"token_len" env:"TOKEN_LEN" env-default:"32"`
	CacheExpirationSeconds int    `yaml:"cache_expiration_seconds" env:"CACHE_EXPIRATION_SECONDS" env-default:"3600"`
	MigrationsPath         string `yaml:"migrations_path" env:"MIGRATIONS_PATH"`
}

type ShortenerConfig struct {
	UpstreamNames string `yaml:"upstream_names" env:"UPSTREAM_NAMES"`
	UpstreamPorts string `yaml:"upstream_ports" env:"UPSTREAM_PORTS"`
	Timeouts      int    `yaml:"timeouts" env:"TIMEOUTS"`
}

type Config struct {
	Redis       redis.Config      `yaml:"redis" env-prefix:"REDIS_"`
	Postgres    postgres.Config   `yaml:"postgres" env-prefix:"POSTGRES_"`
	UserService UserServiceConfig `yaml:"user_service" env-prefix:"USER_SERVICE_"`
	Shortener   ShortenerConfig   `yaml:"shortener" env-prefix:"SHORTENER_"`
}

func New() (Config, error) {
	var cfg Config
	// local workdir - xlink/user_service
	if err := cleanenv.ReadConfig("../configs/config.yaml", &cfg); err != nil {
		fmt.Println(err.Error())
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to read env vars: %v", err)
		}
	}
	fmt.Println(cfg.Redis)
	fmt.Println(cfg.Postgres)
	fmt.Println(cfg.UserService)
	fmt.Println(cfg.Shortener)
	return cfg, nil
}
