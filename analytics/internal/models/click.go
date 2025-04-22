package models

import "time"

type Click struct {
	LinkID     string    `json:"link_id"`
	LinkOwner  string    `json:"link_owner"`
	ClickedAt  time.Time `json:"clicked_at"`
	Referrer   string    `json:"referrer"`
	Region     string    `json:"region"`
	Browser    string    `json:"browser"`
	DeviceType string    `json:"device_type"`
	IPAddress  string    `json:"ip_address"`
	IsUnique   bool      `json:"is_unique"`
}

type RedirectClickInfo struct {
	ShortLink    string    `json:"short_link"`
	ClickedAt    time.Time `json:"clicked_at"`
	Referrer     string    `json:"referrer"`
	IPAddress    string    `json:"ip_address"`
	VisitorToken string    `json:"visitor_token"`
	Browser      string    `json:"browser"`
	DeviceType   string    `json:"device_type"`
	Os           string    `json:"os"`
}
