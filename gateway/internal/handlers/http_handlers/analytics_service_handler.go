package http_handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"xlink/common/gen/analytics"
	"xlink/gateway/internal/handlers"
	"xlink/gateway/internal/handlers/helpers"
	"xlink/gateway/internal/services"
)

type AnalyticsServiceHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsServiceHandler(analyticsService *services.AnalyticsService) *AnalyticsServiceHandler {
	return &AnalyticsServiceHandler{analyticsService: analyticsService}
}

func getClicksRequestFromRequest(ctx *fiber.Ctx) (*analytics.GetClicksRequest, error) {
	linkOwner := ctx.Get(handlers.UserIdKey)
	if linkOwner == "" {
		return nil, helpers.NotAuthenticatedError(ctx, errors.New("not authenticated (must use auth middleware)"))
	}

	shortLink := ctx.Params("shortLink")

	startDate, err := helpers.ParseDateField(ctx, "start_date")
	if err != nil {
		return nil, helpers.InvalidDateBadRequest(ctx, "end_date")
	}

	endDate, err := helpers.ParseDateField(ctx, "end_date")
	if err != nil {
		return nil, helpers.InvalidDateBadRequest(ctx, "end_date")
	}

	request := &analytics.GetClicksRequest{
		LinkOwner: linkOwner,
		ShortLink: shortLink,
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	return request, nil
}

func (h *AnalyticsServiceHandler) GetClicksByCountry(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByCountry(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get country stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *AnalyticsServiceHandler) GetClicksByRegion(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByRegion(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get region stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *AnalyticsServiceHandler) GetClicksByBrowser(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByBrowser(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get browser stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *AnalyticsServiceHandler) GetClicksByOS(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByOS(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get OS stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *AnalyticsServiceHandler) GetClicksByDeviceType(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByDeviceType(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get device type stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *AnalyticsServiceHandler) GetClicksByHour(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByHour(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get per-hour stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *AnalyticsServiceHandler) GetClicksByDate(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByDate(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get per-date stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *AnalyticsServiceHandler) GetClicksByReferrer(ctx *fiber.Ctx) error {
	request, err := getClicksRequestFromRequest(ctx)
	if err != nil {
		return err //nolint:all
	}

	response, err := h.analyticsService.ClicksByReferrer(request)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get referrer stats: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
