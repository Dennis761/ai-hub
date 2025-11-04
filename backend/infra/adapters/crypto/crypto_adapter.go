package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"crypto/sha256"
	"io"

	apikeyports "ai_hub.com/app/core/ports/apikeyports"
	"ai_hub.com/app/infra/config"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

var _ apikeyports.CryptoPort = (*CryptoAdapter)(nil)

type CryptoAdapter struct {
	key []byte
}

func NewCryptoAdapter() *CryptoAdapter {
	secret := []byte(strings.TrimSpace(config.Env.KeyEncryptSecret))
	if len(secret) == 0 {
		secret = []byte("change-me-in-env-please")
	}
	key := make([]byte, chacha20poly1305.KeySize) // 32
	rd := hkdf.New(sha256.New, secret, nil, []byte("api-key-encryption"))
	if _, err := io.ReadFull(rd, key); err != nil {
	}
	return &CryptoAdapter{key: key}
}

func (a *CryptoAdapter) Encrypt(plaintext string) (string, error) {
	aead, err := chacha20poly1305.NewX(a.key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce: %w", err)
	}

	ct := aead.Seal(nil, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(nonce) + ":" + hex.EncodeToString(ct), nil
}

func (a *CryptoAdapter) Decrypt(ciphertext string) (string, error) {
	sep := strings.IndexByte(ciphertext, ':')
	if sep <= 0 || sep >= len(ciphertext)-1 {
		return "", errors.New(`bad format, want "nonceHex:cipherHex"`)
	}
	nonceHex := ciphertext[:sep]
	ctHex := ciphertext[sep+1:]

	nonce, err := hex.DecodeString(nonceHex)
	if err != nil {
		return "", fmt.Errorf("nonce hex: %w", err)
	}
	ct, err := hex.DecodeString(ctHex)
	if err != nil {
		return "", fmt.Errorf("cipher hex: %w", err)
	}

	if len(nonce) != chacha20poly1305.NonceSizeX {
		return "", errors.New("invalid nonce size")
	}

	aead, err := chacha20poly1305.NewX(a.key)
	if err != nil {
		return "", err
	}

	pt, err := aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}
