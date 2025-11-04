package adminports

import (
	"context"

	"time"
)

type Mailer interface {
	SendVerificationCode(ctx context.Context, email string, code string, expiresAt time.Time) error

	SendResetCode(ctx context.Context, email string, code string, expiresAt time.Time) error
}
