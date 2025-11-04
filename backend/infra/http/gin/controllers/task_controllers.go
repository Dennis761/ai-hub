package controllers

import (
	"net/http"

	taskapp "ai_hub.com/app/core/app/task/taskapp"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	createHandler    *taskapp.CreateTaskHandler
	updateHandler    *taskapp.UpdateTaskHandler
	deleteHandler    *taskapp.DeleteTaskHandler
	getByProjectHdlr *taskapp.GetTasksByProjectHandler
}

func NewTaskController(
	createH *taskapp.CreateTaskHandler,
	updateH *taskapp.UpdateTaskHandler,
	deleteH *taskapp.DeleteTaskHandler,
	getByProjectH *taskapp.GetTasksByProjectHandler,
) *TaskController {
	return &TaskController{
		createHandler:    createH,
		updateHandler:    updateH,
		deleteHandler:    deleteH,
		getByProjectHdlr: getByProjectH,
	}
}

// POST /api/tasks
func (ctrl *TaskController) Create(c *gin.Context) {
	type reqBody struct {
		ID          *string `json:"_id"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		ProjectID   string  `json:"projectId"`
		APIMethod   string  `json:"apiMethod"`
		Status      *string `json:"status"`
	}

	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	adminID := c.GetString("userID")

	res, err := ctrl.createHandler.Create(c, taskapp.CreateTaskCommand{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		APIMethod:   req.APIMethod,
		Status:      req.Status,
		CreatedBy:   adminID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// PATCH /api/tasks/:id
func (ctrl *TaskController) Update(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		APIMethod   *string `json:"apiMethod"`
		Status      *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := ctrl.updateHandler.Update(c, taskapp.UpdateTaskCommand{
		ID:          id,
		AdminID:     adminID,
		Name:        req.Name,
		Description: req.Description,
		APIMethod:   req.APIMethod,
		Status:      req.Status,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// src/infra/http/gin/controllers/task_controller.go
func (ctrl *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	if err := ctrl.deleteHandler.Delete(c, taskapp.DeleteTaskCommand{
		ID:      id,
		AdminID: adminID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"ok": true})
}

// GET /api/tasks/:projectId/tasks
func (ctrl *TaskController) GetByProject(c *gin.Context) {
	projectID := c.Param("projectId")
	adminID := c.GetString("userID")

	items, err := ctrl.getByProjectHdlr.GetTasksByProject(
		c,
		taskapp.GetTasksByProjectQuery{
			ProjectID: projectID,
			AdminID:   adminID,
		},
	)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, items)
}
