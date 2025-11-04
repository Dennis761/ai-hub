package promptapp

type ReorderPromptsCommand struct {
	Items   []string `json:"items"`
	AdminID string   `json:"adminID"`
}
