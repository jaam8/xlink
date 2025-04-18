package storage

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"xlink/shortener/internal/models"
)

type ShortenerStorageRepositoryPostgres struct {
	PostgresPool *pgxpool.Pool
}

func LinkSelectQuery(filter squirrel.Eq) (string, []interface{}, error) {
	sql, args, err := squirrel.Select(
		"id", "user_id", "group_id", "generated",
		"short_link", "url", "created_at", "expire_at",
	).
		From("schema_name.urls").
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
		Scan(&link.Id, &link.UserId, &link.GroupId, &link.Generated,
			&link.ShortLink, &link.Url, &link.CreatedAt, &link.ExpireAt)

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

	var link models.Link = models.Link{}

	_ = s.PostgresPool.QueryRow(context.Background(), sql, args...).
		Scan(&link.Id, &link.UserId, &link.GroupId, &link.Generated,
			&link.ShortLink, &link.Url, &link.CreatedAt, &link.ExpireAt)

	if link.Id == uuid.Nil {
		return models.Link{}, fmt.Errorf("link not found with shortUrl '%s'", shortUrl)
	}

	return link, nil
}

func (s *ShortenerStorageRepositoryPostgres) CreateLink(newLink *models.Link) (models.Link, error) {
	sql, args, err := squirrel.Insert("schema_name.urls").
		Columns("user_id", "group_id", "generated",
			"short_link", "url", "expire_at").
		Values(newLink.UserId, newLink.GroupId, newLink.Generated, newLink.ShortLink, newLink.Url, newLink.ExpireAt).
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

	sql, args, err := squirrel.Update("schema_name.urls").
		Where(squirrel.Eq{"id": newLinkWithExistingId.Id}).
		Set("user_id", newLinkWithExistingId.UserId).
		Set("group_id", newLinkWithExistingId.GroupId).
		Set("generated", newLinkWithExistingId.Generated).
		Set("short_link", newLinkWithExistingId.ShortLink).
		Set("url", newLinkWithExistingId.Url).
		Set("expire_at", newLinkWithExistingId.ExpireAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

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

	sql, args, err := squirrel.Delete("schema_name.urls").
		Where(squirrel.Eq{"id": linkId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	_, err = s.PostgresPool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("couldn't delete an existing link: %w", err)
	}

	return nil
}
