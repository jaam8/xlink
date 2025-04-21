package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"xlink/common/kafka"
	"xlink/common/postgres"
	"xlink/common/redis"
)

type ShortenerConfig struct {
	KafkaTopic                   string `yaml:"kafka_topic" env-prefix:"KAFKA_TOPIC"`
	KafkaNumPartitions           int    `yaml:"kafka_num_partitions" env-default:"1"`
	KafkaReplicationFactor       int    `yaml:"kafka_replication_factor" env-default:"1"`
	RedisDB                      int    `yaml:"redis_db" env:"REDIS_DB" env-default:"1"`
	GRPCPort                     int    `yaml:"grpc_port" env:"GRPC_PORT"`
	ExpirationSeconds            int    `yaml:"expiration_seconds" env:"EXPIRATION_SECONDS" env-default:"500"`
	DefaultLinkExpirationMinutes int    `yaml:"default_link_expiration_minutes" env:"DEFAULT_LINK_EXPIRATION_MINUTES" env-default:"14400"`
	MigrationsPath               string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-default:"file:///app/migrations"`
}

type Config struct {
	Redis     redis.Config    `yaml:"redis" env-prefix:"REDIS_"`
	Postgres  postgres.Config `yaml:"postgres" env-prefix:"POSTGRES_"`
	Kafka     kafka.Config    `yaml:"kafka" env-prefix:"KAFKA_"`
	Shortener ShortenerConfig `yaml:"shortener" env-prefix:"SHORTENER_"`
}

func New() (Config, error) {
	var cfg Config
	// docker workdir - app/
	// local workdir - xlink/shortener
	if err := cleanenv.ReadConfig("../configs/config.yaml", &cfg); err != nil {
		fmt.Println(err.Error())
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to read env vars: %v", err)
		}
	}

	return cfg, nil
}
