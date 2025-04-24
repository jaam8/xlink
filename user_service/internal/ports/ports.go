package ports

type UsersCacheRepository interface {
	CheckToken(userId string, token string) (bool, error)
	SetToken(userId string, token string) error

	GetRole(userId string) (string, bool, bool, error)
	SetRole(userId string, isStaff bool, isAdmin bool) error
}

type UserStorageRepository interface {
	CheckToken(userId string, token string) (bool, error)
	RefreshToken(userId string, token string) (string, error)

	CreateUser(telegramId *int64, isStaff *bool, isAdmin *bool) (string, string, error)      // userId, token, err
	GetUser(userId string) (string, string, *int64, error)                                   // userId, role, tgId, err
	GetUserIDByToken(token string) (string, bool, error)                                     // userId, status, err
	GetUserIDByTgId(tgId int64) (string, bool, error)                                        // userId, status, err
	UpdateUser(userId string, telegramId *int64, isStaff *bool, isAdmin *bool) (bool, error) // status, err
	DeleteUser(userId string) (bool, error)                                                  // status, err
	GetRole(userId string) (string, bool, bool, error)                                       // role, isStaff, isAdmin, err
	GetTokenByTgId(tgId int64) (string, error)                                               // token, err
}

type ShortenerRepository interface {
	GetLinksCountByUserId(userId string) (int32, error)
}
