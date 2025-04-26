package statistics_data

type StatisticsData struct {
	Stats []Stat
}

type Stat struct {
	Date  string
	Items []Item
}

type Item struct {
	Clicks       uint64
	UniqueClicks uint64
	ParamValue   string // e.g. "Russia" or "12:00"
}
