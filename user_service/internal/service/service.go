package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"xlink/common/gen/user_service"
	"xlink/common/logger"
	"xlink/user_service/internal/ports"
)

const (
	userIdKey = "user_id"
	tokenKey  = "token"
	statusKey = "status"
)

type Service struct {
	user_service.UserServiceServer
	cacheRepo     ports.UsersCacheRepository
	storageRepo   ports.UserStorageRepository
	shortenerRepo ports.ShortenerRepository
}

func New(
	cacheRepo ports.UsersCacheRepository,
	storageRepo ports.UserStorageRepository,
	shortenerRepo ports.ShortenerRepository,
) *Service {
	return &Service{
		cacheRepo:     cacheRepo,
		storageRepo:   storageRepo,
		shortenerRepo: shortenerRepo,
	}
}

func (s *Service) CheckToken(ctx context.Context, req *user_service.TokenCheckRequest) (*user_service.TokenCheckResponse, error) {
	if len(req.UserId) == 0 {
		return &user_service.TokenCheckResponse{Status: false}, fmt.Errorf("user_id is empty")
	}
	if len(req.Token) == 0 {
		return &user_service.TokenCheckResponse{Status: false}, fmt.Errorf("token is empty")
	}

	// try to get from cache
	tokenCorrect, err := s.cacheRepo.CheckToken(req.UserId, req.Token)
	if err != nil {

		// if cache doesn't work, check in DB
		tokenCorrect, err = s.storageRepo.CheckToken(req.UserId, req.Token)
		if err != nil {
			// both didn't work(
			logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't create new token for userId",
				zap.String(tokenKey, req.Token),
				zap.String(userIdKey, req.UserId),
				zap.Error(err),
			)
			return &user_service.TokenCheckResponse{Status: false}, fmt.Errorf("couldn't check token '%s' for userId %s: %v", req.Token, req.UserId, err)
		}

		// cache if not cached
		if tokenCorrect {
			err = s.cacheRepo.SetToken(req.UserId, req.Token)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).
					Error(ctx, "couldn't set token for userId",
						zap.String(userIdKey, req.UserId),
						zap.String(tokenKey, req.Token),
						zap.Error(err))
			}
		}
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "checked token",
		zap.String(userIdKey, req.UserId),
		zap.String(tokenKey, req.Token),
		zap.Bool(statusKey, tokenCorrect),
	)
	return &user_service.TokenCheckResponse{Status: tokenCorrect}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req *user_service.RefreshTokenRequest) (*user_service.RefreshTokenResponse, error) {
	if len(req.UserId) == 0 {
		return nil, fmt.Errorf("user_id is empty")
	}
	if len(req.Token) == 0 {
		return nil, fmt.Errorf("token is empty")
	}
	newToken, err := s.storageRepo.RefreshToken(req.UserId, req.Token)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't refresh token",
			zap.String(userIdKey, req.UserId),
			zap.String(tokenKey, req.Token),
			zap.Error(err),
		)
		return nil, fmt.Errorf("couldn't get userId by token '%s': %v", req.Token, err)
	}

	// invalidate & write actual cache
	err = s.cacheRepo.SetToken(req.UserId, newToken)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't cache refreshed token",
			zap.String(userIdKey, req.UserId),
			zap.String(tokenKey, req.Token),
			zap.Error(err),
		)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "refreshed token for user", zap.String(userIdKey, req.UserId))
	return &user_service.RefreshTokenResponse{Token: newToken}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	userId, token, err := s.storageRepo.CreateUser(req.TgId, req.IsStaff, req.IsAdmin)
	if err != nil {
		return &user_service.CreateUserResponse{}, fmt.Errorf("couldn't create user: %v", err)
	}
	return &user_service.CreateUserResponse{
		UserId: userId,
		Token:  token,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, req *user_service.GetUserRequest) (*user_service.GetUserResponse, error) {
	var storageErr, shortenerErr error
	var userId, role string
	var telegramId *int64
	var linkCount int32

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		userId, role, telegramId, storageErr = s.storageRepo.GetUser(req.UserId)
	}()

	go func() {
		defer waitGroup.Done()
		linkCount, shortenerErr = s.shortenerRepo.GetLinksCountByUserId(req.UserId)
	}()

	waitGroup.Wait()

	if storageErr != nil {
		return &user_service.GetUserResponse{},
			fmt.Errorf("couldn't get user with id='%s': %v",
				req.UserId, storageErr)
	}
	if shortenerErr != nil {
		return &user_service.GetUserResponse{},
			fmt.Errorf("couldn't get user links count (id='%s'): %v",
				userId, shortenerErr)
	}

	return &user_service.GetUserResponse{
		UserId:    userId,
		Role:      role,
		TgId:      telegramId,
		LinkCount: linkCount,
	}, nil
}

func (s *Service) GetUserIDByToken(ctx context.Context, req *user_service.GetUserIDByTokenRequest) (*user_service.GetUserIDByTokenResponse, error) {
	userId, status, err := s.storageRepo.GetUserIDByToken(req.Token)
	if err != nil {
		return &user_service.GetUserIDByTokenResponse{}, fmt.Errorf("couldn't get user id by token: %v", err)
	}

	return &user_service.GetUserIDByTokenResponse{
		UserId: userId,
		Status: status,
	}, nil
}

func (s *Service) GetUserIDByTgID(ctx context.Context, req *user_service.GetUserIDByTgIDRequest) (*user_service.GetUserIDByTgIDResponse, error) {
	userId, status, err := s.storageRepo.GetUserIDByTgId(req.TgId)
	if err != nil {
		return &user_service.GetUserIDByTgIDResponse{}, fmt.Errorf("couldn't get user id by tgId='%d': %v", req.TgId, err)
	}

	return &user_service.GetUserIDByTgIDResponse{
		UserId: userId,
		Status: status,
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, req *user_service.UpdateUserRequest) (*user_service.UpdateUserResponse, error) {
	status, err := s.storageRepo.UpdateUser(req.UserId, req.TgId, req.IsStaff, req.IsAdmin)
	if err != nil {
		return &user_service.UpdateUserResponse{}, fmt.Errorf("couldn't update user: %v", err)
	}

	return &user_service.UpdateUserResponse{
		Status: status,
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *user_service.DeleteUserRequest) (*user_service.DeleteUserResponse, error) {
	status, err := s.storageRepo.DeleteUser(req.UserId)
	if err != nil {
		return &user_service.DeleteUserResponse{}, fmt.Errorf("couldn't delete user: %v", err)
	}

	return &user_service.DeleteUserResponse{Status: status}, nil
}

func (s *Service) GetRole(ctx context.Context, req *user_service.GetRoleRequest) (*user_service.GetRoleResponse, error) {
	// try to get from cache
	role, isStaff, isAdmin, err := s.cacheRepo.GetRole(req.UserId)

	// couldn't get from cache
	if err != nil {
		role, isStaff, isAdmin, err = s.storageRepo.GetRole(req.UserId)

		// couldn't get from postgres
		if err != nil {
			return &user_service.GetRoleResponse{}, fmt.Errorf("couldn't get role for userId='%s': %v", req.UserId, err)
		}

		go func() {
			_ = s.cacheRepo.SetRole(req.UserId, isStaff, isAdmin) //nolint:all
		}()
	}

	return &user_service.GetRoleResponse{
		Role:    role,
		IsStaff: isStaff,
		IsAdmin: isAdmin,
	}, nil
}
