package promptapp

import (
	"context"
	"fmt"
	"time"

	"ai_hub.com/app/core/app/prompt/shared/access"
	prompttmpl "ai_hub.com/app/core/app/prompt/shared/prompttemplate"
	"ai_hub.com/app/core/app/prompt/shared/textcleaner"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/domain/promptdomain"
	"ai_hub.com/app/core/domain/taskdomain"

	"ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/core/ports/promptports"
	"ai_hub.com/app/core/ports/taskports"
)

// runs prompt and saves result to history
type RunPromptHandler struct {
	promptReadRepo  promptports.PromptReadRepository
	promptWriteRepo promptports.PromptWriteRepository
	uow             promptports.UnitOfWorkPort
	billingRunner   promptports.LLMBillingRunnerPort
	taskReadRepo    taskports.TaskReadRepository
	projectReadRepo projectports.ProjectReadRepository
}

func NewRunPromptHandler(
	promptReadRepo promptports.PromptReadRepository,
	promptWriteRepo promptports.PromptWriteRepository,
	uow promptports.UnitOfWorkPort,
	billingRunner promptports.LLMBillingRunnerPort,
	taskReadRepo taskports.TaskReadRepository,
	projectReadRepo projectports.ProjectReadRepository,
) *RunPromptHandler {
	return &RunPromptHandler{
		promptReadRepo:  promptReadRepo,
		promptWriteRepo: promptWriteRepo,
		uow:             uow,
		billingRunner:   billingRunner,
		taskReadRepo:    taskReadRepo,
		projectReadRepo: projectReadRepo,
	}
}

func (h *RunPromptHandler) Run(ctx context.Context, cmd RunPromptCommand) (*RunPromptResult, error) {
	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate IDs
		promptID, err := promptdomain.NewPromptID(cmd.ID)
		if err != nil {
			return nil, err
		}
		if cmd.RunByAdminID == "" {
			return nil, promptdomain.MissingParameters([]string{"runByAdminID"})
		}

		// load prompt aggregate
		prompt, err := h.promptReadRepo.FindByID(txCtx, promptID)
		if err != nil {
			return nil, err
		}
		if prompt == nil {
			return nil, promptdomain.PromptNotFound()
		}
		prim := prompt.ToPrimitives()

		// resolve task → project
		taskID, err := taskdomain.NewTaskID(prim.TaskID)
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

		// check access
		ownerID := project.OwnerID().Value()
		adminAccess := project.AdminAccess()
		if !access.EnsureAccess(ownerID, adminAccess, cmd.RunByAdminID) {
			return nil, projectdomain.Forbidden()
		}

		// prevent duplicate history items
		if len(prim.History) > 0 {
			last := prim.History[len(prim.History)-1]
			if last.Prompt == prim.PromptText {
				return nil, promptdomain.SamePromptConsecutive()
			}
		}

		// validate placeholders against api params
		apiMethod := task.APIMethod().Value()
		apiParams := prompttmpl.ParseQueryParamsFromAPIMethod(apiMethod)
		placeholders := prompttmpl.ExtractPlaceholders(prim.PromptText)
		if len(placeholders) == 0 {
			return nil, promptdomain.NoPlaceholdersProvided()
		}
		var missing []string
		for _, ph := range placeholders {
			if _, ok := apiParams[ph]; !ok {
				missing = append(missing, ph)
			}
		}
		if len(missing) > 0 {
			return nil, promptdomain.MissingParameters(missing)
		}

		// render final prompt
		formatted := prompttmpl.FormatPrompt(prim.PromptText, apiParams)

		// run LLM with billing
		runRes, err := h.billingRunner.RunWithBilling(txCtx, promptports.RunWithBillingParams{
			ModelID:         prim.ModelID,
			FormattedPrompt: formatted,
			RunByAdminID:    cmd.RunByAdminID,
		})
		if err != nil {
			return nil, err
		}

		// normalize LLM response
		cleanResponse := textcleaner.CleanText(runRes.ResponseText)

		// recalc allowed version
		maxHist := 0
		for _, item := range prim.History {
			if item.Version > maxHist {
				maxHist = item.Version
			}
		}
		allowedMax := maxHist + 1

		currVer := prompt.Version().Value()
		if currVer > allowedMax {
			cappedVO, err := promptdomain.NewVersionNumber(allowedMax)
			if err != nil {
				return nil, err
			}
			prompt.ForceSetVersion(cappedVO)
			currVer = cappedVO.Value()
		}

		// append history entry
		now := time.Now().UTC()
		entry, err := promptdomain.NewHistoryEntry(
			prim.PromptText,
			&cleanResponse,
			currVer,
			&now,
		)
		if err != nil {
			return nil, err
		}
		prompt.AddHistoryEntry(entry)

		// update response text if changed
		prevResp := ""
		if prim.ResponseText != nil {
			prevResp = *prim.ResponseText
		}
		if prevResp != cleanResponse {
			respVO, err := promptdomain.NewResponseText(&cleanResponse)
			if err != nil {
				return nil, err
			}
			prompt.SetResponseText(respVO)
		}

		// persist aggregate
		updated, err := h.promptWriteRepo.Update(txCtx, prompt)
		if err != nil {
			return nil, err
		}
		updatedPrim := updated.ToPrimitives()

		out := &RunPromptResult{
			ID:               updatedPrim.ID,
			Version:          currVer,
			PromptText:       updatedPrim.PromptText,
			ResponseText:     cleanResponse,
			ModelID:          updatedPrim.ModelID,
			UpdatedAt:        updatedPrim.UpdatedAt,
			Cost:             fmt.Sprintf("%.6f", runRes.Cost),
			BalanceRemaining: fmt.Sprintf("%.6f", runRes.BalanceRemaining),
		}
		return out, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*RunPromptResult), nil
}
