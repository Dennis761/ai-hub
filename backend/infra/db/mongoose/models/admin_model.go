package models

import "time"

type AdminDoc struct {
	ID                      string     `bson:"_id"`
	Email                   string     `bson:"email"`
	Name                    *string    `bson:"name,omitempty"`
	IsVerified              bool       `bson:"isVerified"`
	Password                string     `bson:"password"`
	Role                    string     `bson:"role"`
	VerificationCode        *string    `bson:"verificationCode,omitempty"`
	VerificationCodeExpires *time.Time `bson:"verificationCodeExpires,omitempty"`
	IsResetCodeConfirmed    bool       `bson:"isResetCodeConfirmed"`
	CreatedAt               time.Time  `bson:"createdAt"`
	UpdatedAt               time.Time  `bson:"updatedAt"`
}
