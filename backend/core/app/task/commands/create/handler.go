// src/core/application/task/create_task_handler.go
package taskapp

import (
	"context"
	"time"

	"ai_hub.com/app/core/app/task/shared/access"
	"ai_hub.com/app/core/app/task/shared/cache"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/taskdomain"
	"ai_hub.com/app/core/ports/adminports"
	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/taskports"
)

type CreateTaskHandler struct {
	taskWriteRepo   taskports.TaskWriteRepository
	uow             taskports.UnitOfWorkPort
	projectReadRepo projectports.ProjectReadRepository
	projectCache    projectports.ProjectCachePort
	idGen           adminports.IDGenerator
}

func NewCreateTaskHandler(
	taskWriteRepo taskports.TaskWriteRepository,
	uow taskports.UnitOfWorkPort,
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
	idGen adminports.IDGenerator,
) *CreateTaskHandler {
	return &CreateTaskHandler{
		taskWriteRepo:   taskWriteRepo,
		uow:             uow,
		projectReadRepo: projectReadRepo,
		projectCache:    projectCache,
		idGen:           idGen,
	}
}

func (h *CreateTaskHandler) Create(ctx context.Context, cmd CreateTaskCommand) (*taskdomain.Task, error) {
	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate project and access
		projectIDVO, err := projectdomain.NewProjectID(cmd.ProjectID)
		if err != nil {
			return nil, err
		}
		project, err := h.projectReadRepo.FindByID(txCtx, projectIDVO)
		if err != nil {
			return nil, err
		}
		if project == nil {
			return nil, projectdomain.ProjectNotFound()
		}

		ownerID := project.OwnerID().Value()
		adminAccess := project.AdminAccess()
		if !access.EnsureAccess(ownerID, adminAccess, cmd.CreatedBy) {
			return nil, projectdomain.Forbidden()
		}

		// build IDs and VOs
		taskID, err := taskdomain.NewTaskID(h.idGen.NewID())
		if err != nil {
			return nil, err
		}

		nameVO, err := taskdomain.NewTaskName(cmd.Name)
		if err != nil {
			return nil, err
		}

		var descVO taskdomain.TaskDescription
		if cmd.Description == nil {
			descVO, err = taskdomain.NewTaskDescription(nil)
		} else {
			descVO, err = taskdomain.NewTaskDescription(cmd.Description)
		}
		if err != nil {
			return nil, err
		}

		projRef, err := taskdomain.NewTaskProjectID(projectIDVO.Value())
		if err != nil {
			return nil, err
		}

		apiMethodVO, err := taskdomain.NewAPIMethod(cmd.APIMethod)
		if err != nil {
			return nil, err
		}

		creatorVO, err := taskdomain.NewTaskCreatorID(cmd.CreatedBy)
		if err != nil {
			return nil, err
		}

		now := time.Now().UTC()
		if cmd.Now != nil {
			now = cmd.Now.UTC()
		}

		var statusPtr *string
		if cmd.Status != nil {
			stVO, err := taskdomain.NewTaskStatus(*cmd.Status)
			if err != nil {
				return nil, err
			}
			s := stVO.Value()
			statusPtr = &s
		}

		// create aggregate
		task, err := taskdomain.Create(taskdomain.CreateProps{
			ID:          taskID,
			Name:        nameVO,
			Description: &descVO,
			ProjectID:   projRef,
			APIMethod:   apiMethodVO,
			Status:      statusPtr,
			CreatedBy:   creatorVO,
			Now:         &now,
		})
		if err != nil {
			return nil, err
		}

		// persist
		created, err := h.taskWriteRepo.Create(txCtx, task)
		if err != nil {
			return nil, err
		}
		return created, nil
	})
	if err != nil {
		return nil, err
	}

	createdTask := res.(*taskdomain.Task)

	// best-effort cache update
	_ = cache.UpdateProjectCache(ctx, cache.ProjectCacheDeps{
		ProjectReadRepo: h.projectReadRepo,
		ProjectCache:    h.projectCache,
	}, createdTask.ProjectID().Value())

	return createdTask, nil
}
