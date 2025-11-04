package taskapp

import (
	"context"

	"ai_hub.com/app/core/app/task/shared/access"
	"ai_hub.com/app/core/app/task/shared/cache"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/taskdomain"
	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/taskports"
)

type DeleteTaskHandler struct {
	taskReadRepo    taskports.TaskReadRepository
	taskWriteRepo   taskports.TaskWriteRepository
	uow             taskports.UnitOfWorkPort
	projectReadRepo projectports.ProjectReadRepository
	projectCache    projectports.ProjectCachePort
}

func NewDeleteTaskHandler(
	taskReadRepo taskports.TaskReadRepository,
	taskWriteRepo taskports.TaskWriteRepository,
	uow taskports.UnitOfWorkPort,
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
) *DeleteTaskHandler {
	return &DeleteTaskHandler{
		taskReadRepo:    taskReadRepo,
		taskWriteRepo:   taskWriteRepo,
		uow:             uow,
		projectReadRepo: projectReadRepo,
		projectCache:    projectCache,
	}
}

func (h *DeleteTaskHandler) Delete(ctx context.Context, cmd DeleteTaskCommand) error {
	var projectID string

	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate task id
		taskID, err := taskdomain.NewTaskID(cmd.ID)
		if err != nil {
			return nil, err
		}

		// load task aggregate
		task, err := h.taskReadRepo.FindByID(txCtx, taskID)
		if err != nil {
			return nil, err
		}
		// idempotent: no task → success
		if task == nil {
			return struct{}{}, nil
		}

		// check access via project
		projectID = task.ProjectID().Value()

		projectVO, err := projectdomain.NewProjectID(projectID)
		if err != nil {
			return nil, err
		}
		project, err := h.projectReadRepo.FindByID(txCtx, projectVO)
		if err != nil {
			return nil, err
		}
		if project == nil {
			return nil, taskdomain.TaskNotFound()
		}

		ownerID := project.OwnerID().Value()
		adminAccess := project.AdminAccess()
		if !access.EnsureAccess(ownerID, adminAccess, cmd.AdminID) {
			return nil, projectdomain.Forbidden()
		}

		// delete aggregate
		if err := h.taskWriteRepo.Delete(txCtx, taskID); err != nil {
			return nil, err
		}

		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	// best-effort cache update
	if projectID != "" {
		_ = cache.UpdateProjectCache(ctx, cache.ProjectCacheDeps{
			ProjectReadRepo: h.projectReadRepo,
			ProjectCache:    h.projectCache,
		}, projectID)
	}

	return nil
}
