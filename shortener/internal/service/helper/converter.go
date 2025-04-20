package helper

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"xlink/common/gen/shortener"
	"xlink/shortener/internal/models"
)

func LinkResponseFromLinkModel(link models.Link) *shortener.Link {
	var GroupId *string
	if link.GroupId != nil {
		str := link.GroupId.String()
		GroupId = &str
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

	generated = request.GetGenerated()

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

func LinkModelFromLinkRequestWithId(request LinkBodyRequestWithId, expireAtDefault time.Time) (*models.Link, error) {
	var err error
	var inputData *models.Link
	var id uuid.UUID
	inputData, err = LinkModelFromLinkRequest(request, expireAtDefault)
	if err != nil {
		return nil, fmt.Errorf("error while validating link with id: %v", err)
	}
	id, err = GetValidatedId(request)
	if err != nil {
		return nil, fmt.Errorf("error while validating link with id: %v", err)
	}

	return &models.Link{
		Id:        id,
		UserId:    inputData.UserId,
		GroupId:   inputData.GroupId,
		Generated: inputData.Generated,
		ShortLink: inputData.ShortLink,
		Url:       inputData.Url,
		CreatedAt: inputData.CreatedAt,
		ExpireAt:  inputData.ExpireAt,
	}, nil
}
