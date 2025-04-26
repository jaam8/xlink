package handler

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strings"
	"time"
)

func (h *Handler) ChooseMetricsToRenderHandler(ctx *th.Context, update telego.Update) error {
	userID := update.Message.From.ID

	shortLink := strings.TrimPrefix(update.CallbackQuery.Data, "show-metrics-")

	// Инициализируем выборы для юзера, если их ещё нет
	if _, ok := h.userMetricSelections[userID]; !ok {
		h.userMetricSelections[userID] = make(map[string]bool)
	}

	// Инициализируем выбор самой ссылки для юзера
	if _, ok := h.userShortLinkSelections[userID]; ok {
		h.userShortLinkSelections[userID] = shortLink
	}

	// Шлем стартовое сообщение с клавиатурой
	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      tu.ID(userID),
		Text:        "Выбери метрики для отображения:",
		ReplyMarkup: h.renderMetricKeyboard(userID),
	})
	return err
}

func (h *Handler) renderMetricKeyboard(userID int64) *telego.InlineKeyboardMarkup {
	h.mu.Lock()
	selections := h.userMetricSelections[userID]
	h.mu.Unlock()

	btn := func(metricKey, label string) telego.InlineKeyboardButton {
		selected := selections[metricKey]
		if selected {
			label = "✅ " + label
		}
		return tu.InlineKeyboardButton(label).WithCallbackData(metricKey)
	}

	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(btn("clicks-by-country", "Клики по странам")),
		tu.InlineKeyboardRow(btn("clicks-by-region", "Клики по регионам")),
		tu.InlineKeyboardRow(btn("clicks-by-browser", "Клики по браузерам")),
		tu.InlineKeyboardRow(btn("clicks-by-device_type", "Клики по устройствам")),
		tu.InlineKeyboardRow(btn("clicks-by-os", "Клики по ОС")),
		tu.InlineKeyboardRow(btn("clicks-by-referrers", "Клики по реферам")),
		tu.InlineKeyboardRow(btn("clicks-by-hour", "Клики по часам")),
		tu.InlineKeyboardRow(btn("clicks-by-date", "Клики по датам")),
		tu.InlineKeyboardRow(btn("clicks-done", "Готово")),
	)
}

func (h *Handler) HandleMetricSelection(ctx *th.Context, update telego.Update) error {
	userID := update.CallbackQuery.From.ID
	metric := update.CallbackQuery.Data

	if _, ok := h.userMetricSelections[userID]; !ok {
		h.mu.Lock()
		h.userMetricSelections[userID] = make(map[string]bool)
		h.mu.Unlock()
	}

	if metric == "clicks-done" {
		selected := h.userMetricSelections[userID]

		selectedStrings := make([]string, 0, len(selected))
		for metricName := range selected {
			selectedStrings = append(selectedStrings, metricName)
		}

		_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(userID),
			Text: fmt.Sprintf(
				"Ты выбрал: %v. Ввети букву d, дату начала и окончания выборки через пробел (как `d 2025-01-01 %s`)",
				selectedStrings,
				time.Now().Format(time.DateOnly),
			),
		})
		if err != nil {
			return err
		}

		// delete(h.userMetricSelections, userID)
		return nil
	}

	// Переключение выбранной метрики
	h.mu.Lock()
	selections := h.userMetricSelections[userID]
	h.mu.Unlock()

	selections[metric] = !selections[metric]

	// Обновляем клаву
	_, err := h.Bot.EditMessageReplyMarkup(ctx, &telego.EditMessageReplyMarkupParams{
		ChatID:      tu.ID(userID),
		MessageID:   update.CallbackQuery.Message.GetMessageID(),
		ReplyMarkup: h.renderMetricKeyboard(userID),
	})
	if err != nil {
		return err
	}

	// Обязательно ответить на коллбэк, чтобы крутилка исчезла
	_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	return nil
}

func (h *Handler) ChooseDateToRenderHandler(ctx *th.Context, update telego.Update) error {
	userID := update.CallbackQuery.From.ID
	dataSplit := strings.Split(update.Message.Text, " ")

	var err error

	//region dates
	var startDate, endDate time.Time
	startDate, err = time.Parse(time.DateOnly, dataSplit[1])
	if err != nil {
		_, _ = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(userID),
			Text:   fmt.Errorf("invalid startDate: %w", err).Error(),
		})
	}

	endDate, err = time.Parse(time.DateOnly, dataSplit[2])
	if err != nil {
		_, _ = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(userID),
			Text:   fmt.Errorf("invalid endDate: %w", err).Error(),
		})
	}
	//endregion

	//region mutex
	var selected map[string]bool
	var shortLink string
	var ok bool

	h.mu.Lock()

	selected, ok = h.userMetricSelections[userID]
	if !ok {
		h.mu.Unlock()
		_, _ = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(userID),
			Text:   "не выбраны метрики!",
		})
	}
	shortLink, ok = h.userShortLinkSelections[userID]
	if !ok {
		h.mu.Unlock()
		_, _ = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(userID),
			Text:   "не выбрана ссылка!",
		})
	}

	h.mu.Unlock()
	//endregion

	//region token
	token, err := h.user.GetTokenByTgID(userID)
	if err != nil {
		_, _ = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(userID),
			Text:   "Нет доступа",
		})
	}
	//endregion

	for clickByParamShortLinkString, value := range selected {
		if !value {
			continue
		}

		param := strings.Split(clickByParamShortLinkString, "-")[2]

		imageUrl := h.renderer.GetImageUrl(shortLink, token, param, startDate, endDate)
		h.SendImage(ctx, userID, imageUrl)
	}

	return nil
}
