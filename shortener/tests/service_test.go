package tests

import (
	"context"
	"testing"
	"time"
	"xlink/common/gen/shortener"
	"xlink/shortener/internal/models"
	"xlink/shortener/internal/service"
	"xlink/shortener/internal/service/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

type mokShortenerSenderRepository struct{ mock.Mock }

func (m *mokShortenerSenderRepository) SendRedirectInfo() {}

const (
	testdefaultLinkExpireTime = 1 * time.Second
)

func TestCreateNewLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx := context.Background()

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testCreateLinkRequest := shortener.CreateLinkRequest{
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	testCreateLinkRequest4 := shortener.CreateLinkRequest{
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: false,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	expectedModel := models.Link{
		Id:        uuid.New(),
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: true,
		ShortLink: utils.GenerateShortURL(),
		TargetUrl: "http://qwertysdijvnisdnc",
		CreatedAt: time.Now(),
		ExpireAt:  expireAtTime,
	}

	expectedModel4 := models.Link{
		Id:        uuid.New(),
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: testCreateLinkRequest4.Generated,
		ShortLink: utils.GenerateShortURL(),
		TargetUrl: "http://qwertysdijvnisdnc",
		CreatedAt: time.Now(),
		ExpireAt:  expireAtTime,
	}

	testShortenerStorageRepository.On("CreateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel, nil).Times(3)
	testShortenerStorageRepository.On("CreateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel4, nil)

	testCreateNewLinkResponse1, err := s.CreateNewLink(ctx, &testCreateLinkRequest)
	if err != nil {
		t.Fatalf("testCreateNewLinkResponse is wrong, got:%v", err)
	}
	testCreateNewLinkResponse2, err := s.CreateNewLink(ctx, &testCreateLinkRequest)
	if err != nil {
		t.Fatalf("testCreateNewLinkResponse is wrong, got:%v", err)
	}
	testCreateNewLinkResponse3, err := s.CreateNewLink(ctx, &testCreateLinkRequest)
	if err != nil {
		t.Fatalf("testCreateNewLinkResponse is wrong, got:%v", err)
	}
	testCreateNewLinkResponse4, err := s.CreateNewLink(ctx, &testCreateLinkRequest4)
	if err != nil {
		t.Fatalf("testCreateNewLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Id.String(), testCreateNewLinkResponse1.Id)
	assert.Equal(t, expectedModel.UserId.String(), testCreateNewLinkResponse1.UserId)
	assert.Equal(t, expectedModel.GroupId.String(), *testCreateNewLinkResponse1.GroupId)
	assert.Equal(t, expectedModel.Generated, testCreateNewLinkResponse1.Generated)
	assert.Equal(t, expectedModel.ShortLink, testCreateNewLinkResponse1.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testCreateNewLinkResponse1.Url)
	assert.Equal(t, expectedModel.CreatedAt.String(), testCreateNewLinkResponse1.CreatedAt)
	assert.Equal(t, expectedModel.ExpireAt.String(), testCreateNewLinkResponse1.ExpireAt)

	assert.Equal(t, expectedModel.Id.String(), testCreateNewLinkResponse2.Id)
	assert.Equal(t, expectedModel.UserId.String(), testCreateNewLinkResponse2.UserId)
	assert.Equal(t, expectedModel.GroupId.String(), *testCreateNewLinkResponse2.GroupId)
	assert.Equal(t, expectedModel.Generated, testCreateNewLinkResponse2.Generated)
	assert.Equal(t, expectedModel.ShortLink, testCreateNewLinkResponse2.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testCreateNewLinkResponse2.Url)
	assert.Equal(t, expectedModel.CreatedAt.String(), testCreateNewLinkResponse2.CreatedAt)
	assert.Equal(t, expectedModel.ExpireAt.String(), testCreateNewLinkResponse2.ExpireAt)

	assert.Equal(t, expectedModel.Id.String(), testCreateNewLinkResponse3.Id)
	assert.Equal(t, expectedModel.UserId.String(), testCreateNewLinkResponse3.UserId)
	assert.Equal(t, expectedModel.GroupId.String(), *testCreateNewLinkResponse3.GroupId)
	assert.Equal(t, expectedModel.Generated, testCreateNewLinkResponse3.Generated)
	assert.Equal(t, expectedModel.ShortLink, testCreateNewLinkResponse3.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testCreateNewLinkResponse3.Url)
	assert.Equal(t, expectedModel.CreatedAt.String(), testCreateNewLinkResponse3.CreatedAt)
	assert.Equal(t, expectedModel.ExpireAt.String(), testCreateNewLinkResponse3.ExpireAt)

	assert.Equal(t, expectedModel4.Id.String(), testCreateNewLinkResponse4.Id)
	assert.Equal(t, expectedModel4.UserId.String(), testCreateNewLinkResponse4.UserId)
	assert.Equal(t, expectedModel4.GroupId.String(), *testCreateNewLinkResponse4.GroupId)
	assert.Equal(t, expectedModel4.Generated, testCreateNewLinkResponse4.Generated)
	assert.Equal(t, expectedModel4.ShortLink, testCreateNewLinkResponse4.ShortLink)
	assert.Equal(t, expectedModel4.TargetUrl, testCreateNewLinkResponse4.Url)
	assert.Equal(t, expectedModel4.CreatedAt.String(), testCreateNewLinkResponse4.CreatedAt)
	assert.Equal(t, expectedModel4.ExpireAt.String(), testCreateNewLinkResponse4.ExpireAt)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestUpdateLink(t *testing.T) {

	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx := context.Background()

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	createdAtStr := "2025-04-16T11:28:06+03:00"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testUpdateLinkRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: true,
		ShortLink: "http://qwerty",
		TargetUrl: "http://qwertysdijvnisdnc",
		CreatedAt: createdAtTime,
		ExpireAt:  expireAtTime,
	}

	testShortenerStorageRepository.On("UpdateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel, nil)

	testUpdateLinkResponse, err := s.UpdateLink(ctx, &testUpdateLinkRequest)
	if err != nil {
		t.Fatalf("testUpdateLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Id.String(), testUpdateLinkResponse.Id)
	assert.Equal(t, expectedModel.UserId.String(), testUpdateLinkResponse.UserId)
	assert.Equal(t, expectedModel.GroupId.String(), *testUpdateLinkResponse.GroupId)
	assert.Equal(t, expectedModel.Generated, testUpdateLinkResponse.Generated)
	assert.Equal(t, expectedModel.ShortLink, testUpdateLinkResponse.ShortLink)
	assert.Equal(t, expectedModel.TargetUrl, testUpdateLinkResponse.Url)
	assert.Equal(t, expectedModel.CreatedAt.String(), testUpdateLinkResponse.CreatedAt)
	assert.Equal(t, expectedModel.ExpireAt.String(), testUpdateLinkResponse.ExpireAt)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestDeleteLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx := context.Background()

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	createdAtStr := "2025-04-16T11:28:06+03:00"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testDeleteLinkRequest := shortener.DeleteLinkRequest{
		Id: idStr,
	}

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: true,
		ShortLink: "http://qwerty",
		TargetUrl: "http://qwertysdijvnisdnc",
		CreatedAt: createdAtTime,
		ExpireAt:  expireAtTime,
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

	ctx := context.Background()

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	createdAtStr := "2025-04-16T11:28:06+03:00"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testRedirectRequest := shortener.Url{
		Url: "http://qwerty",
	}

	expectedModel := models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: true,
		ShortLink: "http://qwerty",
		TargetUrl: "http://qwertysdijvnisdnc",
		CreatedAt: createdAtTime,
		ExpireAt:  expireAtTime,
	}

	testShortenerCacheRepository.On("GetUrl", mock.AnythingOfType("string")).Return(expectedModel.TargetUrl, nil)
	//testShortenerSenderRepository.On("SendRedirectInfo").Return()

	testRedirectResponse, err := s.Redirect(ctx, &testRedirectRequest)
	if err != nil {
		t.Fatalf("testUpdateLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.TargetUrl, testRedirectResponse.Url)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestGetLinksCountByUserId(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx := context.Background()

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
