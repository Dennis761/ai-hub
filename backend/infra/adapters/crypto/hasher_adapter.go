package crypto

import (
	"ai_hub.com/app/core/ports/adminports"
	"ai_hub.com/app/core/ports/apikeyports"
	"golang.org/x/crypto/bcrypt"
)

var (
	_ adminports.Hasher  = (*BcryptHasher)(nil)
	_ apikeyports.Hasher = (*BcryptHasher)(nil)
)

type BcryptHasher struct {
	rounds int
}

func NewBcryptHasher(rounds int) *BcryptHasher {
	if rounds <= 0 {
		rounds = 10
	}
	return &BcryptHasher{rounds: rounds}
}

func (h *BcryptHasher) Hash(plaintext string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), h.rounds)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func (h *BcryptHasher) Compare(plaintext, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
