package ports

import (
	"context"
	"github.com/google/uuid"
	"time"
	"xlink/analytics/internal/models"
)

type CacheAdapter interface {
	CheckVisitorToken(visitorToken, shortLink string) (bool, error)
	SetVisitorToken(visitorToken, shortLink string) error
}

type StorageAdapter interface {
	SaveClicks(ctx context.Context, clicks []*models.Click) error
	GetClicksByCountry(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.CountryDayStats, error)
	GetClicksByRegion(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.RegionDayStats, error)
	GetClicksByReferrer(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.ReferrerDayStats, error)
	GetClicksByBrowser(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.BrowserDayStats, error)
	GetClicksByOS(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.OSDayStats, error)
	GetClicksByDeviceType(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.DeviceDayStats, error)
	GetClicksByHour(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.HourDayStats, error)
	GetClicksByDate(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.DateStat, error)
}

type ConsumerAdapter interface {
	ConsumeClickEvent(ctx context.Context) (*models.ClickEvent, error)
}

type ShortenerAdapter interface {
	GetLinkOwner(shortLink string) (uuid.UUID, error)
}
