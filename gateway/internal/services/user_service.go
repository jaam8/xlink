package services

import (
	"fmt"
	"time"
	"xlink/common/callers"
	"xlink/common/gen/user_service"
	"xlink/gateway/internal/ports"
)

type UserService struct {
	UserServiceRepo *ports.UserServiceRepository
	MaxRetries      uint
	BaseDelay       time.Duration
}

func NewUserService(userServiceRepo ports.UserServiceRepository, maxRetries uint, baseDelay time.Duration) *UserService {
	return &UserService{
		UserServiceRepo: &userServiceRepo,
		MaxRetries:      maxRetries,
		BaseDelay:       baseDelay,
	}
}

func (s *UserService) CreateUser(request *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	resultChan := make(chan *user_service.CreateUserResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).CreateUser(request)
		if err != nil {
			return fmt.Errorf("error in timeout CreateUser caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call CreateUser: %v", err)
	}

	response := <-resultChan
	close(resultChan)

	return response, nil
}

func (s *UserService) GetUser(request *user_service.GetUserRequest) (*user_service.GetUserResponse, error) {
	resultChan := make(chan *user_service.GetUserResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).GetUser(request)
		if err != nil {
			return fmt.Errorf("error in timeout GetUser caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetUser: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}

func (s *UserService) GetUserIDByToken(request *user_service.GetUserIDByTokenRequest) (*user_service.GetUserIDByTokenResponse, error) {
	resultChan := make(chan *user_service.GetUserIDByTokenResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).GetUserIDByToken(request)
		if err != nil {
			return fmt.Errorf("error in timeout GetUserIDByToken caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetUserIDByToken: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}

func (s *UserService) GetUserIDByTgID(request *user_service.GetUserIDByTgIDRequest) (*user_service.GetUserIDByTgIDResponse, error) {
	resultChan := make(chan *user_service.GetUserIDByTgIDResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).GetUserIDByTgID(request)
		if err != nil {
			return fmt.Errorf("error in timeout GetUserIDByTgID caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetUserIDByTgID: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}

func (s *UserService) UpdateUser(request *user_service.UpdateUserRequest) (*user_service.UpdateUserResponse, error) {
	resultChan := make(chan *user_service.UpdateUserResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).UpdateUser(request)
		if err != nil {
			return fmt.Errorf("error in timeout UpdateUser caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call UdpateUser: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}

func (s *UserService) CheckToken(request *user_service.TokenCheckRequest) (*user_service.TokenCheckResponse, error) {
	resultChan := make(chan *user_service.TokenCheckResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).CheckToken(request)
		if err != nil {
			return fmt.Errorf("error in timeout CheckToken caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call CheckToken: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}

func (s *UserService) RefreshToken(request *user_service.RefreshTokenRequest) (*user_service.RefreshTokenResponse, error) {
	resultChan := make(chan *user_service.RefreshTokenResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).RefreshToken(request)
		if err != nil {
			return fmt.Errorf("error in timeout RefreshToken caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call RefreshToken: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}

func (s *UserService) DeleteUser(request *user_service.DeleteUserRequest) (*user_service.DeleteUserResponse, error) {
	resultChan := make(chan *user_service.DeleteUserResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).DeleteUser(request)
		if err != nil {
			return fmt.Errorf("error in timeout DeleteUser caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call DeleteUser: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}

func (s *UserService) GetRole(request *user_service.GetRoleRequest) (*user_service.GetRoleResponse, error) {
	resultChan := make(chan *user_service.GetRoleResponse, 1)

	err := callers.Retry(func() error {
		response, err := (*s.UserServiceRepo).GetRole(request)
		if err != nil {
			return fmt.Errorf("error in timeout GetRole caller: %v", err)
		}
		resultChan <- response
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetRole: %v", err)
	}

	response := <-resultChan
	close(resultChan)
	return response, nil
}
