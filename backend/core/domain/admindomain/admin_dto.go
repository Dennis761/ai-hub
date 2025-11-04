package admindomain

import "time"

type AdminPrimitives struct {
	ID                     string  `json:"_id"`
	Email                  string  `json:"email"`
	Name                   *string `json:"name"`
	IsVerified             bool    `json:"isVerified"`
	PasswordHash           string  `json:"passwordHash"`
	Role                   string  `json:"role"`
	VerificationCode       *string `json:"verificationCode"`
	VerificationCodeExpiry *string `json:"verificationCodeExpires"`
	IsResetCodeConfirmed   bool    `json:"isResetCodeConfirmed"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}

func (a *Admin) ToPrimitives() AdminPrimitives {
	var primitiveVerificationCode *string
	if value := a.verificationCode.Value(); value != nil {
		primitiveVerificationCode = value
	}

	var primitiveVerificationExpiry *string
	if t := a.verificationCodeExpiry.Value(); t != nil {
		formatted := t.UTC().Format(time.RFC3339)
		primitiveVerificationExpiry = &formatted
	}

	return AdminPrimitives{
		ID:                     a.id.Value(),
		Email:                  a.email.Value(),
		Name:                   a.name.Value(),
		IsVerified:             a.isVerified.Value(),
		PasswordHash:           a.passwordHash.ExposeForPersistence(),
		Role:                   a.role.Value(),
		VerificationCode:       primitiveVerificationCode,
		VerificationCodeExpiry: primitiveVerificationExpiry,
		IsResetCodeConfirmed:   a.isResetCodeConfirmed.Value(),
		CreatedAt:              a.createdAt.UTC().Format(time.RFC3339),
		UpdatedAt:              a.updatedAt.UTC().Format(time.RFC3339),
	}
}
