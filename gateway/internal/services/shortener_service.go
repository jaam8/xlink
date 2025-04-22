package services

import (
	"fmt"
	"time"
	"xlink/common/callers"
	"xlink/common/gen/shortener"
	"xlink/gateway/internal/ports"
)

type ShortenerService struct {
	ShortenerServiceRepo *ports.ShortenerServiceRepository
	MaxRetries           uint
	BaseDelay            time.Duration
}

func NewShortenerService(shortenerServiceRepo ports.ShortenerServiceRepository, maxRetries uint, baseDelay time.Duration) *ShortenerService {
	return &ShortenerService{
		ShortenerServiceRepo: &shortenerServiceRepo,
		MaxRetries:           maxRetries,
		BaseDelay:            baseDelay,
	}
}

func (s *ShortenerService) Redirect(request *shortener.RedirectRequest) (*shortener.RedirectResponse, error) {
	resultChan := make(chan *shortener.RedirectResponse)

	err := callers.Retry(func() error {
		response, err := (*s.ShortenerServiceRepo).Redirect(request)
		if err != nil {
			return fmt.Errorf("error in timeout Redirect caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call Redirect: %v", err)
	}

	response := <-resultChan

	return response, nil
}

func (s *ShortenerService) GetLink(request *shortener.GetLinkRequest) (*shortener.Link, error) {
	resultChan := make(chan *shortener.Link)

	err := callers.Retry(func() error {
		response, err := (*s.ShortenerServiceRepo).GetLink(request)
		if err != nil {
			return fmt.Errorf("error in timeout GetLink caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetLink: %v", err)
	}

	response := <-resultChan

	return response, nil
}

func (s *ShortenerService) CreateNewLink(request *shortener.CreateLinkRequest) (*shortener.Link, error) {
	resultChan := make(chan *shortener.Link)

	err := callers.Retry(func() error {
		response, err := (*s.ShortenerServiceRepo).CreateNewLink(request)
		if err != nil {
			return fmt.Errorf("error in timeout CreateNewLink caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call CreateNewLink: %v", err)
	}

	response := <-resultChan

	return response, nil
}

func (s *ShortenerService) UpdateLink(request *shortener.UpdateLinkRequest) (*shortener.Link, error) {
	resultChan := make(chan *shortener.Link)

	err := callers.Retry(func() error {
		response, err := (*s.ShortenerServiceRepo).UpdateLink(request)
		if err != nil {
			return fmt.Errorf("error in timeout UpdateLink caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call UpdateLink: %v", err)
	}

	response := <-resultChan

	return response, nil
}

func (s *ShortenerService) DeleteLink(request *shortener.DeleteLinkRequest) (*shortener.DeleteLinkResponse, error) {
	resultChan := make(chan *shortener.DeleteLinkResponse)

	err := callers.Retry(func() error {
		response, err := (*s.ShortenerServiceRepo).DeleteLink(request)
		if err != nil {
			return fmt.Errorf("error in timeout DeleteLink caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call DeleteLink: %v", err)
	}

	response := <-resultChan

	return response, nil
}

func (s *ShortenerService) GetLinksCountByUserId(request *shortener.GetLinksCountByUserIdRequest) (*shortener.GetLinksCountByUserIdResponse, error) {
	resultChan := make(chan *shortener.GetLinksCountByUserIdResponse)

	err := callers.Retry(func() error {
		response, err := (*s.ShortenerServiceRepo).GetLinksCountByUserId(request)
		if err != nil {
			return fmt.Errorf("error in timeout GetLinksCountByUserId caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetLinksCountByUserId: %v", err)
	}

	response := <-resultChan

	return response, nil
}
