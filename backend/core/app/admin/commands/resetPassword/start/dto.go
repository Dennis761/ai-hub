package adminapp

type StartPasswordResetCommand struct {
	Email string `json:"email"`
}
