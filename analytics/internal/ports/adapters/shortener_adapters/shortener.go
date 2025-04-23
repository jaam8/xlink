package shortener_adapters

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"time"
	"xlink/common/callers"
	"xlink/common/gen/shortener"
)

type ShortenerAdapter struct {
	Address     string
	DialOptions []grpc.DialOption
	Timeout     time.Duration
}

func NewShortenerAdapter(
	address string,
	options []grpc.DialOption,
	timeout time.Duration,
) *ShortenerAdapter {
	return &ShortenerAdapter{
		Address:     address,
		Timeout:     timeout,
		DialOptions: options,
	}
}

func (s *ShortenerAdapter) GetGRPCClient() (*grpc.ClientConn, *shortener.ShortenerServiceClient, error) {
	grpcConn, err := grpc.NewClient(s.Address, s.DialOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to shortener by gRPC: %v", err)
	}

	client := shortener.NewShortenerServiceClient(grpcConn)

	return grpcConn, &client, nil
}

func (s *ShortenerAdapter) GetLinkOwner(shortLink string) (uuid.UUID, error) {
	conn, clientPointer, err := s.GetGRPCClient()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("cannot connect to shortener by gRPC: %v", err)
	}

	defer conn.Close() //nolint

	resultChan := make(chan string, 1)

	err = callers.Timeout(func() error {
		request := &shortener.GetLinkOwnerByShortLinkRequest{ShortLink: shortLink}
		response, grpcErr := (*clientPointer).GetLinkOwnerByShortLink(context.Background(), request)
		if grpcErr != nil {
			return fmt.Errorf("error in timeout gRPC caller: %v", grpcErr)
		}
		resultChan <- response.GetLinkOwner()
		return nil
	}, s.Timeout)

	if err != nil {
		log.Printf("error in timeout gRPC caller: %v", err)
		return uuid.UUID{}, fmt.Errorf("couldn't get shortener.GetLinkOwnerByShortLinkRequest gRPC response: %v", err)
	}

	responseLinkOwner := <-resultChan
	close(resultChan)
	uuidLinkOwner, err := uuid.Parse(responseLinkOwner)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("couldn't parse link_owner: %v", err)
	}
	return uuidLinkOwner, nil
}
