package handler

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
	"xlink/tg_bot/internal/ports"
)

type Handler struct {
	user      ports.UserServiceAdapter
	shortener ports.ShortenerAdapter
	analytics ports.AnalyticsAdapter
	cache     ports.CacheAdapter
	// key = tgID, value = map[metrics]bool
	// сделать потокобезопасным
	userMetricSelections map[int64]map[string]bool
	basePath             string
	Bot                  *telego.Bot
}

func NewHandler(user ports.UserServiceAdapter, shortener ports.ShortenerAdapter,
	analytics ports.AnalyticsAdapter, cache ports.CacheAdapter,
	basePath string, bot *telego.Bot) *Handler {
	var handler Handler
	handler.user = user
	handler.shortener = shortener
	handler.analytics = analytics
	handler.userMetricSelections = make(map[int64]map[string]bool)
	handler.cache = cache
	handler.basePath = basePath
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
- /start — запустить бота
- /help — показать это сообщение
- /menu — выбрать метрики для отображения`})
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

func (h *Handler) MenuHandler(ctx *th.Context, update telego.Update) error {
	chatID := tu.ID(update.Message.From.ID)
	token, err := h.cache.GetUserToken(strconv.Itoa(int(chatID.ID)))
	if err != nil {
		token, err = h.user.GetTokenByTgID(chatID.ID)
		if err != nil {
			return err
		}
		err = h.cache.SetUserToken(strconv.Itoa(int(chatID.ID)), token)
		if err != nil {
			return err
		}
	}

	links, err := h.shortener.GetUserLinks(token)
	if err != nil {
		h.SendMessage(ctx, chatID.ID, "Что то пошло не так"+err.Error())
		return err
	}

	keyboard := [][]telego.InlineKeyboardButton{
		[]telego.InlineKeyboardButton{
			tu.InlineKeyboardButton("Создать ссылку").WithCallbackData("create-link"),
		},
	}

	for _, link := range links {
		keyboard = append(keyboard, []telego.InlineKeyboardButton{
			tu.InlineKeyboardButton(fmt.Sprintf("%s/l/%s", h.basePath, link)).
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
