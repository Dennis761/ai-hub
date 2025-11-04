package taskdomain

import "time"

type CreateProps struct {
	ID          TaskID
	Name        TaskName
	Description *TaskDescription
	ProjectID   TaskProjectID
	APIMethod   APIMethod
	Status      *string
	CreatedBy   TaskCreatorID
	Now         *time.Time
}

func Create(props CreateProps) (Task, error) {
	currentTime := time.Now().UTC()
	if props.Now != nil {
		currentTime = props.Now.UTC()
	}

	resolvedStatus := "active"
	if props.Status != nil {
		resolvedStatus = *props.Status
	}
	if _, ok := allowedStatuses[resolvedStatus]; !ok {
		return Task{}, InvalidStatus()
	}

	var descVO TaskDescription
	if props.Description != nil {
		descVO = *props.Description
	} else {
		empty, err := NewTaskDescription(nil)
		if err != nil {
			return Task{}, err
		}
		descVO = empty
	}

	return Task{
		id:          props.ID,
		name:        props.Name,
		description: descVO,
		projectID:   props.ProjectID,
		apiMethod:   props.APIMethod,
		status:      resolvedStatus,
		createdBy:   props.CreatedBy,
		createdAt:   currentTime,
		updatedAt:   currentTime,
	}, nil
}

type RestoreProps struct {
	ID          TaskID
	Name        TaskName
	Description TaskDescription
	ProjectID   TaskProjectID
	APIMethod   APIMethod
	Status      string
	CreatedBy   TaskCreatorID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func Restore(props RestoreProps) (Task, error) {
	if _, ok := allowedStatuses[props.Status]; !ok {
		return Task{}, InvalidStatus()
	}
	return Task{
		id:          props.ID,
		name:        props.Name,
		description: props.Description,
		projectID:   props.ProjectID,
		apiMethod:   props.APIMethod,
		status:      props.Status,
		createdBy:   props.CreatedBy,
		createdAt:   props.CreatedAt.UTC(),
		updatedAt:   props.UpdatedAt.UTC(),
	}, nil
}
