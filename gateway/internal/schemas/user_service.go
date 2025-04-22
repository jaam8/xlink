package schemas

type CreateUserSchema struct {
	TgId *int64 `json:"tg_id,omitempty"`
}

type CreateUserSchemaAdmin struct {
	TgId    *int64 `json:"tg_id,omitempty"`
	IsStaff *bool  `json:"is_staff,omitempty"`
	IsAdmin *bool  `json:"is_admin,omitempty"`
}

type UserIdByTokenSchema struct {
	Token string `json:"token"`
}

type UserIdByTgIdSchema struct {
	TgId int64 `json:"tg_id"`
}

type UpdateUserSchema struct {
	TgId int64 `json:"tg_id"`
}

type UpdateUserSchemaAdmin struct {
	TgId    *int64 `json:"tg_id,omitempty"`
	IsStaff *bool  `json:"is_staff,omitempty"`
	IsAdmin *bool  `json:"is_admin,omitempty"`
}

type TokenCheckRequest struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type RefreshTokenSchema struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}
