package schemas

type CreateLinkSchema struct {
	ShortLink *string `json:"short_link,omitempty"`
	TargetUrl string  `json:"target_url"`
}

type CreateLinkSchemaAdmin struct {
	UserId    string  `json:"user_id"`
	ShortLink *string `json:"short_link,omitempty"`
	TargetUrl string  `json:"target_url"`
}

type UpdateLinkSchema struct {
	Regenerate bool    `json:"regenerate"`
	ShortLink  *string `json:"short_link,omitempty"`
	TargetUrl  *string `json:"target_url,omitempty"`
	ExpireAt   string  `json:"expire_at"`
}
