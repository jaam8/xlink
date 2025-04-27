package handler

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
	"sync"
	"time"
	"xlink/common/callers"
	"xlink/tg_bot/internal/ports"
)

type Handler struct {
	user      ports.UserServiceAdapter
	shortener ports.ShortenerAdapter
	analytics ports.AnalyticsAdapter
	cache     ports.CacheAdapter
	renderer  ports.RendererAdapter
	mu        sync.Mutex
	// key = tgID, value = map[metrics]bool
	userMetricSelections    map[int64]map[string]bool
	userShortLinkSelections map[int64]string
	basePath                string
	gatewayServerUrl        string
	baseRetryDelay          time.Duration
	maxRetries              uint
	Bot                     *telego.Bot
}

func NewHandler(
	user ports.UserServiceAdapter,
	shortener ports.ShortenerAdapter,
	analytics ports.AnalyticsAdapter,
	cache ports.CacheAdapter,
	renderer ports.RendererAdapter,
	basePath string,
	gatewayServerUrl string,
	baseRetryDelay time.Duration,
	maxRetries uint,
	bot *telego.Bot,
) *Handler {
	var handler Handler
	handler.user = user
	handler.shortener = shortener
	handler.analytics = analytics
	handler.userMetricSelections = make(map[int64]map[string]bool)
	handler.userShortLinkSelections = make(map[int64]string)
	handler.cache = cache
	handler.renderer = renderer
	handler.basePath = basePath
	handler.gatewayServerUrl = gatewayServerUrl
	handler.baseRetryDelay = baseRetryDelay
	handler.maxRetries = maxRetries
	handler.Bot = bot
	return &handler
}

func (h *Handler) StartHandler(ctx *th.Context, update telego.Update) error {
	chatID := tu.ID(update.Message.Chat.ID)

	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Зарегестироваться").WithCallbackData("sign-in"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Авторизоваться").WithCallbackData("log-in"),
		),
	)

	// Отправка сообщения с инлайн-кнопками
	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      chatID,
		Text:        "Выберите кнопку:",
		ReplyMarkup: inlineKeyboard})

	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) HelpHandler(ctx *th.Context, update telego.Update) error {
	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(update.Message.Chat.ID),
		Text: `Вот, что я могу:
/start - регистрация и авторизация
/help - список команд
/menu - управление ссылками`})
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) SendMessage(ctx *th.Context, chatID int64, text string) {
	_, _ = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(chatID),
		Text:   text,
	})
}

func (h *Handler) SendImage(ctx *th.Context, chatID int64, url string) {
	_, _ = h.Bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID: tu.ID(chatID),
		Photo:  telego.InputFile{URL: url},
	})
}

func (h *Handler) MenuHandler(ctx *th.Context, update telego.Update) error {
	chatID := tu.ID(update.Message.From.ID)
	token, err := h.cache.GetUserToken(strconv.Itoa(int(chatID.ID)))
	if err != nil || token == "" {
		token, err = h.user.GetTokenByTgID(chatID.ID)
		if err != nil {
			h.SendMessage(ctx, chatID.ID, "для начала авторизуйтесь или зарегестрируйтесь")
			return err
		}
	}

	links, err := h.shortener.GetUserLinks(token)
	if err != nil {
		err = callers.Retry(func() error {
			links, err = h.shortener.GetUserLinks(token)
			if err != nil {
				return err
			}
			return nil
		}, h.maxRetries, h.baseRetryDelay)
		if err != nil {
			h.SendMessage(ctx, chatID.ID, "Что то пошло не так, попробуйте еще раз\n (Докер на локалке в 90% случаев не тянет, ловит истекшие таймауты)")
			return err
		}
	}

	keyboard := [][]telego.InlineKeyboardButton{
		[]telego.InlineKeyboardButton{
			tu.InlineKeyboardButton("Создать ссылку").WithCallbackData("create-link"),
		},
	}

	for _, link := range links {
		keyboard = append(keyboard, []telego.InlineKeyboardButton{
			tu.InlineKeyboardButton(fmt.Sprintf("%s/l/%s", "localhost", link)).
				WithCallbackData("link-" + link)},
		)
	}

	inlineKeyboard := tu.InlineKeyboard(keyboard...)
	_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      chatID,
		Text:        "Выберите действие:",
		ReplyMarkup: inlineKeyboard})

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) LinkMenuHandler(ctx *th.Context, update telego.Update) error {
	chatID := tu.ID(update.CallbackQuery.From.ID)
	shortLink := strings.TrimPrefix(update.CallbackQuery.Data, "link-")

	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Обновить ссылку").WithCallbackData("update-link-"+shortLink),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Удалить ссылку").WithCallbackData("delete-link-"+shortLink),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Показать метрики").WithCallbackData("show-metrics-"+shortLink),
		),
	)
	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      chatID,
		Text:        "Выберите действие:",
		ReplyMarkup: inlineKeyboard,
	})
	if err != nil {
		return err
	}
	return nil
}
