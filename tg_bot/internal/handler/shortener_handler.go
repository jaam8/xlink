package handler

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
	"xlink/common/callers"
)

type CreateLinkData struct {
	TargetUrl string
	ShortLink *string
}

var createLinkData = &CreateLinkData{}

func (h *Handler) CreateLinkHandler(ctx *th.Context, firstUpdate telego.Update) error {
	chatID := tu.ID(firstUpdate.CallbackQuery.From.ID)

	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: chatID,
		Text:   "Напишите ссылку для сокращения",
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) ChooseLinkType(ctx *th.Context, update telego.Update) error {
	chatID := update.Message.From.ID
	createLinkData.TargetUrl = update.Message.Text
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Сделать кастомной").WithCallbackData("do-custom-link"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Сгенерировать").WithCallbackData("do-generate-link"),
		),
	)
	_, err := h.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      tu.ID(chatID),
		Text:        "Сгенерировать ссылку или сделать кастомной?\n Кастмоная ссылка должна быть на английском длинной до 10 символов",
		ReplyMarkup: inlineKeyboard,
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DoCustomLinkFinal(ctx *th.Context, update telego.Update) error {
	chatID := update.Message.From.ID
	text := update.Message.Text[2:]
	createLinkData.ShortLink = &text
	token, err := h.cache.GetUserToken(strconv.Itoa(int(chatID)))
	if err != nil || token == "" {
		token, err = h.user.GetTokenByTgID(chatID)
		if err != nil {
			return err
		}
		err = h.cache.SetUserToken(strconv.Itoa(int(chatID)), token)
		if err != nil {
			return err
		}
	}
	ShortLink, TargetURL, CreatedAt, ExpireAt, err := h.shortener.CreateLink(
		token, createLinkData.TargetUrl, createLinkData.ShortLink)
	if err != nil {
		err = callers.Retry(func() error {
			ShortLink, TargetURL, CreatedAt, ExpireAt, err = h.shortener.CreateLink(
				token, createLinkData.TargetUrl, createLinkData.ShortLink)
			if err != nil {
				return err
			}
			return nil
		}, h.maxRetries, h.baseRetryDelay)
		h.SendMessage(ctx, chatID, "Что то пошло не так, попробуйте еще раз\n (Докер на локалке в 90% случаев не тянет, ловит истекшие таймауты)")
		return err
	}
	_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	createLinkData.ShortLink = nil
	createLinkData.TargetUrl = ""
	resultLink := fmt.Sprintf("%s/l/%s", h.gatewayServerUrl, ShortLink)
	msg := fmt.Sprintf(
		`Короткая ссылка: <a href="%[1]s">%[1]s</a>
Целевой ресурс: %s
Создана: %s
Истекает: %s`,
		resultLink, TargetURL, CreatedAt, ExpireAt,
	)

	_, err = h.Bot.SendMessage(ctx,
		tu.Message(tu.ID(chatID), msg).
			WithParseMode(telego.ModeHTML),
	)

	if err != nil {
		return err
	}
	return nil
}

// callback = "do-custom-link"
func (h *Handler) DoCustomLink(ctx *th.Context, update telego.Update) error {
	h.SendMessage(ctx, update.CallbackQuery.From.ID, "Напишите кастомную ссылку в формате `l <link>`")
	_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	return nil
}

// callback = "do-generate-link"
func (h *Handler) DoGenerateLink(ctx *th.Context, update telego.Update) error {
	chatID := update.CallbackQuery.From.ID
	_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	var token string
	token, err := h.cache.GetUserToken(strconv.Itoa(int(chatID)))
	if err != nil || token == "" {
		token, err = h.user.GetTokenByTgID(chatID)
		if err != nil {
			return err
		}
		err = h.cache.SetUserToken(strconv.Itoa(int(chatID)), token)
		if err != nil {
			return err
		}
	}
	ShortLink, TargetURL, CreatedAt, ExpireAt, err := h.shortener.CreateLink(
		token, createLinkData.TargetUrl, createLinkData.ShortLink)
	if err != nil {
		return err
	}
	createLinkData.ShortLink = nil
	createLinkData.TargetUrl = ""
	resultLink := fmt.Sprintf("%s/l/%s", h.gatewayServerUrl, ShortLink)
	msg := fmt.Sprintf(
		`Короткая ссылка: <a href="%[1]s">%[1]s</a>
Целевой ресурс: %s
Создана: %s
Истекает: %s`,
		resultLink, TargetURL, CreatedAt, ExpireAt,
	)

	_, err = h.Bot.SendMessage(ctx,
		tu.Message(tu.ID(chatID), msg).
			WithParseMode(telego.ModeHTML),
	)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DeleteLinkHandler(ctx *th.Context, update telego.Update) error {
	chatID := tu.ID(update.CallbackQuery.From.ID)
	shortLink := strings.TrimPrefix(update.CallbackQuery.Data, "delete-link-")
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
	fmt.Println(shortLink)
	fmt.Println(token)
	err = h.shortener.DeleteLink(token, shortLink)
	if err != nil {
		return err
	}
	h.SendMessage(ctx, chatID.ID, "Ссылка удалена")
	_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	return nil
}

//func (h *Handler) UpdateLinkHandler(ctx *th.Context, update telego.Update) error {}

//func (h *Handler) GetUserLinksHandler(ctx *th.Context, update telego.Update) error {
//	chatID := tu.ID(update.CallbackQuery.From.ID)
//	token, err := h.cache.GetUserToken(strconv.Itoa(int(chatID.ID)))
//	if err != nil {
//		token, err = h.user.GetTokenByTgID(chatID.ID)
//		if err != nil {
//			return err
//		}
//		err = h.cache.SetUserToken(strconv.Itoa(int(chatID.ID)), token)
//		if err != nil {
//			return err
//		}
//	}
//
//	userLinks, err := h.shortener.GetUserLinks(token)
//	if err != nil {
//		return err
//	}
//
//	for _, link := range userLinks {
//		_, err = h.Bot.SendMessage(ctx,
//			tu.Messagef(tu.ID(chatID.ID),
//				"Короткая ссылка: %s/l/%s\nЦелевой ресурс: %s\nСоздана: %s\nИстекает: %s",
//				h.basePath, link.ShortLink, link.TargetUrl, link.CreatedAt, link.ExpireAt))
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//
//}
