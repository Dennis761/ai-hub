package promptapp

import (
	"context"

	"ai_hub.com/app/core/app/prompt/shared/access"
	"ai_hub.com/app/core/app/prompt/shared/cache"
	prompttmpl "ai_hub.com/app/core/app/prompt/shared/prompttemplate"
	"ai_hub.com/app/core/app/prompt/shared/textcleaner"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"

	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/promptports"
	"ai_hub.com/app/core/ports/taskports"
)

type UpdatePromptHandler struct {
	promptReadRepo  promptports.PromptReadRepository
	promptWriteRepo promptports.PromptWriteRepository
	uow             promptports.UnitOfWorkPort
	taskReadRepo    taskports.TaskReadRepository
	projectReadRepo projectports.ProjectReadRepository
	projectCache    projectports.ProjectCachePort
}

func NewUpdatePromptHandler(
	promptReadRepo promptports.PromptReadRepository,
	promptWriteRepo promptports.PromptWriteRepository,
	uow promptports.UnitOfWorkPort,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
) *UpdatePromptHandler {
	return &UpdatePromptHandler{
		promptReadRepo:  promptReadRepo,
		promptWriteRepo: promptWriteRepo,
		uow:             uow,
		taskReadRepo:    taskReadRepo,
		projectReadRepo: projectReadRepo,
		projectCache:    projectCache,
	}
}

func (h *UpdatePromptHandler) Update(ctx context.Context, cmd UpdatePromptCommand) (*promptdomain.Prompt, error) {
	var projectID string

	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
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

		// resolve task → project → access
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

		ownerID := project.OwnerID().Value()
		adminAccess := project.AdminAccess()
		if !access.EnsureAccess(ownerID, adminAccess, cmd.AdminID) {
			return nil, projectdomain.Forbidden()
		}

		// apply updates
		prim := prompt.ToPrimitives()
		changed := false
		bump := false
		var newPromptText string

		// name
		if cmd.Name != nil {
			nameVO, err := promptdomain.NewPromptName(*cmd.Name)
			if err != nil {
				return nil, err
			}
			if prim.Name != nameVO.Value() {
				prompt.Rename(nameVO)
				changed = true
			}
		}

		// model
		if cmd.ModelID != nil {
			modelVO, err := promptdomain.NewModelRefID(*cmd.ModelID)
			if err != nil {
				return nil, err
			}
			if prim.ModelID != modelVO.Value() {
				prompt.RebindModel(modelVO)
				changed = true
			}
		}

		// execution order
		if cmd.ExecutionOrder != nil {
			orderVO, err := promptdomain.NewExecutionOrder(*cmd.ExecutionOrder)
			if err != nil {
				return nil, err
			}
			if prim.ExecutionOrder != orderVO.Value() {
				prompt.SetExecutionOrder(orderVO)
				changed = true
			}
		}

		// prompt text (bump)
		if cmd.PromptText != nil {
			textVO, err := promptdomain.NewPromptText(*cmd.PromptText)
			if err != nil {
				return nil, err
			}
			if prim.PromptText != textVO.Value() {
				newPromptText = textVO.Value()
				if err := prompt.SetPromptText(textVO, true); err != nil {
					return nil, err
				}
				changed = true
				bump = true
			}
		}

		// response text (no bump)
		if cmd.ResponseText != nil {
			cleaned := textcleaner.CleanText(*cmd.ResponseText)
			respVO, err := promptdomain.NewResponseText(&cleaned)
			if err != nil {
				return nil, err
			}
			curr := prim.ResponseText
			next := respVO.Value()
			same := (curr == nil && next == nil) ||
				(curr != nil && next != nil && *curr == *next)
			if !same {
				prompt.SetResponseText(respVO)
				changed = true
			}
		}

		// nothing to update
		if !changed {
			return prompt, nil
		}

		// validate placeholders when text changed
		if bump {
			apiParams := prompttmpl.ParseQueryParamsFromAPIMethod(task.APIMethod().Value())
			ph := prompttmpl.ExtractPlaceholders(newPromptText)
			if len(ph) == 0 {
				return nil, promptdomain.NoPlaceholdersProvided()
			}

			var missing []string
			for _, p := range ph {
				if _, ok := apiParams[p]; !ok {
					missing = append(missing, p)
				}
			}
			if len(missing) > 0 {
				return nil, promptdomain.MissingParameters(missing)
			}

			// keep version in sync with history
			nowPrim := prompt.ToPrimitives()
			maxHist := 0
			for _, hst := range nowPrim.History {
				if hst.Version > maxHist {
					maxHist = hst.Version
				}
			}
			allowedMax := maxHist + 1
			if nowPrim.Version > allowedMax {
				capped, err := promptdomain.NewVersionNumber(allowedMax)
				if err != nil {
					return nil, err
				}
				prompt.ForceSetVersion(capped)
			}
		}

		// persist
		updated, err := h.promptWriteRepo.Update(txCtx, prompt)
		if err != nil {
			return nil, err
		}

		return updated, nil
	})
	if err != nil {
		return nil, err
	}

	// best-effort cache refresh
	if projectID != "" {
		_ = cache.UpdateProjectCache(ctx, cache.ProjectCacheDeps{
			ProjectReadRepo: h.projectReadRepo,
			ProjectCache:    h.projectCache,
		}, projectID)
	}

	return res.(*promptdomain.Prompt), nil
}
