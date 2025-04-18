package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"time"
	"xlink/common/logger"
	"xlink/common/postgres"
	"xlink/common/redis"
	"xlink/shortener/internal/config"
	"xlink/shortener/internal/ports/adapters/cache"
	"xlink/shortener/internal/ports/adapters/sender"
	"xlink/shortener/internal/ports/adapters/storage"
	"xlink/shortener/internal/runner"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to load config", zap.Error(err))
	}

	redisCfg := cfg.Redis
	postgresCfg := cfg.Postgres

	redisClient, err := redis.NewRedisClient(ctx, redisCfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to Redis database: %w", err))
	}

	postgresClient, err := postgres.New(ctx, postgresCfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to Postgres database: %w", err))
	}

	// postgres migration
	err = postgres.Migrate(ctx, postgresCfg, postgresCfg.MigrationsPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to migrate postgres database: %w", err))
	}

	repositoryRedis := cache.NewShortenerCacheRepositoryRedis(
		redisClient, time.Duration(cfg.ExpirationSeconds)*time.Second,
	)
	repositoryPostgres := storage.NewShortenerStorageRepositoryPostgres(postgresClient)
	senderRepositoryMock := sender.NewShortenerSenderRepositoryMock()

	grpcServer, err := runner.CreateGRPC(
		repositoryRedis, repositoryPostgres, senderRepositoryMock,
		time.Minute*time.Duration(cfg.DefaultLinkExpirationMinutes),
	)
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	go runner.RunGRPC(ctx, grpcServer, cfg.GRPCPort)

	<-ctx.Done()

	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "server stopped")
}
