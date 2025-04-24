package handler

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (h *Handler) LoginHandler(ctx *th.Context, update telego.Update) error {
	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(update.CallbackQuery.From.ID),
		Text:   "Введите ваш API ключ",
	})
	apiKey := update.Message.Text

	if err != nil {
		return err
	}
	// send user_id to redis
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
			ChatID:      tu.ID(update.CallbackQuery.From.ID),
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
				tu.InlineKeyboardButton("Попробовать еще раз").WithCallbackData("log-in"),
			),
		)
		_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:      tu.ID(update.CallbackQuery.From.ID),
			Text:        "Кажется, вы пытаетесь зайти в чужой аккаунт",
			ReplyMarkup: inlineKeyboard,
		})
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (h *Handler) SignInHandler(ctx *th.Context, update telego.Update) error {
	// send user_id to redis
	_, token, err := h.user.CreateUser(&update.CallbackQuery.From.ID)
	if err != nil {
		return err
	}
	_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(update.CallbackQuery.From.ID),
		Text:   fmt.Sprintf("Ваш API ключ: `%s`\nСохраните его и никому не показывайте", token),
	})
	if err != nil {
		return err
	}
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
