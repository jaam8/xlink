package ports

import (
	"time"
	"xlink/tg_bot/internal/models"
)

type ShortenerAdapter interface {
	CreateLink(userToken, targetUrl string, shortLink *string) (string, string, string, string, error)
	UpdateLink(userToken, shortLink, targetURL string, expireAt time.Time, regenerate bool,
	) (string, string, string, string, error)
	GetUserLinks(userToken string) ([]string, error)
	DeleteLink(userToken, shortLink string) error
}

type UserServiceAdapter interface {
	CreateUser(tgID *int64) (string, string, error)
	RefreshToken(userID, token string) (string, error)
	GetTokenByTgID(tgID int64) (string, error)
	LoginUser(token string) (string, *int64, error)
	GetUserProfile(userID string) (*int64, string, int64, error)
	SetTgID(userID string, tgID int64) error
}

type CacheAdapter interface {
	GetUserToken(tgID string) (string, error)
	SetUserToken(tgID, userToken string) error
}

type AnalyticsAdapter interface {
	ClicksByCountry(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.CountryDayStat], error)
	ClicksByRegion(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.RegionDayStat], error)
	ClicksByBrowser(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.BrowserDayStat], error)
	ClicksByDevice(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.DeviceDayStat], error)
	ClicksByOS(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.OSDayStat], error)
	ClicksByReferrer(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.ReferrerDayStat], error)
	ClicksByHour(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.HourDayStat], error)
	ClicksByDate(userToken, shortLink, startDate, endDate string) ([]models.DateStat, error)
}

type RendererAdapter interface {
	RenderChart(chartType string) (string, error)
}
