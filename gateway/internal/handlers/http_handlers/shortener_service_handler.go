package http_handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mileusna/useragent"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"xlink/common/gen/shortener"
	"xlink/common/gen/user_service"
	"xlink/common/logger"
	"xlink/gateway/internal/handlers"
	"xlink/gateway/internal/handlers/helpers"
	"xlink/gateway/internal/schemas"
	"xlink/gateway/internal/services"
)

type ShortenerServiceHandler struct {
	shortenerService    *services.ShortenerService
	userService         *services.UserService
	unknownRefererValue string
}

func NewShortenerServiceHandler(
	shortenerService *services.ShortenerService,
	userService *services.UserService,
	unknownRefererValue string,
) *ShortenerServiceHandler {
	return &ShortenerServiceHandler{
		shortenerService:    shortenerService,
		userService:         userService,
		unknownRefererValue: unknownRefererValue,
	}
}

func (h *ShortenerServiceHandler) getLinkIdByShortLinkParameter(ctx *fiber.Ctx) (string, error) {
	shortLink := ctx.Params("shortLink")
	if len(shortLink) == 0 {
		return "", helpers.BadRequest(ctx, errors.New("invalid shortLink: must be a non-empty string"))
	}

	responseLinkId, err := h.shortenerService.GetLinkIdByShortLink(
		&shortener.GetLinkIdByShortLinkRequest{ShortLink: shortLink})
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener response", zap.Error(err))
		return "", helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response"))
	}

	linkId := responseLinkId.LinkId
	return linkId, nil
}

func (h *ShortenerServiceHandler) Redirect(ctx *fiber.Ctx) error {
	shortLink := ctx.Params("shortLink")
	clickedAt := timestamppb.New(time.Now())

	referrer := ctx.Get("HTTP_REFERER")
	if referrer == "" {
		referrer = h.unknownRefererValue
	}

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
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener_service response", zap.Error(err))
		return helpers.NotFoundError(ctx, "couldn't get shortener_service response")
	}

	targetLink := response.TargetUrl

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "redirected link",
			zap.String("shortLink", shortLink),
			zap.String("targetLink", targetLink))

	return ctx.Render("redirect", fiber.Map{
		"TargetLink": targetLink,
	})
}

func (h *ShortenerServiceHandler) CreateNewLink(ctx *fiber.Ctx) error {
	var body schemas.CreateLinkSchema
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	userIdValue := ctx.Context().Value(handlers.UserIdKey)
	if userIdValue == nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "unauthorized (Use auth middleware before Owner only middleware!!!)"})
	}
	userId := userIdValue.(string)

	request := &shortener.CreateLinkRequest{
		UserId:    userId,
		ShortLink: body.ShortLink,
		TargetUrl: body.TargetUrl,
	}

	response, err := h.shortenerService.CreateNewLink(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "created link (administrator)", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusCreated).JSON(schemas.Link{
		LinkId:    response.LinkId,
		UserId:    response.UserId,
		ShortLink: response.ShortLink,
		TargetUrl: response.TargetUrl,
		CreatedAt: response.CreatedAt.AsTime().Format(time.RFC3339),
		ExpireAt:  response.ExpireAt.AsTime().Format(time.RFC3339),
	})
}

func (h *ShortenerServiceHandler) CreateNewLinkAdmin(ctx *fiber.Ctx) error {
	var body schemas.CreateLinkSchemaAdmin
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	if _, err := helpers.ParseUUID(body.UserId); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid user id: %v", err))
	}

	if _, err := h.userService.GetUser(&user_service.GetUserRequest{UserId: body.UserId}); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("couldn't find user with given id='%s': %v",
			body.UserId, err))
	}

	request := &shortener.CreateLinkRequest{
		UserId:    body.UserId,
		ShortLink: body.ShortLink,
		TargetUrl: body.TargetUrl,
	}

	response, err := h.shortenerService.CreateNewLink(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "created link (administrator)", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusCreated).JSON(schemas.Link{
		LinkId:    response.LinkId,
		UserId:    response.UserId,
		ShortLink: response.ShortLink,
		TargetUrl: response.TargetUrl,
		CreatedAt: response.CreatedAt.AsTime().Format(time.RFC3339),
		ExpireAt:  response.ExpireAt.AsTime().Format(time.RFC3339),
	})
}

func (h *ShortenerServiceHandler) UpdateLink(ctx *fiber.Ctx) error {
	linkId, err := h.getLinkIdByShortLinkParameter(ctx)
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("couldn't get link id by short link: %v", err))
	}

	var body schemas.UpdateLinkSchema
	if err = ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	var expireAt time.Time
	expireAt, err = helpers.ParseDateTime(body.ExpireAt)
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid expire_at: %v", err))
	}

	request := &shortener.UpdateLinkRequest{
		LinkId:     linkId,
		UserId:     ctx.Get(handlers.UserIdKey),
		Regenerate: body.Regenerate,
		ShortLink:  body.ShortLink,
		TargetUrl:  body.TargetUrl,
		ExpireAt:   timestamppb.New(expireAt),
	}

	var response *shortener.Link
	response, err = h.shortenerService.UpdateLink(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "updated link", zap.String("id", response.LinkId))
	return ctx.Status(fiber.StatusCreated).JSON(schemas.Link{
		LinkId:    response.LinkId,
		UserId:    response.UserId,
		ShortLink: response.ShortLink,
		TargetUrl: response.TargetUrl,
		CreatedAt: response.CreatedAt.AsTime().Format(time.RFC3339),
		ExpireAt:  response.ExpireAt.AsTime().Format(time.RFC3339),
	})
}

func (h *ShortenerServiceHandler) DeleteLink(ctx *fiber.Ctx) error {
	linkId, err := h.getLinkIdByShortLinkParameter(ctx)
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("couldn't get link id by short link: %v", err))
	}

	request := &shortener.DeleteLinkRequest{LinkId: linkId}

	var response *shortener.DeleteLinkResponse
	response, err = h.shortenerService.DeleteLink(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response"))
	}

	if !response.Status {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "link deletion was unsuccessful", zap.String("id", linkId))
		return helpers.BadRequest(ctx, fmt.Errorf("link deletion was unsuccessful: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "deleted link",
			zap.String("id", linkId),
			zap.Bool("status", response.Status))
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ShortenerServiceHandler) MyLinks(ctx *fiber.Ctx) error {
	userIdValue := ctx.Context().Value(handlers.UserIdKey)
	if userIdValue == nil {
		return helpers.NotAuthenticatedError(ctx, errors.New("unauthorized (Use auth middleware before Owner only middleware!!!)"))
	}
	userId := userIdValue.(string)

	response, err := h.shortenerService.GetLinks(&shortener.GetLinksRequest{UserId: userId})
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "got user's links", zap.String("id", userId))

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *ShortenerServiceHandler) GetLinksByUserId(ctx *fiber.Ctx) error {
	userId, err := helpers.ParseUUIDField(ctx, "userId")
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("couldn't parse user id: %v", err))
	}

	response, err := h.shortenerService.GetLinks(&shortener.GetLinksRequest{UserId: userId.String()})
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get shortener response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get shortener response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "got user's links", zap.String("id", userId.String()))

	return ctx.Status(fiber.StatusOK).JSON(response)
}
