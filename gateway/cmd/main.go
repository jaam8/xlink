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
	"xlink/gateway/internal/config"
	"xlink/gateway/internal/handlers/http_handlers"
	"xlink/gateway/internal/handlers/middlewares"
	"xlink/gateway/internal/ports/adapters/user_service_adapters"
	"xlink/gateway/internal/services"
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

	var usersGrpcPool *pool.GrpcPool

	userServiceAddress := fmt.Sprintf("%s:%s", mainConfig.UpstreamNames.UserService, mainConfig.UpstreamPorts.UserService)

	usersGrpcPool, err = pool.NewGrpcPool(ctx, pool.Config{
		Address:        userServiceAddress,
		MaxConnections: mainConfig.MaxConnections,
		MinConnections: mainConfig.MinConnections,
		DialOptions:    []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	})
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "couldn't create grpc pool for user service",
			zap.Int("MaxConnections", mainConfig.MaxConnections),
			zap.Int("MinConnections", mainConfig.MinConnections),
			zap.String("Address", userServiceAddress),
			zap.Error(err))
	}

	//region repos
	userServiceRepo := user_service_adapters.NewUserServiceRepositoryGRPC(usersGrpcPool)

	userService := services.NewUserService(
		userServiceRepo,
		mainConfig.MaxRetries,
		time.Millisecond*time.Duration(mainConfig.BaseRetryDelayMilliseconds),
	)
	//endregion repos

	//region handlers
	userServiceHandler := http_handlers.NewUserServiceHandler(userService)
	//endregion handlers

	//region routing
	app := fiber.New()

	//region api
	apiGroup := app.Group("/api")
	apiGroup.Use(middlewares.LoggerMiddleware())

	//region v1
	v1Group := apiGroup.Group("/v1")

	//region user v1
	userGroup := v1Group.Group("/user")
	userGroup.Post("/create", userServiceHandler.CreateUser)             //
	userGroup.Get("/:id", userServiceHandler.GetUser)                    // staff | admin
	userGroup.Post("/get/by-token", userServiceHandler.GetUserIDByToken) // admin
	userGroup.Post("/get/by-tg-id", userServiceHandler.GetUserIdByTgId)  // staff | admin
	userGroup.Patch("/:id", userServiceHandler.UpdateUser)               //
	userGroup.Post("/token/check", userServiceHandler.CheckToken)        // admin
	userGroup.Post("/token/refresh", userServiceHandler.RefreshToken)    //
	userGroup.Delete("/:id", userServiceHandler.DeleteUser)              // staff | admin
	userGroup.Get("/role/:id", userServiceHandler.GetRole)               // staff | admin
	//endregion user v1
	//endregion v1
	//endregion api
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
