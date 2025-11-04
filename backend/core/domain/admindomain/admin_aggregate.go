package admindomain

import "time"

// Admin — Aggregate Root.
type Admin struct {
	id                     AdminID
	email                  AdminEmail
	name                   AdminName
	isVerified             IsVerified
	passwordHash           PasswordHash
	role                   AdminRole
	verificationCode       VerificationCode
	verificationCodeExpiry VerificationExpiry
	isResetCodeConfirmed   IsResetCodeConfirmed
	createdAt              time.Time
	updatedAt              time.Time
}

func (a *Admin) touch() {
	a.updatedAt = time.Now().UTC()
}
