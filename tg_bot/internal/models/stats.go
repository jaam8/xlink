package models

type Stat[T any] struct {
	Date    string `json:"date"` // "2025-04-20"
	DayStat []T    `json:"stats"`
}

type CountryDayStat struct {
	Country      string `json:"country"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}

type RegionDayStat struct {
	Region       string `json:"region"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}

type BrowserDayStat struct {
	Browser      string `json:"browser"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}

type DeviceDayStat struct {
	DeviceType   string `json:"device_type"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}

type OSDayStat struct {
	OS           string `json:"os"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}

type ReferrerDayStat struct {
	Referrer     string `json:"referrer"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}

type HourDayStat struct {
	Hour         uint32 `json:"hour"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}

type DateStat struct {
	Date         string `json:"date"`
	Clicks       uint64 `json:"clicks"`
	UniqueClicks uint64 `json:"unique_clicks"`
}
