package security

import "unicode"

type PasswordValidationResult struct {
	OK      bool
	TooWeak bool
}

type SimplePasswordPolicy struct {
	requireStrength bool
}

func NewSimpleProjectPasswordPolicy(requireStrength bool) *SimplePasswordPolicy {
	return &SimplePasswordPolicy{
		requireStrength: requireStrength,
	}
}

// Validate expects already trimmed password (VO does it).
func (p *SimplePasswordPolicy) Validate(plain string) PasswordValidationResult {
	res := PasswordValidationResult{OK: true}

	if p.requireStrength {
		hasUpper := false
		hasLower := false
		hasDigit := false
		hasSpecial := false

		for _, r := range plain {
			switch {
			case r >= 'A' && r <= 'Z':
				hasUpper = true
			case r >= 'a' && r <= 'z':
				hasLower = true
			case r >= '0' && r <= '9':
				hasDigit = true
			default:
				if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
					hasSpecial = true
				}
			}
		}

		if !(hasUpper && hasLower && hasDigit && hasSpecial) {
			res.OK = false
			res.TooWeak = true
		}
	}

	return res
}
