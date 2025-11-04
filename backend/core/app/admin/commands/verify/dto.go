package adminapp

type VerifyAdminCommand struct {
	ID   string `json:"_id"`
	Code string `json:"code"`
}
