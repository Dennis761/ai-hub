package projectdomain

import "time"

// Api key — Aggregate Root.
type Project struct {
	id          ProjectID
	name        ProjectName
	status      string
	apiKey      ProjectAPIKey
	ownerID     OwnerID
	adminAccess map[string]struct{}
	createdAt   time.Time
	updatedAt   time.Time
}

func (p *Project) touch() {
	p.updatedAt = time.Now().UTC()
}
