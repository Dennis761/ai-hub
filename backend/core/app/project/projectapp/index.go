// src/core/app/project/commands/commands.go
package projectapp

import (
	create "ai_hub.com/app/core/app/project/commands/create"
	deletecmd "ai_hub.com/app/core/app/project/commands/delete"
	join "ai_hub.com/app/core/app/project/commands/join"
	update "ai_hub.com/app/core/app/project/commands/update"

	getmyprojects "ai_hub.com/app/core/app/project/queries/getmyprojects"
	getprojectbyid "ai_hub.com/app/core/app/project/queries/getprojectbyid"
)

// command aliases for project module

// create
type CreateProjectHandler = create.CreateProjectHandler
type CreateProjectCommand = create.CreateProjectCommand

var NewCreateProjectHandler = create.NewCreateProjectHandler

// update
type UpdateProjectHandler = update.UpdateProjectHandler
type UpdateProjectCommand = update.UpdateProjectCommand

var NewUpdateProjectHandler = update.NewUpdateProjectHandler

// delete
type DeleteProjectHandler = deletecmd.DeleteProjectHandler
type DeleteProjectCommand = deletecmd.DeleteProjectCommand

var NewDeleteProjectHandler = deletecmd.NewDeleteProjectHandler

// join
type JoinProjectByNameHandler = join.JoinProjectByNameHandler
type JoinProjectByNameCommand = join.JoinProjectByNameCommand

var NewJoinProjectByNameHandler = join.NewJoinProjectByNameHandler

// query: all projects of a user
type GetMyProjectsHandler = getmyprojects.GetMyProjectsHandler
type GetMyProjectsQuery = getmyprojects.GetMyProjectsQuery

var NewGetMyProjectsHandler = getmyprojects.NewGetMyProjectsHandler

// query: project by id
type GetProjectByIDHandler = getprojectbyid.GetProjectByIDHandler
type GetProjectByIDQuery = getprojectbyid.GetProjectByIDQuery

var NewGetProjectByIDHandler = getprojectbyid.NewGetProjectByIDHandler
