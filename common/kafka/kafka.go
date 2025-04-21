package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
	"xlink/common/logger"
)

type Config struct {
	Host           string   `yaml:"host" env:"HOST" env-default:"kafka"`
	Port           uint16   `yaml:"port" env:"PORT" env-default:"9092"`
	Brokers        []string `yaml:"brokers" env:"BROKERS" env-separator:","`
	Topic          string   `yaml:"topic" env:"TOPIC"`
	GroupID        string   `yaml:"group_id" env:"GROUP_ID"`
	MinBytes       int      `yaml:"min_bytes" env:"MIN_BYTES" env-default:"10"`
	MaxBytes       int      `yaml:"max_bytes" env:"MAX_BYTES" env-default:"1048576"` // 1MB
	MaxWaitMs      int      `yaml:"max_wait_ms" env:"MAX_WAIT_MS" env-default:"500"`
	CommitInterval int      `yaml:"commit_interval_ms" env:"COMMIT_INTERVAL_MS" env-default:"1000"`
}

func NewReader(ctx context.Context, cfg Config, topic, groupID string) *kafka.Reader {
	l := logger.GetOrCreateLoggerFromCtx(ctx)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		MaxWait:        time.Duration(cfg.MaxWaitMs) * time.Millisecond,
		CommitInterval: time.Duration(cfg.CommitInterval) * time.Millisecond,
	})
	l.Info(ctx, "connected to Kafka topic",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
		zap.String("group_id", groupID),
	)
	return r
}

func NewWriter(ctx context.Context, cfg Config, topic string) *kafka.Writer {
	l := logger.GetOrCreateLoggerFromCtx(ctx)
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
	}

	l.Info(ctx, "created Kafka writer",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
	)
	return w
}
