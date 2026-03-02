package repository

// RepositoryContainer возвращает репозитории пользователей и секретов.
type RepositoryContainer interface {
	UserRepository() UserRepository
	SecretRepository() SecretRepository
}
