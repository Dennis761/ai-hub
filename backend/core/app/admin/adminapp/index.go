package adminapp

import (
	delete "ai_hub.com/app/core/app/admin/commands/delete"
	login "ai_hub.com/app/core/app/admin/commands/login"
	register "ai_hub.com/app/core/app/admin/commands/register"
	rename "ai_hub.com/app/core/app/admin/commands/rename"
	change "ai_hub.com/app/core/app/admin/commands/resetPassword/change"
	confirm "ai_hub.com/app/core/app/admin/commands/resetPassword/confirm"
	start "ai_hub.com/app/core/app/admin/commands/resetPassword/start"
	verify "ai_hub.com/app/core/app/admin/commands/verify"
)

// command aliases for admin module

// registration
type CreateAdminHandler = register.CreateAdminHandler
type CreateAdminCommand = register.CreateAdminCommand

var NewCreateAdminHandler = register.NewCreateAdminHandler

// verification
type VerifyAdminHandler = verify.VerifyAdminHandler
type VerifyAdminCommand = verify.VerifyAdminCommand

var NewVerifyAdminHandler = verify.NewVerifyAdminHandler

// password reset flow
type StartPasswordResetHandler = start.StartPasswordResetHandler
type StartPasswordResetCommand = start.StartPasswordResetCommand

var NewStartPasswordResetHandler = start.NewStartPasswordResetHandler

type ConfirmResetCodeHandler = confirm.ConfirmResetCodeHandler
type ConfirmResetCodeCommand = confirm.ConfirmResetCodeCommand

var NewConfirmResetCodeHandler = confirm.NewConfirmResetCodeHandler

type ChangePasswordWithCodeHandler = change.ChangePasswordWithCodeHandler
type ChangePasswordWithCodeCommand = change.ChangePasswordWithCodeCommand

var NewChangePasswordWithCodeHandler = change.NewChangePasswordWithCodeHandler

// rename admin
type RenameAdminHandler = rename.RenameAdminHandler
type RenameAdminCommand = rename.RenameAdminCommand

var NewRenameAdminHandler = rename.NewRenameAdminHandler

// delete admin
type DeleteAdminHandler = delete.DeleteAdminHandler
type DeleteAdminCommand = delete.DeleteAdminCommand

var NewDeleteAdminHandler = delete.NewDeleteAdminHandler

// login
type LoginAdminHandler = login.LoginAdminHandler
type LoginAdminCommand = login.LoginAdminCommand
type LoginAdminResponse = login.LoginAdminResponse
type LoginAdminResult = login.LoginAdminResult

var NewLoginAdminHandler = login.NewLoginAdminHandler
