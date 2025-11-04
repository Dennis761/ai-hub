package admindomain

import "time"

type NewAdminProps struct {
	ID           AdminID
	Email        AdminEmail
	Name         AdminName
	PasswordHash PasswordHash
	Role         AdminRole
	CurrentTime  time.Time
}

// NewAdmin creates a new Admin aggregate.
func NewAdmin(props NewAdminProps) (*Admin, error) {
	admin := &Admin{
		id:                     props.ID,
		email:                  props.Email,
		name:                   props.Name,
		isVerified:             NewIsVerified(false),
		passwordHash:           props.PasswordHash,
		role:                   props.Role,
		verificationCode:       emptyVerificationCode(),
		verificationCodeExpiry: emptyVerificationExpiry(),
		isResetCodeConfirmed:   NewIsResetCodeConfirmed(false),
		createdAt:              props.CurrentTime,
		updatedAt:              props.CurrentTime,
	}
	return admin, nil
}

type RestoreProps struct {
	ID                     AdminID
	Email                  AdminEmail
	Name                   AdminName
	IsVerified             IsVerified
	PasswordHash           PasswordHash
	Role                   AdminRole
	VerificationCode       VerificationCode
	VerificationCodeExpiry VerificationExpiry
	IsResetCodeConfirmed   IsResetCodeConfirmed
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// Restore rebuilds an Admin aggregate from stored data.
func Restore(props RestoreProps) *Admin {
	return &Admin{
		id:                     props.ID,
		email:                  props.Email,
		name:                   props.Name,
		isVerified:             props.IsVerified,
		passwordHash:           props.PasswordHash,
		role:                   props.Role,
		verificationCode:       props.VerificationCode,
		verificationCodeExpiry: props.VerificationCodeExpiry,
		isResetCodeConfirmed:   props.IsResetCodeConfirmed,
		createdAt:              props.CreatedAt.UTC(),
		updatedAt:              props.UpdatedAt.UTC(),
	}
}
