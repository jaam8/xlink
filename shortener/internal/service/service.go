package service

import (
	"context"
	"fmt"
	"shortener/internal/models"
	"shortener/internal/ports"
	"shortener/internal/service/helper"
	"shortener/pkg/api/shortener"
	"time"

	"github.com/chempik1234/common-chempik-pkg-golang/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	shortener.ShortenerServiceServer
	cachingRepo           ports.ShortenerCacheRepository
	storageRepo           ports.ShortenerStorageRepository
	senderRepo            ports.ShortenerSenderRepository
	defaultLinkExpireTime time.Duration
}

func New(
	cachingRepo ports.ShortenerCacheRepository,
	storageRepo ports.ShortenerStorageRepository,
	senderRepo ports.ShortenerSenderRepository,
	defaultLinkExpireTime time.Duration,
) *Service {
	return &Service{
		cachingRepo:           cachingRepo,
		storageRepo:           storageRepo,
		senderRepo:            senderRepo,
		defaultLinkExpireTime: defaultLinkExpireTime,
	}
}

func (s *Service) CreateNewLink(ctx context.Context, request *shortener.CreateLinkRequest) (*shortener.Link, error) {
	inputData, err := helper.LinkModelFromLinkRequest(request, time.Now().Add(s.defaultLinkExpireTime))
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while validating link: %v", err)
	}

	newLink, err := s.storageRepo.CreateLink(inputData)

	return helper.LinkResponseFromLinkModel(newLink), nil
}

func (s *Service) UpdateLink(ctx context.Context, request *shortener.UpdateLinkRequest) (*shortener.Link, error) {
	inputData, err := helper.LinkModelFromLinkRequestWithId(request, time.Now().Add(s.defaultLinkExpireTime))
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while validating link: %v", err)
	}

	newLink, err := s.storageRepo.UpdateLink(inputData)

	return helper.LinkResponseFromLinkModel(newLink), nil
}

func (s *Service) DeleteLink(ctx context.Context, request *shortener.DeleteLinkRequest) (*shortener.DeleteLinkResponse, error) {
	id, err := helper.GetValidatedId(request)
	if err != nil {
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while getting id: %v", err)
	}

	var link models.Link
	link, err = s.storageRepo.GetLinkById(id)
	if err != nil {
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while getting link: %v", err)
	}
	cacheKeyToDelete := link.ShortLink

	err = s.storageRepo.DeleteLink(id)
	if err != nil {
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while deleting link from storage: %v", err)
	}

	err = s.cachingRepo.DeleteUrl(cacheKeyToDelete)
	if err != nil {
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while deleting link from caching: %v", err)
	}

	return &shortener.DeleteLinkResponse{Status: true}, nil
}

func (s *Service) Redirect(ctx context.Context, request *shortener.Url) (*shortener.Url, error) {
	var originalUrl string
	var shortUrl = request.Url
	var err error

	originalUrl, err = s.cachingRepo.GetUrl(shortUrl)
	if err != nil {
		// if it's not in cache, then we get in from relational DB
		var link models.Link
		link, err = s.storageRepo.GetLinkByShortUrl(shortUrl)
		if err != nil {
			return &shortener.Url{}, fmt.Errorf("error while getting link: %v", err)
		}

		originalUrl = link.Url

		// we better cache the link, so we won't have to visit DB too often
		go func() {
			cacheErr := s.cachingRepo.SetUrl(shortUrl, originalUrl)
			if cacheErr != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't cache link", zap.String("key", shortUrl), zap.String("value", originalUrl))
			}
		}()
	}

	go func() {
		s.senderRepo.SendRedirectInfo()
	}()

	return &shortener.Url{Url: originalUrl}, nil
}
