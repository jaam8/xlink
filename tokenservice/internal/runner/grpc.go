package runner

import (
	"context"
	"fmt"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/logger"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/transport/grpc/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"tokenservice/internal/ports"
	"tokenservice/internal/service"
	"tokenservice/pkg/api/tokenservice"
)

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

func CreateGRPC(tokensRepo ports.TokensRepository) (*grpc.Server, error) {
	grpcSrv := service.New(tokensRepo)
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AddLogMiddleware))
	tokenservice.RegisterTokenServiceServer(server, grpcSrv)
	return server, nil
}
