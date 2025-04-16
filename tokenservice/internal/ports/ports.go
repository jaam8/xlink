package ports

type TokensRepository interface {
	Check(userId string, token string) (bool, error)
	GetUserIdByToken(token string) (string, error)
	Create(userId string) (string, error)
	Delete(userId string) error
}
