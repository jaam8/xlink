package handler

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
)

func (h *Handler) LoginHandler(ctx *th.Context, firstUpdate telego.Update) error {
	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(firstUpdate.CallbackQuery.From.ID),
		Text:   "Введите ваш API ключ",
	})
	if err != nil {
		return err
	}

	// wait for user input
	for {
		updates, err := h.Bot.GetUpdates(ctx, &telego.GetUpdatesParams{
			Offset: firstUpdate.UpdateID + 1,
		})
		if err != nil {
			return err
		}

		for _, update := range updates {
			if update.Message != nil {
				apiKey := update.Message.Text

				_, tgID, err := h.user.LoginUser(apiKey)
				if err != nil {
					return err
				}
				if tgID == nil {
					inlineKeyboard := tu.InlineKeyboard(
						tu.InlineKeyboardRow(
							tu.InlineKeyboardButton("Привязать тг аккаунт").WithCallbackData("set-tg-id"),
						),
					)
					_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
						ChatID:      tu.ID(update.Message.From.ID),
						Text:        "К вашему профилю не привязан тг аккаунт",
						ReplyMarkup: inlineKeyboard,
					})
					if err != nil {
						return err
					}
					return nil
				}
				if *tgID != update.Message.From.ID {
					inlineKeyboard := tu.InlineKeyboard(
						tu.InlineKeyboardRow(
							tu.InlineKeyboardButton("Зарегестироваться").WithCallbackData("sign-in"),
						),
						tu.InlineKeyboardRow(
							tu.InlineKeyboardButton("Попробовать другой токен").WithCallbackData("log-in"),
						),
					)
					_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
						ChatID:      tu.ID(update.Message.From.ID),
						Text:        "Кажется, вы пытаетесь зайти в чужой аккаунт",
						ReplyMarkup: inlineKeyboard,
					})
					if err != nil {
						return err
					}
					return nil
				}
				err = h.cache.SetUserToken(strconv.Itoa(int(update.Message.From.ID)), apiKey)
				if err != nil {
					return err
				}
				_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
					ChatID: tu.ID(update.Message.From.ID),
					Text:   "Вы успешно авторизовались",
				})
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
}

func (h *Handler) SignInHandler(ctx *th.Context, update telego.Update) error {
	_, token, err := h.user.CreateUser(&update.CallbackQuery.From.ID)
	if err != nil {
		return err
	}
	err = h.cache.SetUserToken(strconv.Itoa(int(update.CallbackQuery.From.ID)), token)
	_, err = h.Bot.SendMessage(ctx,
		tu.Messagef(tu.ID(update.CallbackQuery.From.ID),
			"Ваш API ключ: `%s`\nСохраните его и никому не показывайте", token,
		).WithParseMode(telego.ModeMarkdownV2),
	)
	if err != nil {
		return err
	}
	_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	return nil

}

func (h *Handler) SetTgIDHandler(ctx *th.Context, update telego.Update) error {
	userID := "something get"
	err := h.user.SetTgID(userID, update.CallbackQuery.From.ID)
	if err != nil {
		_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(update.CallbackQuery.From.ID),
			Text:   "Ошибка привязки аккаунта, попробуйте позже",
		})
		if err != nil {
			return err
		}
		return nil
	}
	_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(update.CallbackQuery.From.ID),
		Text:   "Ваш аккаунт успешно привязан",
	})
	return nil
}
