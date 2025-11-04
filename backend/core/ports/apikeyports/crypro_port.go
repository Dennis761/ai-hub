package apikeyports

type CryptoPort interface {
	Encrypt(plaintext string) (string, error)

	Decrypt(cipher string) (string, error)
}
