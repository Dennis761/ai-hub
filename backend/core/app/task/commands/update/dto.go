package taskapp

type UpdateTaskCommand struct {
	ID          string
	AdminID     string
	Name        *string
	Description *string
	APIMethod   *string
	Status      *string
}
