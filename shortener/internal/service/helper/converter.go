package helper

import (
	"fmt"
	"github.com/google/uuid"
	"shortener/internal/models"
	"shortener/pkg/api/shortener"
	"time"
)

func LinkResponseFromLinkModel(link models.Link) *shortener.Link {
	var GroupId *string = nil
	if link.GroupId != nil {
		*GroupId = link.GroupId.String()
	}

	return &shortener.Link{
		Id:        link.Id.String(),
		UserId:    link.UserId.String(),
		GroupId:   GroupId,
		Generated: link.Generated,
		ShortLink: link.ShortLink,
		Url:       link.Url,
		CreatedAt: link.CreatedAt.String(),
		ExpireAt:  link.ExpireAt.String(),
	}
}

func LinkModelFromLinkRequest(request LinkBodyRequest, expireAtDefault time.Time) (*models.Link, error) {
	var userId uuid.UUID
	var groupId *uuid.UUID
	var err error
	var expireAt time.Time
	var generated bool
	var shortLink, url string

	userId, err = GetValidatedUserId(request)
	if err != nil {
		return nil, fmt.Errorf("error while getting user id: %v", err)
	}

	// group id is optional
	groupId, err = GetValidatedGroupId(request, nil)
	if err != nil {
		return nil, fmt.Errorf("error while getting group id: %v", err)
	}

	// expireAt is optional too
	expireAt, err = GetValidatedExpireAt(request, expireAtDefault)
	if err != nil {
		return nil, fmt.Errorf("error while getting expire_at: %v", err)
	}

	err = ValidateStringNotEmpty(request.GetShortLink())
	if err != nil {
		return nil, fmt.Errorf("error while validating short link: %v", err)
	}
	shortLink = request.GetShortLink()

	err = ValidateUrl(request.GetUrl())
	if err != nil {
		return nil, fmt.Errorf("error while validating url: %v", err)
	}
	url = request.GetUrl()

	inputData := &models.Link{
		UserId:    userId,
		GroupId:   groupId,
		Generated: generated,
		ShortLink: shortLink,
		Url:       url,
		ExpireAt:  expireAt,
	}

	return inputData, nil
}
