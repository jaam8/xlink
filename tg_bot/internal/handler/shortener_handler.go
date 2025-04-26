package handler

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
)

type CreateLinkData struct {
	TargetUrl string
	ShortLink *string
}

var createLinkData = &CreateLinkData{}

//func (h *Handler) CreateLinkHandler(ctx *th.Context, firstUpdate telego.Update) error {
//	chatID := tu.ID(firstUpdate.CallbackQuery.From.ID)
//
//	var targetUrl string
//	var shortLink *string
//	var err error
//	var stop = false
//
//	_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
//		ChatID: chatID,
//		Text:   "Напишите ссылку для сокращения",
//	})
//	if err != nil {
//		return err
//	}
//	updateID := firstUpdate.UpdateID
//	for {
//		updates, err := h.Bot.GetUpdates(ctx, &telego.GetUpdatesParams{
//			Offset:  updateID + 1,
//			Timeout: 3,
//		})
//		if err != nil {
//			return err
//		}
//		log.Println(updateID)
//		log.Println(updates)
//		if stop {
//			log.Println("40 stop")
//			break
//		}
//		var update telego.Update
//		if len(updates) > 0 {
//			update = updates[0]
//		}
//		switch {
//		case update.Message != nil && strings.HasPrefix(update.Message.Text, "http"):
//			targetUrl = update.Message.Text
//			inlineKeyboard := tu.InlineKeyboard(
//				tu.InlineKeyboardRow(
//					tu.InlineKeyboardButton("Сделать кастомной").WithCallbackData("do-custom-link"),
//				),
//				tu.InlineKeyboardRow(
//					tu.InlineKeyboardButton("Сгенерировать").WithCallbackData("do-generate-link"),
//				),
//			)
//			_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
//				ChatID:      chatID,
//				Text:        `Сгенерировать ссылку или сделать кастомной?\n Кастмоная ссылка должна быть на английском длинной до 10 символов`,
//				ReplyMarkup: inlineKeyboard,
//			})
//			if err != nil {
//				return err
//			}
//			continue
//		case update.CallbackQuery != nil && update.CallbackQuery.Data == "do-custom-link":
//			_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
//				ChatID: chatID,
//				Text:   "Напишите кастомную ссылку",
//			})
//			if err != nil {
//				return err
//			}
//			_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
//				CallbackQueryID: update.CallbackQuery.ID,
//			})
//			continue
//		case update.CallbackQuery != nil && update.CallbackQuery.Data == "do-generate-link":
//			shortLink = nil
//			_ = h.Bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
//				CallbackQueryID: update.CallbackQuery.ID,
//			})
//			stop = true
//			break
//		case update.Message != nil:
//			if len(update.Message.Text) <= 10 && helper.IsValidShortLink(update.Message.Text) {
//				shortLink = &update.Message.Text
//				stop = true
//				break
//			} else {
//				_, err = h.Bot.SendMessage(ctx, &telego.SendMessageParams{
//					ChatID: chatID,
//					Text:   "Кастомная ссылка должна быть на английском и длинной до 10 символов",
//				})
//				if err != nil {
//					return err
//				}
//				continue
//			}
//		}
//	}
//
//	var token string
//	token, err = h.cache.GetUserToken(strconv.Itoa(int(chatID.ID)))
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
//	ShortLink, TargetURL, CreatedAt, ExpireAt, err := h.shortener.CreateLink(token, targetUrl, shortLink)
//	if err != nil {
//		return err
//	}
//	_, err = h.Bot.SendMessage(ctx,
//		tu.Messagef(tu.ID(chatID.ID),
//			"Короткая ссылка: %s/l/%s\nЦелевой ресурс: %s\nСоздана: %s\nИстекает: %s",
//			h.basePath, ShortLink, TargetURL, CreatedAt, ExpireAt))
//	if err != nil {
//		return err
//	}
//	return nil
//}

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
		Text:        `Сгенерировать ссылку или сделать кастомной?\n Кастмоная ссылка должна быть на английском длинной до 10 символов`,
		ReplyMarkup: inlineKeyboard,
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DoCustomLinkFinal(ctx *th.Context, update telego.Update) error {
	chatID := update.Message.From.ID
	createLinkData.ShortLink = nil
	var token string
	token, err := h.cache.GetUserToken(strconv.Itoa(int(chatID)))
	if err != nil {
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
	_, err = h.Bot.SendMessage(ctx,
		tu.Messagef(tu.ID(chatID),
			"Короткая ссылка: %s/l/%s\nЦелевой ресурс: %s\nСоздана: %s\nИстекает: %s",
			h.basePath, ShortLink, TargetURL, CreatedAt, ExpireAt))
	if err != nil {
		return err
	}
	return nil
}

// callback = "do-custom-link"
func (h *Handler) DoCustomLink(ctx *th.Context, update telego.Update) error {
	h.SendMessage(ctx, update.CallbackQuery.From.ID, "Напишите кастомную ссылку")
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
	if err != nil {
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
	_, err = h.Bot.SendMessage(ctx,
		tu.Messagef(tu.ID(chatID),
			"Короткая ссылка: %s/l/%s\nЦелевой ресурс: %s\nСоздана: %s\nИстекает: %s",
			h.basePath, ShortLink, TargetURL, CreatedAt, ExpireAt))
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

	err = h.shortener.DeleteLink(token, shortLink)
	if err != nil {
		return err
	}

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
