package services

import (
	"fmt"
	"time"
	"xlink/common/callers"
	"xlink/common/gen/analytics"
	"xlink/gateway/internal/ports"
)

type AnalyticsService struct {
	AnalyticsServiceRepo *ports.AnalyticsServiceRepository
	MaxRetries           uint
	BaseDelay            time.Duration
}

func NewAnalyticsService(analyticsServiceRepo ports.AnalyticsServiceRepository, maxRetries uint, baseDelay time.Duration) *AnalyticsService {
	return &AnalyticsService{
		AnalyticsServiceRepo: &analyticsServiceRepo,
		MaxRetries:           maxRetries,
		BaseDelay:            baseDelay,
	}
}

func (s *AnalyticsService) ClicksByCountry(request *analytics.GetClicksRequest) (*analytics.ClicksByCountryResponse, error) {
	resultChan := make(chan *analytics.ClicksByCountryResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByCountry(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByCountry caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByCountry: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *AnalyticsService) ClicksByRegion(request *analytics.GetClicksRequest) (*analytics.ClicksByRegionResponse, error) {
	resultChan := make(chan *analytics.ClicksByRegionResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByRegion(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByRegion caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByRegion: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *AnalyticsService) ClicksByBrowser(request *analytics.GetClicksRequest) (*analytics.ClicksByBrowserResponse, error) {
	resultChan := make(chan *analytics.ClicksByBrowserResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByBrowser(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByBrowser caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByBrowser: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *AnalyticsService) ClicksByOS(request *analytics.GetClicksRequest) (*analytics.ClicksByOSResponse, error) {
	resultChan := make(chan *analytics.ClicksByOSResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByOS(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByOS caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByOS: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *AnalyticsService) ClicksByDeviceType(request *analytics.GetClicksRequest) (*analytics.ClicksByDeviceTypeResponse, error) {
	resultChan := make(chan *analytics.ClicksByDeviceTypeResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByDeviceType(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByDeviceType caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByDeviceType: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *AnalyticsService) ClicksByHour(request *analytics.GetClicksRequest) (*analytics.ClicksByHourResponse, error) {
	resultChan := make(chan *analytics.ClicksByHourResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByHour(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByHour caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByHour: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *AnalyticsService) ClicksByDate(request *analytics.GetClicksRequest) (*analytics.ClicksByDateResponse, error) {
	resultChan := make(chan *analytics.ClicksByDateResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByDate(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByDate caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByDate: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *AnalyticsService) ClicksByReferrer(request *analytics.GetClicksRequest) (*analytics.ClicksByReferrerResponse, error) {
	resultChan := make(chan *analytics.ClicksByReferrerResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.AnalyticsServiceRepo).ClicksByReferrer(request)
		if err != nil {
			return fmt.Errorf("error in retry ClicksByReferrer caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ClicksByReferrer: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}
