package adminports

type Hasher interface {
	Hash(plaintext string) (string, error)

	Compare(plaintext, hash string) (bool, error)
}
