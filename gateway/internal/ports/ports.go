package ports

import "xlink/common/gen/user_service"

type UserServiceRepository interface {
	CreateUser(request *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error)
	GetUser(request *user_service.GetUserRequest) (*user_service.GetUserResponse, error)
	GetUserIDByToken(request *user_service.GetUserIDByTokenRequest) (*user_service.GetUserIDByTokenResponse, error)
	GetUserIDByTgID(request *user_service.GetUserIDByTgIDRequest) (*user_service.GetUserIDByTgIDResponse, error)
	UpdateUser(request *user_service.UpdateUserRequest) (*user_service.UpdateUserResponse, error)
	CheckToken(request *user_service.TokenCheckRequest) (*user_service.TokenCheckResponse, error)
	RefreshToken(request *user_service.RefreshTokenRequest) (*user_service.RefreshTokenResponse, error)
	GetRole(request *user_service.GetRoleRequest) (*user_service.GetRoleResponse, error)
	DeleteUser(request *user_service.DeleteUserRequest) (*user_service.DeleteUserResponse, error)
}
