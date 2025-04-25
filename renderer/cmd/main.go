package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"time"
	"xlink/common/grpc/pool"
	"xlink/common/logger"
	"xlink/renderer/internal/config"
	"xlink/renderer/internal/handlers/http_handlers"
	"xlink/renderer/internal/handlers/middlewares"
	"xlink/renderer/internal/ports/adapters/analytics_service_adapters"
	"xlink/renderer/internal/ports/adapters/drawer_adapters"
	"xlink/renderer/internal/services"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, _ = logger.New(ctx)

	mainConfig, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "couldn't load configs", zap.Error(err))
	}

	//region analytics grpc pool
	var analyticsGrpcPool *pool.GrpcPool

	analyticsServiceAddress := fmt.Sprintf("%s:%s", mainConfig.UpstreamNames.Analytics, mainConfig.UpstreamPorts.Analytics)

	analyticsGrpcPool, err = pool.NewGrpcPool(ctx, pool.Config{
		Address:        analyticsServiceAddress,
		MaxConnections: mainConfig.GrpcPool.MaxConnections,
		MinConnections: mainConfig.GrpcPool.MinConnections,
		DialOptions:    []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	})
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "couldn't create grpc pool for analytics service",
			zap.Int("MaxConnections", mainConfig.GrpcPool.MaxConnections),
			zap.Int("MinConnections", mainConfig.GrpcPool.MinConnections),
			zap.String("Address", analyticsServiceAddress),
			zap.Error(err))
	}
	//endregion

	//region repos
	analyticsServiceRepo := analytics_service_adapters.NewAnalyticsServiceRepositoryGRPC(analyticsGrpcPool)
	analyticsService := services.NewAnalyticsService(
		analyticsServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		time.Millisecond*time.Duration(mainConfig.GrpcPool.BaseRetryDelayMilliseconds),
	)

	drawerRepo := drawer_adapters.NewDrawerRepositoryEcharts()
	drawerService := services.NewDrawerService(drawerRepo)
	//endregion repos

	//region handlers
	analyticsServiceHandler := http_handlers.NewRendererHandler(analyticsService, drawerService)
	//endregion handlers

	//region middlewares
	loggingMiddleware := middlewares.LoggerMiddleware()
	//endregion middlewares

	//region routing
	app := fiber.New()
	app.Use(loggingMiddleware)
	app.Get("/image", analyticsServiceHandler.Image)
	//endregion routing

	go func() {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "listening on port", zap.Int("port", mainConfig.HTTPPort))

		listenError := app.Listen(fmt.Sprintf(":%d", mainConfig.HTTPPort))

		if listenError != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "error while running the application", zap.Error(listenError))
		}
	}()

	<-ctx.Done()

	err = app.Shutdown()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "error shutting down", zap.Error(err))
	}
}
