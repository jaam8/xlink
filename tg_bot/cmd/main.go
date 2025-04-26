package main

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"time"
	"xlink/tg_bot/internal/config"
	"xlink/tg_bot/internal/handler"
	"xlink/tg_bot/internal/ports/adapters/analytics_adapter"
	"xlink/tg_bot/internal/ports/adapters/shortener_adapter"
	"xlink/tg_bot/internal/ports/adapters/user_service_adapter"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telego.NewBot(cfg.BotConfig.BotToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	botCfg := cfg.BotConfig
	userServiceCfg := cfg.UserService
	fmt.Println(botCfg)
	fmt.Println(userServiceCfg)

	userAdapter := user_service_adapter.NewUserServiceAdapter(
		fmt.Sprintf("%s:%d", userServiceCfg.UpstreamNames, userServiceCfg.UpstreamPorts),
		//fmt.Sprintf("%s/%s", botCfg.BaseAPIURL, "api/v1/user"),
		fmt.Sprintf("%s/%s", "http://nginx:80", "api/v1/user/"),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		time.Millisecond*time.Duration(userServiceCfg.Timeouts),
	)

	shortenerAdapter := shortener_adapter.NewShortenerAdapter(
		fmt.Sprintf("%s/%s", botCfg.BaseAPIURL, "api/v1/link"),
		time.Millisecond*time.Duration(botCfg.Timeouts),
	)

	analyticsAdapter := analytics_adapter.NewAnalyticsAdapter(
		fmt.Sprintf("%s/%s", botCfg.BaseAPIURL, "api/v1/analytics"),
		time.Millisecond*time.Duration(botCfg.Timeouts),
	)

	h := handler.NewHandler(userAdapter, shortenerAdapter, analyticsAdapter, bot)
	updates, _ := h.Bot.UpdatesViaLongPolling(ctx, nil)
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}

	bh.Handle(h.StartHandler, th.CommandEqual("start"))
	bh.Handle(h.HelpHandler, th.CommandEqual("help"))
	bh.Handle(h.LoginHandler, th.CallbackDataEqual("log-in"))
	bh.Handle(h.SignInHandler, th.CallbackDataEqual("sign-in"))
	bh.Handle(h.SetTgIDHandler, th.CallbackDataEqual("set-tg-id"))
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
