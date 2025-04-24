package tests

import (
	"testing"
	"time"
	"xlink/common/gen/shortener"
	"xlink/shortener/internal/models"
	"xlink/shortener/internal/service/helper"
	"xlink/shortener/internal/service/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestLinkResponseFromLinkModel(t *testing.T) {

	gen := true
	shortLink := "qwerty"
	targetUrl := "https://syubfugsebfuyegfbyu"
	expireAt := time.Now()

	model := models.Link{
		Id:        uuid.New(),
		UserId:    uuid.New(),
		Generated: &gen,
		ShortLink: &shortLink,
		TargetUrl: &targetUrl,
		CreatedAt: time.Now(),
		ExpireAt:  &expireAt,
	}

	resp := helper.LinkResponseFromLinkModel(model)

	assert.Equal(t, model.Id.String(), resp.LinkId)
	assert.Equal(t, model.UserId.String(), resp.UserId)
	assert.Equal(t, *model.ShortLink, resp.ShortLink)
	assert.Equal(t, *model.TargetUrl, resp.TargetUrl)
	assert.Equal(t, timestamppb.New(model.CreatedAt), resp.CreatedAt)
	assert.Equal(t, timestamppb.New(expireAt), resp.ExpireAt)
}

func TestLinkModelFromLinkCreateRequest_with_ready_shortlink(t *testing.T) {

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"

	testCreateLinkRequest := shortener.CreateLinkRequest{
		UserId:    userIDStr,
		ShortLink: &shortLinkStr,
		TargetUrl: "http://qwertysdijvnisdnc",
	}

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	expireAtStr := "2025-04-16T11:28:07+03:00"
	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	targetUrl := "http://qwertysdijvnisdnc"

	generated := false

	expectedModel := models.Link{
		UserId:    userUUID,
		ShortLink: &shortLinkStr,
		TargetUrl: &targetUrl,
		ExpireAt:  &expireAtTime,
		Generated: &generated,
	}

	testResponse, err := helper.LinkModelFromLinkCreateRequest(&testCreateLinkRequest, expireAtTime)
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel.UserId, testResponse.UserId)
	assert.Equal(t, expectedModel.ShortLink, testResponse.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testResponse.TargetUrl)
	assert.Equal(t, expireAtTime, *testResponse.ExpireAt)
	assert.Equal(t, *expectedModel.Generated, *testResponse.Generated)

}

func TestLinkModelFromLinkCreateRequest_without_ready_shortlink(t *testing.T) {

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := ""

	testCreateLinkRequest := shortener.CreateLinkRequest{
		UserId:    userIDStr,
		ShortLink: &shortLinkStr,
		TargetUrl: "http://qwertysdijvnisdnc",
	}

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	expireAtStr := "2025-04-16T11:28:07+03:00"
	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	generated := true

	targetUrl := "http://qwertysdijvnisdnc"

	shortLinkForModel := "nsuidfhuyefbuye"

	expectedModel := models.Link{
		UserId:    userUUID,
		ShortLink: &shortLinkForModel,
		TargetUrl: &targetUrl,
		ExpireAt:  &expireAtTime,
		Generated: &generated,
	}

	testResponse, err := helper.LinkModelFromLinkCreateRequest(&testCreateLinkRequest, expireAtTime)
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel.UserId, testResponse.UserId)
	assert.NotEqual(t, *expectedModel.ShortLink, *testResponse.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testResponse.TargetUrl)
	assert.Equal(t, expireAtTime, *testResponse.ExpireAt)
	assert.Equal(t, *expectedModel.Generated, *testResponse.Generated)

}

func TestLinkModelFromLinkUpdateRequest_with_regenarating(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"
	targerUrl := "http://qwertysdijvnisdnc"
	expireAtStr := "2025-05-16T11:28:07+03:00"

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	expireAtTimestamppb := timestamppb.New(expireAtTime)
	require.NoError(t, err)

	regenerated := true

	testUpdateLinkRequest := shortener.UpdateLinkRequest{
		LinkId:     idStr,
		UserId:     userIDStr,
		Regenerate: regenerated,
		ShortLink:  &shortLinkStr,
		TargetUrl:  &targerUrl,
		ExpireAt:   expireAtTimestamppb,
	}

	createdAtStr := "2025-04-16T11:28:06+03:00"
	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)

	newShortLink := utils.GenerateShortURL()

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		Generated: &regenerated,
		ShortLink: &newShortLink,
		TargetUrl: &targerUrl,
		CreatedAt: createdAtTime,
		ExpireAt:  &expireAtTime,
	}

	testResponse, err := helper.LinkModelFromLinkUpdateRequest(&testUpdateLinkRequest)
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel.Id, testResponse.Id)
	assert.Equal(t, expectedModel.UserId, testResponse.UserId)
	assert.NotEqual(t, *expectedModel.ShortLink, *testResponse.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testResponse.TargetUrl)
	assert.Equal(t, expireAtTime.UTC(), *testResponse.ExpireAt)
	assert.Equal(t, *expectedModel.Generated, *testResponse.Generated)

}
func TestLinkModelFromLinkUpdateRequest_without_regenarating(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"
	targetUrl := "http://qwertysdijvnisdnc"
	expireAtStr := "2025-05-16T11:28:07+03:00"

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	expireAtTimestamppb := timestamppb.New(expireAtTime)
	require.NoError(t, err)

	regenerated := false

	testUpdateLinkRequest := shortener.UpdateLinkRequest{
		LinkId:     idStr,
		UserId:     userIDStr,
		Regenerate: regenerated,
		ShortLink:  &shortLinkStr,
		TargetUrl:  &targetUrl,
		ExpireAt:   expireAtTimestamppb,
	}

	createdAtStr := "2025-04-16T11:28:06+03:00"
	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		Generated: &regenerated,
		ShortLink: &shortLinkStr,
		TargetUrl: &targetUrl,
		CreatedAt: createdAtTime,
		ExpireAt:  &expireAtTime,
	}

	testResponse, err := helper.LinkModelFromLinkUpdateRequest(&testUpdateLinkRequest)
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel.Id, testResponse.Id)
	assert.Equal(t, expectedModel.UserId, testResponse.UserId)
	assert.Equal(t, *expectedModel.ShortLink, *testResponse.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testResponse.TargetUrl)
	assert.Equal(t, expireAtTime.UTC(), *testResponse.ExpireAt)
	assert.Equal(t, *expectedModel.Generated, *testResponse.Generated)

}

func TestRedirectRequestToClick(t *testing.T) {
	shortLink := "abc123"
	clickedAt := timestamppb.Now()
	referrer := "https://example.com"
	ipAdress := "127.0.0.1"
	visitorToken := "segryureh747reh9er9h"
	userAgent := shortener.UserAgent{Browser: "chrome", DeviceType: "mobile", Os: "android"}

	req := &shortener.RedirectRequest{
		ShortLink:    shortLink,
		ClickedAt:    clickedAt,
		Referrer:     referrer,
		IpAddress:    ipAdress,
		VisitorToken: visitorToken,
		UserAgent:    &userAgent,
	}

	expectedModel := models.Click{
		ShortLink:    shortLink,
		ClickedAt:    clickedAt.AsTime(),
		Referrer:     referrer,
		IPAddress:    ipAdress,
		VisitorToken: visitorToken,
		Browser:      userAgent.Browser,
		DeviceType:   userAgent.DeviceType,
		Os:           userAgent.Os,
	}

	resp, err := helper.RedirectRequestToClick(req)
	if err != nil {
		t.Errorf("failed to make click model:%v", err)
	}

	assert.Equal(t, expectedModel.ShortLink, resp.ShortLink)
	assert.Equal(t, expectedModel.ClickedAt, resp.ClickedAt)
	assert.Equal(t, expectedModel.Referrer, resp.Referrer)
	assert.Equal(t, expectedModel.IPAddress, resp.IPAddress)
	assert.Equal(t, expectedModel.VisitorToken, resp.VisitorToken)
	assert.Equal(t, expectedModel.Browser, resp.Browser)
	assert.Equal(t, expectedModel.DeviceType, resp.DeviceType)
	assert.Equal(t, expectedModel.Os, resp.Os)
}
