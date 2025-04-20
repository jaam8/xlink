package helper

import (
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"xlink/common/gen/shortener"
	"xlink/shortener/internal/models"
	"xlink/shortener/internal/service/utils"
)

func LinkResponseFromLinkModel(link models.Link) *shortener.Link {
	return &shortener.Link{
		LinkId:    link.Id.String(),
		UserId:    link.UserId.String(),
		ShortLink: *link.ShortLink,
		TargetUrl: *link.Url,
		CreatedAt: timestamppb.New(*link.CreatedAt),
		ExpireAt:  timestamppb.New(*link.ExpireAt),
	}
}

func LinkModelFromLinkCreateRequest(request LinkCreateRequest, expiration time.Time) (*models.Link, error) {
	userId, err := GetValidatedUserId(request)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse create link request: invalid user id: %v", err)
	}
	shortLink := request.GetShortLink()
	targetUrl := request.GetTargetUrl()

	generated := false

	if len(shortLink) == 0 {
		shortLink = utils.GenerateShortURL()
		generated = true
	}

	return &models.Link{
		UserId:    &userId,
		ShortLink: &shortLink,
		Url:       &targetUrl,
		ExpireAt:  &expiration,
		Generated: &generated,
	}, nil
}

func LinkModelFromLinkUpdateRequest(request LinkUpdateRequest) (*models.Link, error) {
	var err error
	var linkId, userId uuid.UUID

	linkId, err = GetValidatedId(request)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse create link request: invalid link id: %v", err)
	}
	userId, err = GetValidatedUserId(request)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse create link request: invalid user id: %v", err)
	}

	var shortLink *string
	var generatedValue bool

	shortLinkText := request.GetShortLink()

	if request.GetRegenerate() {
		generatedValue = true
		shortLinkText = utils.GenerateShortURL()
	} else if len(shortLinkText) > 0 {
		generatedValue = false
	}
	generated := &generatedValue
	shortLink = &shortLinkText

	targetUrlText := request.GetTargetUrl()
	targetUrl := &targetUrlText

	expireAt := request.GetExpireAt().AsTime()

	return &models.Link{
		Id:        &linkId,
		UserId:    &userId,
		Generated: generated,
		ShortLink: shortLink,
		Url:       targetUrl,
		ExpireAt:  &expireAt,
	}, nil
}
