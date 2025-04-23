package tests

import (
	"testing"
	"time"
	"xlink/analytics/internal/models"
	"xlink/analytics/internal/service/helper"
	"xlink/common/gen/analytics"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClickEventToClick(t *testing.T) {

	linkOwner := uuid.New()
	clickedAt := time.Now()

	clickEvent := &models.ClickEvent{
		ShortLink:    "short123",
		ClickedAt:    clickedAt,
		Referrer:     "https://example.com",
		IPAddress:    "192.168.0.1",
		VisitorToken: "some-token",
		Browser:      "Chrome",
		DeviceType:   "Mobile",
		Os:           "Android",
	}

	country := "RU"
	region := "Moscow"

	click, err := helper.ClickEventToClick(clickEvent, linkOwner, country, region, true)
	assert.NoError(t, err)
	assert.Equal(t, linkOwner, click.LinkOwner)
	assert.Equal(t, "short123", click.ShortLink)
	assert.Equal(t, clickedAt, click.ClickedAt)
	assert.Equal(t, "https://example.com", click.Referrer)
	assert.Equal(t, "192.168.0.1", click.IPAddress)
	assert.Equal(t, "Chrome", click.Browser)
	assert.Equal(t, "Mobile", click.DeviceType)
	assert.Equal(t, "Android", click.Os)
	assert.Equal(t, "RU", click.Country)
	assert.Equal(t, "Moscow", click.Region)
	assert.False(t, click.IsUnique)

	click2, err := helper.ClickEventToClick(clickEvent, linkOwner, country, region, false)
	assert.NoError(t, err)
	assert.True(t, click2.IsUnique)
}

func TestToClicksByCountryResponse(t *testing.T) {

	stats := []models.CountryDayStats{
		{
			Date: time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Stats: []models.CountryClicks{
				{Country: "RU", Clicks: 100, UniqueClicks: 80},
				{Country: "US", Clicks: 50, UniqueClicks: 40},
			},
		},
	}

	expected := []*analytics.CountryDayStats{
		{
			Date: stats[0].Date,
			Stats: []*analytics.CountryClicks{
				{Country: "RU", Clicks: 100, UniqueClicks: 80},
				{Country: "US", Clicks: 50, UniqueClicks: 40},
			},
		},
	}

	result := helper.ToClicksByCountryResponse(stats)

	assert.Equal(t, expected, result)
}

func TestToClicksByRegionResponse(t *testing.T) {

	stats := []models.RegionDayStats{
		{
			Date: time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Stats: []models.RegionClicks{
				{Region: "Moscow", Clicks: 120, UniqueClicks: 90},
				{Region: "California", Clicks: 75, UniqueClicks: 60},
			},
		},
	}

	expected := []*analytics.RegionDayStats{
		{
			Date: stats[0].Date,
			Stats: []*analytics.RegionClicks{
				{Region: "Moscow", Clicks: 120, UniqueClicks: 90},
				{Region: "California", Clicks: 75, UniqueClicks: 60},
			},
		},
	}

	result := helper.ToClicksByRegionResponse(stats)

	assert.Equal(t, expected, result)
}

func TestToClicksByBrowserResponse(t *testing.T) {
	stats := []models.BrowserDayStats{
		{
			Date: time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Stats: []models.BrowserClicks{
				{Browser: "Chrome", Clicks: 200, UniqueClicks: 150},
				{Browser: "Firefox", Clicks: 100, UniqueClicks: 80},
			},
		},
	}

	expected := []*analytics.BrowserDayStats{
		{
			Date: stats[0].Date,
			Stats: []*analytics.BrowserClicks{
				{Browser: "Chrome", Clicks: 200, UniqueClicks: 150},
				{Browser: "Firefox", Clicks: 100, UniqueClicks: 80},
			},
		},
	}

	result := helper.ToClicksByBrowserResponse(stats)

	assert.Equal(t, expected, result)
}

func TestToClicksByDeviceTypeResponse(t *testing.T) {
	stats := []models.DeviceDayStats{
		{
			Date: time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Stats: []models.DeviceClicks{
				{DeviceType: "Mobile", Clicks: 300, UniqueClicks: 250},
				{DeviceType: "Desktop", Clicks: 150, UniqueClicks: 100},
			},
		},
	}

	expected := []*analytics.DeviceDayStats{
		{
			Date: stats[0].Date,
			Stats: []*analytics.DeviceClicks{
				{DeviceType: "Mobile", Clicks: 300, UniqueClicks: 250},
				{DeviceType: "Desktop", Clicks: 150, UniqueClicks: 100},
			},
		},
	}

	result := helper.ToClicksByDeviceTypeResponse(stats)

	assert.Equal(t, expected, result)
}

func TestToClicksByOSResponse(t *testing.T) {

	stats := []models.OSDayStats{
		{
			Date: time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Stats: []models.OSClicks{
				{OS: "iOS", Clicks: 500, UniqueClicks: 400},
				{OS: "Android", Clicks: 300, UniqueClicks: 250},
			},
		},
	}

	expected := []*analytics.OSDayStats{
		{
			Date: stats[0].Date,
			Stats: []*analytics.OSClicks{
				{Os: "iOS", Clicks: 500, UniqueClicks: 400},
				{Os: "Android", Clicks: 300, UniqueClicks: 250},
			},
		},
	}

	result := helper.ToClicksByOSResponse(stats)

	assert.Equal(t, expected, result)
}

func TestToClicksByHourResponse(t *testing.T) {

	stats := []models.HourDayStats{
		{
			Date: time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Stats: []models.HourStat{
				{Hour: 9, Clicks: 100, UniqueClicks: 80},
				{Hour: 10, Clicks: 150, UniqueClicks: 120},
			},
		},
	}

	expected := []*analytics.HourDayStats{
		{
			Date: stats[0].Date,
			Stats: []*analytics.HourStat{
				{Hour: 9, Clicks: 100, UniqueClicks: 80},
				{Hour: 10, Clicks: 150, UniqueClicks: 120},
			},
		},
	}

	result := helper.ToClicksByHourResponse(stats)

	assert.Equal(t, expected, result)
}

func TestToClicksByReferrerResponse(t *testing.T) {

	stats := []models.ReferrerDayStats{
		{
			Date: time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Stats: []models.ReferrerClicks{
				{Referrer: "google.com", Clicks: 100, UniqueClicks: 80},
				{Referrer: "facebook.com", Clicks: 150, UniqueClicks: 120},
			},
		},
	}

	expected := []*analytics.ReferrerDayStats{
		{
			Date: stats[0].Date,
			Stats: []*analytics.ReferrerClicks{
				{Referrer: "google.com", Clicks: 100, UniqueClicks: 80},
				{Referrer: "facebook.com", Clicks: 150, UniqueClicks: 120},
			},
		},
	}

	result := helper.ToClicksByReferrerResponse(stats)

	assert.Equal(t, expected, result)
}

func TestToClicksByDateResponse(t *testing.T) {

	stats := []models.DateStat{
		{
			Date:         time.Date(2025, 4, 23, 0, 0, 0, 0, time.UTC).String(),
			Clicks:       200,
			UniqueClicks: 150,
		},
		{
			Date:         time.Date(2025, 4, 24, 0, 0, 0, 0, time.UTC).String(),
			Clicks:       300,
			UniqueClicks: 250,
		},
	}

	expected := []*analytics.DateStat{
		{
			Date:         stats[0].Date,
			Clicks:       200,
			UniqueClicks: 150,
		},
		{
			Date:         stats[1].Date,
			Clicks:       300,
			UniqueClicks: 250,
		},
	}

	result := helper.ToClicksByDateResponse(stats)

	assert.Equal(t, expected, result)
}
