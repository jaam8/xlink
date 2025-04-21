package helper

import "google.golang.org/protobuf/types/known/timestamppb"

type LinkCreateRequest interface {
	LinkRequestOnlyUserId
	LinkRequestOnlyShortLink
	LinkRequestOnlyTargetUrl
}

type LinkUpdateRequest interface {
	LinkRequestOnlyLinkId
	LinkRequestOnlyUserId
	LinkRequestOnlyRegenerate
	LinkRequestOnlyShortLink
	LinkRequestOnlyTargetUrl
	LinkRequestOnlyExpireAt
}

type LinkRequestOnlyLinkId interface {
	GetLinkId() string
}

type LinkRequestOnlyUserId interface {
	GetUserId() string
}

type LinkRequestOnlyShortLink interface {
	GetShortLink() string
}

type LinkRequestOnlyTargetUrl interface {
	GetTargetUrl() string
}

type LinkRequestOnlyRegenerate interface {
	GetRegenerate() bool
}

type LinkRequestOnlyExpireAt interface {
	GetExpireAt() *timestamppb.Timestamp
}
