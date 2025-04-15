package helper

type LinkBodyRequest interface {
	GetUserId() string
	GetGroupId() string
	GetGenerated() bool
	GetShortLink() string
	GetUrl() string
	GetExpireAt() string
}

type LinkBodyRequestWithId interface {
	GetId() string
}
