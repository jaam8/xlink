package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
	"xlink/common/gen/shortener"
	"xlink/common/logger"
	"xlink/shortener/internal/models"
	"xlink/shortener/internal/ports"
	"xlink/shortener/internal/service/helper"
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

func (s *Service) GetLink(ctx context.Context, request *shortener.GetLinkRequest) (*shortener.Link, error) {
	var err error
	var linkId uuid.UUID
	linkId, err = helper.GetValidatedId(request)
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while validating link_id: %w", err)
	}

	var link models.Link
	link, err = s.storageRepo.GetLinkById(linkId)
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while getting link by id: %w", err)
	}

	return helper.LinkResponseFromLinkModel(link), nil
}

func (s *Service) CreateNewLink(ctx context.Context, request *shortener.CreateLinkRequest) (*shortener.Link, error) {
	var err error
	var inputData *models.Link
	inputData, err = helper.LinkModelFromLinkCreateRequest(request, time.Now().Add(s.defaultLinkExpireTime))
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while validating link: %v", err)
	}

	var newLink models.Link
	newLink, err = s.storageRepo.CreateLink(inputData)
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while creating new link: %w", err)
	}

	return helper.LinkResponseFromLinkModel(newLink), nil
}

func (s *Service) UpdateLink(ctx context.Context, request *shortener.UpdateLinkRequest) (*shortener.Link, error) {
	var err error
	var inputData *models.Link
	inputData, err = helper.LinkModelFromLinkUpdateRequest(request)
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while validating link: %v", err)
	}

	var newLink models.Link
	newLink, err = s.storageRepo.UpdateLink(inputData)
	if err != nil {
		return &shortener.Link{}, fmt.Errorf("error while updating link by id: %w", err)
	}

	return helper.LinkResponseFromLinkModel(newLink), nil
}

func (s *Service) DeleteLink(ctx context.Context, request *shortener.DeleteLinkRequest) (*shortener.DeleteLinkResponse, error) {
	var err error
	var id uuid.UUID
	id, err = helper.GetValidatedId(request)
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

	err = s.cachingRepo.DeleteUrl(*cacheKeyToDelete)
	if err != nil {
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while deleting link from caching: %v", err)
	}

	return &shortener.DeleteLinkResponse{Status: true}, nil
}

func (s *Service) Redirect(ctx context.Context, request *shortener.RedirectRequest) (*shortener.RedirectResponse, error) {
	var err error
	var shortLink = request.ShortLink
	var targetUrl string

	targetUrl, err = s.cachingRepo.GetUrl(shortLink)
	if err != nil {
		// if it's not in cache, then we get in from relational DB
		var link models.Link
		link, err = s.storageRepo.GetLinkByShortUrl(shortLink)
		if err != nil {
			return &shortener.RedirectResponse{}, fmt.Errorf("error while getting link: %v", err)
		}

		targetUrl = link.TargetUrl

		// we better cache the link, so we won't have to visit DB too often
		go func() {
			cacheErr := s.cachingRepo.SetUrl(shortLink, targetUrl)
			if cacheErr != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx,
					"couldn't cache link",
					zap.String("key", shortLink),
					zap.String("value", targetUrl))
			}
		}()
	}

	go func() {
		var click *models.Click
		click, err = helper.RedirectRequestToClick(request)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx,
				"failed to parse click",
				zap.Error(err),
				zap.String("short_link", shortLink),
				zap.String("visitor_token", request.GetVisitorToken()),
				zap.Time("clicked_at", request.GetClickedAt().AsTime()))
		}
		err = s.senderRepo.SendClick(ctx, click)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx,
				"failed to send click to kafka",
				zap.Error(err),
				zap.String("short_link", shortLink),
				zap.String("visitor_token", request.GetVisitorToken()),
				zap.Time("clicked_at", request.GetClickedAt().AsTime()))
		}
	}()

	return &shortener.RedirectResponse{TargetUrl: targetUrl}, nil
}

func (s *Service) GetLinksCountByUserId(ctx context.Context, request *shortener.GetLinksCountByUserIdRequest) (*shortener.GetLinksCountByUserIdResponse, error) {
	var err error
	var userId uuid.UUID
	userId, err = helper.GetValidatedUserId(request)
	if err != nil {
		return &shortener.GetLinksCountByUserIdResponse{Count: 0},
			fmt.Errorf("error while getting user id: %v", err)
	}

	var count int32
	count, err = s.storageRepo.GetLinksCountByUserId(userId)
	if err != nil {
		return &shortener.GetLinksCountByUserIdResponse{Count: 0},
			fmt.Errorf("error while getting links count by id: %v", err)
	}

	return &shortener.GetLinksCountByUserIdResponse{Count: count}, nil
}

func (s *Service) GetLinkOwnerByShortLink(ctx context.Context, request *shortener.GetLinkOwnerByShortLinkRequest) (*shortener.GetLinkOwnerByShortLinkResponse, error) {
	shortLink, err := helper.ValidateNotEmptyStr(request.ShortLink)
	if err != nil {
		return &shortener.GetLinkOwnerByShortLinkResponse{}, fmt.Errorf("error while validating link: %v", err)
	}

	var userId string
	userId, err = s.storageRepo.GetLinkOwnerByShortLink(shortLink)
	if err != nil {
		return &shortener.GetLinkOwnerByShortLinkResponse{}, fmt.Errorf("error while getting user id: %v", err)
	}

	return &shortener.GetLinkOwnerByShortLinkResponse{UserId: userId}, nil
}
