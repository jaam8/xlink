package tests

import (
	"context"
	"testing"
	"time"
	"xlink/common/gen/shortener"
	"xlink/common/logger"
	"xlink/shortener/internal/models"
	"xlink/shortener/internal/service"
	"xlink/shortener/internal/service/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mokShortenerCacheRepository struct{ mock.Mock }

func (m *mokShortenerCacheRepository) GetUrl(shortUrl string) (string, error) {
	args := m.Called(shortUrl)
	return args.Get(0).(string), args.Error(1)
}
func (m *mokShortenerCacheRepository) SetUrl(shortUrl string, url string) error { return nil }
func (m *mokShortenerCacheRepository) DeleteUrl(shortUrl string) error {
	args := m.Called(shortUrl)
	return args.Error(0)
}

type mokShortenerStorageRepository struct{ mock.Mock }

func (m *mokShortenerStorageRepository) GetLinkById(linkId uuid.UUID) (models.Link, error) {
	args := m.Called(linkId)
	return args.Get(0).(models.Link), args.Error(1)
}
func (m *mokShortenerStorageRepository) GetLinkByShortUrl(shortUrl string) (models.Link, error) {
	return models.Link{}, nil
}
func (m *mokShortenerStorageRepository) CreateLink(newLink *models.Link) (models.Link, error) {
	args := m.Called(newLink)
	return args.Get(0).(models.Link), args.Error(1)
}

func (m *mokShortenerStorageRepository) UpdateLink(newLinkWithExistingId *models.Link) (models.Link, error) {
	args := m.Called(newLinkWithExistingId)
	return args.Get(0).(models.Link), args.Error(1)
}
func (m *mokShortenerStorageRepository) DeleteLink(linkId uuid.UUID) error {
	args := m.Called(linkId)
	return args.Error(0)
}

func (m *mokShortenerStorageRepository) GetLinksCountByUserId(userId uuid.UUID) (int32, error) {
	args := m.Called(userId)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mokShortenerStorageRepository) GetLinkOwnerByShortLink(shortLink string) (string, error) {
	args := m.Called(shortLink)
	return args.Get(0).(string), args.Error(1)
}

func (m *mokShortenerStorageRepository) GetLinkIdByShortLink(shortLink string) (string, error) {
	args := m.Called(shortLink)
	return args.Get(0).(string), args.Error(1)
}

type mokShortenerSenderRepository struct{ mock.Mock }

func (m *mokShortenerSenderRepository) SendClick(ctx context.Context, click *models.Click) error {
	args := m.Called(ctx, click)
	return args.Error(0)
}

const (
	testdefaultLinkExpireTime = 1 * time.Second
)

func TestGetLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"
	targetUrl := "http://qwertysdijvnisdnc"
	createdAtStr := "2025-04-16T11:28:06+03:00"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	generated := false

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testGetLinkRequest := shortener.GetLinkRequest{
		LinkId: idStr,
	}

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		Generated: &generated,
		ShortLink: &shortLinkStr,
		TargetUrl: &targetUrl,
		CreatedAt: createdAtTime,
		ExpireAt:  &expireAtTime,
	}

	testShortenerStorageRepository.On("GetLinkById", mock.AnythingOfType("uuid.UUID")).Return(expectedModel, nil).Once()

	testGetLinkResponse, err := s.GetLink(ctx, &testGetLinkRequest)
	if err != nil {
		t.Fatalf("testGetLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Id.String(), testGetLinkResponse.LinkId)
	assert.Equal(t, expectedModel.UserId.String(), testGetLinkResponse.UserId)
	assert.Equal(t, *expectedModel.ShortLink, testGetLinkResponse.ShortLink)
	assert.Equal(t, *expectedModel.TargetUrl, testGetLinkResponse.TargetUrl)
	assert.Equal(t, timestamppb.New(expectedModel.CreatedAt), testGetLinkResponse.CreatedAt)
	assert.Equal(t, timestamppb.New(*expectedModel.ExpireAt), testGetLinkResponse.ExpireAt)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestCreateNewLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"
	targetUrl := "http://qwertysdijvnisdnc"

	testCreateLinkRequest := shortener.CreateLinkRequest{
		UserId:    userIDStr,
		ShortLink: &shortLinkStr,
		TargetUrl: "http://qwertysdijvnisdnc",
	}

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	expireAtStr := "2025-04-16T11:28:07+03:00"
	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	generated := false

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		Generated: &generated,
		ShortLink: &shortLinkStr,
		TargetUrl: &targetUrl,
		CreatedAt: time.Now(),
		ExpireAt:  &expireAtTime,
	}

	testShortenerStorageRepository.On("CreateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel, nil)

	testCreateNewLinkResponse, err := s.CreateNewLink(ctx, &testCreateLinkRequest)
	if err != nil {
		t.Fatalf("testCreateNewLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Id.String(), testCreateNewLinkResponse.LinkId)
	assert.Equal(t, expectedModel.UserId.String(), testCreateNewLinkResponse.UserId)
	assert.Equal(t, *expectedModel.ShortLink, testCreateNewLinkResponse.ShortLink)
	assert.Equal(t, *expectedModel.TargetUrl, testCreateNewLinkResponse.TargetUrl)
	assert.Equal(t, timestamppb.New(expectedModel.CreatedAt), testCreateNewLinkResponse.CreatedAt)
	assert.Equal(t, timestamppb.New(*expectedModel.ExpireAt), testCreateNewLinkResponse.ExpireAt)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestUpdateLink(t *testing.T) {

	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"
	targetUrl := "http://qwertysdijvnisdnc"
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

	newShortLink := utils.GenerateShortURL()

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		Generated: &regenerated,
		ShortLink: &newShortLink,
		TargetUrl: &targetUrl,
		CreatedAt: createdAtTime,
		ExpireAt:  &expireAtTime,
	}

	testShortenerStorageRepository.On("UpdateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel, nil)

	testUpdateLinkResponse, err := s.UpdateLink(ctx, &testUpdateLinkRequest)
	if err != nil {
		t.Fatalf("testUpdateLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Id.String(), testUpdateLinkResponse.LinkId)
	assert.Equal(t, expectedModel.UserId.String(), testUpdateLinkResponse.UserId)
	assert.Equal(t, *expectedModel.ShortLink, testUpdateLinkResponse.ShortLink)
	assert.Equal(t, *expectedModel.TargetUrl, testUpdateLinkResponse.TargetUrl)
	assert.Equal(t, timestamppb.New(expectedModel.CreatedAt), testUpdateLinkResponse.CreatedAt)
	assert.Equal(t, timestamppb.New(*expectedModel.ExpireAt), testUpdateLinkResponse.ExpireAt)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestDeleteLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"
	targetUrl := "http://qwertysdijvnisdnc"
	createdAtStr := "2025-04-16T11:28:06+03:00"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	generated := false

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testDeleteLinkRequest := shortener.DeleteLinkRequest{
		LinkId: idStr,
	}

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		Generated: &generated,
		ShortLink: &shortLinkStr,
		TargetUrl: &targetUrl,
		CreatedAt: createdAtTime,
		ExpireAt:  &expireAtTime,
	}

	testShortenerStorageRepository.On("GetLinkById", mock.AnythingOfType("uuid.UUID")).Return(expectedModel, nil).Once()
	testShortenerStorageRepository.On("DeleteLink", mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
	testShortenerCacheRepository.On("DeleteUrl", mock.AnythingOfType("string")).Return(nil).Once()

	testDeleteLinkResponse, err := s.DeleteLink(ctx, &testDeleteLinkRequest)
	if err != nil {
		t.Fatalf("testDeleteLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, true, testDeleteLinkResponse.Status)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestRedirect(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

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

	targetUrl := "https://ignuyruyrnfucufnwbrewunygb"

	testShortenerCacheRepository.On("GetUrl", shortLink).Return(targetUrl, nil).Once()
	testShortenerSenderRepository.On("SendClick", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Redirect(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, targetUrl, resp.TargetUrl)

	testShortenerCacheRepository.AssertExpectations(t)
}

func TestGetLinksCountByUserId(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"

	testGetLinksCountByUserIdRequest := shortener.GetLinksCountByUserIdRequest{
		UserId: idStr,
	}

	var expectedModel int32 = 7

	testShortenerStorageRepository.On("GetLinksCountByUserId", mock.AnythingOfType("uuid.UUID")).Return(expectedModel, nil)

	testGetLinksCountByUserIdResponse, err := s.GetLinksCountByUserId(ctx, &testGetLinksCountByUserIdRequest)
	if err != nil {
		t.Fatalf("testUpdateLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel, testGetLinksCountByUserIdResponse.Count)

	testShortenerStorageRepository.AssertExpectations(t)
}

func TestGetLinkOwnerByShortLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	shortLink := "qwerty"

	req := shortener.GetLinkOwnerByShortLinkRequest{ShortLink: shortLink}

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"

	testShortenerStorageRepository.On("GetLinkOwnerByShortLink", mock.AnythingOfType("string")).Return(userIDStr, nil)

	resp, err := s.GetLinkOwnerByShortLink(ctx, &req)
	if err != nil {
		t.Errorf("cant get link owner by short link:%v", err)
	}

	assert.Equal(t, userIDStr, resp.LinkOwner)

}

func TestGetLinkIdByShortLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	shortLink := "qwerty"

	req := shortener.GetLinkIdByShortLinkRequest{ShortLink: shortLink}

	linkId := "f9e71cb4-e1e1-4721-8eef-887338db7247"

	testShortenerStorageRepository.On("GetLinkIdByShortLink", mock.AnythingOfType("string")).Return(linkId, nil)

	resp, err := s.GetLinkIdByShortLink(ctx, &req)
	if err != nil {
		t.Errorf("cant get link owner by short link:%v", err)
	}

	assert.Equal(t, linkId, resp.LinkId)

}
