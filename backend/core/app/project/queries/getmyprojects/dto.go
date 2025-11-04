package projectapp

type GetMyProjectsQuery struct {
	AdminID string `json:"adminID"`
}

type ProjectListItem struct {
	ID     string `json:"_id"`
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
	Status string `json:"status"`
}

type GetMyProjectsResult struct {
	Projects []ProjectListItem `json:"projects"`
}
