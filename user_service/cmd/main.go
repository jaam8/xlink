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
	"xlink/common/logger"
	"xlink/common/postgres"
	"xlink/common/redis"
	"xlink/user_service/internal/config"
	"xlink/user_service/internal/ports/adapters/cache"
	"xlink/user_service/internal/ports/adapters/shortener_adapters"
	"xlink/user_service/internal/ports/adapters/storage"
	"xlink/user_service/internal/runner"
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

	cacheRepo := cache.NewUserCacheRepositoryRedis(redisClient, time.Second*time.Duration(cfg.CacheExpirationSeconds))
	storageRepo := storage.NewUserStorageRepositoryPostgres(postgresClient, int8(cfg.TokenLength))

	shortenerRepo := shortener_adapters.NewShortenerRepositoryGRPC(
		fmt.Sprintf("%s:%s", cfg.UpstreamNames.Shortener, cfg.UpstreamPorts.Shortener),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		time.Millisecond*time.Duration(cfg.Timeouts.Shortener),
	)

	grpcServer, err := runner.CreateGRPC(cacheRepo, storageRepo, shortenerRepo)
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	go runner.RunGRPC(ctx, grpcServer, cfg.GRPCPort)

	<-ctx.Done()

	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "server stopped")
}
