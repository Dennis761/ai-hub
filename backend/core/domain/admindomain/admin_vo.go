package admindomain

import (
	"regexp"
	"strings"
	"time"
)

const (
	emailMaxLength = 254
	nameMinLength  = 6
	nameMaxLength  = 120

	// stored hash (already hashed)
	minPasswordHashLen = 20

	// plain password limits (before hashing)
	plainPasswordMinLen = 8
	plainPasswordMaxLen = 100
)

var (
	emailRegex            = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	verificationCodeRegex = regexp.MustCompile(`^\d{6}$`)
)

// ===================== AdminID =====================

type AdminID struct {
	value string
}

func NewAdminID(raw string) (AdminID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return AdminID{}, InvalidAdminID()
	}
	return AdminID{value: trimmed}, nil
}

func (a AdminID) Value() string { return a.value }

// ===================== AdminEmail =====================

type AdminEmail struct {
	value string
}

func NewAdminEmail(raw string) (AdminEmail, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" || !emailRegex.MatchString(normalized) {
		return AdminEmail{}, InvalidEmailFormat()
	}
	if len(normalized) > emailMaxLength {
		return AdminEmail{}, EmailTooLong(emailMaxLength)
	}
	return AdminEmail{value: normalized}, nil
}

func (e AdminEmail) Value() string { return e.value }

// ===================== AdminName =====================

type AdminName struct {
	value *string
}

func NewAdminName(raw string) (AdminName, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return AdminName{value: nil}, nil
	}
	if len(trimmed) < nameMinLength {
		return AdminName{}, NameTooShort(nameMinLength)
	}
	if len(trimmed) > nameMaxLength {
		return AdminName{}, NameTooLong(nameMaxLength)
	}
	copied := trimmed
	return AdminName{value: &copied}, nil
}

func (n AdminName) Value() *string { return n.value }

// ===================== PasswordHash =====================

type PasswordHash struct {
	hash string
}

func NewPasswordHashFromString(raw string) (PasswordHash, error) {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) < minPasswordHashLen {
		return PasswordHash{}, InvalidPasswordHash()
	}
	return PasswordHash{hash: trimmed}, nil
}

func (p PasswordHash) ExposeForPersistence() string { return p.hash }

// ===================== AdminRole =====================

type AdminRole struct {
	value string
}

func NewAdminRole(raw string) (AdminRole, error) {
	if raw != "admin" {
		return AdminRole{}, InvalidRole()
	}
	return AdminRole{value: raw}, nil
}

func (r AdminRole) Value() string { return r.value }

// ===================== VerificationCode =====================

type VerificationCode struct {
	value *string
}

func NewVerificationCode(raw *string) (VerificationCode, error) {
	if raw == nil {
		return emptyVerificationCode(), nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return emptyVerificationCode(), nil
	}
	if !verificationCodeRegex.MatchString(trimmed) {
		return VerificationCode{}, InvalidVerificationCode()
	}
	copied := trimmed
	return VerificationCode{value: &copied}, nil
}

func (vc VerificationCode) Value() *string { return vc.value }

func emptyVerificationCode() VerificationCode {
	return VerificationCode{value: nil}
}

// ===================== VerificationExpiry =====================

type VerificationExpiry struct {
	value *time.Time
}

func NewVerificationExpiry(t *time.Time) (VerificationExpiry, error) {
	if t == nil {
		return emptyVerificationExpiry(), nil
	}
	if t.IsZero() {
		return VerificationExpiry{}, InvalidVerificationExpiry()
	}
	utc := t.UTC()
	return VerificationExpiry{value: &utc}, nil
}

func NewVerificationExpiryFromString(s *string) (VerificationExpiry, error) {
	if s == nil {
		return emptyVerificationExpiry(), nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return emptyVerificationExpiry(), nil
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return VerificationExpiry{}, InvalidVerificationExpiry()
	}
	parsed = parsed.UTC()
	return VerificationExpiry{value: &parsed}, nil
}

func (ve VerificationExpiry) Value() *time.Time { return ve.value }

func (ve VerificationExpiry) IsExpired(now time.Time) bool {
	if ve.value == nil {
		return false
	}
	ref := now
	if ref.IsZero() {
		ref = time.Now().UTC()
	}
	return !ve.value.After(ref)
}

func emptyVerificationExpiry() VerificationExpiry {
	return VerificationExpiry{value: nil}
}

// ===================== IsVerified =====================

type IsVerified struct {
	value bool
}

func NewIsVerified(raw bool) IsVerified { return IsVerified{value: raw} }
func (v IsVerified) Value() bool        { return v.value }

// ===================== IsResetCodeConfirmed =====================

type IsResetCodeConfirmed struct {
	value bool
}

func NewIsResetCodeConfirmed(raw bool) IsResetCodeConfirmed {
	return IsResetCodeConfirmed{value: raw}
}

func (v IsResetCodeConfirmed) Value() bool { return v.value }

// ===================== PlainPassword =====================

type PlainPassword struct {
	value string
}

func NewPlainPassword(raw string) (PlainPassword, error) {
	trimmed := strings.TrimSpace(raw)
	l := len(trimmed)

	if l < plainPasswordMinLen {
		return PlainPassword{}, PasswordTooShort(plainPasswordMinLen)
	}
	if l > plainPasswordMaxLen {
		return PlainPassword{}, PasswordTooLong(plainPasswordMaxLen)
	}

	return PlainPassword{value: trimmed}, nil
}

func (p PlainPassword) Value() string {
	return p.value
}
