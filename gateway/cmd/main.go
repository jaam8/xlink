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
	"xlink/gateway/internal/ports/adapters/analytics_service_adapters"
	"xlink/gateway/internal/ports/adapters/shortener_service_adapters"
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

	//region user grpc pool
	var usersGrpcPool *pool.GrpcPool

	userServiceAddress := fmt.Sprintf("%s:%s", mainConfig.UpstreamNames.UserService, mainConfig.UpstreamPorts.UserService)

	usersGrpcPool, err = pool.NewGrpcPool(ctx, pool.Config{
		Address:        userServiceAddress,
		MaxConnections: mainConfig.GrpcPool.MaxConnections,
		MinConnections: mainConfig.GrpcPool.MinConnections,
		DialOptions:    []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	})
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "couldn't create grpc pool for user service",
			zap.Int("MaxConnections", mainConfig.GrpcPool.MaxConnections),
			zap.Int("MinConnections", mainConfig.GrpcPool.MinConnections),
			zap.String("Address", userServiceAddress),
			zap.Error(err))
	}
	//endregion

	//region shortener grpc pool
	var shortenerGrpcPool *pool.GrpcPool

	shortenerServiceAddress := fmt.Sprintf("%s:%s", mainConfig.UpstreamNames.Shortener, mainConfig.UpstreamPorts.Shortener)

	shortenerGrpcPool, err = pool.NewGrpcPool(ctx, pool.Config{
		Address:        shortenerServiceAddress,
		MaxConnections: mainConfig.GrpcPool.MaxConnections,
		MinConnections: mainConfig.GrpcPool.MinConnections,
		DialOptions:    []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	})
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "couldn't create grpc pool for shortener service",
			zap.Int("MaxConnections", mainConfig.GrpcPool.MaxConnections),
			zap.Int("MinConnections", mainConfig.GrpcPool.MinConnections),
			zap.String("Address", shortenerServiceAddress),
			zap.Error(err))
	}
	//endregion

	//region shortener grpc pool
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
	userServiceRepo := user_service_adapters.NewUserServiceRepositoryGRPC(usersGrpcPool)
	shortenerServiceRepo := shortener_service_adapters.NewShortenerServiceRepositoryGRPC(shortenerGrpcPool)
	analyticsServiceRepo := analytics_service_adapters.NewAnalyticsServiceRepositoryGRPC(analyticsGrpcPool)

	userService := services.NewUserService(
		userServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		time.Millisecond*time.Duration(mainConfig.GrpcPool.BaseRetryDelayMilliseconds),
	)
	shortenerService := services.NewShortenerService(
		shortenerServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		time.Millisecond*time.Duration(mainConfig.GrpcPool.BaseRetryDelayMilliseconds),
	)
	analyticsService := services.NewAnalyticsService(
		analyticsServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		time.Millisecond*time.Duration(mainConfig.GrpcPool.BaseRetryDelayMilliseconds),
	)
	//endregion repos

	//region handlers
	userServiceHandler := http_handlers.NewUserServiceHandler(userService)
	shortenerServiceHandler := http_handlers.NewShortenerServiceHandler(shortenerService, userService)
	analyticsServiceHandler := http_handlers.NewAnalyticsServiceHandler(analyticsService)
	//endregion handlers

	//region middlewares
	loggingMiddleware := middlewares.LoggerMiddleware()
	authMiddleware := middlewares.AuthMiddleware(userService)
	isAdminMiddleware := middlewares.RoleMiddleware(false, true, userService)
	isStaffMiddleware := middlewares.RoleMiddleware(true, false, userService)
	//endregion middlewares

	//region routing
	app := fiber.New()

	//region api
	apiGroup := app.Group("/api")
	apiGroup.Use(loggingMiddleware)

	//region v1
	v1Group := apiGroup.Group("/v1")

	//region user v1
	userGroup := v1Group.Group("/user")
	userAdminGroup := userGroup.Group("/admin")
	userStaffGroup := userGroup.Group("/staff")

	userAdminGroup.Use(isAdminMiddleware)
	userStaffGroup.Use(isStaffMiddleware)

	userGroup.Post("/create", userServiceHandler.CreateUser)          //
	userGroup.Patch("/:id", userServiceHandler.UpdateUser)            //
	userGroup.Post("/token/refresh", userServiceHandler.RefreshToken) //
	userGroup.Post("/token/login", userServiceHandler.Login)

	userStaffGroup.Get("/:id", userServiceHandler.GetUser)                         // staff | admin
	userStaffGroup.Get("/get/by-tg-id/:tg_id", userServiceHandler.GetUserIdByTgId) // staff | admin
	userStaffGroup.Delete("/:id", userServiceHandler.DeleteUser)                   // staff | admin
	userStaffGroup.Get("/role/:id", userServiceHandler.GetRole)                    // staff | admin

	userAdminGroup.Post("/create", userServiceHandler.CreateUserAdmin)        // admin
	userAdminGroup.Patch("/update/:id", userServiceHandler.UpdateUserAdmin)   // admin
	userAdminGroup.Delete("/delete/:id", userServiceHandler.DeleteUserAdmin)  // admin
	userAdminGroup.Post("/get/by-token", userServiceHandler.GetUserIDByToken) // admin
	userAdminGroup.Post("/token/check", userServiceHandler.CheckToken)        // admin
	//endregion user v1

	//region shortener v1
	shortenerGroup := v1Group.Group("/s")

	shortenerCRUDGroup := shortenerGroup.Group("/crud")
	shortenerCRUDGroup.Use(authMiddleware)

	shortenerAdminGroup := shortenerCRUDGroup.Group("/admin")
	shortenerAdminGroup.Use(isAdminMiddleware)

	shortenerOwnerOnlyGroup := shortenerCRUDGroup.Group("/owner")
	shortenerOwnerOnlyGroup.Use(middlewares.ShortenerOwnerOnlyMiddleware("id", shortenerService))

	shortenerGroup.Get("/:shortLink", shortenerServiceHandler.Redirect)        //
	shortenerCRUDGroup.Post("/", shortenerServiceHandler.CreateNewLink)        // authenticated
	shortenerOwnerOnlyGroup.Put("/:id", shortenerServiceHandler.UpdateLink)    // owner
	shortenerOwnerOnlyGroup.Delete("/:id", shortenerServiceHandler.DeleteLink) // owner
	shortenerAdminGroup.Put("/:id", shortenerServiceHandler.UpdateLink)        // admin
	shortenerAdminGroup.Delete("/:id", shortenerServiceHandler.DeleteLink)     // admin
	//endregion shortener v1

	//region analytics v1
	analyticsGroup := v1Group.Group("/analytics")
	analyticsGroup.Use(authMiddleware)

	analyticsGroup.Get("/by-country", analyticsServiceHandler.GetClicksByCountry)
	analyticsGroup.Get("/by-region", analyticsServiceHandler.GetClicksByRegion)
	analyticsGroup.Get("/by-browser", analyticsServiceHandler.GetClicksByBrowser)
	analyticsGroup.Get("/by-os", analyticsServiceHandler.GetClicksByOS)
	analyticsGroup.Get("/by-device-type", analyticsServiceHandler.GetClicksByDeviceType)
	analyticsGroup.Get("/by-hour", analyticsServiceHandler.GetClicksByHour)
	analyticsGroup.Get("/by-date", analyticsServiceHandler.GetClicksByDate)
	analyticsGroup.Get("/by-referrer", analyticsServiceHandler.GetClicksByReferrer)
	//endregion analytics v1

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
