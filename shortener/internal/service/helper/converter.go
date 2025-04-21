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
		TargetUrl: link.TargetUrl,
		CreatedAt: timestamppb.New(link.CreatedAt),
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
		UserId:    userId,
		ShortLink: &shortLink,
		TargetUrl: targetUrl,
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

	targetUrl := request.GetTargetUrl()

	expireAt := request.GetExpireAt().AsTime()

	return &models.Link{
		Id:        linkId,
		UserId:    userId,
		Generated: generated,
		ShortLink: shortLink,
		TargetUrl: targetUrl,
		ExpireAt:  &expireAt,
	}, nil
}

func RedirectRequestToClick(request *shortener.RedirectRequest) (*models.Click, error) {
	shortLink, err := ValidateShortLink(request.GetShortLink())
	if err != nil {
		return nil, err
	}
	clickedAt := request.GetClickedAt().AsTime()
	referrer, err := ValidateNotEmptyStr(request.GetReferrer())
	if err != nil {
		return nil, fmt.Errorf("invalid referrer: %w", err)
	}
	ipAddress, err := ValidateIPAddress(request.GetIpAddress())
	if err != nil {
		return nil, err
	}
	visitorToken, err := ValidateNotEmptyStr(request.GetVisitorToken())
	if err != nil {
		return nil, fmt.Errorf("invalid visitor_token: %w", err)
	}
	userAgent := request.GetUserAgent()
	if userAgent == nil {
		return nil, fmt.Errorf("invalid user agent: user agent is nil")
	}
	browser, err := ValidateNotEmptyStr(userAgent.GetBrowser())
	if err != nil {
		return nil, fmt.Errorf("invalid browser: %w", err)
	}
	deviceType, err := ValidateNotEmptyStr(userAgent.GetDeviceType())
	if err != nil {
		return nil, fmt.Errorf("invalid device type: %w", err)
	}
	os, err := ValidateNotEmptyStr(userAgent.GetOs())
	if err != nil {
		return nil, fmt.Errorf("invalid os: %w", err)
	}
	return &models.Click{
		ShortLink:    shortLink,
		ClickedAt:    clickedAt,
		Referrer:     referrer,
		IPAddress:    ipAddress,
		VisitorToken: visitorToken,
		Browser:      browser,
		DeviceType:   deviceType,
		Os:           os,
	}, nil
}
