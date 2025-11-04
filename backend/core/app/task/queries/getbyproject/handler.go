package taskapp

import (
	"context"

	"ai_hub.com/app/core/app/task/shared/access"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/taskdomain"
	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/taskports"
)

type GetTasksByProjectHandler struct {
	taskReadRepo    taskports.TaskReadRepository
	projectReadRepo projectports.ProjectReadRepository
}

func NewGetTasksByProjectHandler(
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
) *GetTasksByProjectHandler {
	return &GetTasksByProjectHandler{
		taskReadRepo:    taskReadRepo,
		projectReadRepo: projectReadRepo,
	}
}

func (h *GetTasksByProjectHandler) GetTasksByProject(
	ctx context.Context,
	q GetTasksByProjectQuery,
) ([]TaskListItem, error) {
	// validate project id
	projectRef, err := taskdomain.NewTaskProjectID(q.ProjectID)
	if err != nil {
		return nil, err
	}

	// check project and access
	projectID, err := projectdomain.NewProjectID(projectRef.Value())
	if err != nil {
		return nil, err
	}

	project, err := h.projectReadRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, taskdomain.TaskNotFound()
	}

	projectPrim := project.ToPrimitives()
	if !access.EnsureAccess(projectPrim.OwnerID, projectPrim.AdminAccess, q.AdminID) {
		return nil, projectdomain.Forbidden()
	}

	// load tasks
	tasks, err := h.taskReadRepo.FindAllByProject(ctx, projectRef)
	if err != nil {
		return nil, err
	}

	// map to DTO
	out := make([]TaskListItem, 0, len(tasks))
	for _, task := range tasks {
		p := task.ToPrimitives()
		out = append(out, TaskListItem{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			ProjectID:   p.ProjectID,
			APIMethod:   p.APIMethod,
			Status:      p.Status,
			CreatedBy:   p.CreatedBy,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return out, nil
}
