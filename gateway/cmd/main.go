package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	_ "github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"time"
	"xlink/common/prometheus"

	"xlink/common/grpc/pool"
	"xlink/common/logger"
	"xlink/gateway/internal/config"
	"xlink/gateway/internal/handlers/http_handlers"
	"xlink/gateway/internal/handlers/middlewares"
	"xlink/gateway/internal/ports/adapters/analytics_service_adapters"
	"xlink/gateway/internal/ports/adapters/renderer_service_adapters"
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
	rendererServiceRepo := renderer_service_adapters.NewRendererServiceRepositoryHTTP(
		"http",
		mainConfig.UpstreamNames.FileGenerator,
		mainConfig.UpstreamPorts.FileGenerator,
		time.Duration(mainConfig.Timeouts.FileGenerator)*time.Millisecond,
	)

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
	rendererService := services.NewRendererService(
		rendererServiceRepo,
		mainConfig.GrpcPool.MaxRetries,
		time.Millisecond*time.Duration(mainConfig.GrpcPool.BaseRetryDelayMilliseconds),
	)
	//endregion repos

	//region handlers
	userServiceHandler := http_handlers.NewUserServiceHandler(userService)
	shortenerServiceHandler := http_handlers.NewShortenerServiceHandler(shortenerService, userService, mainConfig.Gateway.UnknownRefererValue)
	analyticsServiceHandler := http_handlers.NewAnalyticsServiceHandler(analyticsService)
	rendererHandler := http_handlers.NewRendererServiceHandler(rendererService)
	//endregion handlers

	//region middlewares
	loggingMiddleware := middlewares.LoggerMiddleware()
	authMiddleware := middlewares.AuthMiddleware(userService)
	authMiddlewareTokenParam := middlewares.AuthMiddlewareTokenParam(userService, "token")
	isAdminMiddleware := middlewares.RoleMiddleware(false, true, userService)
	isStaffMiddleware := middlewares.RoleMiddleware(true, false, userService)
	shortLinkOwnerOnlyMiddleware := middlewares.ShortenerOwnerOnlyMiddleware("shortLink", shortenerService)
	//endregion middlewares

	//region html
	htmlEngine := html.New("./web/html", ".html")
	//endregion html

	prometheus.InitMetrics()

	//region routing
	app := fiber.New(fiber.Config{Views: htmlEngine})
	app.Use(loggingMiddleware)
	app.Use(middlewares.MetricsMiddlewareFiber())

	//region api
	apiGroup := app.Group("/api")

	//region v1
	v1Group := apiGroup.Group("/v1")

	//region user v1
	userGroup := v1Group.Group("/user")
	{
		userGroup.Post("/create", userServiceHandler.CreateUser)
		userGroup.Patch("/update/:id", userServiceHandler.UpdateUser)
		userGroup.Post("/refresh", userServiceHandler.RefreshToken)
		userGroup.Post("/login", userServiceHandler.Login)

		userAdminGroup := userGroup.Group("/admin")
		userAdminGroup.Use(isAdminMiddleware)
		{
			userAdminGroup.Post("/create", userServiceHandler.CreateUserAdmin)        // admin
			userAdminGroup.Patch("/update/:id", userServiceHandler.UpdateUserAdmin)   // admin
			userAdminGroup.Delete("/delete/:id", userServiceHandler.DeleteUserAdmin)  // admin
			userAdminGroup.Post("/get-by-token", userServiceHandler.GetUserIDByToken) // admin
			userAdminGroup.Post("/token-check", userServiceHandler.CheckToken)        // admin
		}

		userStaffGroup := userGroup.Group("/staff")
		userStaffGroup.Use(isStaffMiddleware)
		{
			userStaffGroup.Get("/get/:id", userServiceHandler.GetUser)                     // staff | admin
			userStaffGroup.Get("/get-by-tg-id/:tg_id", userServiceHandler.GetUserIdByTgId) // staff | admin
			userStaffGroup.Delete("/delete/:id", userServiceHandler.DeleteUser)            // staff | admin
			userStaffGroup.Get("/role/:id", userServiceHandler.GetRole)                    // staff | admin
		}

		userAuthedGroup := userGroup.Group("")
		userAuthedGroup.Use(authMiddleware)
		{
			userAuthedGroup.Get("/profile", userServiceHandler.Profile)
		}
	}
	//endregion user v1

	//region shortener v1
	shortenerGroup := v1Group.Group("/link")

	shortenerAuthenticatedGroup := shortenerGroup.Group("")
	shortenerAuthenticatedGroup.Use(authMiddleware)
	{
		shortenerAuthenticatedGroup.Get("/my-links", shortenerServiceHandler.MyLinks)      // authenticated
		shortenerAuthenticatedGroup.Post("/create", shortenerServiceHandler.CreateNewLink) // authenticated

		shortenerAdminGroup := shortenerAuthenticatedGroup.Group("/admin")
		shortenerAdminGroup.Use(isAdminMiddleware)
		{
			shortenerAdminGroup.Get("/links/:userId", shortenerServiceHandler.GetLinksByUserId)     // admin
			shortenerAdminGroup.Post("/create/:userId", shortenerServiceHandler.CreateNewLinkAdmin) // admin
			shortenerAdminGroup.Put("/update/:shortLink", shortenerServiceHandler.UpdateLink)       // admin
			shortenerAdminGroup.Delete("/delete/:shortLink", shortenerServiceHandler.DeleteLink)    // admin
		}

		shortenerOwnerOnlyGroup := shortenerAuthenticatedGroup.Group("")
		shortenerOwnerOnlyGroup.Use(shortLinkOwnerOnlyMiddleware)
		{
			shortenerOwnerOnlyGroup.Put("/update/:shortLink", shortenerServiceHandler.UpdateLink)    // owner
			shortenerOwnerOnlyGroup.Delete("/delete/:shortLink", shortenerServiceHandler.DeleteLink) // owner
		}
	}
	app.Get("/l/:shortLink", shortenerServiceHandler.Redirect) //
	//endregion shortener v1

	//region analytics v1
	analyticsGroup := v1Group.Group("/analytics")
	analyticsGroup.Use(authMiddleware)
	{
		analyticsGroup.Get("/by-country", analyticsServiceHandler.GetClicksByCountry)
		analyticsGroup.Get("/by-region", analyticsServiceHandler.GetClicksByRegion)
		analyticsGroup.Get("/by-browser", analyticsServiceHandler.GetClicksByBrowser)
		analyticsGroup.Get("/by-os", analyticsServiceHandler.GetClicksByOS)
		analyticsGroup.Get("/by-device-type", analyticsServiceHandler.GetClicksByDeviceType)
		analyticsGroup.Get("/by-hour", analyticsServiceHandler.GetClicksByHour)
		analyticsGroup.Get("/by-date", analyticsServiceHandler.GetClicksByDate)
		analyticsGroup.Get("/by-referrer", analyticsServiceHandler.GetClicksByReferrer)
	}
	//endregion analytics v1

	//region renderer v1
	rendererGroup := v1Group.Group("/img/:shortLink")
	rendererGroup.Use(authMiddlewareTokenParam, shortLinkOwnerOnlyMiddleware)
	{
		rendererGroup.Get("/", rendererHandler.Image)
	}
	for _, b := range app.GetRoutes() {
		fmt.Println(b)
	}
	//endregion renderer v1

	//region static
	app.Static("/static/", "./web/static")
	//endregion static

	//endregion v1
	//endregion api

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
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
