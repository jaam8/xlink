package shortener_adapters

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
	"xlink/common/callers"
	"xlink/common/gen/shortener"
)

type ShortenerRepositoryGRPC struct {
	Address     string
	DialOptions []grpc.DialOption
	Timeout     time.Duration
}

func NewShortenerRepositoryGRPC(
	address string,
	options []grpc.DialOption,
	timeout time.Duration,
) *ShortenerRepositoryGRPC {
	return &ShortenerRepositoryGRPC{
		Address:     address,
		Timeout:     timeout,
		DialOptions: options,
	}
}

func (s *ShortenerRepositoryGRPC) GetGRPCClient() (*grpc.ClientConn, *shortener.ShortenerServiceClient, error) {
	grpcConn, err := grpc.NewClient(s.Address, s.DialOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to shortener by gRPC: %v", err)
	}

	client := shortener.NewShortenerServiceClient(grpcConn)

	return grpcConn, &client, nil
}

func (s *ShortenerRepositoryGRPC) GetLinksCountByUserId(userId string) (int32, error) {
	conn, clientPointer, err := s.GetGRPCClient()
	if err != nil {
		return 0, fmt.Errorf("cannot connect to shortener by gRPC: %v", err)
	}

	defer conn.Close() //nolint:all

	resultChan := make(chan int32)

	err = callers.Timeout(func() error {
		request := &shortener.GetLinksCountByUserIdRequest{UserId: userId}
		response, grpcErr := (*clientPointer).GetLinksCountByUserId(context.Background(), request)
		if grpcErr != nil {
			return fmt.Errorf("error in timeout gRPC caller: %v", grpcErr)
		}
		resultChan <- response.Count
		return nil
	}, s.Timeout)

	if err != nil {
		return 0, fmt.Errorf("couldn't get shortener.GetLinksCountByUserIdRequest gRPC response: %v", err)
	}
	linkCount := <-resultChan
	return linkCount, nil
}
