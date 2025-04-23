package models

// ========================== Country ==========================

type CountryDayStats struct {
	Date  string          `ch:"date"`
	Stats []CountryClicks `ch:"-"`
}

type CountryClicks struct {
	Country      string `ch:"country"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}

// ========================== Region ==========================

type RegionDayStats struct {
	Date  string         `ch:"date"`
	Stats []RegionClicks `ch:"-"`
}

type RegionClicks struct {
	Region       string `ch:"region"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}

// ========================== Browser ==========================

type BrowserDayStats struct {
	Date  string          `ch:"date"`
	Stats []BrowserClicks `ch:"-"`
}

type BrowserClicks struct {
	Browser      string `ch:"browser"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}

// ========================== Device Type ==========================

type DeviceDayStats struct {
	Date  string         `ch:"date"`
	Stats []DeviceClicks `ch:"-"`
}

type DeviceClicks struct {
	DeviceType   string `ch:"device_type"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}

// ========================== OS ==========================

type OSDayStats struct {
	Date  string     `ch:"date"`
	Stats []OSClicks `ch:"-"`
}

type OSClicks struct {
	OS           string `ch:"os"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}

// ========================== Referrer ==========================

type ReferrerDayStats struct {
	Date  string           `ch:"date"`
	Stats []ReferrerClicks `ch:"-"`
}

type ReferrerClicks struct {
	Referrer     string `ch:"referrer"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}

// ========================== Hour ==========================

type HourDayStats struct {
	Date  string     `ch:"date"`
	Stats []HourStat `ch:"-"`
}

type HourStat struct {
	Hour         uint32 `ch:"hour"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}

// ========================== Date ==========================

type DateStat struct {
	Date         string `ch:"date"`
	Clicks       uint64 `ch:"clicks"`
	UniqueClicks uint64 `ch:"unique_clicks"`
}
