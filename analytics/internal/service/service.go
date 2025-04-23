package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
	"xlink/analytics/internal/models"
	"xlink/analytics/internal/ports"
	"xlink/analytics/internal/service/helper"
	"xlink/common/gen/analytics"
	"xlink/common/logger"
)

type Service struct {
	analytics.AnalyticsServiceServer
	cacheAdapter     ports.CacheAdapter
	storageAdapter   ports.StorageAdapter
	consumerAdapter  ports.ConsumerAdapter
	shortenerAdapter ports.ShortenerAdapter
}

func New(
	cacheAdapter ports.CacheAdapter,
	storageAdapter ports.StorageAdapter,
	consumerAdapter ports.ConsumerAdapter,
	shortenerAdapter ports.ShortenerAdapter,
) *Service {
	return &Service{
		cacheAdapter:     cacheAdapter,
		storageAdapter:   storageAdapter,
		consumerAdapter:  consumerAdapter,
		shortenerAdapter: shortenerAdapter,
	}
}

func (s *Service) HandleConsumer(ctx context.Context, batchSize int, flushTimeout time.Duration) {
	ticker := time.NewTicker(flushTimeout)
	defer ticker.Stop()

	var batch []*models.Click

	flushBatch := func() {
		if len(batch) == 0 {
			return
		}
		if err := s.storageAdapter.SaveClicks(ctx, batch); err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx,
				"failed to save clicks batch to storage: %v",
				zap.Error(err))
		}
		batch = nil
	}

	for {
		select {
		case <-ctx.Done():
			flushBatch()
			logger.GetLoggerFromCtx(ctx).Info(ctx, "stop handling kafka consumer")
			return
		case <-ticker.C:
			flushBatch()
		default:
			event, err := s.consumerAdapter.ConsumeClickEvent(ctx)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx,
					"failed to consume click event: %v",
					zap.Error(err))
				continue
			}
			if event == nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx,
					"empty click event")
				continue
			}
			linkOwner, err := s.shortenerAdapter.GetLinkOwner(event.ShortLink)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx,
					"failed to get link owner",
					zap.String("short_link", event.ShortLink),
					zap.String("link_owner", linkOwner.String()),
					zap.Error(err))
				continue
			}

			hasToken, err := s.cacheAdapter.CheckVisitorToken(event.VisitorToken, event.ShortLink)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx,
					"failed to check visitor token: %v",
					zap.Error(err))
				continue
			}

			if !hasToken {
				if err = s.cacheAdapter.SetVisitorToken(event.VisitorToken, event.ShortLink); err != nil {
					logger.GetLoggerFromCtx(ctx).Error(ctx,
						"failed to set visitor token: %v",
						zap.Error(err))
					continue
				}
			}

			region, country, err := helper.ParseRegionAndCountry(event.IPAddress)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx,
					"failed to parse IP info",
					zap.Error(err))
				continue
			}

			click, err := helper.ClickEventToClick(event, linkOwner, country, region, hasToken)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx,
					"failed to convert click_event to click: %v",
					zap.Error(err))
				continue
			}

			batch = append(batch, click)

			if len(batch) >= batchSize {
				flushBatch()
			}
		}
	}
}

func (s *Service) ClicksByCountry(
	ctx context.Context, request *analytics.GetClicksRequest,
) (*analytics.ClicksByCountryResponse, error) {
	response := &analytics.ClicksByCountryResponse{}
	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"ClicksByCountry", zap.Any("request", request))
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"ClicksByCountry", zap.Any("startDate", startDate), zap.Any("endDate", endDate))

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByCountry(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by country: %v", zap.Error(err))
		return response, err
	}
	response.Data = helper.ToClicksByCountryResponse(rows)
	return response, nil
}

func (s *Service) ClicksByRegion(
	ctx context.Context, request *analytics.GetClicksRequest,
) (*analytics.ClicksByRegionResponse, error) {
	response := &analytics.ClicksByRegionResponse{}
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByRegion(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by region: %v", zap.Error(err))
		return response, err
	}
	response.Data = helper.ToClicksByRegionResponse(rows)
	return response, nil
}

func (s *Service) ClicksByBrowser(ctx context.Context, request *analytics.GetClicksRequest) (*analytics.ClicksByBrowserResponse, error) {
	response := &analytics.ClicksByBrowserResponse{}
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByBrowser(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by browser: %v", zap.Error(err))
		return response, err
	}
	response.Data = helper.ToClicksByBrowserResponse(rows)
	return response, nil
}

func (s *Service) ClicksByOS(ctx context.Context, request *analytics.GetClicksRequest) (*analytics.ClicksByOSResponse, error) {
	response := &analytics.ClicksByOSResponse{}
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByOS(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by os: %v", zap.Error(err))
		return response, err
	}
	response.Data = helper.ToClicksByOSResponse(rows)
	return response, nil
}

func (s *Service) ClicksByDeviceType(ctx context.Context, request *analytics.GetClicksRequest) (*analytics.ClicksByDeviceTypeResponse, error) {
	response := &analytics.ClicksByDeviceTypeResponse{}
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByDeviceType(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by device type: %v", zap.Error(err))
		return response, err
	}
	response.Data = helper.ToClicksByDeviceTypeResponse(rows)
	return response, nil
}

func (s *Service) ClicksByHour(ctx context.Context, request *analytics.GetClicksRequest) (*analytics.ClicksByHourResponse, error) {
	response := &analytics.ClicksByHourResponse{}
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByHour(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by hour: %v", zap.Error(err))
		return response, err
	}
	response.Stats = helper.ToClicksByHourResponse(rows)
	return response, nil
}

func (s *Service) ClicksByDate(ctx context.Context, request *analytics.GetClicksRequest) (*analytics.ClicksByDateResponse, error) {
	response := &analytics.ClicksByDateResponse{}
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByDate(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by date: %v", zap.Error(err))
		return response, err
	}
	response.Stats = helper.ToClicksByDateResponse(rows)
	return response, nil
}

func (s *Service) ClicksByReferrer(ctx context.Context, request *analytics.GetClicksRequest) (*analytics.ClicksByReferrerResponse, error) {
	response := &analytics.ClicksByReferrerResponse{}
	startDate, endDate, err := helper.ValidateRequestDates(ctx, request)
	if err != nil {
		return response, err
	}

	realLinkOwner, err := s.shortenerAdapter.GetLinkOwner(request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get link owner: %v", zap.Error(err))
		return response, err
	}

	reqLinkOwner, err := helper.ValidateRequestLinkOwner(ctx, request)
	if err != nil {
		return response, err
	}

	if realLinkOwner != reqLinkOwner {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"link owner mismatch: %v", zap.Error(err))
		return response, fmt.Errorf("link_owner mismatch exepted %v, got %v", realLinkOwner, reqLinkOwner)
	}

	rows, err := s.storageAdapter.GetClicksByReferrer(startDate, endDate, realLinkOwner, request.GetShortLink())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to get clicks by referrer: %v", zap.Error(err))
		return response, err
	}
	response.Data = helper.ToClicksByReferrerResponse(rows)
	return response, nil
}
