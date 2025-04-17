package ports

type TokensRepository interface {
	Check(userId string, token string) (bool, error)
	Create(userId string) (string, error)
	Refresh(userId string, token string) (string, error)
	Delete(userId string) (bool, error)
}
