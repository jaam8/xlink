package http_handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mileusna/useragent"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"xlink/common/gen/shortener"
	"xlink/common/logger"
	"xlink/gateway/internal/handlers"
	"xlink/gateway/internal/handlers/helpers"
	"xlink/gateway/internal/schemas"
	"xlink/gateway/internal/services"
)

type ShortenerServiceHandler struct {
	shortenerService *services.ShortenerService
}

func NewShortenerServiceHandler(shortenerService *services.ShortenerService) *ShortenerServiceHandler {
	return &ShortenerServiceHandler{shortenerService: shortenerService}
}

func (h *ShortenerServiceHandler) Redirect(ctx *fiber.Ctx) error {
	shortLink := ctx.Params("shortLink")
	clickedAt := timestamppb.New(time.Now())
	referrer := ctx.Get("HTTP_REFERER")

	ipAddress := ctx.Get("X-Forwarded-For")
	if len(ipAddress) == 0 {
		ipAddress = ctx.Get("X-Real-IP")
		if len(ipAddress) == 0 {
			ipAddress = ctx.IP()
		}
	}

	visitorToken := ctx.Cookies("xlinkVisitor")

	userAgentText := ctx.Get("User-Agent")
	userAgent := useragent.Parse(userAgentText)
	browser := userAgent.Name
	deviceType := userAgent.Device
	userOs := userAgent.OS

	request := &shortener.RedirectRequest{
		ShortLink:    shortLink,
		ClickedAt:    clickedAt,
		Referrer:     referrer,
		IpAddress:    ipAddress,
		VisitorToken: visitorToken,
		UserAgent: &shortener.UserAgent{
			Browser:    browser,
			DeviceType: deviceType,
			Os:         userOs,
		},
	}

	response, err := h.shortenerService.Redirect(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.Context()).
			Error(ctx.Context(), "couldn't get shortener_service response", zap.Error(err))
		return helpers.NotFoundError(ctx, "couldn't get shortener_service response")
	}

	targetLink := response.TargetUrl

	logger.GetOrCreateLoggerFromCtx(ctx.Context()).
		Info(ctx.Context(), "redirected link",
			zap.String("shortLink", shortLink),
			zap.String("targetLink", targetLink))

	//TODO: return HTML that redirects itself
	return ctx.Redirect(targetLink, 200)
}

func (h *ShortenerServiceHandler) CreateNewLink(ctx *fiber.Ctx) error {
	var body schemas.CreateLinkSchema
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err).Error())
	}

	request := &shortener.CreateLinkRequest{
		UserId:    body.UserId,
		ShortLink: body.ShortLink,
		TargetUrl: body.TargetUrl,
	}

	response, err := h.shortenerService.CreateNewLink(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.Context()).
			Error(ctx.Context(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.Context()).
		Info(ctx.Context(), "created user", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *ShortenerServiceHandler) UpdateLink(ctx *fiber.Ctx) error {
	linkId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, "invalid link id: must be a valid uuid")
	}

	var body schemas.UpdateLinkSchema
	if err = ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err).Error())
	}

	var expireAt time.Time
	expireAt, err = helpers.ParseDateTime(body.ExpireAt)
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid expire_at: %v", err).Error())
	}

	request := &shortener.UpdateLinkRequest{
		LinkId:     linkId.String(),
		UserId:     ctx.Get(handlers.UserIdKey),
		Regenerate: body.Regenerate,
		ShortLink:  body.ShortLink,
		TargetUrl:  body.TargetUrl,
		ExpireAt:   timestamppb.New(expireAt),
	}

	var response *shortener.Link
	response, err = h.shortenerService.UpdateLink(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.Context()).
			Error(ctx.Context(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.Context()).
		Info(ctx.Context(), "updated link", zap.String("id", response.LinkId))
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *ShortenerServiceHandler) DeleteLink(ctx *fiber.Ctx) error {
	linkId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, "invalid link id: must be a valid uuid")
	}

	linkIdText := linkId.String()

	request := &shortener.DeleteLinkRequest{LinkId: linkIdText}

	var response *shortener.DeleteLinkResponse
	response, err = h.shortenerService.DeleteLink(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.Context()).
			Error(ctx.Context(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response"))
	}

	if !response.Status {
		logger.GetOrCreateLoggerFromCtx(ctx.Context()).
			Error(ctx.Context(), "link deletion was unsuccessful", zap.String("id", linkIdText))
		return helpers.BadRequest(ctx, "link deletion was unsuccessful")
	}

	logger.GetOrCreateLoggerFromCtx(ctx.Context()).
		Info(ctx.Context(), "deleted link",
			zap.String("id", linkIdText),
			zap.Bool("status", response.Status))
	return ctx.SendStatus(fiber.StatusNoContent)
}
