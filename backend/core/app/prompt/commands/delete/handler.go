package promptapp

import (
	"context"

	"ai_hub.com/app/core/app/prompt/shared/access"
	"ai_hub.com/app/core/app/prompt/shared/cache"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"
	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/promptports"
	"ai_hub.com/app/core/ports/taskports"
)

type DeletePromptHandler struct {
	promptReadRepo   promptports.PromptReadRepository
	promptWriteRepo  promptports.PromptWriteRepository
	uow              promptports.UnitOfWorkPort
	taskReadRepo     taskports.TaskReadRepository
	projectReadRepo  projectports.ProjectReadRepository
	projectCachePort projectports.ProjectCachePort
}

func NewDeletePromptHandler(
	promptReadRepo promptports.PromptReadRepository,
	promptWriteRepo promptports.PromptWriteRepository,
	uow promptports.UnitOfWorkPort,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
	projectCachePort projectports.ProjectCachePort,
) *DeletePromptHandler {
	return &DeletePromptHandler{
		promptReadRepo:   promptReadRepo,
		promptWriteRepo:  promptWriteRepo,
		uow:              uow,
		taskReadRepo:     taskReadRepo,
		projectReadRepo:  projectReadRepo,
		projectCachePort: projectCachePort,
	}
}

func (h *DeletePromptHandler) Delete(
	ctx context.Context,
	cmd DeletePromptCommand,
) error {
	var projectID string

	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate id
		promptID, err := promptdomain.NewPromptID(cmd.ID)
		if err != nil {
			return nil, err
		}

		// load prompt
		prompt, err := h.promptReadRepo.FindByID(txCtx, promptID)
		if err != nil {
			return nil, err
		}
		// idempotent delete
		if prompt == nil {
			return struct{}{}, nil
		}

		// load task
		taskID := prompt.TaskID().Value()
		if taskID == "" {
			return nil, promptdomain.TaskNotFound()
		}
		taskVO, err := taskdomain.NewTaskID(taskID)
		if err != nil {
			return nil, err
		}
		task, err := h.taskReadRepo.FindByID(txCtx, taskVO)
		if err != nil {
			return nil, err
		}
		if task == nil {
			return nil, promptdomain.TaskNotFound()
		}

		// load project & ACL
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
			return nil, promptdomain.TaskNotFound()
		}
		if !access.EnsureAccess(
			project.OwnerID().Value(),
			project.AdminAccess(),
			cmd.AdminID,
		) {
			return nil, projectdomain.Forbidden()
		}

		// delete prompt
		if err := h.promptWriteRepo.Delete(txCtx, promptID); err != nil {
			return nil, err
		}

		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	// best-effort cache update
	if projectID != "" && h.projectCachePort != nil {
		_ = cache.UpdateProjectCache(ctx, cache.ProjectCacheDeps{
			ProjectReadRepo: h.projectReadRepo,
			ProjectCache:    h.projectCachePort,
		}, projectID)
	}

	return nil
}
