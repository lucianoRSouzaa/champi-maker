package port

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
