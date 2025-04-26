package ports

import (
	"time"
	"xlink/tg_bot/internal/models"
)

type ShortenerAdapter interface {
	CreateLink(shortLink *string, targetUrl string) (string, string, time.Time, time.Time, error)
	UpdateLink(shortLink, targetURL string, expireAt time.Time, regenerate bool,
	) (string, string, time.Time, time.Time, error)
	DeleteLink(shortLink string) error
}

type UserServiceAdapter interface {
	CreateUser(tgID *int64) (string, string, error)
	RefreshToken(userID, token string) (string, error)
	GetTokenByTgID(tgID int64) (string, error)
	LoginUser(token string) (string, *int64, error)
	GetUserProfile(userID string) (*int64, string, int64, error)
	SetTgID(userID string, tgID int64) error
}

type AnalyticsAdapter interface {
	ClicksByCountry(linkOwner, shortLink, startDate, endDate string) (models.Stat[models.CountryDayStat], error)
	ClicksByRegion(linkOwner, shortLink, startDate, endDate string) (models.Stat[models.RegionDayStat], error)
	ClicksByBrowser(linkOwner, shortLink, startDate, endDate string) (models.Stat[models.BrowserDayStat], error)
	ClicksByDevice(linkOwner, shortLink, startDate, endDate string) (models.Stat[models.DeviceDayStat], error)
	ClicksByOS(linkOwner, shortLink, startDate, endDate string) (models.Stat[models.OSDayStat], error)
	ClicksByReferrer(linkOwner, shortLink, startDate, endDate string) (models.Stat[models.ReferrerDayStat], error)
	ClicksByHour(linkOwner, shortLink, startDate, endDate string) (models.Stat[models.HourDayStat], error)
	ClicksByDate(linkOwner, shortLink, startDate, endDate string) (models.DateStat, error)
}
