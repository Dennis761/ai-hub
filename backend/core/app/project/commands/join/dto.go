package projectapp

type JoinProjectByNameCommand struct {
	Name    string `json:"name"`
	APIKey  string `json:"apiKey"`
	AdminID string `json:"adminID"`
}
