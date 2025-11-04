package adminapp

type ConfirmResetCodeCommand struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
