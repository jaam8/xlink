package models

import "time"

type Click struct {
	LinkID     string
	LinkOwner  string
	ClickedAt  time.Time
	Referrer   string
	Region     string
	Browser    string
	DeviceType string
	IPAddress  string
	IsUnique   bool
}
