package mappers

import (
	"time"

	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/infra/db/mongoose/models"
)

func AdminFromDoc(doc *models.AdminDoc) (*admindomain.Admin, error) {
	if doc == nil {
		return nil, nil
	}

	// ID
	id, err := admindomain.NewAdminID(doc.ID)
	if err != nil {
		return nil, err
	}

	// Email
	email, err := admindomain.NewAdminEmail(doc.Email)
	if err != nil {
		return nil, err
	}

	// Name
	name, err := admindomain.NewAdminName(*doc.Name)
	if err != nil {
		return nil, err
	}

	// Password hash
	pass, err := admindomain.NewPasswordHashFromString(doc.Password)
	if err != nil {
		return nil, err
	}

	// Role
	role, err := admindomain.NewAdminRole(doc.Role)
	if err != nil {
		return nil, err
	}

	// Verification code
	vc, err := admindomain.NewVerificationCode(doc.VerificationCode)
	if err != nil {
		return nil, err
	}

	// Verification expiry
	ve, err := admindomain.NewVerificationExpiry(doc.VerificationCodeExpires)
	if err != nil {
		return nil, err
	}

	restore := admindomain.RestoreProps{
		ID:                     id,
		Email:                  email,
		Name:                   name,
		IsVerified:             admindomain.NewIsVerified(doc.IsVerified),
		PasswordHash:           pass,
		Role:                   role,
		VerificationCode:       vc,
		VerificationCodeExpiry: ve,
		IsResetCodeConfirmed:   admindomain.NewIsResetCodeConfirmed(doc.IsResetCodeConfirmed),
		CreatedAt:              doc.CreatedAt,
		UpdatedAt:              doc.UpdatedAt,
	}

	return admindomain.Restore(restore), nil
}

func AdminToPersistence(entity *admindomain.Admin) (*models.AdminDoc, error) {
	p := entity.ToPrimitives()

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &models.AdminDoc{
		ID:                      p.ID,
		Email:                   p.Email,
		Name:                    p.Name,
		IsVerified:              p.IsVerified,
		Password:                p.PasswordHash,
		Role:                    p.Role,
		VerificationCode:        p.VerificationCode,
		VerificationCodeExpires: parseRFC3339Ptr(p.VerificationCodeExpiry),
		IsResetCodeConfirmed:    p.IsResetCodeConfirmed,
		CreatedAt:               createdAt.UTC(),
		UpdatedAt:               updatedAt.UTC(),
	}, nil
}

func parseRFC3339Ptr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	trim := *s
	if trim == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, trim)
	if err != nil {
		return nil
	}
	u := parsed.UTC()
	return &u
}
