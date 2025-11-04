package adminports

import "time"

type CodeGenerator interface {
	GenerateVerificationCode() string

	VerificationTTL() time.Duration
}
