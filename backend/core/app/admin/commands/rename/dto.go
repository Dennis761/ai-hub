package adminapp

type RenameAdminCommand struct {
	ID   string  `json:"_id"`
	Name *string `json:"name,omitempty"`
}
