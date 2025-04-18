package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"xlink/common/logger"
	"xlink/common/redis"
	"xlink/user_service/internal/config"
	"xlink/user_service/internal/ports/adapters"
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

	redisClient, err := redis.NewRedisClient(ctx, redisCfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to Redis database: %w", err))
	}

	tokensRepo := adapters.NewTokensRepositoryRedis(redisClient, int8(cfg.TokenLength))

	grpcServer, err := runner.CreateGRPC(tokensRepo)
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	go runner.RunGRPC(ctx, grpcServer, cfg.GRPCPort)

	<-ctx.Done()

	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "server stopped")
}
