package promptapp

import (
	"context"

	"ai_hub.com/app/core/app/prompt/shared/access"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"

	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/promptports"
	"ai_hub.com/app/core/ports/taskports"
)

type GetPromptByIDHandler struct {
	promptReadRepo  promptports.PromptReadRepository
	taskReadRepo    taskports.TaskReadRepository
	projectReadRepo projectports.ProjectReadRepository
}

func NewGetPromptByIDHandler(
	promptReadRepo promptports.PromptReadRepository,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
) *GetPromptByIDHandler {
	return &GetPromptByIDHandler{
		promptReadRepo:  promptReadRepo,
		taskReadRepo:    taskReadRepo,
		projectReadRepo: projectReadRepo,
	}
}

func (h *GetPromptByIDHandler) GetPromptByID(ctx context.Context, q GetPromptByIDQuery) (*GetPromptByIDResult, error) {
	// validate input
	if q.ID == "" || q.AdminID == "" {
		missing := make([]string, 0, 2)
		if q.ID == "" {
			missing = append(missing, "_id")
		}
		if q.AdminID == "" {
			missing = append(missing, "adminID")
		}
		return nil, promptdomain.MissingParameters(missing)
	}

	promptID, err := promptdomain.NewPromptID(q.ID)
	if err != nil {
		return nil, err
	}

	// load prompt aggregate
	prompt, err := h.promptReadRepo.FindByID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	if prompt == nil {
		return nil, promptdomain.PromptNotFound()
	}

	// access via task → project
	taskID, err := taskdomain.NewTaskID(prompt.TaskID().Value())
	if err != nil {
		return nil, err
	}
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

	// project to DTO
	p := prompt.ToPrimitives()

	hist := make([]HistoryItem, 0, len(p.History))
	for _, it := range p.History {
		hist = append(hist, HistoryItem{
			Prompt:    it.Prompt,
			Response:  it.Response,
			Version:   it.Version,
			CreatedAt: it.CreatedAt,
		})
	}

	return &GetPromptByIDResult{
		ID:             p.ID,
		TaskID:         p.TaskID,
		Name:           p.Name,
		ModelID:        p.ModelID,
		PromptText:     p.PromptText,
		ResponseText:   p.ResponseText,
		History:        hist,
		ExecutionOrder: p.ExecutionOrder,
		Version:        p.Version,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}, nil
}
