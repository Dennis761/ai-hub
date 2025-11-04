package adminapp

type LoginAdminResult struct {
	ID         string  `json:"_id"`
	Email      *string `json:"email,omitempty"`
	Name       *string `json:"name,omitempty"`
	Role       *string `json:"role,omitempty"`
	IsVerified bool    `json:"isVerified"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

type LoginAdminResponse struct {
	Token string           `json:"token"`
	Admin LoginAdminResult `json:"admin"`
}
