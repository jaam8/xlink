package services

import (
	"fmt"
	"time"
	"xlink/common/callers"
	"xlink/gateway/internal/ports"
)

type RendererService struct {
	rendererServiceRepo ports.RendererServiceRepository
	MaxRetries          uint
	BaseDelay           time.Duration
}

func NewRendererService(rendererServiceRepo ports.RendererServiceRepository, maxRetries uint, baseDelay time.Duration) *RendererService {
	return &RendererService{
		rendererServiceRepo: rendererServiceRepo,
		MaxRetries:          maxRetries,
		BaseDelay:           baseDelay,
	}
}

func (s *RendererService) Generate(shortLink string, param string, startDate time.Time, endDate time.Time, linkOwner string) ([]byte, error) {
	resultChan := make(chan []byte, 1)

	err := callers.Retry(func() error {
		response, err := s.rendererServiceRepo.Generate(shortLink, param, startDate, endDate, linkOwner)
		if err != nil {
			return fmt.Errorf("error in retry Renderer caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call HTTP image generate: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}
