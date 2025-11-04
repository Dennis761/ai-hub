// core/app/apikey/commands/commands.go
package apikeyapp

import (
	createcmd "ai_hub.com/app/core/app/apikey/commands/create"
	delcmd "ai_hub.com/app/core/app/apikey/commands/delete"
	updatecmd "ai_hub.com/app/core/app/apikey/commands/update"

	qbyid "ai_hub.com/app/core/app/apikey/queries/getkeybyid"
	qbyowner "ai_hub.com/app/core/app/apikey/queries/getkeysbyowner"
)

// command aliases for api key module

// create
type CreateAPIKeyHandler = createcmd.CreateAPIKeyHandler
type CreateAPIKeyCommand = createcmd.CreateAPIKeyCommand

var NewCreateAPIKeyHandler = createcmd.NewCreateAPIKeyHandler

// update
type UpdateAPIKeyHandler = updatecmd.UpdateAPIKeyHandler
type UpdateAPIKeyCommand = updatecmd.UpdateAPIKeyCommand

var NewUpdateAPIKeyHandler = updatecmd.NewUpdateAPIKeyHandler

// delete
type DeleteAPIKeyHandler = delcmd.DeleteAPIKeyHandler
type DeleteAPIKeyCommand = delcmd.DeleteAPIKeyCommand

var NewDeleteAPIKeyHandler = delcmd.NewDeleteAPIKeyHandler

// query: by id
type GetKeyByIDHandler = qbyid.GetKeyByIDHandler
type GetKeyByIDQuery = qbyid.GetKeyByIDQuery

var NewGetKeyByIDHandler = qbyid.NewGetKeyByIDHandler

// query: by owner
type GetKeysByOwnerHandler = qbyowner.GetKeysByOwnerHandler
type GetKeysByOwnerQuery = qbyowner.GetKeysByOwnerQuery

var NewGetKeysByOwnerHandler = qbyowner.NewGetKeysByOwnerHandler
