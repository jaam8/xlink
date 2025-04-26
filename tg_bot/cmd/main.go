package main

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"time"
	"xlink/common/logger"
	"xlink/common/redis"
	"xlink/tg_bot/internal/config"
	"xlink/tg_bot/internal/handler"
	"xlink/tg_bot/internal/ports/adapters/analytics_adapter"
	"xlink/tg_bot/internal/ports/adapters/cache"
	"xlink/tg_bot/internal/ports/adapters/shortener_adapter"
	"xlink/tg_bot/internal/ports/adapters/user_service_adapter"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to load config", zap.Error(err))
	}

	bot, err := telego.NewBot(cfg.BotConfig.BotToken, telego.WithDefaultLogger(false, true))
	if err != nil {
		log.Fatal(err)
	}

	botCfg := cfg.BotConfig
	userServiceCfg := cfg.UserService
	redisCfg := cfg.RedisConfig

	redisClient, err := redis.NewRedisClient(ctx, redisCfg, botCfg.RedisDB)

	redisAdapter := cache.NewRedisAdapter(redisClient, time.Duration(botCfg.ExpirationSeconds)*time.Second)

	userAdapter := user_service_adapter.NewUserServiceAdapter(
		fmt.Sprintf("%s:%d", userServiceCfg.UpstreamNames, userServiceCfg.UpstreamPorts),
		fmt.Sprintf("%s/%s", botCfg.BaseAPIURL, "api/v1/user/"),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		time.Millisecond*time.Duration(userServiceCfg.Timeouts),
	)

	shortenerAdapter := shortener_adapter.NewShortenerAdapter(
		fmt.Sprintf("%s/%s", botCfg.BaseAPIURL, "api/v1/link/"),
		time.Millisecond*time.Duration(botCfg.Timeouts),
	)

	analyticsAdapter := analytics_adapter.NewAnalyticsAdapter(
		fmt.Sprintf("%s/%s", botCfg.BaseAPIURL, "api/v1/analytics/"),
		time.Millisecond*time.Duration(botCfg.Timeouts),
	)

	h := handler.NewHandler(userAdapter, shortenerAdapter, analyticsAdapter,
		redisAdapter, cfg.BotConfig.BaseAPIURL, bot)

	updates, _ := h.Bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		Timeout: 3,
	})
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}

	bh.Handle(h.StartHandler, th.CommandEqual("start"))
	bh.Handle(h.HelpHandler, th.CommandEqual("help"))
	bh.Handle(h.LoginHandler, th.CallbackDataEqual("log-in"))
	bh.Handle(h.MenuHandler, th.CommandEqual("menu"))
	bh.Handle(h.LinkMenuHandler, th.CallbackDataPrefix("link-"))
	bh.Handle(h.SignInHandler, th.CallbackDataEqual("sign-in"))
	bh.Handle(h.SetTgIDHandler, th.CallbackDataEqual("set-tg-id"))
	bh.Handle(h.HandleMetricSelection, th.CallbackDataPrefix("clicks-"))
	bh.Handle(h.ChooseMetricsToRenderHandler, th.CallbackDataEqual("show-metrics"))
	bh.Handle(h.DeleteLinkHandler, th.CallbackDataPrefix("delete-link"))
	bh.Handle(h.CreateLinkHandler, th.CallbackDataEqual("create-link"))
	bh.Handle(h.DoCustomLink, th.CallbackDataEqual("do-custom-link"))
	bh.Handle(h.DoGenerateLink, th.CallbackDataEqual("do-generate-link"))
	bh.Handle(h.ChooseLinkType)
	err = bh.Start()
	if err != nil {
		log.Fatal(err)
	}

	select {
	case <-ctx.Done():
		err = bh.Stop()
		if err != nil {
			log.Fatal(err)
		}
		return
	default:
	}
}
