package admindomain

func (a *Admin) Rename(newName AdminName) {
	a.name = newName
	a.touch()
}

func (a *Admin) ChangePassword(newHash PasswordHash) {
	a.passwordHash = newHash
	a.touch()
}

func (a *Admin) Verify() {
	if !a.isVerified.Value() {
		a.isVerified = NewIsVerified(true)
		a.touch()
	}
}

func (a *Admin) SetVerificationCode(code VerificationCode, expiry VerificationExpiry) {
	a.verificationCode = code
	a.verificationCodeExpiry = expiry
	a.touch()
}

func (a *Admin) ClearVerificationCode() {
	a.verificationCode = emptyVerificationCode()
	a.verificationCodeExpiry = emptyVerificationExpiry()
	a.touch()
}

func (a *Admin) ConfirmResetCode() {
	a.isResetCodeConfirmed = NewIsResetCodeConfirmed(true)
	a.touch()
}

func (a *Admin) ResetConfirmation() {
	a.isResetCodeConfirmed = NewIsResetCodeConfirmed(false)
	a.touch()
}

func (a *Admin) SetRole(role AdminRole) {
	a.role = role
	a.touch()
}
