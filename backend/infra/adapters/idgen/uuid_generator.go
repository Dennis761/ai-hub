package idgen

import (
	"ai_hub.com/app/core/ports/adminports"
	"github.com/google/uuid"
)

type UUIDGenerator struct{}

var _ adminports.IDGenerator = (*UUIDGenerator)(nil)

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) NewID() string {
	return uuid.NewString()
}
