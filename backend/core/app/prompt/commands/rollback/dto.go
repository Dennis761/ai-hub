package promptapp

type RollbackPromptCommand struct {
	ID      string `json:"_id"`
	Version int    `json:"version"`
	AdminID string `json:"adminID"`
}
