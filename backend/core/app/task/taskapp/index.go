package taskapp

import (
	create "ai_hub.com/app/core/app/task/commands/create"
	delete "ai_hub.com/app/core/app/task/commands/delete"
	update "ai_hub.com/app/core/app/task/commands/update"
	getbyproject "ai_hub.com/app/core/app/task/queries/getbyproject"
)

// Type aliases for task command and query handlers.

// create
type CreateTaskHandler = create.CreateTaskHandler
type CreateTaskCommand = create.CreateTaskCommand

var NewCreateTaskHandler = create.NewCreateTaskHandler

// update
type UpdateTaskHandler = update.UpdateTaskHandler
type UpdateTaskCommand = update.UpdateTaskCommand

var NewUpdateTaskHandler = update.NewUpdateTaskHandler

// delete
type DeleteTaskHandler = delete.DeleteTaskHandler
type DeleteTaskCommand = delete.DeleteTaskCommand

var NewDeleteTaskHandler = delete.NewDeleteTaskHandler

// query: get by project
type GetTasksByProjectHandler = getbyproject.GetTasksByProjectHandler
type GetTasksByProjectQuery = getbyproject.GetTasksByProjectQuery

var NewGetTasksByProjectHandler = getbyproject.NewGetTasksByProjectHandler
