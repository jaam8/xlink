package unitTests

import (
	"context"
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

	testShortenerStorageRepository.On("UpdateLink", mock.AnythingOfType("*models.Link")).Return(expectedModel, nil)

	testUpdateLinkResponse, err := s.UpdateLink(ctx, &testUpdateLinkRequest)
	if err != nil {
		t.Fatalf("testUpdateLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Id.String(), testUpdateLinkResponse.Id)
	assert.Equal(t, expectedModel.UserId.String(), testUpdateLinkResponse.UserId)
	assert.Equal(t, expectedModel.GroupId.String(), *testUpdateLinkResponse.GroupId)
	assert.Equal(t, expectedModel.ShortLink, testUpdateLinkResponse.ShortLink)
	assert.Equal(t, expectedModel.Url, testUpdateLinkResponse.Url)

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
	expireDuration := 1 * time.Second
	expireAtTime := time.Now().Add(expireDuration)

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
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
		Url:       "http://qwertysdijvnisdnc",
		CreatedAt: time.Now(),
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
	expireDuration := 1 * time.Second
	expireAtTime := time.Now().Add(expireDuration)

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
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
		Url:       "http://qwertysdijvnisdnc",
		CreatedAt: time.Now(),
		ExpireAt:  expireAtTime,
	}

	testShortenerCacheRepository.On("GetUrl", mock.AnythingOfType("string")).Return(expectedModel.Url, nil)
	//testShortenerSenderRepository.On("SendRedirectInfo").Return()

	testRedirectResponse, err := s.Redirect(ctx, &testRedirectRequest)
	if err != nil {
		t.Fatalf("testUpdateLinkResponse is wrong, got:%v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedModel.Url, testRedirectResponse.Url)

	testShortenerStorageRepository.AssertExpectations(t)

}

func TestGetValidatedId(t *testing.T) {
	testRequest := shortener.DeleteLinkRequest{Id: "f9e71cb4-e1e1-4721-8eef-806338db2222"}

	testResponse, err := helper.GetValidatedId(&testRequest)
	if err != nil {
		t.Fatalf("Unable to get id, %v", err)
	}

	expected, err := uuid.Parse("f9e71cb4-e1e1-4721-8eef-806338db2222")
	if err != nil {
		t.Fatalf("Unable to parse id, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestGetValidatedUserId(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireDuration := 1 * time.Second
	expireAtTime := time.Now().Add(expireDuration).Format(time.RFC3339)

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtTime,
	}

	testResponse, err := helper.GetValidatedUserId(&testRequest)
	if err != nil {
		t.Fatalf("Unable to get id, %v", err)
	}

	expected, err := uuid.Parse(userIDStr)
	if err != nil {
		t.Fatalf("Unable to parse id, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestGetValidatedGroupId(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireDuration := 1 * time.Second
	expireAtTime := time.Now().Add(expireDuration).Format(time.RFC3339)

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtTime,
	}

	testResponse, err := helper.GetValidatedGroupId(&testRequest, nil)
	if err != nil {
		t.Fatalf("Unable to get id, %v", err)
	}

	parsed, err := uuid.Parse(groupIDStr)
	expected := &parsed
	if err != nil {
		t.Fatalf("Unable to parse id, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestGetValidatedExpireAt(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtTime := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtTime,
	}

	testResponse, err := helper.GetValidatedExpireAt(&testRequest, time.Now())
	if err != nil {
		t.Fatalf("Unable to get ExpireAt, %v", err)
	}

	expected, err := time.Parse(time.RFC3339, expireAtTime)
	if err != nil {
		t.Fatalf("Unable to parse ExpireAt, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestValidatedStringNotEmpty(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtTime := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtTime,
	}

	testerror := helper.ValidateStringNotEmpty(testRequest.GetShortLink())
	if testerror != nil {
		t.Fatalf("ValidateStringNotEmpty error, %v", testerror)
	}

	assert.NoError(t, testerror)

}

func TestValidatedUrl(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtTime := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtTime,
	}

	testerror := helper.ValidateUrl(testRequest.GetUrl())
	if testerror != nil {
		t.Fatalf("ValidateUrl error, %v", testerror)
	}

	assert.NoError(t, testerror)

}

func TestLinkModelFromLinkRequest(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtTime := "2025-04-16T11:28:07+03:00"
	expireAtAsTime, _ := time.Parse(time.RFC3339, expireAtTime)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtTime,
	}

	expectedModel := &models.Link{
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: false,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  expireAtAsTime,
	}

	testResponse, err := helper.LinkModelFromLinkRequest(&testRequest, time.Now())
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel, testResponse)

}

func TestLinkModelFromLinkRequestWithId(t *testing.T) {
	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtTime := "2025-04-16T11:28:07+03:00"
	expireAtAsTime, _ := time.Parse(time.RFC3339, expireAtTime)
	var CreatedAt time.Time

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtTime,
	}

	expectedModel := &models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: false,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		CreatedAt: CreatedAt,
		ExpireAt:  expireAtAsTime,
	}

	testResponse, err := helper.LinkModelFromLinkRequestWithId(&testRequest, time.Now())
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel, testResponse)

}
