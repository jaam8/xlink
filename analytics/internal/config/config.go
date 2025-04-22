package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"xlink/common/clickhouse"
	"xlink/common/kafka"
	"xlink/common/redis"
)

type AnalyticsConfig struct {
	KafkaTopic     string `yaml:"kafka_topic" env-prefix:"KAFKA_TOPIC"`
	KafkaGroupID   string `yaml:"kafka_group_id" env-prefix:"KAFKA_GROUP"`
	RedisDB        int    `yaml:"redis_db" env:"REDIS_DB" env-default:"1"`
	GRPCPort       int    `yaml:"grpc_port" env:"GRPC_PORT"`
	Timezone       string `yaml:"timezone" env:"TIMEZONE"`
	BatchSize      int    `yaml:"batch_size" env:"BATCH_SIZE"`
	FlushTimeout   int    `yaml:"flush_timeout" env:"FLUSH_TIMEOUT"`
	MigrationsPath string `yaml:"migrations_path" env:"MIGRATIONS_PATH"`
}

type Config struct {
	Redis      redis.Config      `yaml:"redis" env-prefix:"REDIS_"`
	ClickHouse clickhouse.Config `yaml:"clickhouse" env-prefix:"CLICKHOUSE_"`
	Kafka      kafka.Config      `yaml:"kafka" env-prefix:"KAFKA_"`
	Analytics  AnalyticsConfig   `yaml:"analytics" env-prefix:"ANALYTICS_"`
}

func New() (Config, error) {
	var cfg Config
	// docker workdir - app/
	// local workdir - xlink/analytics
	if err := cleanenv.ReadConfig("../configs/config.yaml", &cfg); err != nil {
		fmt.Println(err.Error())
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to read env vars: %v", err)
		}
	}

	return cfg, nil
}
