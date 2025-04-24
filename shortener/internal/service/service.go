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
		logger.GetLoggerFromCtx(ctx).Info(ctx, "couldn't receive a link: invalid linkId",
			zap.Error(err),
		)
		return &shortener.Link{}, fmt.Errorf("error while validating link_id: %w", err)
	}

	var link models.Link
	link, err = s.storageRepo.GetLinkById(linkId)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't receive a link from storage",
			zap.Error(err),
		)
		return &shortener.Link{}, fmt.Errorf("error while getting link by id: %w", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "received a link",
		zap.String(helper.LinkIdKey, request.LinkId),
	)
	return helper.LinkResponseFromLinkModel(link), nil
}

func (s *Service) CreateNewLink(ctx context.Context, request *shortener.CreateLinkRequest) (*shortener.Link, error) {
	var err error
	var inputData *models.Link
	inputData, err = helper.LinkModelFromLinkCreateRequest(request, time.Now().Add(s.defaultLinkExpireTime))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "couldn't create a link: invalid input data",
			zap.Error(err),
		)
		return &shortener.Link{}, fmt.Errorf("error while validating link: %v", err)
	}

	var newLink models.Link
	newLink, err = s.storageRepo.CreateLink(inputData)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't create a link in storage",
			zap.Error(err),
		)
		return &shortener.Link{}, fmt.Errorf("error while creating new link: %w", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "created a new link",
		zap.String(helper.LinkIdKey, newLink.Id.String()),
	)
	return helper.LinkResponseFromLinkModel(newLink), nil
}

func (s *Service) UpdateLink(ctx context.Context, request *shortener.UpdateLinkRequest) (*shortener.Link, error) {
	var err error
	var inputData *models.Link
	inputData, err = helper.LinkModelFromLinkUpdateRequest(request)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "couldn't update a link: invalid input",
			zap.Error(err),
		)
		return &shortener.Link{}, fmt.Errorf("error while validating link: %v", err)
	}

	var newLink models.Link
	newLink, err = s.storageRepo.UpdateLink(inputData)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't update a link in storage",
			zap.Error(err),
		)
		return &shortener.Link{}, fmt.Errorf("error while updating link by id: %w", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "updated a link",
		zap.String(helper.LinkIdKey, request.LinkId),
	)
	return helper.LinkResponseFromLinkModel(newLink), nil
}

func (s *Service) DeleteLink(ctx context.Context, request *shortener.DeleteLinkRequest) (*shortener.DeleteLinkResponse, error) {
	var err error
	var id uuid.UUID
	id, err = helper.GetValidatedId(request)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "couldn't delete a link: invalid linkId",
			zap.Error(err),
		)
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while getting id: %v", err)
	}

	var link models.Link
	link, err = s.storageRepo.GetLinkById(id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't delete a link: error while checking it's existence in storage",
			zap.Error(err),
		)
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while getting link: %v", err)
	}
	cacheKeyToDelete := link.ShortLink

	err = s.storageRepo.DeleteLink(id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't delete a link: error while deleting from storage",
			zap.Error(err),
		)
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while deleting link from storage: %v", err)
	}

	err = s.cachingRepo.DeleteUrl(*cacheKeyToDelete)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't delete a link: error while invalidating cache",
			zap.Error(err),
		)
		return &shortener.DeleteLinkResponse{Status: false}, fmt.Errorf("error while deleting link from caching: %v", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "received link",
		zap.String(helper.LinkIdKey, request.LinkId),
	)
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
			logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't perform a redirect: error while getting target url",
				zap.Error(err),
			)
			return &shortener.RedirectResponse{}, fmt.Errorf("error while getting link: %v", err)
		}

		targetUrl = *link.TargetUrl

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

	logger.GetLoggerFromCtx(ctx).Info(ctx, "performed a redirect",
		zap.String(helper.ShortUrlKey, request.ShortLink),
		zap.String(helper.TargetUrlKey, targetUrl),
		zap.Error(err),
	)
	return &shortener.RedirectResponse{TargetUrl: targetUrl}, nil
}

func (s *Service) GetLinksCountByUserId(ctx context.Context, request *shortener.GetLinksCountByUserIdRequest) (*shortener.GetLinksCountByUserIdResponse, error) {
	var err error
	var userId uuid.UUID
	userId, err = helper.GetValidatedUserId(request)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "couldn't get links count by user id: invalid userId",
			zap.Error(err),
		)
		return &shortener.GetLinksCountByUserIdResponse{Count: 0},
			fmt.Errorf("error while getting user id: %v", err)
	}

	var count int32
	count, err = s.storageRepo.GetLinksCountByUserId(userId)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't delete a link: error while querying the storage",
			zap.Error(err),
		)
		return &shortener.GetLinksCountByUserIdResponse{Count: 0},
			fmt.Errorf("error while getting links count by id: %v", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "got links count by user id",
		zap.String(helper.UserIdKey, request.UserId),
		zap.Int32(helper.CountKey, count),
		zap.Error(err),
	)
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

	logger.GetLoggerFromCtx(ctx).Info(ctx, "got link owner by short link",
		zap.String(helper.ShortUrlKey, request.ShortLink),
		zap.String(helper.UserIdKey, userId),
		zap.Error(err),
	)
	return &shortener.GetLinkOwnerByShortLinkResponse{LinkOwner: userId}, nil
}

func (s *Service) GetLinkIdByShortLink(ctx context.Context, request *shortener.GetLinkIdByShortLinkRequest) (*shortener.GetLinkIdByShortLinkResponse, error) {
	shortLink, err := helper.ValidateNotEmptyStr(request.ShortLink)
	if err != nil {
		return &shortener.GetLinkIdByShortLinkResponse{}, fmt.Errorf("error while validating link: %v", err)
	}

	var linkId string
	linkId, err = s.storageRepo.GetLinkIdByShortLink(shortLink)
	if err != nil {
		return &shortener.GetLinkIdByShortLinkResponse{}, fmt.Errorf("error while getting link id: %v", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "got link id by short link",
		zap.String(helper.ShortUrlKey, request.ShortLink),
		zap.String(helper.LinkIdKey, linkId),
		zap.Error(err),
	)
	return &shortener.GetLinkIdByShortLinkResponse{LinkId: linkId}, nil
}
