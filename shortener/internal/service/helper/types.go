package helper

type LinkBodyRequest interface {
	LinkBodyRequestOnlyUserId
	LinkBodyRequestOnlyGroupId
	GetGenerated() bool
	GetShortLink() string
	GetUrl() string
	LinkBodyRequestOnlyExpireAt
}

type LinkBodyRequestOnlyId interface {
	GetId() string
}

type LinkBodyRequestOnlyUserId interface {
	GetUserId() string
}

type LinkBodyRequestOnlyGroupId interface {
	GetGroupId() string
}

type LinkBodyRequestOnlyExpireAt interface {
	GetExpireAt() string
}

type LinkBodyRequestWithId interface {
	LinkBodyRequest
	LinkBodyRequestOnlyId
}
