package promptapp

import (
	"context"
	"encoding/json"

	"ai_hub.com/app/core/app/prompt/shared/access"
	"ai_hub.com/app/core/app/prompt/shared/cache"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"

	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/promptports"
	"ai_hub.com/app/core/ports/taskports"
)

type ReorderPromptsHandler struct {
	promptReadRepo   promptports.PromptReadRepository
	promptWriteRepo  promptports.PromptWriteRepository
	uow              promptports.UnitOfWorkPort
	taskReadRepo     taskports.TaskReadRepository
	projectReadRepo  projectports.ProjectReadRepository
	projectCachePort projectports.ProjectCachePort
}

func NewReorderPromptsHandler(
	promptReadRepo promptports.PromptReadRepository,
	promptWriteRepo promptports.PromptWriteRepository,
	uow promptports.UnitOfWorkPort,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
) *ReorderPromptsHandler {
	return &ReorderPromptsHandler{
		promptReadRepo:   promptReadRepo,
		promptWriteRepo:  promptWriteRepo,
		uow:              uow,
		taskReadRepo:     taskReadRepo,
		projectReadRepo:  projectReadRepo,
		projectCachePort: projectCache,
	}
}

func (h *ReorderPromptsHandler) ReorderPrompts(ctx context.Context, cmd ReorderPromptsCommand) error {
	return h.reorder(ctx, cmd.Items, cmd.AdminID)
}

func (h *ReorderPromptsHandler) ReorderPromptsRaw(ctx context.Context, adminID string, raw []byte) error {
	ids, err := parseIDsFlexible(raw)
	if err != nil {
		return err
	}
	return h.reorder(ctx, ids, adminID)
}

func (h *ReorderPromptsHandler) reorder(ctx context.Context, ids []string, adminID string) error {
	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate input
		if len(ids) == 0 {
			return nil, promptdomain.MissingParameters([]string{"items"})
		}
		if adminID == "" {
			return nil, promptdomain.MissingParameters([]string{"adminID"})
		}

		ordered := uniquePreserveOrder(ids)
		if len(ordered) == 0 {
			return nil, promptdomain.MissingParameters([]string{"items"})
		}

		// build prompt IDs
		promptIDs := make([]promptdomain.PromptID, 0, len(ordered))
		for _, id := range ordered {
			vo, err := promptdomain.NewPromptID(id)
			if err != nil {
				return nil, err
			}
			promptIDs = append(promptIDs, vo)
		}

		// load all prompts
		prompts, err := h.promptReadRepo.FindByIDs(txCtx, promptIDs)
		if err != nil {
			return nil, err
		}
		if len(prompts) != len(ordered) {
			return nil, promptdomain.PromptNotFound()
		}

		// ensure all prompts belong to the same task
		byPromptID := make(map[string]*promptdomain.Prompt, len(prompts))
		var commonTaskID string
		for _, promptAgg := range prompts {
			promptID := promptAgg.ID().Value()
			taskID := promptAgg.TaskID().Value()

			byPromptID[promptID] = promptAgg

			if commonTaskID == "" {
				commonTaskID = taskID
				continue
			}
			if commonTaskID != taskID {
				return nil, promptdomain.InvalidExecutionOrder()
			}
		}

		// check access: task → project
		tid, err := taskdomain.NewTaskID(commonTaskID)
		if err != nil {
			return nil, err
		}
		taskAgg, err := h.taskReadRepo.FindByID(txCtx, tid)
		if err != nil {
			return nil, err
		}
		if taskAgg == nil {
			return nil, promptdomain.TaskNotFound()
		}

		projectID := taskAgg.ProjectID().Value()
		pid, err := projectdomain.NewProjectID(projectID)
		if err != nil {
			return nil, err
		}
		projectAgg, err := h.projectReadRepo.FindByID(txCtx, pid)
		if err != nil {
			return nil, err
		}
		if projectAgg == nil {
			return nil, promptdomain.TaskNotFound()
		}

		ownerID := projectAgg.OwnerID().Value()
		adminAccess := projectAgg.AdminAccess()
		if !access.EnsureAccess(ownerID, adminAccess, adminID) {
			return nil, projectdomain.Forbidden()
		}

		// update order 1..N
		updated := make([]*promptdomain.Prompt, 0, len(ordered))
		for i, id := range ordered {
			promptAgg := byPromptID[id]
			orderVO, err := promptdomain.NewExecutionOrder(i + 1)
			if err != nil {
				return nil, err
			}
			promptAgg.SetExecutionOrder(orderVO)
			updated = append(updated, promptAgg)
		}

		// persist changes
		if err := h.promptWriteRepo.UpdateMany(txCtx, updated); err != nil {
			return nil, err
		}

		// best-effort cache refresh
		_ = cache.UpdateProjectCache(txCtx, cache.ProjectCacheDeps{
			ProjectReadRepo: h.projectReadRepo,
			ProjectCache:    h.projectCachePort,
		}, projectID)

		return struct{}{}, nil
	})
	return err
}

func parseIDsFlexible(body []byte) ([]string, error) {
	// try plain []string
	var arrStr []string
	if err := json.Unmarshal(body, &arrStr); err == nil && len(arrStr) > 0 {
		return arrStr, nil
	}

	// try []any
	var arrAny []any
	if err := json.Unmarshal(body, &arrAny); err == nil && len(arrAny) > 0 {
		return extractIDs(arrAny), nil
	}

	// try {"items": ...}
	var wrap struct {
		Items json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &wrap); err == nil && len(wrap.Items) > 0 {
		var itemsStr []string
		if err2 := json.Unmarshal(wrap.Items, &itemsStr); err2 == nil && len(itemsStr) > 0 {
			return itemsStr, nil
		}
		var itemsAny []any
		if err3 := json.Unmarshal(wrap.Items, &itemsAny); err3 == nil && len(itemsAny) > 0 {
			return extractIDs(itemsAny), nil
		}
	}

	return nil, promptdomain.MissingParameters([]string{"items"})
}

func extractIDs(items []any) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		switch v := it.(type) {
		case string:
			if v != "" {
				out = append(out, v)
			}
		case map[string]any:
			if s, ok := v["_id"].(string); ok && s != "" {
				out = append(out, s)
			}
		}
	}
	return out
}

func uniquePreserveOrder(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
