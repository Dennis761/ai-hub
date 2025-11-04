package projectdomain

import "time"

type CreateProps struct {
	ID          ProjectID
	Name        ProjectName
	APIKey      ProjectAPIKey
	OwnerID     OwnerID
	AdminAccess []AccessAdminID
	Status      *string
	Now         *time.Time
}

// Create builds a new project aggregate.
func Create(props CreateProps) (Project, error) {
	status := "active"
	if props.Status != nil {
		status = *props.Status
	}
	if _, ok := allowedProjectStatuses[status]; !ok {
		return Project{}, InvalidStatus()
	}

	accessSet := make(map[string]struct{}, len(props.AdminAccess))
	for _, adm := range props.AdminAccess {
		accessSet[adm.Value()] = struct{}{}
	}

	return Project{
		id:          props.ID,
		name:        props.Name,
		status:      status,
		apiKey:      props.APIKey,
		ownerID:     props.OwnerID,
		adminAccess: accessSet,
		createdAt:   *props.Now,
		updatedAt:   *props.Now,
	}, nil
}

type RestoreProps struct {
	ID          ProjectID
	Name        ProjectName
	Status      string
	APIKey      ProjectAPIKey
	OwnerID     OwnerID
	AdminAccess []AccessAdminID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Restore rebuilds a project aggregate from stored data.
func Restore(props RestoreProps) (Project, error) {
	if _, ok := allowedProjectStatuses[props.Status]; !ok {
		return Project{}, InvalidStatus()
	}

	accessSet := make(map[string]struct{}, len(props.AdminAccess))
	for _, adm := range props.AdminAccess {
		accessSet[adm.Value()] = struct{}{}
	}

	return Project{
		id:          props.ID,
		name:        props.Name,
		status:      props.Status,
		apiKey:      props.APIKey,
		ownerID:     props.OwnerID,
		adminAccess: accessSet,
		createdAt:   props.CreatedAt.UTC(),
		updatedAt:   props.UpdatedAt.UTC(),
	}, nil
}
