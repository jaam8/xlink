package user_service_adapters

import (
	"context"
	"fmt"
	"xlink/common/gen/user_service"
	"xlink/common/grpc/pool"
)

type UserServiceRepositoryGRPC struct {
	grpcPool *pool.GrpcPool
}

func NewUserServiceRepositoryGRPC(grpcPool *pool.GrpcPool) *UserServiceRepositoryGRPC {
	return &UserServiceRepositoryGRPC{
		grpcPool: grpcPool,
	}
}

func (s *UserServiceRepositoryGRPC) CreateUser(request *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.CreateUser(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in CreateUser grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) GetUser(request *user_service.GetUserRequest) (*user_service.GetUserResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.GetUser(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in GetUser grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) GetUserIDByToken(request *user_service.GetUserIDByTokenRequest) (*user_service.GetUserIDByTokenResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.GetUserIDByToken(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in GetUserIDByToken grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) GetUserIDByTgID(request *user_service.GetUserIDByTgIDRequest) (*user_service.GetUserIDByTgIDResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.GetUserIDByTgID(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in GetUserIDByTgUD grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) UpdateUser(request *user_service.UpdateUserRequest) (*user_service.UpdateUserResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.UpdateUser(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in UpdateUser grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) CheckToken(request *user_service.TokenCheckRequest) (*user_service.TokenCheckResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.CheckToken(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in CheckToken grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) RefreshToken(request *user_service.RefreshTokenRequest) (*user_service.RefreshTokenResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.RefreshToken(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in RefreshToken grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) GetRole(request *user_service.GetRoleRequest) (*user_service.GetRoleResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.GetRole(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in GetRole grpc: %v", grpcErr)
	}

	return response, nil
}

func (s *UserServiceRepositoryGRPC) DeleteUser(request *user_service.DeleteUserRequest) (*user_service.DeleteUserResponse, error) {
	conn, err := s.grpcPool.GetConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get conn from pool: %v", err)
	}
	defer conn.Close() //nolint:all

	client := user_service.NewUserServiceClient(conn)

	response, grpcErr := client.DeleteUser(context.Background(), request)
	if grpcErr != nil {
		return nil, fmt.Errorf("error in DeleteUser grpc: %v", grpcErr)
	}

	return response, nil
}
