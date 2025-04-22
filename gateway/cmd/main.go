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

	//region repos
	userServiceRepo := user_service_adapters.NewUserServiceRepositoryGRPC(usersGrpcPool)
	shortenerServiceRepo := shortener_service_adapters.NewShortenerServiceRepositoryGRPC(shortenerGrpcPool)

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
	//endregion repos

	//region handlers
	userServiceHandler := http_handlers.NewUserServiceHandler(userService)
	shortenerServiceHandler := http_handlers.NewShortenerServiceHandler(shortenerService, userService)
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
	userAdminGroup := userGroup.Group("/admin")
	userStaffGroup := userGroup.Group("/staff")

	userAdminGroup.Use(middlewares.RoleMiddleware(false, true, userService))
	userStaffGroup.Use(middlewares.RoleMiddleware(true, false, userService))

	userGroup.Post("/create", userServiceHandler.CreateUser)          //
	userGroup.Patch("/:id", userServiceHandler.UpdateUser)            //
	userGroup.Post("/token/refresh", userServiceHandler.RefreshToken) //

	userStaffGroup.Get("/:id", userServiceHandler.GetUser)                   // staff | admin
	userStaffGroup.Post("/get/by-tg-id", userServiceHandler.GetUserIdByTgId) // staff | admin
	userStaffGroup.Delete("/:id", userServiceHandler.DeleteUser)             // staff | admin
	userStaffGroup.Get("/role/:id", userServiceHandler.GetRole)              // staff | admin

	userAdminGroup.Post("/create", userServiceHandler.CreateUserAdmin)        // admin
	userAdminGroup.Patch("/update/:id", userServiceHandler.UpdateUserAdmin)   // admin
	userAdminGroup.Delete("/delete/:id", userServiceHandler.DeleteUserAdmin)  // admin
	userAdminGroup.Post("/get/by-token", userServiceHandler.GetUserIDByToken) // admin
	userAdminGroup.Post("/token/check", userServiceHandler.CheckToken)        // admin
	//endregion user v1

	//region shortener v1
	shortenerGroup := v1Group.Group("/s")

	shortenerCRUDGroup := shortenerGroup.Group("/crud")
	shortenerCRUDGroup.Use(middlewares.AuthMiddleware(userService))

	shortenerAdminGroup := shortenerCRUDGroup.Group("/admin")
	shortenerAdminGroup.Use(middlewares.RoleMiddleware(false, true, userService))

	shortenerOwnerOnlyGroup := shortenerCRUDGroup.Group("/owner")
	shortenerOwnerOnlyGroup.Use(middlewares.ShortenerOwnerOnlyMiddleware("id", shortenerService))

	shortenerGroup.Get("/:shortLink", shortenerServiceHandler.Redirect)        //
	shortenerCRUDGroup.Post("/", shortenerServiceHandler.CreateNewLink)        // authenticated
	shortenerOwnerOnlyGroup.Put("/:id", shortenerServiceHandler.UpdateLink)    // owner
	shortenerOwnerOnlyGroup.Delete("/:id", shortenerServiceHandler.DeleteLink) // owner
	shortenerAdminGroup.Put("/:id", shortenerServiceHandler.UpdateLink)        // admin
	shortenerAdminGroup.Delete("/:id", shortenerServiceHandler.DeleteLink)     // admin
	//endregion shortener v1

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
