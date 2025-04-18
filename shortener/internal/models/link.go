package models

import (
	"github.com/google/uuid"
	"time"
)

type Link struct {
	Id        uuid.UUID  `json:"id"`
	UserId    uuid.UUID  `json:"user_id"`
	GroupId   *uuid.UUID `json:"group_id"`
	Generated bool       `json:"generated"`
	ShortLink string     `json:"short_link"`
	Url       string     `json:"url"`
	CreatedAt time.Time  `json:"created_at"`
	ExpireAt  time.Time  `json:"expire_at"`
}
