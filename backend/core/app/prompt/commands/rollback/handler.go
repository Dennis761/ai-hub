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

type RollbackPromptHandler struct {
	promptReadRepo  promptports.PromptReadRepository
	promptWriteRepo promptports.PromptWriteRepository
	uow             promptports.UnitOfWorkPort
	taskReadRepo    taskports.TaskReadRepository
	projectReadRepo projectports.ProjectReadRepository
	projectCache    projectports.ProjectCachePort
}

func NewRollbackPromptHandler(
	promptReadRepo promptports.PromptReadRepository,
	promptWriteRepo promptports.PromptWriteRepository,
	uow promptports.UnitOfWorkPort,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
) *RollbackPromptHandler {
	return &RollbackPromptHandler{
		promptReadRepo:  promptReadRepo,
		promptWriteRepo: promptWriteRepo,
		uow:             uow,
		taskReadRepo:    taskReadRepo,
		projectReadRepo: projectReadRepo,
		projectCache:    projectCache,
	}
}

func (h *RollbackPromptHandler) Rollback(ctx context.Context, cmd RollbackPromptCommand) error {
	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate input
		missing := make([]string, 0, 3)
		if cmd.ID == "" {
			missing = append(missing, "_id")
		}
		if cmd.AdminID == "" {
			missing = append(missing, "adminID")
		}
		if cmd.Version < 1 {
			return nil, promptdomain.InvalidVersionNumber()
		}
		if len(missing) > 0 {
			return nil, promptdomain.MissingParameters(missing)
		}

		// load prompt aggregate
		promptID, err := promptdomain.NewPromptID(cmd.ID)
		if err != nil {
			return nil, err
		}
		prompt, err := h.promptReadRepo.FindByID(txCtx, promptID)
		if err != nil {
			return nil, err
		}
		if prompt == nil {
			return nil, promptdomain.PromptNotFound()
		}

		// resolve access through task → project
		taskID, err := taskdomain.NewTaskID(prompt.TaskID().Value())
		if err != nil {
			return nil, err
		}
		task, err := h.taskReadRepo.FindByID(txCtx, taskID)
		if err != nil {
			return nil, err
		}
		if task == nil {
			return nil, promptdomain.TaskNotFound()
		}
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
			return nil, promptdomain.TaskNotFound()
		}

		ownerID := project.OwnerID().Value()
		adminAccess := project.AdminAccess()
		if !access.EnsureAccess(ownerID, adminAccess, cmd.AdminID) {
			return nil, projectdomain.Forbidden()
		}

		// find history entry by version
		history := prompt.History()
		targetIdx := -1
		for i, snap := range history {
			if snap.Version() == cmd.Version {
				targetIdx = i
				break
			}
		}
		if targetIdx < 0 {
			return nil, promptdomain.HistoryIndexOutOfRange()
		}

		// rollback aggregate
		if err := prompt.RollbackToIndex(targetIdx); err != nil {
			return nil, promptdomain.HistoryIndexOutOfRange()
		}

		// persist updated aggregate
		saved, err := h.promptWriteRepo.Update(txCtx, prompt)
		if err != nil {
			return nil, err
		}

		// refresh cache (best effort)
		if saved != nil {
			_ = cache.UpdateProjectCache(txCtx, cache.ProjectCacheDeps{
				ProjectReadRepo: h.projectReadRepo,
				ProjectCache:    h.projectCache,
			}, projectID)
		}

		return struct{}{}, nil
	})

	return err
}
