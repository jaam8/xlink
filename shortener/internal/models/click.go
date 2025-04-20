package models

import "time"

type Click struct {
	ShortLink    string    `json:"short_link"`
	ClickedAt    time.Time `json:"clicked_at"`
	Referrer     string    `json:"referrer"`
	IPAddress    string    `json:"ip_address"`
	VisitorToken string    `json:"visitor_token"`
	Browser      string    `json:"browser"`
	DeviceType   string    `json:"device_type"`
	Os           string    `json:"os"`
}
