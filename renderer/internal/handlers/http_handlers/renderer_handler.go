package http_handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"xlink/common/gen/analytics"
	"xlink/common/logger"
	"xlink/renderer/internal/handlers/helpers"
	"xlink/renderer/internal/services"
)

type RendererHandler struct {
	analyticsService *services.AnalyticsService
}

func NewRendererHandler(analyticsService *services.AnalyticsService) *RendererHandler {
	return &RendererHandler{analyticsService: analyticsService}
}

func (h *RendererHandler) Image(ctx *fiber.Ctx) error {
	var err error

	var shortLink string
	shortLink, err = helpers.ParseNotEmptyStringField(ctx, "short_link")
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid query param 'short_link': %w", err))
	}

	var startDate time.Time
	startDate, err = helpers.ParseDateField(ctx, "start_date")
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid query param 'start_date': %w", err))
	}

	var endDate time.Time
	endDate, err = helpers.ParseDateField(ctx, "end_date")
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid query param 'end_date': %w", err))
	}

	var linkOwner string
	linkOwner, err = helpers.ParseNotEmptyStringField(ctx, "link_owner")
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid query param 'link_owner': %w", err))
	}

	var param string
	param, err = helpers.ParseNotEmptyStringField(ctx, "param")
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid query param 'param': %w", err))
	}

	request := &analytics.GetClicksRequest{
		ShortLink: shortLink,
		LinkOwner: linkOwner,
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	response, err := h.analyticsService.ClicksByHour(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get analytics response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get analytics response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "generated image",
			zap.String("shortLink", shortLink),
			zap.String("linkOwner", linkOwner),
			zap.String("param", param),
			zap.Time("startDate", startDate),
			zap.Time("endDate", endDate))
	return ctx.Status(fiber.StatusCreated).JSON(response)
}
