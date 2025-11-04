package controllers

import (
	"net/http"

	adminapp "ai_hub.com/app/core/app/admin/adminapp"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	createHandler         *adminapp.CreateAdminHandler
	verifyHandler         *adminapp.VerifyAdminHandler
	startResetHandler     *adminapp.StartPasswordResetHandler
	confirmResetHandler   *adminapp.ConfirmResetCodeHandler
	changePasswordHandler *adminapp.ChangePasswordWithCodeHandler
	renameHandler         *adminapp.RenameAdminHandler
	deleteHandler         *adminapp.DeleteAdminHandler
	loginHandler          *adminapp.LoginAdminHandler
}

func NewAdminController(
	createH *adminapp.CreateAdminHandler,
	verifyH *adminapp.VerifyAdminHandler,
	startResetH *adminapp.StartPasswordResetHandler,
	confirmResetH *adminapp.ConfirmResetCodeHandler,
	changePwH *adminapp.ChangePasswordWithCodeHandler,
	renameH *adminapp.RenameAdminHandler,
	deleteH *adminapp.DeleteAdminHandler,
	loginH *adminapp.LoginAdminHandler,
) *AdminController {
	return &AdminController{
		createHandler:         createH,
		verifyHandler:         verifyH,
		startResetHandler:     startResetH,
		confirmResetHandler:   confirmResetH,
		changePasswordHandler: changePwH,
		renameHandler:         renameH,
		deleteHandler:         deleteH,
		loginHandler:          loginH,
	}
}

// POST /api/admins/login
func (ctrl *AdminController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	result, err := ctrl.loginHandler.Login(c, adminapp.LoginAdminCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// POST /api/admins/register
func (ctrl *AdminController) Register(c *gin.Context) {
	var cmd adminapp.CreateAdminCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.Error(err)
		return
	}

	admin, err := ctrl.createHandler.CreateAdmin(c, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"admin": admin})
}

// POST /api/admins/:id/verify
func (ctrl *AdminController) Verify(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := ctrl.verifyHandler.Verify(c, adminapp.VerifyAdminCommand{
		ID:   id,
		Code: req.Code,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/admins/reset/start
func (ctrl *AdminController) StartReset(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := ctrl.startResetHandler.StartPasswordReset(c, adminapp.StartPasswordResetCommand{
		Email: req.Email,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/admins/reset/confirm
func (ctrl *AdminController) ConfirmReset(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := ctrl.confirmResetHandler.ConfirmResetCode(c, adminapp.ConfirmResetCodeCommand{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/admins/reset/change
func (ctrl *AdminController) ChangePassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := ctrl.changePasswordHandler.ChangePasswordWithCode(c, adminapp.ChangePasswordWithCodeCommand{
		Email:       req.Email,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// PATCH /api/admins/:id/rename
func (ctrl *AdminController) Rename(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name *string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	result, err := ctrl.renameHandler.RenameAdmin(c, adminapp.RenameAdminCommand{
		ID:   id,
		Name: req.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// DELETE /api/admins/:id
func (ctrl *AdminController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.deleteHandler.DeleteAdmin(c, adminapp.DeleteAdminCommand{
		ID: id,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"ok": true})
}
