package ports

import (
	"context"
	"github.com/google/uuid"
	"xlink/shortener/internal/models"
)

type ShortenerCacheRepository interface {
	GetUrl(shortUrl string) (string, error)
	SetUrl(shortUrl string, url string) error
	DeleteUrl(shortUrl string) error
}

type ShortenerStorageRepository interface {
	GetLinks(userId uuid.UUID) ([]*models.Link, error)
	GetLinkById(linkId uuid.UUID) (models.Link, error)
	GetLinkByShortUrl(shortUrl string) (models.Link, error)
	CreateLink(newLink *models.Link) (models.Link, error)
	UpdateLink(newLinkWithExistingId *models.Link) (models.Link, error)
	DeleteLink(linkId uuid.UUID) error
	GetLinksCountByUserId(userId uuid.UUID) (int32, error)
	GetLinkOwnerByShortLink(shortLink string) (string, error)
	GetLinkIdByShortLink(shortLink string) (string, error)
}

type ShortenerSenderRepository interface {
	SendClick(ctx context.Context, click *models.Click) error
}
