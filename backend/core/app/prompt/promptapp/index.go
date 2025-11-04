package promptapp

import (
	create "ai_hub.com/app/core/app/prompt/commands/create"
	delete "ai_hub.com/app/core/app/prompt/commands/delete"
	reorder "ai_hub.com/app/core/app/prompt/commands/reorder"
	rollback "ai_hub.com/app/core/app/prompt/commands/rollback"
	run "ai_hub.com/app/core/app/prompt/commands/run"
	update "ai_hub.com/app/core/app/prompt/commands/update"

	getbyid "ai_hub.com/app/core/app/prompt/queries/getbyid"
	getbytask "ai_hub.com/app/core/app/prompt/queries/getbytask"
)

// prompt module command/query aliases

// create
type CreatePromptHandler = create.CreatePromptHandler
type CreatePromptCommand = create.CreatePromptCommand

var NewCreatePromptHandler = create.NewCreatePromptHandler

// delete
type DeletePromptHandler = delete.DeletePromptHandler
type DeletePromptCommand = delete.DeletePromptCommand

var NewDeletePromptHandler = delete.NewDeletePromptHandler

// reorder
type ReorderPromptsHandler = reorder.ReorderPromptsHandler
type ReorderPromptsCommand = reorder.ReorderPromptsCommand

var NewReorderPromptsHandler = reorder.NewReorderPromptsHandler

// rollback
type RollbackPromptHandler = rollback.RollbackPromptHandler
type RollbackPromptCommand = rollback.RollbackPromptCommand

var NewRollbackPromptHandler = rollback.NewRollbackPromptHandler

// run
type RunPromptHandler = run.RunPromptHandler
type RunPromptCommand = run.RunPromptCommand

var NewRunPromptHandler = run.NewRunPromptHandler

// update
type UpdatePromptHandler = update.UpdatePromptHandler
type UpdatePromptCommand = update.UpdatePromptCommand

var NewUpdatePromptHandler = update.NewUpdatePromptHandler

// queries
type GetPromptByIDHandler = getbyid.GetPromptByIDHandler
type GetPromptByIDQuery = getbyid.GetPromptByIDQuery

var NewGetPromptByIDHandler = getbyid.NewGetPromptByIDHandler

type GetPromptsByTaskHandler = getbytask.GetPromptsByTaskHandler
type GetPromptsByTaskQuery = getbytask.GetPromptsByTaskQuery

var NewGetPromptsByTaskHandler = getbytask.NewGetPromptsByTaskHandler
