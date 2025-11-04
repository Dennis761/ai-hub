package adminapp

type CreateAdminCommand struct {
	ID       string `json:"_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
