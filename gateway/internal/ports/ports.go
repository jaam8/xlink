package ports

import (
	"time"
	"xlink/common/gen/analytics"
	"xlink/common/gen/shortener"
	"xlink/common/gen/user_service"
)

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

type ShortenerServiceRepository interface {
	Redirect(request *shortener.RedirectRequest) (*shortener.RedirectResponse, error)
	GetLinks(request *shortener.GetLinksRequest) (*shortener.GetLinksResponse, error)
	GetLink(request *shortener.GetLinkRequest) (*shortener.Link, error)
	CreateNewLink(request *shortener.CreateLinkRequest) (*shortener.Link, error)
	UpdateLink(request *shortener.UpdateLinkRequest) (*shortener.Link, error)
	DeleteLink(request *shortener.DeleteLinkRequest) (*shortener.DeleteLinkResponse, error)
	GetLinksCountByUserId(request *shortener.GetLinksCountByUserIdRequest) (*shortener.GetLinksCountByUserIdResponse, error)
	GetLinkIdByShortLink(request *shortener.GetLinkIdByShortLinkRequest) (*shortener.GetLinkIdByShortLinkResponse, error)
}

type AnalyticsServiceRepository interface {
	ClicksByCountry(request *analytics.GetClicksRequest) (*analytics.ClicksByCountryResponse, error)
	ClicksByRegion(request *analytics.GetClicksRequest) (*analytics.ClicksByRegionResponse, error)
	ClicksByBrowser(request *analytics.GetClicksRequest) (*analytics.ClicksByBrowserResponse, error)
	ClicksByOS(request *analytics.GetClicksRequest) (*analytics.ClicksByOSResponse, error)
	ClicksByDeviceType(request *analytics.GetClicksRequest) (*analytics.ClicksByDeviceTypeResponse, error)
	ClicksByHour(request *analytics.GetClicksRequest) (*analytics.ClicksByHourResponse, error)
	ClicksByDate(request *analytics.GetClicksRequest) (*analytics.ClicksByDateResponse, error)
	ClicksByReferrer(request *analytics.GetClicksRequest) (*analytics.ClicksByReferrerResponse, error)
}

type RendererServiceRepository interface {
	Generate(shortLink string, param string, startDate time.Time, endDate time.Time, linkOwner string) ([]byte, error)
}
