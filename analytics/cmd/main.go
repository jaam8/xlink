package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"time"
	"xlink/analytics/internal/config"
	"xlink/analytics/internal/ports/adapters/cache"
	"xlink/analytics/internal/ports/adapters/consumer"
	"xlink/analytics/internal/ports/adapters/shortener_adapters"
	"xlink/analytics/internal/ports/adapters/storage"
	"xlink/analytics/internal/server"
	"xlink/common/clickhouse"
	"xlink/common/kafka"
	"xlink/common/logger"
	"xlink/common/redis"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx,
			"failed to load config", zap.Error(err))
	}

	redisCfg := cfg.Redis
	clickHouseCfg := cfg.ClickHouse
	kafkaCfg := cfg.Kafka
	analyticsCfg := cfg.Analytics
	shortenerCfg := cfg.Shortener

	redisClient, err := redis.NewRedisClient(ctx, redisCfg, analyticsCfg.RedisDB)
	//defer redisClient.Close()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to Redis database: %v", err))
	}
	clickhouseClient, err := clickhouse.New(ctx, clickHouseCfg)
	//nolint
	defer clickhouseClient.Close()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to ClickHouse database: %v", err))
	}

	kafkaConsumer := kafka.NewReader(ctx,
		kafkaCfg, analyticsCfg.KafkaTopic, analyticsCfg.KafkaGroupID)
	//nolint
	defer kafkaConsumer.Close()
	log.Printf("success connect to kafka at %s:%d", kafkaCfg.Host, kafkaCfg.Port)

	//clickhouse migration
	err = clickhouse.Migrate(clickHouseCfg, analyticsCfg.MigrationsPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to migrate clickhouse database: %w", err))
	}

	redisAdapter := cache.NewRedisAdapter(redisClient, analyticsCfg.Timezone)
	clickhouseAdapter := storage.NewClickHouseAdapter(clickhouseClient)
	kafkaAdapter := consumer.NewKafkaAdapter(kafkaConsumer)
	shortenerAdapter := shortener_adapters.NewShortenerAdapter(
		fmt.Sprintf("%s:%s", shortenerCfg.UpstreamNames, shortenerCfg.UpstreamPorts),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		time.Millisecond*time.Duration(shortenerCfg.Timeouts),
	)

	Server := server.NewService(redisAdapter, clickhouseAdapter, kafkaAdapter, shortenerAdapter)

	grpcServer, err := server.CreateGRPC(Server)

	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	go server.RunGRPC(ctx, grpcServer, analyticsCfg.GRPCPort)

	go Server.HandleConsumer(ctx, analyticsCfg.BatchSize,
		time.Second*time.Duration(analyticsCfg.FlushTimeout),
	)

	<-ctx.Done()

	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "ANALYTICS server stopped")
}
