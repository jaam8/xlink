package models

import (
	"github.com/google/uuid"
	"time"
)

type Click struct {
	LinkOwner  uuid.UUID `json:"link_owner" ch:"link_owner"`
	ShortLink  string    `json:"short_link" ch:"short_link"`
	ClickedAt  time.Time `json:"clicked_at" ch:"clicked_at"`
	Referrer   string    `json:"referrer" ch:"referrer"`
	Country    string    `json:"country" ch:"country"`
	Region     string    `json:"region" ch:"region"`
	Browser    string    `json:"browser" ch:"browser"`
	DeviceType string    `json:"device_type" ch:"device_type"`
	Os         string    `json:"os" ch:"os"`
	IPAddress  string    `json:"ip_address" ch:"ip_address"`
	IsUnique   bool      `json:"is_unique" ch:"is_unique"`
}

type ClickEvent struct {
	ShortLink    string    `json:"short_link"`
	ClickedAt    time.Time `json:"clicked_at"`
	Referrer     string    `json:"referrer"`
	IPAddress    string    `json:"ip_address"`
	VisitorToken string    `json:"visitor_token"`
	Browser      string    `json:"browser"`
	DeviceType   string    `json:"device_type"`
	Os           string    `json:"os"`
}
