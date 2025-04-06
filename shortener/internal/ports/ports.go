package ports

import (
	"github.com/google/uuid"
	"shortener/internal/models"
)

type ShortenerCacheRepository interface {
	GetUrl(shortUrl string) (string, error)
	SetUrl(shortUrl string, url string) error
	DeleteUrl(shortUrl string) error
}

type ShortenerStorageRepository interface {
	GetLinkById(linkId uuid.UUID) (models.Link, error)
	GetLinkByShortUrl(shortUrl string) (models.Link, error)
	CreateLink(newLink *models.Link) (models.Link, error)
	UpdateLink(newLinkWithExistingId *models.Link) (models.Link, error)
	DeleteLink(linkId uuid.UUID) error
}

type ShortenerSenderRepository interface {
	SendRedirectInfo()
}
