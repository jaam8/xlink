package storage

import (
	"context"
	"fmt"
	"xlink/shortener/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortenerStorageRepositoryPostgres struct {
	PostgresPool *pgxpool.Pool
}

func LinkSelectQuery(filter squirrel.Eq) (string, []interface{}, error) {
	sql, args, err := squirrel.Select(
		"id", "user_id", "group_id", "generated",
		"short_link", "url", "created_at", "expire_at",
	).
		From("shortener.urls").
		Where(filter).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	return sql, args, err
}

func NewShortenerStorageRepositoryPostgres(pool *pgxpool.Pool) *ShortenerStorageRepositoryPostgres {
	return &ShortenerStorageRepositoryPostgres{
		PostgresPool: pool,
	}
}

func (s *ShortenerStorageRepositoryPostgres) GetLinkById(linkId uuid.UUID) (models.Link, error) {
	sql, args, err := LinkSelectQuery(squirrel.Eq{"id": linkId})

	if err != nil {
		return models.Link{}, fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	var link = models.Link{}

	_ = s.PostgresPool.QueryRow(context.Background(), sql, args...).
		Scan(&link.Id, &link.UserId, &link.Generated,
			&link.ShortLink, &link.TargetUrl, &link.CreatedAt, &link.ExpireAt)

	if link.Id == uuid.Nil {
		return models.Link{}, fmt.Errorf("link not found with id '%s'", linkId)
	}

	return link, nil
}

func (s *ShortenerStorageRepositoryPostgres) GetLinkByShortUrl(shortUrl string) (models.Link, error) {
	sql, args, err := LinkSelectQuery(squirrel.Eq{"short_link": shortUrl})

	if err != nil {
		return models.Link{}, fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	var link = models.Link{}

	_ = s.PostgresPool.QueryRow(context.Background(), sql, args...).
		Scan(&link.Id, &link.UserId, &link.Generated,
			&link.ShortLink, &link.TargetUrl, &link.CreatedAt, &link.ExpireAt)

	if link.Id == uuid.Nil {
		return models.Link{}, fmt.Errorf("link not found with shortUrl '%s'", shortUrl)
	}

	return link, nil
}

func (s *ShortenerStorageRepositoryPostgres) CreateLink(newLink *models.Link) (models.Link, error) {
	sql, args, err := squirrel.Insert("shortener.urls").
		Columns("user_id", "generated",
			"short_link", "url", "expire_at").
		Values(newLink.UserId, newLink.Generated, newLink.ShortLink, newLink.TargetUrl, newLink.ExpireAt).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return models.Link{}, fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	var link = *newLink
	_ = s.PostgresPool.QueryRow(context.Background(), sql, args...).Scan(&link.Id, &link.CreatedAt)

	return link, nil
}

func (s *ShortenerStorageRepositoryPostgres) UpdateLink(newLinkWithExistingId *models.Link) (models.Link, error) {
	_, err := s.GetLinkById(newLinkWithExistingId.Id)
	if err != nil {
		return models.Link{}, fmt.Errorf("couldn't update an existing link: %w", err)
	}

	updateBuilder := squirrel.Update("shortener.urls").
		Where(squirrel.Eq{"id": newLinkWithExistingId.Id}).
		Set("user_id", newLinkWithExistingId.UserId)

	if newLinkWithExistingId.Generated != nil {
		updateBuilder = updateBuilder.Set("generated", newLinkWithExistingId.Generated)
	}

	if newLinkWithExistingId.ShortLink != nil {
		updateBuilder = updateBuilder.Set("short_link", newLinkWithExistingId.ShortLink)
	}

	if newLinkWithExistingId.ExpireAt != nil {
		updateBuilder = updateBuilder.Set("expire_at", newLinkWithExistingId.ExpireAt)
	}

	if newLinkWithExistingId.TargetUrl != nil {
		updateBuilder = updateBuilder.Set("url", newLinkWithExistingId.TargetUrl)
	}

	sql, args, err := updateBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return models.Link{}, fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	_, err = s.PostgresPool.Exec(context.Background(), sql, args...)
	if err != nil {
		return models.Link{}, fmt.Errorf("couldn't update an existing link: %w", err)
	}

	return *newLinkWithExistingId, nil
}

func (s *ShortenerStorageRepositoryPostgres) DeleteLink(linkId uuid.UUID) error {
	_, err := s.GetLinkById(linkId)
	if err != nil {
		return fmt.Errorf("couldn't delete an existing link: %w", err)
	}

	sql, args, err := squirrel.Delete("shortener.urls").
		Where(squirrel.Eq{"id": linkId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	_, err = s.PostgresPool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("couldn't delete an existing link: %w", err)
	}

	return nil
}

func (s *ShortenerStorageRepositoryPostgres) GetLinksCountByUserId(userId uuid.UUID) (int32, error) {
	var count int32
	sql, args, err := squirrel.Select("count(*)").
		From("shortener.urls").
		Where(squirrel.Eq{"user_id": userId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return count, fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	err = s.PostgresPool.QueryRow(context.Background(), sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("couldn't get links count by user id '%s': %w", userId, err)
	}
	return count, nil
}

func (s *ShortenerStorageRepositoryPostgres) GetLinkOwnerByShortLink(shortLink string) (string, error) {
	var userId string
	sql, args, err := squirrel.Select("user_id").
		From("shortener.urls").
		Where(squirrel.Eq{"short_link": shortLink}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	err = s.PostgresPool.QueryRow(context.Background(), sql, args...).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("couldn't get link owner by short link '%s': %w", shortLink, err)
	}
	return userId, nil
}

func (s *ShortenerStorageRepositoryPostgres) GetLinkIdByShortLink(shortLink string) (string, error) {
	var linkId string
	sql, args, err := squirrel.Select("id").
		From("shortener.urls").
		Where(squirrel.Eq{"short_link": shortLink}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("couldn't build an SQL query: %w", err)
	}

	err = s.PostgresPool.QueryRow(context.Background(), sql, args...).Scan(&linkId)
	if err != nil {
		return "", fmt.Errorf("couldn't get link id by shortLink='%s': %w", shortLink, err)
	}
	return linkId, nil
}
