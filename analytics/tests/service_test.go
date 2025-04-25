package tests

import (
	"context"
	"testing"
	"time"
	"xlink/analytics/internal/models"
	"xlink/analytics/internal/service"
	"xlink/common/gen/analytics"
	"xlink/common/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockCacheAdapter struct {
	mock.Mock
}

func (m *MockCacheAdapter) CheckVisitorToken(visitorToken, shortLink string) (bool, error) {
	args := m.Called(visitorToken, shortLink)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheAdapter) SetVisitorToken(visitorToken, shortLink string) error {
	args := m.Called(visitorToken, shortLink)
	return args.Error(0)
}

// StorageAdapter mock
type MockStorageAdapter struct {
	mock.Mock
}

func (m *MockStorageAdapter) SaveClicks(ctx context.Context, clicks []*models.Click) error {
	args := m.Called(ctx, clicks)
	return args.Error(0)
}

func (m *MockStorageAdapter) GetClicksByCountry(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.CountryDayStats, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.CountryDayStats), args.Error(1)
}

func (m *MockStorageAdapter) GetClicksByRegion(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.RegionDayStats, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.RegionDayStats), args.Error(1)
}

func (m *MockStorageAdapter) GetClicksByReferrer(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.ReferrerDayStats, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.ReferrerDayStats), args.Error(1)
}

func (m *MockStorageAdapter) GetClicksByBrowser(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.BrowserDayStats, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.BrowserDayStats), args.Error(1)
}

func (m *MockStorageAdapter) GetClicksByOS(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.OSDayStats, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.OSDayStats), args.Error(1)
}

func (m *MockStorageAdapter) GetClicksByDeviceType(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.DeviceDayStats, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.DeviceDayStats), args.Error(1)
}

func (m *MockStorageAdapter) GetClicksByHour(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.HourDayStats, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.HourDayStats), args.Error(1)
}

func (m *MockStorageAdapter) GetClicksByDate(startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string) ([]models.DateStat, error) {
	args := m.Called(startDate, endDate, linkOwner, shortLink)
	return args.Get(0).([]models.DateStat), args.Error(1)
}

// ConsumerAdapter mock
type MockConsumerAdapter struct {
	mock.Mock
}

func (m *MockConsumerAdapter) ConsumeClickEvent(ctx context.Context) (*models.ClickEvent, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.ClickEvent), args.Error(1)
}

// ShortenerAdapter mock
type MockShortenerAdapter struct {
	mock.Mock
}

func (m *MockShortenerAdapter) GetLinkOwner(shortLink string) (uuid.UUID, error) {
	args := m.Called(shortLink)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func TestHandleConsumer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, err := logger.New(ctx)
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	event := &models.ClickEvent{
		ShortLink:    "abc123",
		ClickedAt:    time.Now(),
		Referrer:     "google.com",
		IPAddress:    "8.8.8.8",
		VisitorToken: "token123",
		Browser:      "Chrome",
		DeviceType:   "Mobile",
		Os:           "Android",
	}
	linkOwner := uuid.New()

	mockConsumer.On("ConsumeClickEvent", ctx).Return(event, nil)
	mockShortener.On("GetLinkOwner", event.ShortLink).Return(linkOwner, nil)
	mockCache.On("CheckVisitorToken", event.VisitorToken, event.ShortLink).Return(false, nil)
	mockCache.On("SetVisitorToken", event.VisitorToken, event.ShortLink).Return(nil)

	mockStorage.On("SaveClicks", mock.Anything, mock.AnythingOfType("[]*models.Click")).Return(nil).Run(func(args mock.Arguments) {
		cancel()
	})

	service.HandleConsumer(ctx, 1, time.Second)

	require.Equal(t, context.Canceled, ctx.Err())
}

func TestClicksByCountry(t *testing.T) {

	mockStorage := new(MockStorageAdapter)
	mockShortener := new(MockShortenerAdapter)
	mockCache := new(MockCacheAdapter)
	mockConsumer := new(MockConsumerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	linkOwner := uuid.New()
	shortLink := "abc123"

	req := analytics.GetClicksRequest{
		ShortLink: shortLink,
		StartDate: timestamppb.New(time.Now().Add(-24 * time.Hour)),
		EndDate:   timestamppb.New(time.Now()),
		LinkOwner: linkOwner.String(),
	}

	mockShortener.On("GetLinkOwner", shortLink).Return(linkOwner, nil)

	expectedData := []models.CountryDayStats{
		{
			Date: "2025-04-24",
			Stats: []models.CountryClicks{
				{Country: "RU", Clicks: 42},
				{Country: "US", Clicks: 13},
			},
		},
	}

	mockStorage.
		On("GetClicksByCountry", mock.Anything, mock.Anything, linkOwner, shortLink).
		Return(expectedData, nil)

	resp, err := service.ClicksByCountry(ctx, &req)

	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "2025-04-24", resp.Data[0].Date)
	require.Len(t, resp.Data[0].Stats, 2)
	require.Equal(t, "RU", resp.Data[0].Stats[0].Country)
	require.Equal(t, uint64(42), resp.Data[0].Stats[0].Clicks)
	require.Equal(t, "US", resp.Data[0].Stats[1].Country)
	require.Equal(t, uint64(13), resp.Data[0].Stats[1].Clicks)
}

func TestClicksByRegion(t *testing.T) {

	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	linkOwner := uuid.New()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

	request := &analytics.GetClicksRequest{
		ShortLink: "abc123",
		LinkOwner: linkOwner.String(),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	mockShortener.On("GetLinkOwner", request.GetShortLink()).Return(linkOwner, nil)

	mockStorage.On("GetClicksByRegion", &startDate, &endDate, linkOwner, request.GetShortLink()).Return([]models.RegionDayStats{
		{
			Date: "2024-01-01",
			Stats: []models.RegionClicks{
				{Region: "Moscow", Clicks: 42},
				{Region: "California", Clicks: 13},
			},
		},
	}, nil)

	resp, err := service.ClicksByRegion(ctx, request)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "2024-01-01", resp.Data[0].Date)
	require.Equal(t, "Moscow", resp.Data[0].Stats[0].Region)
	require.Equal(t, uint64(42), resp.Data[0].Stats[0].Clicks)
	require.Equal(t, "California", resp.Data[0].Stats[1].Region)
	require.Equal(t, uint64(13), resp.Data[0].Stats[1].Clicks)
}

func TestClicksByBrowser(t *testing.T) {

	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	linkOwner := uuid.New()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	request := &analytics.GetClicksRequest{
		ShortLink: "abc123",
		LinkOwner: linkOwner.String(),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	mockShortener.On("GetLinkOwner", request.GetShortLink()).Return(linkOwner, nil)

	mockStorage.On("GetClicksByBrowser", &startDate, &endDate, linkOwner, request.GetShortLink()).Return([]models.BrowserDayStats{
		{
			Date: "2024-01-01",
			Stats: []models.BrowserClicks{
				{Browser: "Chrome", Clicks: 100},
				{Browser: "Firefox", Clicks: 25},
			},
		},
	}, nil)

	resp, err := service.ClicksByBrowser(ctx, request)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)

	require.Equal(t, "2024-01-01", resp.Data[0].Date)
	require.Len(t, resp.Data[0].Stats, 2)
	require.Equal(t, "Chrome", resp.Data[0].Stats[0].Browser)
	require.Equal(t, uint64(100), resp.Data[0].Stats[0].Clicks)
	require.Equal(t, "Firefox", resp.Data[0].Stats[1].Browser)
	require.Equal(t, uint64(25), resp.Data[0].Stats[1].Clicks)
}

func TestClicksByOS_Success(t *testing.T) {

	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	linkOwner := uuid.New()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	request := &analytics.GetClicksRequest{
		ShortLink: "abc123",
		LinkOwner: linkOwner.String(),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	mockShortener.On("GetLinkOwner", request.GetShortLink()).Return(linkOwner, nil)

	mockStorage.On("GetClicksByOS", &startDate, &endDate, linkOwner, request.GetShortLink()).Return([]models.OSDayStats{
		{
			Date: "2024-01-01",
			Stats: []models.OSClicks{
				{OS: "Android", Clicks: 120},
				{OS: "iOS", Clicks: 80},
			},
		},
	}, nil)

	resp, err := service.ClicksByOS(ctx, request)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)

	require.Equal(t, "2024-01-01", resp.Data[0].Date)
	require.Len(t, resp.Data[0].Stats, 2)
	require.Equal(t, "Android", resp.Data[0].Stats[0].Os)
	require.Equal(t, uint64(120), resp.Data[0].Stats[0].Clicks)
	require.Equal(t, "iOS", resp.Data[0].Stats[1].Os)
	require.Equal(t, uint64(80), resp.Data[0].Stats[1].Clicks)
}

func TestClicksByDeviceType(t *testing.T) {

	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	linkOwner := uuid.New()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	request := &analytics.GetClicksRequest{
		ShortLink: "abc123",
		LinkOwner: linkOwner.String(),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	mockShortener.On("GetLinkOwner", request.GetShortLink()).Return(linkOwner, nil)

	mockStorage.On("GetClicksByDeviceType", &startDate, &endDate, linkOwner, request.GetShortLink()).Return([]models.DeviceDayStats{
		{
			Date: "2024-01-01",
			Stats: []models.DeviceClicks{
				{DeviceType: "Mobile", Clicks: 150},
				{DeviceType: "Desktop", Clicks: 200},
			},
		},
	}, nil)

	resp, err := service.ClicksByDeviceType(ctx, request)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)

	require.Equal(t, "2024-01-01", resp.Data[0].Date)
	require.Len(t, resp.Data[0].Stats, 2)
	require.Equal(t, "Mobile", resp.Data[0].Stats[0].DeviceType)
	require.Equal(t, uint64(150), resp.Data[0].Stats[0].Clicks)
	require.Equal(t, "Desktop", resp.Data[0].Stats[1].DeviceType)
	require.Equal(t, uint64(200), resp.Data[0].Stats[1].Clicks)
}

func TestClicksByHour(t *testing.T) {

	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	linkOwner := uuid.New()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 23, 59, 59, 0, time.UTC)

	request := &analytics.GetClicksRequest{
		ShortLink: "abc123",
		LinkOwner: linkOwner.String(),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	mockShortener.On("GetLinkOwner", request.GetShortLink()).Return(linkOwner, nil)

	mockStorage.On("GetClicksByHour", &startDate, &endDate, linkOwner, request.GetShortLink()).Return([]models.HourDayStats{
		{
			Date: "2025-01-01",
			Stats: []models.HourStat{
				{Hour: 0, Clicks: 100, UniqueClicks: 80},
				{Hour: 1, Clicks: 150, UniqueClicks: 120},
				{Hour: 2, Clicks: 120, UniqueClicks: 90},
			},
		},
	}, nil)

	resp, err := service.ClicksByHour(ctx, request)
	require.NoError(t, err)
	require.Len(t, resp.Stats, 1) // Один день в результате

	// Проверка данных для первого дня
	require.Equal(t, "2025-01-01", resp.Stats[0].Date)
	require.Len(t, resp.Stats[0].Stats, 3) // Три часа

	// Проверка данных для первого часа
	require.Equal(t, uint32(0), resp.Stats[0].Stats[0].Hour)
	require.Equal(t, uint64(100), resp.Stats[0].Stats[0].Clicks)
	require.Equal(t, uint64(80), resp.Stats[0].Stats[0].UniqueClicks)

	// Проверка данных для второго часа
	require.Equal(t, uint32(1), resp.Stats[0].Stats[1].Hour)
	require.Equal(t, uint64(150), resp.Stats[0].Stats[1].Clicks)
	require.Equal(t, uint64(120), resp.Stats[0].Stats[1].UniqueClicks)

	// Проверка данных для третьего часа
	require.Equal(t, uint32(2), resp.Stats[0].Stats[2].Hour)
	require.Equal(t, uint64(120), resp.Stats[0].Stats[2].Clicks)
	require.Equal(t, uint64(90), resp.Stats[0].Stats[2].UniqueClicks)
}

func TestClicksByDate(t *testing.T) {
	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	linkOwner := uuid.New()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 23, 59, 59, 0, time.UTC)

	request := &analytics.GetClicksRequest{
		ShortLink: "abc123",
		LinkOwner: linkOwner.String(),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	mockShortener.On("GetLinkOwner", request.GetShortLink()).Return(linkOwner, nil)

	mockStorage.On("GetClicksByDate", &startDate, &endDate, linkOwner, request.GetShortLink()).Return([]models.DateStat{
		{
			Date:         "2025-01-01",
			Clicks:       100,
			UniqueClicks: 80},
	}, nil)

	resp, err := service.ClicksByDate(ctx, request)
	require.NoError(t, err)
	require.Len(t, resp.Stats, 1) // Один день в результате

	// Проверка данных для первого дня
	require.Equal(t, "2025-01-01", resp.Stats[0].Date)
	require.Equal(t, uint64(100), resp.Stats[0].Clicks)
	require.Equal(t, uint64(80), resp.Stats[0].UniqueClicks)
}

func TestClicksByReferrer(t *testing.T) {
	mockCache := new(MockCacheAdapter)
	mockStorage := new(MockStorageAdapter)
	mockConsumer := new(MockConsumerAdapter)
	mockShortener := new(MockShortenerAdapter)

	service := service.New(mockCache, mockStorage, mockConsumer, mockShortener)

	ctx, err := logger.New(context.Background())
	require.NoError(t, err)

	linkOwner := uuid.New()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

	request := &analytics.GetClicksRequest{
		ShortLink: "abc123",
		LinkOwner: linkOwner.String(),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	mockShortener.On("GetLinkOwner", request.GetShortLink()).Return(linkOwner, nil)

	mockStorage.On("GetClicksByReferrer", &startDate, &endDate, linkOwner, request.GetShortLink()).Return([]models.ReferrerDayStats{
		{
			Date: "2025-01-01",
			Stats: []models.ReferrerClicks{
				{Referrer: "google.com", Clicks: 120, UniqueClicks: 100},
				{Referrer: "twitter.com", Clicks: 80, UniqueClicks: 75},
			},
		},
	}, nil)

	resp, err := service.ClicksByReferrer(ctx, request)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)

	require.Equal(t, "2025-01-01", resp.Data[0].Date)
	require.Len(t, resp.Data[0].Stats, 2)

	require.Equal(t, "google.com", resp.Data[0].Stats[0].Referrer)
	require.Equal(t, uint64(120), resp.Data[0].Stats[0].Clicks)
	require.Equal(t, uint64(100), resp.Data[0].Stats[0].UniqueClicks)

	require.Equal(t, "twitter.com", resp.Data[0].Stats[1].Referrer)
	require.Equal(t, uint64(80), resp.Data[0].Stats[1].Clicks)
	require.Equal(t, uint64(75), resp.Data[0].Stats[1].UniqueClicks)
}
