package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"ai_hub.com/app/core/ports/adminports"
)

var _ adminports.CodeGenerator = (*SimpleCodeGenerator)(nil)

// SimpleCodeGenerator generates 6-digit verification codes and TTL.
type SimpleCodeGenerator struct {
	ttl time.Duration
}

func NewSimpleCodeGenerator(ttl time.Duration) *SimpleCodeGenerator {
	return &SimpleCodeGenerator{ttl: ttl}
}

func (g *SimpleCodeGenerator) GenerateVerificationCode() string {
	const max = 1_000_000

	if n, err := rand.Int(rand.Reader, big.NewInt(max)); err == nil {
		return fmt.Sprintf("%06d", n.Int64())
	}

	now := time.Now().UnixNano() % max
	return fmt.Sprintf("%06d", now)
}

func (g *SimpleCodeGenerator) VerificationTTL() time.Duration {
	return g.ttl
}
