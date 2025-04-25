package ports

import (
	"xlink/common/gen/analytics"
)

type AnalyticsServiceRepository interface {
	ClicksByCountry(request *analytics.GetClicksRequest) (*analytics.ClicksByCountryResponse, error)
	ClicksByRegion(request *analytics.GetClicksRequest) (*analytics.ClicksByRegionResponse, error)
	ClicksByBrowser(request *analytics.GetClicksRequest) (*analytics.ClicksByBrowserResponse, error)
	ClicksByOS(request *analytics.GetClicksRequest) (*analytics.ClicksByOSResponse, error)
	ClicksByDeviceType(request *analytics.GetClicksRequest) (*analytics.ClicksByDeviceTypeResponse, error)
	ClicksByHour(request *analytics.GetClicksRequest) (*analytics.ClicksByHourResponse, error)
	ClicksByDate(request *analytics.GetClicksRequest) (*analytics.ClicksByDateResponse, error)
	ClicksByReferrer(request *analytics.GetClicksRequest) (*analytics.ClicksByReferrerResponse, error)
}
