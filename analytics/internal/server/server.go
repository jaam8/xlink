package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"xlink/analytics/internal/ports"
	"xlink/analytics/internal/service"
	"xlink/common/gen/analytics"
	"xlink/common/grpc/interceptors"
	"xlink/common/logger"
)

func CreateGRPC(grpcSrv *service.Service) (*grpc.Server, error) {
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AddLogMiddleware))
	analytics.RegisterAnalyticsServiceServer(server, grpcSrv)
	return server, nil
}

func NewService(
	cacheAdapter ports.CacheAdapter,
	storageAdapter ports.StorageAdapter,
	consumerAdapter ports.ConsumerAdapter,
	shortenerAdapter ports.ShortenerAdapter,
) *service.Service {
	return service.New(cacheAdapter, storageAdapter, consumerAdapter, shortenerAdapter)
}

func RunGRPC(ctx context.Context, server *grpc.Server, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx,
			"ANALYTICS failed to create listener on port",
			zap.Int("port", port),
			zap.Error(err))
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("ANALYTICS listening at :%d", port))
	if err = server.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx,
			"ANALYTICS failed to serve grpc server",
			zap.Error(err))
	}
}
