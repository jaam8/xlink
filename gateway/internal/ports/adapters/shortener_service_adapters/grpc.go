package shortener_service_adapters

import (
	"context"
	"fmt"
	"xlink/common/gen/shortener"
	"xlink/common/grpc/pool"
)

type ShortenerServiceRepositoryGRPC struct {
	grpcPool *pool.GrpcPool
}

func NewShortenerServiceRepositoryGRPC(grpcPool *pool.GrpcPool) *ShortenerServiceRepositoryGRPC {
	return &ShortenerServiceRepositoryGRPC{
		grpcPool: grpcPool,
	}
}

func (s ShortenerServiceRepositoryGRPC) Redirect(request *shortener.RedirectRequest) (*shortener.RedirectResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := shortener.NewShortenerServiceClient(conn)

	response, grpcErr := client.Redirect(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in Redirect grpc: %v", grpcErr)
	}

	return response, nil
}

func (s ShortenerServiceRepositoryGRPC) CreateNewLink(request *shortener.CreateLinkRequest) (*shortener.Link, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := shortener.NewShortenerServiceClient(conn)

	response, grpcErr := client.CreateNewLink(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in CreateNewLink grpc: %v", grpcErr)
	}

	return response, nil
}

func (s ShortenerServiceRepositoryGRPC) UpdateLink(request *shortener.UpdateLinkRequest) (*shortener.Link, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := shortener.NewShortenerServiceClient(conn)

	response, grpcErr := client.UpdateLink(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in UpdateLink grpc: %v", grpcErr)
	}

	return response, nil
}

func (s ShortenerServiceRepositoryGRPC) DeleteLink(request *shortener.DeleteLinkRequest) (*shortener.DeleteLinkResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := shortener.NewShortenerServiceClient(conn)

	response, grpcErr := client.DeleteLink(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in DeleteLink grpc: %v", grpcErr)
	}

	return response, nil
}

func (s ShortenerServiceRepositoryGRPC) GetLinksCountByUserId(request *shortener.GetLinksCountByUserIdRequest) (*shortener.GetLinksCountByUserIdResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := shortener.NewShortenerServiceClient(conn)

	response, grpcErr := client.GetLinksCountByUserId(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in GetLinksCountByUserId grpc: %v", grpcErr)
	}

	return response, nil
}
