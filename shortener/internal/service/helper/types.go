package helper

type LinkBodyRequest interface {
	GetUserId() string
	GetGroupId() string
	GetGenerated() bool
	GetShortLink() string
	GetUrl() string
	GetExpireAt() string
}

type LinkBodyRequestOnlyId interface {
	GetId() string
}

type LinkBodyRequestWithId interface {
	LinkBodyRequest
	LinkBodyRequestOnlyId
}
