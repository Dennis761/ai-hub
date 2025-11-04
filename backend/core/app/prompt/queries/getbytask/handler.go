package promptapp

import (
	"context"
	"sort"
	"time"

	"ai_hub.com/app/core/app/prompt/shared/access"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"

	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/promptports"
	"ai_hub.com/app/core/ports/taskports"
)

// flat list of prompts for a task
type GetPromptsByTaskResult []PromptListItem

type GetPromptsByTaskHandler struct {
	promptReadRepo  promptports.PromptReadRepository
	taskReadRepo    taskports.TaskReadRepository
	projectReadRepo projectports.ProjectReadRepository
}

func NewGetPromptsByTaskHandler(
	promptReadRepo promptports.PromptReadRepository,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
) *GetPromptsByTaskHandler {
	return &GetPromptsByTaskHandler{
		promptReadRepo:  promptReadRepo,
		taskReadRepo:    taskReadRepo,
		projectReadRepo: projectReadRepo,
	}
}

func (h *GetPromptsByTaskHandler) GetPromptsByTask(
	ctx context.Context,
	q GetPromptsByTaskQuery,
) (GetPromptsByTaskResult, error) {
	// validate input
	missing := make([]string, 0, 2)
	if q.TaskID == "" {
		missing = append(missing, "taskID")
	}
	if q.AdminID == "" {
		missing = append(missing, "adminID")
	}
	if len(missing) > 0 {
		return nil, promptdomain.MissingParameters(missing)
	}

	// normalize task IDs
	taskRef, err := promptdomain.NewTaskRefID(q.TaskID)
	if err != nil {
		return nil, err
	}
	taskID, err := taskdomain.NewTaskID(taskRef.Value())
	if err != nil {
		return nil, err
	}

	// check access via task → project
	task, err := h.taskReadRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, promptdomain.TaskNotFound()
	}

	projectID, err := projectdomain.NewProjectID(task.ProjectID().Value())
	if err != nil {
		return nil, err
	}
	project, err := h.projectReadRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, promptdomain.TaskNotFound()
	}

	projectPrim := project.ToPrimitives()
	if !access.EnsureAccess(projectPrim.OwnerID, projectPrim.AdminAccess, q.AdminID) {
		return nil, projectdomain.Forbidden()
	}

	// load all prompts for task
	prompts, err := h.promptReadRepo.FindAllByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// sort by executionOrder ASC, then createdAt ASC
	sort.SliceStable(prompts, func(i, j int) bool {
		pi := prompts[i].ToPrimitives()
		pj := prompts[j].ToPrimitives()

		if pi.ExecutionOrder != pj.ExecutionOrder {
			return pi.ExecutionOrder < pj.ExecutionOrder
		}

		ti, _ := time.Parse(time.RFC3339, pi.CreatedAt)
		tj, _ := time.Parse(time.RFC3339, pj.CreatedAt)
		return ti.Before(tj)
	})

	// project aggregates to DTO
	items := make([]PromptListItem, 0, len(prompts))
	for _, prompt := range prompts {
		p := prompt.ToPrimitives()

		history := make([]HistoryItem, 0, len(p.History))
		for _, hItem := range p.History {
			history = append(history, HistoryItem{
				Prompt:    hItem.Prompt,
				Response:  *hItem.Response,
				Version:   hItem.Version,
				CreatedAt: hItem.CreatedAt,
			})
		}

		items = append(items, PromptListItem{
			ID:             p.ID,
			TaskID:         p.TaskID,
			Name:           p.Name,
			ModelID:        p.ModelID,
			PromptText:     p.PromptText,
			ResponseText:   p.ResponseText,
			History:        history,
			ExecutionOrder: p.ExecutionOrder,
			Version:        p.Version,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		})
	}

	return items, nil
}
