package mappers

import (
	"time"

	"ai_hub.com/app/core/domain/taskdomain"
)

type TaskDoc struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description *string   `bson:"description,omitempty"`
	ProjectID   string    `bson:"projectId"`
	APIMethod   string    `bson:"apiMethod"`
	Status      string    `bson:"status"`
	CreatedBy   string    `bson:"createdBy"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

func TaskFromDoc(doc *TaskDoc) (*taskdomain.Task, error) {
	if doc == nil {
		return nil, nil
	}
	// ID
	id, err := taskdomain.NewTaskID(doc.ID)
	if err != nil {
		return nil, err
	}
	// Name
	name, err := taskdomain.NewTaskName(doc.Name)
	if err != nil {
		return nil, err
	}
	// Description
	desc, err := taskdomain.NewTaskDescription(doc.Description)
	if err != nil {
		return nil, err
	}
	// ProjectID
	projectID, err := taskdomain.NewTaskProjectID(doc.ProjectID)
	if err != nil {
		return nil, err
	}
	// APIMethod
	apiMethod, err := taskdomain.NewAPIMethod(doc.APIMethod)
	if err != nil {
		return nil, err
	}
	// CreatedBy
	creator, err := taskdomain.NewTaskCreatorID(doc.CreatedBy)
	if err != nil {
		return nil, err
	}

	props := taskdomain.RestoreProps{
		ID:          id,
		Name:        name,
		Description: desc,
		ProjectID:   projectID,
		APIMethod:   apiMethod,
		Status:      doc.Status,
		CreatedBy:   creator,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	t, err := taskdomain.Restore(props)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func TaskToPersistence(entity *taskdomain.Task) (*TaskDoc, error) {
	p := entity.ToPrimitives()

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &TaskDoc{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ProjectID:   p.ProjectID,
		APIMethod:   p.APIMethod,
		Status:      p.Status,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   createdAt.UTC(),
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}
