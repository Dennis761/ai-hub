// src/core/application/task/update_task_handler.go
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

type UpdateTaskHandler struct {
	taskReadRepo    taskports.TaskReadRepository
	taskWriteRepo   taskports.TaskWriteRepository
	uow             taskports.UnitOfWorkPort
	projectReadRepo projectports.ProjectReadRepository
	projectCache    projectports.ProjectCachePort
}

func NewUpdateTaskHandler(
	taskReadRepo taskports.TaskReadRepository,
	taskWriteRepo taskports.TaskWriteRepository,
	uow taskports.UnitOfWorkPort,
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
) *UpdateTaskHandler {
	return &UpdateTaskHandler{
		taskReadRepo:    taskReadRepo,
		taskWriteRepo:   taskWriteRepo,
		uow:             uow,
		projectReadRepo: projectReadRepo,
		projectCache:    projectCache,
	}
}

func (h *UpdateTaskHandler) Update(ctx context.Context, cmd UpdateTaskCommand) (*taskdomain.Task, error) {
	var updatedTask *taskdomain.Task

	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate id
		taskID, err := taskdomain.NewTaskID(cmd.ID)
		if err != nil {
			return nil, err
		}

		// load aggregate
		task, err := h.taskReadRepo.FindByID(txCtx, taskID)
		if err != nil {
			return nil, err
		}
		if task == nil {
			return nil, taskdomain.TaskNotFound()
		}

		// check access via project
		projectID := task.ProjectID().Value()
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

		// apply changes

		// name
		if cmd.Name != nil {
			nameVO, err := taskdomain.NewTaskName(*cmd.Name)
			if err != nil {
				return nil, err
			}
			task.Rename(nameVO)

		}

		// description
		if cmd.Description != nil {
			empty := ""
			descVO, err := taskdomain.NewTaskDescription(&empty)
			if err != nil {
				return nil, err
			}
			task.SetDescription(descVO)
		} else {
			descVO, err := taskdomain.NewTaskDescription(cmd.Description)
			if err != nil {
				return nil, err
			}
			task.SetDescription(descVO)
		}

		// API method
		if cmd.APIMethod != nil {
			methodVO, err := taskdomain.NewAPIMethod(*cmd.APIMethod)
			if err != nil {
				return nil, err
			}
			task.SetAPIMethod(methodVO)
		}

		// status
		if cmd.Status != nil {
			statusVO, err := taskdomain.NewTaskStatus(*cmd.Status)
			if err != nil {
				return nil, err
			}
			if err := task.SetStatus(statusVO.Value()); err != nil {
				return nil, err
			}

		}

		// persist
		updated, err := h.taskWriteRepo.Update(txCtx, *task)
		if err != nil {
			return nil, err
		}
		updatedTask = updated

		return updated, nil
	})
	if err != nil {
		return nil, err
	}

	// best-effort cache update
	if updatedTask != nil {
		_ = cache.UpdateProjectCache(ctx, cache.ProjectCacheDeps{
			ProjectReadRepo: h.projectReadRepo,
			ProjectCache:    h.projectCache,
		}, updatedTask.ProjectID().Value())
	}

	return res.(*taskdomain.Task), nil
}
