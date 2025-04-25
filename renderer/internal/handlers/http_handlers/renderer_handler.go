package http_handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
	"xlink/common/gen/analytics"
	"xlink/common/logger"
	"xlink/renderer/internal/handlers/helpers"
	"xlink/renderer/internal/services"
	"xlink/renderer/internal/statistics_data"
)

type RendererHandler struct {
	analyticsService *services.AnalyticsService
	drawerGenerator  *services.DrawerService
}

func NewRendererHandler(analyticsService *services.AnalyticsService, drawerGenerator *services.DrawerService) *RendererHandler {
	return &RendererHandler{analyticsService: analyticsService, drawerGenerator: drawerGenerator}
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

	inputStats := make([]statistics_data.Stat, 0)

	switch param {
	case "country":
		var response *analytics.ClicksByCountryResponse
		response, err = h.analyticsService.ClicksByCountry(request)

		if err != nil {
			break
		}

		for _, stat := range response.Data {
			items := make([]statistics_data.Item, len(stat.Stats))
			for _, item := range stat.Stats {
				items = append(items, statistics_data.Item{
					Clicks:       item.Clicks,
					UniqueClicks: item.UniqueClicks,
					ParamValue:   item.Country,
				})
			}

			inputStats = append(inputStats, statistics_data.Stat{
				Date:  stat.Date,
				Items: items,
			})
		}
	case "region":
		var response *analytics.ClicksByRegionResponse
		response, err = h.analyticsService.ClicksByRegion(request)

		if err != nil {
			break
		}

		for _, stat := range response.Data {
			items := make([]statistics_data.Item, len(stat.Stats))
			for _, item := range stat.Stats {
				items = append(items, statistics_data.Item{
					Clicks:       item.Clicks,
					UniqueClicks: item.UniqueClicks,
					ParamValue:   item.Region,
				})
			}

			inputStats = append(inputStats, statistics_data.Stat{
				Date:  stat.Date,
				Items: items,
			})
		}
	case "browser":
		var response *analytics.ClicksByBrowserResponse
		response, err = h.analyticsService.ClicksByBrowser(request)

		if err != nil {
			break
		}

		for _, stat := range response.Data {
			items := make([]statistics_data.Item, len(stat.Stats))
			for _, item := range stat.Stats {
				items = append(items, statistics_data.Item{
					Clicks:       item.Clicks,
					UniqueClicks: item.UniqueClicks,
					ParamValue:   item.Browser,
				})
			}

			inputStats = append(inputStats, statistics_data.Stat{
				Date:  stat.Date,
				Items: items,
			})
		}
	case "os":
		var response *analytics.ClicksByOSResponse
		response, err = h.analyticsService.ClicksByOS(request)

		if err != nil {
			break
		}

		for _, stat := range response.Data {
			items := make([]statistics_data.Item, len(stat.Stats))
			for _, item := range stat.Stats {
				items = append(items, statistics_data.Item{
					Clicks:       item.Clicks,
					UniqueClicks: item.UniqueClicks,
					ParamValue:   item.Os,
				})
			}

			inputStats = append(inputStats, statistics_data.Stat{
				Date:  stat.Date,
				Items: items,
			})
		}
	case "device_type":
		var response *analytics.ClicksByDeviceTypeResponse
		response, err = h.analyticsService.ClicksByDeviceType(request)

		if err != nil {
			break
		}

		for _, stat := range response.Data {
			items := make([]statistics_data.Item, len(stat.Stats))
			for _, item := range stat.Stats {
				items = append(items, statistics_data.Item{
					Clicks:       item.Clicks,
					UniqueClicks: item.UniqueClicks,
					ParamValue:   item.DeviceType,
				})
			}

			inputStats = append(inputStats, statistics_data.Stat{
				Date:  stat.Date,
				Items: items,
			})
		}
	case "hour":
		var response *analytics.ClicksByHourResponse
		response, err = h.analyticsService.ClicksByHour(request)

		if err != nil {
			break
		}

		for _, stat := range response.Stats {
			items := make([]statistics_data.Item, len(stat.Stats))
			for _, item := range stat.Stats {
				items = append(items, statistics_data.Item{
					Clicks:       item.Clicks,
					UniqueClicks: item.UniqueClicks,
					ParamValue:   strconv.Itoa(int(item.Hour)),
				})
			}

			inputStats = append(inputStats, statistics_data.Stat{
				Date:  stat.Date,
				Items: items,
			})
		}
	case "date":
		var response *analytics.ClicksByDateResponse
		response, err = h.analyticsService.ClicksByDate(request)

		if err != nil {
			break
		}

		for _, stat := range response.Stats {
			items := []statistics_data.Item{
				{
					Clicks:       stat.Clicks,
					UniqueClicks: stat.UniqueClicks,
					ParamValue:   stat.Date,
				},
			}

			inputStats = append(inputStats, statistics_data.Stat{
				Date:  stat.Date,
				Items: items,
			})
		}
	}

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get analytics response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get analytics response"))
	}

	var imageBytes []byte
	imageBytes, err = h.drawerGenerator.Generate(statistics_data.StatisticsData{Stats: inputStats}, param)

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "generated image",
			zap.String("shortLink", shortLink),
			zap.String("linkOwner", linkOwner),
			zap.String("param", param),
			zap.Time("startDate", startDate),
			zap.Time("endDate", endDate))

	return ctx.Send(imageBytes)
}
