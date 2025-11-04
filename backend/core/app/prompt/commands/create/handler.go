package promptapp

import (
	"context"
	"time"

	"ai_hub.com/app/core/app/prompt/shared/access"
	"ai_hub.com/app/core/app/prompt/shared/cache"
	prompttmpl "ai_hub.com/app/core/app/prompt/shared/prompttemplate"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"

	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/promptports"
	"ai_hub.com/app/core/ports/taskports"
)

type CreatePromptHandler struct {
	promptWriteRepo   promptports.PromptWriteRepository
	uow               promptports.UnitOfWorkPort
	idGen             promptports.IDGenerator
	taskReadRepo      taskports.TaskReadRepository
	projectReadRepo   projectports.ProjectReadRepository
	projectCache      projectports.ProjectCachePort
	billingRunnerPort promptports.LLMBillingRunnerPort
}

func NewCreatePromptHandler(
	promptWriteRepo promptports.PromptWriteRepository,
	uow promptports.UnitOfWorkPort,
	idGen promptports.IDGenerator,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
	billingRunner promptports.LLMBillingRunnerPort,
) *CreatePromptHandler {
	return &CreatePromptHandler{
		promptWriteRepo:   promptWriteRepo,
		uow:               uow,
		idGen:             idGen,
		taskReadRepo:      taskReadRepo,
		projectReadRepo:   projectReadRepo,
		projectCache:      projectCache,
		billingRunnerPort: billingRunner,
	}
}

func (h *CreatePromptHandler) Create(
	ctx context.Context,
	cmd CreatePromptCommand,
) (*promptdomain.Prompt, error) {
	var (
		createdPrompt *promptdomain.Prompt
		projectID     string
	)

	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// check task
		taskID, err := taskdomain.NewTaskID(cmd.TaskID)
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

		taskPrim := task.ToPrimitives()
		projectID = taskPrim.ProjectID

		// check project + ACL
		projectIDVO, err := projectdomain.NewProjectID(projectID)
		if err != nil {
			return nil, err
		}
		project, err := h.projectReadRepo.FindByID(txCtx, projectIDVO)
		if err != nil {
			return nil, err
		}
		if project == nil {
			return nil, promptdomain.TaskNotFound()
		}
		if !access.EnsureAccess(
			project.OwnerID().Value(),
			project.AdminAccess(),
			cmd.CreatedBy,
		) {
			return nil, projectdomain.Forbidden()
		}

		// derive params from task API method
		apiParams := prompttmpl.ParseQueryParamsFromAPIMethod(taskPrim.APIMethod)

		// validate placeholders
		textPlaceholders := prompttmpl.ExtractPlaceholders(cmd.PromptText)
		if len(textPlaceholders) == 0 {
			return nil, promptdomain.NoPlaceholdersProvided()
		}
		var missing []string
		for _, ph := range textPlaceholders {
			if _, ok := apiParams[ph]; !ok {
				missing = append(missing, ph)
			}
		}
		if len(missing) > 0 {
			return nil, promptdomain.MissingParameters(missing)
		}

		// build formatted prompt
		formatted := prompttmpl.FormatPrompt(cmd.PromptText, apiParams)

		// billing check
		if _, err := h.billingRunnerPort.RunWithBilling(txCtx, promptports.RunWithBillingParams{
			ModelID:         cmd.ModelID,
			FormattedPrompt: formatted,
			RunByAdminID:    cmd.CreatedBy,
		}); err != nil {
			return nil, err
		}

		// build aggregate
		prompt, err := h.buildPromptAggregate(cmd)
		if err != nil {
			return nil, err
		}

		// persist
		created, err := h.promptWriteRepo.Create(txCtx, prompt)
		if err != nil {
			return nil, err
		}
		createdPrompt = created

		return struct{}{}, nil
	})
	if err != nil {
		return nil, err
	}

	// best-effort cache update
	_ = cache.UpdateProjectCache(ctx, cache.ProjectCacheDeps{
		ProjectReadRepo: h.projectReadRepo,
		ProjectCache:    h.projectCache,
	}, projectID)

	return createdPrompt, nil
}

func (h *CreatePromptHandler) buildPromptAggregate(
	cmd CreatePromptCommand,
) (*promptdomain.Prompt, error) {
	id := h.idGen.NewID()

	promptID, err := promptdomain.NewPromptID(id)
	if err != nil {
		return nil, err
	}

	taskRef, err := promptdomain.NewTaskRefID(cmd.TaskID)
	if err != nil {
		return nil, err
	}
	nameVO, err := promptdomain.NewPromptName(cmd.Name)
	if err != nil {
		return nil, err
	}
	modelRef, err := promptdomain.NewModelRefID(cmd.ModelID)
	if err != nil {
		return nil, err
	}
	promptText, err := promptdomain.NewPromptText(cmd.PromptText)
	if err != nil {
		return nil, err
	}

	var respPtr *promptdomain.ResponseText
	if cmd.ResponseText != nil {
		resp, err := promptdomain.NewResponseText(cmd.ResponseText)
		if err != nil {
			return nil, err
		}
		respPtr = &resp
	}

	var execOrderPtr *promptdomain.ExecutionOrder
	if cmd.ExecutionOrder != nil {
		execOrder, err := promptdomain.NewExecutionOrder(*cmd.ExecutionOrder)
		if err != nil {
			return nil, err
		}
		execOrderPtr = &execOrder
	}

	now := time.Now().UTC()
	prompt, err := promptdomain.Create(promptdomain.CreateProps{
		ID:             promptID,
		TaskID:         taskRef,
		Name:           nameVO,
		ModelID:        modelRef,
		PromptText:     promptText,
		ResponseText:   respPtr,
		History:        nil,
		ExecutionOrder: execOrderPtr,
		Version:        nil,
		Now:            &now,
	})
	if err != nil {
		return nil, err
	}
	return &prompt, nil
}
