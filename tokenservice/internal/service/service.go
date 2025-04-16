package service

import (
	"context"
	"fmt"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/logger"
	"go.uber.org/zap"
	"tokenservice/internal/ports"
	"tokenservice/pkg/api/tokenservice"
)

const (
	userIdKey = "user_id"
	tokenKey  = "token"
	statusKey = "status"
)

type Service struct {
	tokenservice.TokenServiceServer
	tokensRepo ports.TokensRepository
}

func New(tokensRepo ports.TokensRepository) *Service {
	return &Service{
		tokensRepo: tokensRepo,
	}
}

func (s *Service) CheckToken(ctx context.Context, req *tokenservice.TokenCheckRequest) (*tokenservice.TokenCheckResponse, error) {
	tokenCorrect, err := s.tokensRepo.Check(req.UserId, req.Token)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't create new token for userId",
			zap.String(tokenKey, req.Token),
			zap.String(userIdKey, req.UserId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("couldn't check token '%s' for userId %s: %v", req.Token, req.UserId, err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "checked token",
		zap.String(userIdKey, req.UserId),
		zap.String(tokenKey, req.Token),
		zap.Bool(statusKey, tokenCorrect),
	)
	return &tokenservice.TokenCheckResponse{Status: tokenCorrect}, nil
}

func (s *Service) CreateToken(ctx context.Context, req *tokenservice.TokenRequest) (*tokenservice.TokenCreateResponse, error) {
	if len(req.UserId) == 0 {
		return nil, fmt.Errorf("user id is empty")
	}

	newToken, err := s.tokensRepo.Create(req.UserId)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't create new token for userId",
			zap.String(userIdKey, req.UserId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("couldn't create new token for userId %s: %v", req.UserId, err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "created new token", zap.String(userIdKey, req.UserId))
	return &tokenservice.TokenCreateResponse{Token: newToken}, nil
}

func (s *Service) GetUserIDByToken(ctx context.Context, req *tokenservice.TokenGetUserIDRequest) (*tokenservice.TokenGetUserIDResponse, error) {
	if len(req.Token) == 0 {
		return nil, fmt.Errorf("user id is empty")
	}
	userId, err := s.tokensRepo.GetUserIdByToken(req.Token)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't get userId by token",
			zap.String(tokenKey, req.Token),
			zap.Error(err),
		)
		return nil, fmt.Errorf("couldn't get userId by token '%s': %v", req.Token, err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "got userId by token", zap.String(tokenKey, req.Token))
	return &tokenservice.TokenGetUserIDResponse{UserId: userId}, nil
}

func (s *Service) DeleteToken(ctx context.Context, req *tokenservice.TokenRequest) (*tokenservice.TokenDeleteResponse, error) {
	if len(req.UserId) == 0 {
		return nil, fmt.Errorf("user id is empty")
	}

	if err := s.tokensRepo.Delete(req.UserId); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't delete token for userId",
			zap.String(userIdKey, req.UserId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("couldn't delete token for userId %s: %v", req.UserId, err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "deleted token", zap.String(userIdKey, req.UserId))

	return &tokenservice.TokenDeleteResponse{}, nil
}
