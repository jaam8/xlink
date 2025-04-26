package analytics_service_adapters

import (
	"context"
	"fmt"
	"xlink/common/gen/analytics"
	"xlink/common/grpc/pool"
)

type AnalyticsServiceRepositoryGRPC struct {
	grpcPool *pool.GrpcPool
}

func NewAnalyticsServiceRepositoryGRPC(grpcPool *pool.GrpcPool) *AnalyticsServiceRepositoryGRPC {
	return &AnalyticsServiceRepositoryGRPC{
		grpcPool: grpcPool,
	}
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByCountry(request *analytics.GetClicksRequest) (*analytics.ClicksByCountryResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByCountry(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByCountry grpc: %v", grpcErr)
	}

	return response, nil
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByRegion(request *analytics.GetClicksRequest) (*analytics.ClicksByRegionResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByRegion(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByRegion grpc: %v", grpcErr)
	}

	return response, nil
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByBrowser(request *analytics.GetClicksRequest) (*analytics.ClicksByBrowserResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByBrowser(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByBrowser grpc: %v", grpcErr)
	}

	return response, nil
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByOS(request *analytics.GetClicksRequest) (*analytics.ClicksByOSResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByOS(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByOS grpc: %v", grpcErr)
	}

	return response, nil
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByDeviceType(request *analytics.GetClicksRequest) (*analytics.ClicksByDeviceTypeResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByDeviceType(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByDeviceType grpc: %v", grpcErr)
	}

	return response, nil
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByHour(request *analytics.GetClicksRequest) (*analytics.ClicksByHourResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByHour(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByHour grpc: %v", grpcErr)
	}

	return response, nil
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByDate(request *analytics.GetClicksRequest) (*analytics.ClicksByDateResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByDate(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByDate grpc: %v", grpcErr)
	}

	return response, nil
}

func (s AnalyticsServiceRepositoryGRPC) ClicksByReferrer(request *analytics.GetClicksRequest) (*analytics.ClicksByReferrerResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := analytics.NewAnalyticsServiceClient(conn)

	response, grpcErr := client.ClicksByReferrer(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in ClicksByReferrer grpc: %v", grpcErr)
	}

	return response, nil
}
