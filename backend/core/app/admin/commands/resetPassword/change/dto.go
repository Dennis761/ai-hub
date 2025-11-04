package adminapp

type ChangePasswordWithCodeCommand struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}
