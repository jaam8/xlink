package runner

import (
	"context"
	"fmt"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/logger"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/transport/grpc/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"shortener/internal/ports"
	"shortener/internal/service"
	"shortener/pkg/api/shortener"
	"time"
)

func CreateGRPC(
	cachingRepo ports.ShortenerCacheRepository,
	storageRepo ports.ShortenerStorageRepository,
	senderRepo ports.ShortenerSenderRepository,
	defaultLinkExpireTime time.Duration,
) (*grpc.Server, error) {
	grpcSrv := service.New(cachingRepo, storageRepo, senderRepo, defaultLinkExpireTime)
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AddLogMiddleware))
	shortener.RegisterShortenerServiceServer(server, grpcSrv)
	return server, nil
}

func RunGRPC(ctx context.Context, server *grpc.Server, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "couldn't create listener on port", zap.Int("port", port), zap.Error(err))
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("listening at :%d", port))
	if err = server.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to serve grpc server", zap.Error(err))
	}
}
