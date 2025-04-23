package helper

import (
	"github.com/google/uuid"
	"xlink/analytics/internal/models"
	"xlink/common/gen/analytics"
)

func ClickEventToClick(clickEvent *models.ClickEvent, linkOwner uuid.UUID,
	country, region string, hasToken bool) (*models.Click, error) {

	var click models.Click
	click.LinkOwner = linkOwner
	click.ShortLink = clickEvent.ShortLink
	click.ClickedAt = clickEvent.ClickedAt
	click.Referrer = clickEvent.Referrer
	click.Country = country
	click.Region = region
	click.Browser = clickEvent.Browser
	click.DeviceType = clickEvent.DeviceType
	click.Os = clickEvent.Os
	click.IPAddress = clickEvent.IPAddress
	click.IsUnique = !hasToken

	return &click, nil
}

func ToClicksByCountryResponse(stats []models.CountryDayStats) []*analytics.CountryDayStats {
	respData := make([]*analytics.CountryDayStats, 0, len(stats))

	for _, day := range stats {
		protoClicks := make([]*analytics.CountryClicks, 0, len(day.Stats))
		for _, cc := range day.Stats {
			protoClicks = append(protoClicks, &analytics.CountryClicks{
				Country:      cc.Country,
				Clicks:       cc.Clicks,
				UniqueClicks: cc.UniqueClicks,
			})
		}

		respData = append(respData, &analytics.CountryDayStats{
			Date:  day.Date,
			Stats: protoClicks,
		})
	}

	return respData
}

func ToClicksByRegionResponse(stats []models.RegionDayStats) []*analytics.RegionDayStats {
	respData := make([]*analytics.RegionDayStats, 0, len(stats))

	for _, day := range stats {
		protoClicks := make([]*analytics.RegionClicks, 0, len(day.Stats))
		for _, cc := range day.Stats {
			protoClicks = append(protoClicks, &analytics.RegionClicks{
				Region:       cc.Region,
				Clicks:       cc.Clicks,
				UniqueClicks: cc.UniqueClicks,
			})
		}

		respData = append(respData, &analytics.RegionDayStats{
			Date:  day.Date,
			Stats: protoClicks,
		})
	}

	return respData
}

func ToClicksByBrowserResponse(stats []models.BrowserDayStats) []*analytics.BrowserDayStats {
	respData := make([]*analytics.BrowserDayStats, 0, len(stats))

	for _, day := range stats {
		protoClicks := make([]*analytics.BrowserClicks, 0, len(day.Stats))
		for _, cc := range day.Stats {
			protoClicks = append(protoClicks, &analytics.BrowserClicks{
				Browser:      cc.Browser,
				Clicks:       cc.Clicks,
				UniqueClicks: cc.UniqueClicks,
			})
		}

		respData = append(respData, &analytics.BrowserDayStats{
			Date:  day.Date,
			Stats: protoClicks,
		})
	}

	return respData
}

func ToClicksByDeviceTypeResponse(stats []models.DeviceDayStats) []*analytics.DeviceDayStats {
	respData := make([]*analytics.DeviceDayStats, 0, len(stats))

	for _, day := range stats {
		protoClicks := make([]*analytics.DeviceClicks, 0, len(day.Stats))
		for _, cc := range day.Stats {
			protoClicks = append(protoClicks, &analytics.DeviceClicks{
				DeviceType:   cc.DeviceType,
				Clicks:       cc.Clicks,
				UniqueClicks: cc.UniqueClicks,
			})
		}

		respData = append(respData, &analytics.DeviceDayStats{
			Date:  day.Date,
			Stats: protoClicks,
		})
	}

	return respData
}

func ToClicksByOSResponse(stats []models.OSDayStats) []*analytics.OSDayStats {
	respData := make([]*analytics.OSDayStats, 0, len(stats))

	for _, day := range stats {
		protoClicks := make([]*analytics.OSClicks, 0, len(day.Stats))
		for _, cc := range day.Stats {
			protoClicks = append(protoClicks, &analytics.OSClicks{
				Os:           cc.OS,
				Clicks:       cc.Clicks,
				UniqueClicks: cc.UniqueClicks,
			})
		}

		respData = append(respData, &analytics.OSDayStats{
			Date:  day.Date,
			Stats: protoClicks,
		})
	}

	return respData
}

func ToClicksByHourResponse(stats []models.HourDayStats) []*analytics.HourDayStats {
	respData := make([]*analytics.HourDayStats, 0, len(stats))

	for _, day := range stats {
		protoClicks := make([]*analytics.HourStat, 0, len(day.Stats))
		for _, cc := range day.Stats {
			protoClicks = append(protoClicks, &analytics.HourStat{
				Hour:         cc.Hour,
				Clicks:       cc.Clicks,
				UniqueClicks: cc.UniqueClicks,
			})
		}

		respData = append(respData, &analytics.HourDayStats{
			Date:  day.Date,
			Stats: protoClicks,
		})
	}

	return respData
}

func ToClicksByReferrerResponse(stats []models.ReferrerDayStats) []*analytics.ReferrerDayStats {
	respData := make([]*analytics.ReferrerDayStats, 0, len(stats))

	for _, day := range stats {
		protoClicks := make([]*analytics.ReferrerClicks, 0, len(day.Stats))
		for _, cc := range day.Stats {
			protoClicks = append(protoClicks, &analytics.ReferrerClicks{
				Referrer:     cc.Referrer,
				Clicks:       cc.Clicks,
				UniqueClicks: cc.UniqueClicks,
			})
		}

		respData = append(respData, &analytics.ReferrerDayStats{
			Date:  day.Date,
			Stats: protoClicks,
		})
	}

	return respData
}

func ToClicksByDateResponse(stats []models.DateStat) []*analytics.DateStat {
	respData := make([]*analytics.DateStat, 0, len(stats))

	for _, day := range stats {
		respData = append(respData, &analytics.DateStat{
			Date:         day.Date,
			Clicks:       day.Clicks,
			UniqueClicks: day.UniqueClicks,
		})
	}

	return respData
}
