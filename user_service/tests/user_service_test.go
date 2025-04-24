package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"xlink/common/gen/user_service"
	"xlink/common/logger"
	"xlink/user_service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUsersCacheRepository struct {
	tokens map[string]string
	roles  map[string]string
}

func NewMockUsersCacheRepository() *mockUsersCacheRepository {
	return &mockUsersCacheRepository{
		tokens: make(map[string]string),
		roles:  make(map[string]string),
	}
}

func (m *mockUsersCacheRepository) CheckToken(userId string, token string) (bool, error) {
	storedToken, ok := m.tokens[userId]
	if !ok {
		return false, nil
	}
	return storedToken == token, nil
}

func (m *mockUsersCacheRepository) SetToken(userId string, token string) error {
	m.tokens[userId] = token
	return nil
}

func (m *mockUsersCacheRepository) GetRole(userId string) (string, bool, bool, error) {
	data, ok := m.roles[userId]
	if !ok {
		return "", false, false, fmt.Errorf("roles not found")
	}
	parts := strings.Split(data, ";")
	if len(parts) != 2 {
		return "", false, false, fmt.Errorf("invalid role data format")
	}
	isStaff := parts[0] == "1"
	isAdmin := parts[1] == "1"
	role := getRoleByIsStaffIsAdmin(isStaff, isAdmin)
	return role, isStaff, isAdmin, nil
}

func (m *mockUsersCacheRepository) SetRole(userId string, isStaff bool, isAdmin bool) error {
	roleString := fmt.Sprintf("%s;%s", boolToString(isStaff), boolToString(isAdmin))
	m.roles[userId] = roleString
	return nil
}

// вспомогательные функции
func boolToString(val bool) string {
	if val {
		return "1"
	}
	return "0"
}

func getRoleByIsStaffIsAdmin(isStaff, isAdmin bool) string {
	switch {
	case isAdmin:
		return "admin"
	case isStaff:
		return "staff"
	default:
		return "user"
	}
}

type mockUserStorageRepository struct {
	mock.Mock
}

func (m *mockUserStorageRepository) CheckToken(userId string, token string) (bool, error) {
	args := m.Called(userId, token)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserStorageRepository) RefreshToken(userId string, token string) (string, error) {
	args := m.Called(userId, token)
	return args.String(0), args.Error(1)
}

func (m *mockUserStorageRepository) CreateUser(telegramId *int64, isStaff *bool, isAdmin *bool) (string, string, error) {
	args := m.Called(telegramId, isStaff, isAdmin)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockUserStorageRepository) GetUser(userId string) (string, string, *int64, error) {
	args := m.Called(userId)
	var tgId *int64
	if args.Get(2) != nil {
		tgId = args.Get(2).(*int64)
	}
	return args.String(0), args.String(1), tgId, args.Error(3)
}

func (m *mockUserStorageRepository) GetUserIDByToken(token string) (string, bool, error) {
	args := m.Called(token)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *mockUserStorageRepository) GetUserIDByTgId(tgId int64) (string, bool, error) {
	args := m.Called(tgId)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *mockUserStorageRepository) UpdateUser(userId string, telegramId *int64, isStaff *bool, isAdmin *bool) (bool, error) {
	args := m.Called(userId, telegramId, isStaff, isAdmin)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserStorageRepository) DeleteUser(userId string) (bool, error) {
	args := m.Called(userId)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserStorageRepository) GetRole(userId string) (string, bool, bool, error) {
	args := m.Called(userId)
	return args.String(0), args.Bool(1), args.Bool(2), args.Error(3)
}

func (m *mockUserStorageRepository) GetTokenByTgId(tgId int64) (string, error) {
	args := m.Called(tgId)
	return args.String(0), args.Error(1)
}

type mockShortenerRepository struct {
	mock.Mock
}

func (m *mockShortenerRepository) GetLinksCountByUserId(userId string) (int32, error) {
	args := m.Called(userId)
	return args.Get(0).(int32), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	tgId := int64(123456)
	isStaff := true
	isAdmin := false

	expectedUserID := "userid123"
	expectedToken := "token123"

	storageRepo.On("CreateUser", &tgId, &isStaff, &isAdmin).Return(expectedUserID, expectedToken, nil)
	req := &user_service.CreateUserRequest{
		TgId:    &tgId,
		IsStaff: &isStaff,
		IsAdmin: &isAdmin,
	}

	resp, err := svc.CreateUser(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.UserId != expectedUserID || resp.Token != expectedToken {
		t.Errorf("unexpected response: got userId=%s, token=%s", resp.UserId, resp.Token)
	}
}

func TestGetUser(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	userId := "user123"

	req := &user_service.GetUserRequest{
		UserId: userId,
	}

	expectedUserId := userId
	expectedRole := "Stuff"
	expectedTgId := new(int64)
	*expectedTgId = 123456
	expectedLinkCount := int32(5)

	storageRepo.On("GetUser", userId).Return(expectedUserId, expectedRole, expectedTgId, nil)
	shortenerRepo.On("GetLinksCountByUserId", userId).Return(expectedLinkCount, nil)

	resp, err := svc.GetUser(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, expectedUserId, resp.UserId)
	assert.Equal(t, expectedRole, resp.Role)
	assert.Equal(t, expectedTgId, resp.TgId)
	assert.Equal(t, expectedLinkCount, resp.LinkCount)
}

func TestGetUserIDByToken(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	token := "token321"

	expectedUserID := "user99"
	expectedStatus := true

	storageRepo.On("GetUserIDByToken", token).Return(expectedUserID, expectedStatus, nil)

	req := &user_service.GetUserIDByTokenRequest{Token: token}

	resp, err := svc.GetUserIDByToken(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, resp.UserId, expectedUserID)
	assert.Equal(t, resp.Status, expectedStatus)
}
func TestGetUserIDByTgID(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	tgId := int64(123456)

	expectedUserID := "user99"
	expectedStatus := true

	storageRepo.On("GetUserIDByTgId", tgId).Return(expectedUserID, expectedStatus, nil)

	req := &user_service.GetUserIDByTgIDRequest{TgId: tgId}

	resp, err := svc.GetUserIDByTgID(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, resp.UserId, expectedUserID)
	assert.Equal(t, resp.Status, expectedStatus)
}

func TestUpdateUser(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	userId := "user123"
	tgId := int64(987654321)
	isStaff := true
	isAdmin := false

	req := &user_service.UpdateUserRequest{
		UserId:  userId,
		TgId:    &tgId,
		IsStaff: &isStaff,
		IsAdmin: &isAdmin,
	}

	expectedStatus := true

	storageRepo.On("UpdateUser", userId, &tgId, &isStaff, &isAdmin).Return(expectedStatus, nil)

	resp, err := svc.UpdateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, resp.Status)

	storageRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	userId := "user123"
	expectedStatus := true

	req := &user_service.DeleteUserRequest{UserId: userId}

	storageRepo.On("DeleteUser", userId).Return(expectedStatus, nil)

	resp, err := svc.DeleteUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, resp.Status)

	storageRepo.AssertExpectations(t)
}

func TestGetRole(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	userId := "user123"
	setrolo_error := cacheRepo.SetRole(userId, true, false)
	if setrolo_error != nil {
		t.Errorf("failed to set role:%v", err)
	}

	req := &user_service.GetRoleRequest{UserId: userId}

	resp, err := svc.GetRole(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "staff", resp.Role)
	assert.True(t, resp.IsStaff)
	assert.False(t, resp.IsAdmin)
}

func TestGetTokenByTgId(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	s := service.New(cacheRepo, storageRepo, shortenerRepo)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	tgId := int64(123456)

	req := user_service.GetTokenByTgIdRequest{TgId: tgId}

	token := "h384ntv873hc8wh8c7th3"

	storageRepo.On("GetTokenByTgId", tgId).Return(token, nil)

	resp, err := s.GetTokenByTgId(ctx, &req)
	if err != nil {
		t.Errorf("failed to get response from GetTokenByTgId:%v", err)
	}

	assert.Equal(t, token, resp.Token)
}

func TestCheckToken_EmptyUserId(t *testing.T) {
	svc := service.New(NewMockUsersCacheRepository(), new(mockUserStorageRepository), new(mockShortenerRepository))

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	req := &user_service.TokenCheckRequest{
		UserId: "",
		Token:  "token123",
	}
	resp, err := svc.CheckToken(ctx, req)

	assert.Error(t, err)
	assert.False(t, resp.Status)
	assert.Contains(t, err.Error(), "user_id is empty")
}

func TestCheckToken_EmptyToken(t *testing.T) {
	svc := service.New(NewMockUsersCacheRepository(), new(mockUserStorageRepository), new(mockShortenerRepository))

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	req := &user_service.TokenCheckRequest{
		UserId: "user123",
		Token:  "",
	}
	resp, err := svc.CheckToken(ctx, req)

	assert.Error(t, err)
	assert.False(t, resp.Status)
	assert.Contains(t, err.Error(), "token is empty")
}

func TestCheckToken_FromCache(t *testing.T) {
	cache := NewMockUsersCacheRepository()
	storage := new(mockUserStorageRepository)
	shortener := new(mockShortenerRepository)

	userId := "user123"
	token := "token123"
	err := cache.SetToken(userId, token)
	if err != nil {
		t.Errorf("failed to set role:%v", err)
	}

	svc := service.New(cache, storage, shortener)

	req := &user_service.TokenCheckRequest{
		UserId: userId,
		Token:  token,
	}

	loggerCtx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("Error getting new logger%v", err)
	}

	resp, err := svc.CheckToken(loggerCtx, req)

	assert.NoError(t, err)
	assert.True(t, resp.Status)
}

func TestRefreshToken(t *testing.T) {
	cacheRepo := NewMockUsersCacheRepository()
	storageRepo := new(mockUserStorageRepository)
	shortenerRepo := new(mockShortenerRepository)

	svc := service.New(cacheRepo, storageRepo, shortenerRepo)

	userId := "user123"
	oldToken := "old_token_123"
	newToken := "new_token_456"

	storageRepo.On("RefreshToken", userId, oldToken).Return(newToken, nil)

	loggerCtx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("Error getting new logger%v", err)
	}

	req := &user_service.RefreshTokenRequest{
		UserId: userId,
		Token:  oldToken,
	}

	resp, err := svc.RefreshToken(loggerCtx, req)
	assert.NoError(t, err)
	assert.Equal(t, newToken, resp.Token)

	cachedTokenCorrect, err := cacheRepo.CheckToken(userId, newToken)
	assert.NoError(t, err)
	assert.True(t, cachedTokenCorrect)

	storageRepo.AssertExpectations(t)
}
