package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"xlink/common/logger"
	"xlink/token_service/internal/ports"
	"xlink/token_service/pkg/api/token_service"
)

const (
	userIdKey = "user_id"
	tokenKey  = "token"
	statusKey = "status"
)

type Service struct {
	token_service.TokenServiceServer
	tokensRepo ports.TokensRepository
}

func New(tokensRepo ports.TokensRepository) *Service {
	return &Service{
		tokensRepo: tokensRepo,
	}
}

func (s *Service) CheckToken(ctx context.Context, req *token_service.TokenCheckRequest) (*token_service.TokenStatusResponse, error) {
	if len(req.UserId) == 0 {
		return nil, fmt.Errorf("user_id is empty")
	}
	if len(req.Token) == 0 {
		return nil, fmt.Errorf("token is empty")
	}
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
	return &token_service.TokenStatusResponse{Status: tokenCorrect}, nil
}

func (s *Service) CreateToken(ctx context.Context, req *token_service.TokenRequest) (*token_service.TokenCreateResponse, error) {
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
	return &token_service.TokenCreateResponse{Token: newToken}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req *token_service.TokenCheckRequest) (*token_service.TokenCreateResponse, error) {
	if len(req.UserId) == 0 {
		return nil, fmt.Errorf("user_id is empty")
	}
	if len(req.Token) == 0 {
		return nil, fmt.Errorf("token is empty")
	}
	token, err := s.tokensRepo.Refresh(req.UserId, req.Token)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't get userId by token",
			zap.String(tokenKey, req.Token),
			zap.Error(err),
		)
		return nil, fmt.Errorf("couldn't get userId by token '%s': %v", req.Token, err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "got userId by token", zap.String(tokenKey, req.Token))
	return &token_service.TokenCreateResponse{Token: token}, nil
}

func (s *Service) DeleteToken(ctx context.Context, req *token_service.TokenRequest) (*token_service.TokenStatusResponse, error) {
	if len(req.UserId) == 0 {
		return nil, fmt.Errorf("user_id is empty")
	}
	status, err := s.tokensRepo.Delete(req.UserId)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't delete token for userId",
			zap.String(userIdKey, req.UserId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("couldn't delete token for userId %s: %v", req.UserId, err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "deleted token", zap.String(userIdKey, req.UserId))

	return &token_service.TokenStatusResponse{Status: status}, nil
}
