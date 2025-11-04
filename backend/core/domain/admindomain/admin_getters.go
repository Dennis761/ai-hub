package admindomain

import "time"

func (a *Admin) ID() AdminID                        { return a.id }
func (a *Admin) Email() AdminEmail                  { return a.email }
func (a *Admin) IsVerified() IsVerified             { return a.isVerified }
func (a *Admin) VerificationCode() VerificationCode { return a.verificationCode }
func (a *Admin) VerificationCodeExpiry() VerificationExpiry {
	return a.verificationCodeExpiry
}
func (a *Admin) IsResetCodeConfirmed() IsResetCodeConfirmed {
	return a.isResetCodeConfirmed
}
func (a *Admin) CreatedAt() time.Time { return a.createdAt }
func (a *Admin) UpdatedAt() time.Time { return a.updatedAt }

func (a *Admin) ExposePasswordHash() string {
	return a.passwordHash.ExposeForPersistence()
}
