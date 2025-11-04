package admindomain

import "fmt"

// DomainError represents domain-level errors with a unique code.
type DomainError interface {
	error
	Code() string
}

// --- Basic error types ---

// simpleError handles fixed errors (e.g. Forbidden, InvalidCredentials).
type simpleError struct {
	code string
	msg  string
}

func (e simpleError) Error() string { return e.msg }
func (e simpleError) Code() string  { return e.code }

// limitError handles limit-based errors (e.g. too short, too long).
type limitError struct {
	code   string
	format string
	limit  int
	def    int
}

func (e limitError) Error() string {
	lim := e.limit
	if lim <= 0 {
		lim = e.def
	}
	return fmt.Sprintf(e.format, lim)
}

func (e limitError) Code() string { return e.code }

// --- Resource / Access errors ---

func AdminNotFound() DomainError {
	return simpleError{"AdminNotFound", "Admin not found"}
}

func Forbidden() DomainError {
	return simpleError{"Forbidden", "Forbidden"}
}

func EmailAlreadyUsed() DomainError {
	return simpleError{"EmailAlreadyUsed", "Email already used"}
}

// --- Verification / Reset ---

func VerificationCodeExpired() DomainError {
	return simpleError{"VerificationCodeExpired", "Verification code has expired"}
}

func InvalidVerificationCode() DomainError {
	return simpleError{"InvalidVerificationCode", "Invalid verification code"}
}

func ResetCodeNotConfirmed() DomainError {
	return simpleError{"ResetCodeNotConfirmed", "Reset code has not been confirmed"}
}

// --- Auth / Login ---

func InvalidCredentials() DomainError {
	return simpleError{"InvalidCredentials", "Invalid email or password"}
}

func EmailNotVerified() DomainError {
	return simpleError{"EmailNotVerified", "Email has not been verified"}
}

// --- Roles ---

func InvalidRole() DomainError {
	return simpleError{"InvalidRole", "Invalid role. Allowed: admin"}
}

// --- Validation (Value Objects) ---

func InvalidAdminID() DomainError {
	return simpleError{"InvalidAdminID", "Invalid admin _id"}
}

func InvalidEmailFormat() DomainError {
	return simpleError{"InvalidEmailFormat", "Invalid email format"}
}

func EmailTooLong(max int) DomainError {
	return limitError{"EmailTooLong", "Email is too long (max %d)", max, 254}
}

func NameTooShort(min int) DomainError {
	return limitError{"NameTooShort", "Name is too short (min %d)", min, 6}
}

func NameTooLong(max int) DomainError {
	return limitError{"NameTooLong", "Name is too long (max %d)", max, 120}
}

func InvalidPasswordHash() DomainError {
	return simpleError{"InvalidPasswordHash", "Invalid password hash"}
}

func PasswordTooShort(min int) DomainError {
	return limitError{"PasswordTooShort", "Password is too short (min %d)", min, 8}
}

func PasswordTooLong(max int) DomainError {
	return limitError{"PasswordTooLong", "Password is too long (max %d)", max, 100}
}

func PasswordTooWeak() DomainError {
	return simpleError{"PasswordTooWeak", "Password is too weak (must include upper/lowercase, number, special character)"}
}

func AdminIDRequired() DomainError {
	return simpleError{"AdminIDRequired", "Admin ID is required"}
}

func InvalidVerificationExpiry() DomainError {
	return simpleError{"InvalidVerificationExpiry", "Invalid verification expiry date"}
}
