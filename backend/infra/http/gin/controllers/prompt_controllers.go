package controllers

import (
	"net/http"

	promptapp "ai_hub.com/app/core/app/prompt/promptapp"
	"github.com/gin-gonic/gin"
)

type PromptController struct {
	createHandler    *promptapp.CreatePromptHandler
	updateHandler    *promptapp.UpdatePromptHandler
	deleteHandler    *promptapp.DeletePromptHandler
	rollbackHandler  *promptapp.RollbackPromptHandler
	reorderHandler   *promptapp.ReorderPromptsHandler
	runHandler       *promptapp.RunPromptHandler
	getByTaskHandler *promptapp.GetPromptsByTaskHandler
	getByIDHandler   *promptapp.GetPromptByIDHandler
}

func NewPromptController(
	createH *promptapp.CreatePromptHandler,
	updateH *promptapp.UpdatePromptHandler,
	deleteH *promptapp.DeletePromptHandler,
	rollbackH *promptapp.RollbackPromptHandler,
	reorderH *promptapp.ReorderPromptsHandler,
	runH *promptapp.RunPromptHandler,
	getByTaskH *promptapp.GetPromptsByTaskHandler,
	getByIDH *promptapp.GetPromptByIDHandler,
) *PromptController {
	return &PromptController{
		createHandler:    createH,
		updateHandler:    updateH,
		deleteHandler:    deleteH,
		rollbackHandler:  rollbackH,
		reorderHandler:   reorderH,
		runHandler:       runH,
		getByTaskHandler: getByTaskH,
		getByIDHandler:   getByIDH,
	}
}

// POST /api/prompts
func (ctrl *PromptController) Create(c *gin.Context) {
	type createReq struct {
		ID             *string `json:"_id"`
		TaskID         string  `json:"taskId"`
		Name           string  `json:"name"`
		ModelID        string  `json:"modelId"`
		PromptText     string  `json:"promptText"`
		ResponseText   *string `json:"responseText"`
		ExecutionOrder *int    `json:"executionOrder"`
	}

	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	adminID := c.GetString("userID")

	prompt, err := ctrl.createHandler.Create(c, promptapp.CreatePromptCommand{
		ID:             req.ID,
		TaskID:         req.TaskID,
		Name:           req.Name,
		ModelID:        req.ModelID,
		PromptText:     req.PromptText,
		ResponseText:   req.ResponseText,
		ExecutionOrder: req.ExecutionOrder,
		CreatedBy:      adminID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	p := prompt.ToPrimitives()
	c.JSON(http.StatusCreated, gin.H{
		"_id":            p.ID,
		"taskId":         p.TaskID,
		"name":           p.Name,
		"modelId":        p.ModelID,
		"promptText":     p.PromptText,
		"responseText":   p.ResponseText,
		"history":        p.History,
		"executionOrder": p.ExecutionOrder,
		"version":        p.Version,
		"createdAt":      p.CreatedAt,
		"updatedAt":      p.UpdatedAt,
	})
}

// PATCH /api/prompts/:id
func (ctrl *PromptController) Update(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	var req struct {
		Name           *string `json:"name"`
		ModelID        *string `json:"modelId"`
		PromptText     *string `json:"promptText"`
		ResponseText   *string `json:"responseText"`
		ExecutionOrder *int    `json:"executionOrder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	entity, err := ctrl.updateHandler.Update(c, promptapp.UpdatePromptCommand{
		ID:             id,
		AdminID:        adminID,
		Name:           req.Name,
		ModelID:        req.ModelID,
		PromptText:     req.PromptText,
		ResponseText:   req.ResponseText,
		ExecutionOrder: req.ExecutionOrder,
	})
	if err != nil {
		c.Error(err)
		return
	}

	p := entity.ToPrimitives()
	c.JSON(http.StatusOK, gin.H{
		"_id":            p.ID,
		"taskId":         p.TaskID,
		"name":           p.Name,
		"modelId":        p.ModelID,
		"promptText":     p.PromptText,
		"responseText":   p.ResponseText,
		"history":        p.History,
		"executionOrder": p.ExecutionOrder,
		"version":        p.Version,
		"createdAt":      p.CreatedAt,
		"updatedAt":      p.UpdatedAt,
	})
}

// DELETE /api/prompts/:id
func (ctrl *PromptController) Delete(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	if err := ctrl.deleteHandler.Delete(c, promptapp.DeletePromptCommand{
		ID:      id,
		AdminID: adminID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// POST /api/prompts/:id/rollback
func (ctrl *PromptController) Rollback(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	var req struct {
		Version int `json:"version"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := ctrl.rollbackHandler.Rollback(c, promptapp.RollbackPromptCommand{
		ID:      id,
		Version: req.Version,
		AdminID: adminID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/prompts/reorder
func (ctrl *PromptController) Reorder(c *gin.Context) {
	adminID := c.GetString("userID")

	type reqBody struct {
		Items []string `json:"items"`
	}
	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := ctrl.reorderHandler.ReorderPrompts(c, promptapp.ReorderPromptsCommand{
		Items:   req.Items,
		AdminID: adminID,
	}); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/prompts/:id/run
func (ctrl *PromptController) Run(c *gin.Context) {
	promptID := c.Param("id")
	adminID := c.GetString("userID")

	result, err := ctrl.runHandler.Run(c, promptapp.RunPromptCommand{
		ID:           promptID,
		RunByAdminID: adminID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /api/prompts/task/:taskId
func (ctrl *PromptController) GetByTask(c *gin.Context) {
	taskID := c.Param("taskId")
	adminID := c.GetString("userID")

	result, err := ctrl.getByTaskHandler.GetPromptsByTask(c, promptapp.GetPromptsByTaskQuery{
		TaskID:  taskID,
		AdminID: adminID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /api/prompts/:id
func (ctrl *PromptController) GetByID(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	result, err := ctrl.getByIDHandler.GetPromptByID(c, promptapp.GetPromptByIDQuery{
		ID:      id,
		AdminID: adminID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}
