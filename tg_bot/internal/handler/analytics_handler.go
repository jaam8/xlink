package handler

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (h *Handler) ChooseMetricsToRenderHandler(ctx *th.Context, update telego.Update) error {
	userID := update.Message.From.ID

	// Инициализируем выборы для юзера, если их ещё нет
	if _, ok := h.userMetricSelections[userID]; !ok {
		h.userMetricSelections[userID] = make(map[string]bool)
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
	selections := h.userMetricSelections[userID]

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
		tu.InlineKeyboardRow(btn("clicks-by-device", "Клики по устройствам")),
		tu.InlineKeyboardRow(btn("clicks-by-os", "Клики по ОС")),
		tu.InlineKeyboardRow(btn("clicks-by-referrer", "Клики по реферам")),
		tu.InlineKeyboardRow(btn("clicks-by-hour", "Клики по часам")),
		tu.InlineKeyboardRow(btn("clicks-by-date", "Клики по датам")),
		tu.InlineKeyboardRow(btn("clicks-done", "Готово")),
	)
}

func (h *Handler) HandleMetricSelection(ctx *th.Context, update telego.Update) error {
	userID := update.CallbackQuery.From.ID
	metric := update.CallbackQuery.Data

	if _, ok := h.userMetricSelections[userID]; !ok {
		h.userMetricSelections[userID] = make(map[string]bool)
	}

	if metric == "clicks-done" {
		selected := h.userMetricSelections[userID]
		// Здесь ты обрабатываешь выбранные метрики и чистишь за собой

		_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(userID),
			Text:   fmt.Sprintf("Ты выбрал: %v", selected),
		})
		if err != nil {
			return err
		}
		delete(h.userMetricSelections, userID)
		return nil
	}

	// Переключение выбранной метрики
	selections := h.userMetricSelections[userID]
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
