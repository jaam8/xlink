package unitTests

import (
	"context"
	"fmt"
	"shortener/internal/models"
	"shortener/internal/service"
	"shortener/internal/service/helper"
	"shortener/pkg/api/shortener"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mokShortenerCacheRepository struct{ mock.Mock }

func (m *mokShortenerCacheRepository) GetUrl(shortUrl string) (string, error)   { return "", nil }
func (m *mokShortenerCacheRepository) SetUrl(shortUrl string, url string) error { return nil }
func (m *mokShortenerCacheRepository) DeleteUrl(shortUrl string) error          { return nil }

type mokShortenerStorageRepository struct{ mock.Mock }

func (m *mokShortenerStorageRepository) GetLinkById(linkId uuid.UUID) (models.Link, error) {
	return models.Link{}, nil
}
func (m *mokShortenerStorageRepository) GetLinkByShortUrl(shortUrl string) (models.Link, error) {
	return models.Link{}, nil
}
func (m *mokShortenerStorageRepository) CreateLink(newLink *models.Link) (models.Link, error) {
	args := m.Called(newLink)
	return args.Get(0).(models.Link), args.Error(1)
}

func (m *mokShortenerStorageRepository) UpdateLink(newLinkWithExistingId *models.Link) (models.Link, error) {
	return models.Link{}, nil
}
func (m *mokShortenerStorageRepository) DeleteLink(linkId uuid.UUID) error { return nil }

type mokShortenerSenderRepository struct{ mock.Mock }

func (m *mokShortenerSenderRepository) SendRedirectInfo() {}

const (
	testdefaultLinkExpireTime = 100 * time.Millisecond
)

func TestCreateNewLink(t *testing.T) {
	testShortenerCacheRepository := new(mokShortenerCacheRepository)
	testShortenerStorageRepository := new(mokShortenerStorageRepository)
	testShortenerSenderRepository := new(mokShortenerSenderRepository)

	s := service.New(testShortenerCacheRepository, testShortenerStorageRepository, testShortenerSenderRepository, testdefaultLinkExpireTime)

	ctx := context.Background()

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireDuration := 1 * time.Second
	expireAtTime := time.Now().Add(expireDuration)
	expireAtStr := expireAtTime.Format(time.RFC3339)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	testCreateLinkRequest := shortener.CreateLinkRequest{
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
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		CreatedAt: time.Now(),
		ExpireAt:  expireAtTime,
	}

	testShortenerStorageRepository.On("CreateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel, nil)

	testCreateNewLinkResponse, err := s.CreateNewLink(ctx, &testCreateLinkRequest)
	if err != nil {
		t.Fatalf("testCreateNewLinkResponse is wrong, got:%v", err)
	}

	fmt.Println(*helper.LinkResponseFromLinkModel(expectedModel))

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.UserId.String(), testCreateNewLinkResponse.UserId)
	assert.Equal(t, expectedModel.GroupId.String(), *testCreateNewLinkResponse.GroupId)
	assert.Equal(t, expectedModel.ShortLink, testCreateNewLinkResponse.ShortLink)
	assert.Equal(t, expectedModel.Url, testCreateNewLinkResponse.Url)

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
	expireDuration := 1 * time.Second
	expireAtTime := time.Now().Add(expireDuration)
	expireAtStr := expireAtTime.Format(time.RFC3339)

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
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
		Url:       "http://qwertysdijvnisdnc",
		CreatedAt: time.Now(),
		ExpireAt:  expireAtTime,
	}
	fmt.Println(expectedModel)

	assert.Equal(t, userIDStr, "f9e71cb4-e1e1-4721-8eef-806338db7282")
	fmt.Println(userUUID)

	testShortenerStorageRepository.On("UpdateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel, nil)

	testUpdateLinkResponse, err := s.UpdateLink(ctx, &testUpdateLinkRequest)
	if err != nil {
		t.Fatalf("testUpdateLinkResponse is wrong, got:%v", err)
	}
	fmt.Println(*(&shortener.Link{
		Id:        expectedModel.Id.String(),
		UserId:    expectedModel.UserId.String(),
		GroupId:   &(expectedModel.GroupId.String()),
		Generated: expectedModel.Generated,
		ShortLink: expectedModel.ShortLink,
		Url:       expectedModel.Url,
		CreatedAt: expectedModel.CreatedAt.String(),
		ExpireAt:  expectedModel.ExpireAt.String(),
	}))

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Id.String(), testUpdateLinkResponse.Id)
	assert.Equal(t, expectedModel.UserId.String(), testUpdateLinkResponse.UserId)
	assert.Equal(t, expectedModel.GroupId.String(), *testUpdateLinkResponse.GroupId)
	assert.Equal(t, expectedModel.ShortLink, testUpdateLinkResponse.ShortLink)
	assert.Equal(t, expectedModel.Url, testUpdateLinkResponse.Url)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestDeleteLink(t *testing.T) {

}

func TestRedirect(t *testing.T) {

}
