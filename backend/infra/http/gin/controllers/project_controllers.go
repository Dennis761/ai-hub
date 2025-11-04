package controllers

import (
	"net/http"

	projectapp "ai_hub.com/app/core/app/project/projectapp"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	createHandler  *projectapp.CreateProjectHandler
	updateHandler  *projectapp.UpdateProjectHandler
	deleteHandler  *projectapp.DeleteProjectHandler
	joinHandler    *projectapp.JoinProjectByNameHandler
	getMyHandler   *projectapp.GetMyProjectsHandler
	getByIDHandler *projectapp.GetProjectByIDHandler
}

func NewProjectController(
	createH *projectapp.CreateProjectHandler,
	updateH *projectapp.UpdateProjectHandler,
	deleteH *projectapp.DeleteProjectHandler,
	joinH *projectapp.JoinProjectByNameHandler,
	getMyH *projectapp.GetMyProjectsHandler,
	getByIDH *projectapp.GetProjectByIDHandler,
) *ProjectController {
	return &ProjectController{
		createHandler:  createH,
		updateHandler:  updateH,
		deleteHandler:  deleteH,
		joinHandler:    joinH,
		getMyHandler:   getMyH,
		getByIDHandler: getByIDH,
	}
}

// POST /api/projects
func (ctrl *ProjectController) Create(c *gin.Context) {
	type createProjectReq struct {
		ID          *string  `json:"_id"`
		Name        string   `json:"name"`
		APIKey      string   `json:"apiKey"`
		AdminAccess []string `json:"adminAccess"`
		Status      *string  `json:"status"`
	}

	var req createProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	ownerID := c.GetString("userID")

	result, err := ctrl.createHandler.Create(c, projectapp.CreateProjectCommand{
		ID:          req.ID,
		Name:        req.Name,
		APIKey:      req.APIKey,
		OwnerID:     ownerID,
		AdminAccess: req.AdminAccess,
		Status:      req.Status,
	})
	if err != nil {
		c.Error(err)
		return
	}

	p := result.ToPrimitives()
	c.JSON(http.StatusCreated, gin.H{
		"_id":         p.ID,
		"name":        p.Name,
		"status":      p.Status,
		"ownerId":     p.OwnerID,
		"adminAccess": p.AdminAccess,
		"createdAt":   p.CreatedAt,
		"updatedAt":   p.UpdatedAt,
	})
}

// PATCH /api/projects/:id
func (ctrl *ProjectController) Update(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name   *string `json:"name"`
		Status *string `json:"status"`
		APIKey *string `json:"apiKey"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	adminID := c.GetString("userID")

	updated, err := ctrl.updateHandler.Update(c, projectapp.UpdateProjectCommand{
		ID:      id,
		OwnerID: adminID,
		Name:    req.Name,
		Status:  req.Status,
		APIKey:  req.APIKey,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DELETE /api/projects/:id
func (ctrl *ProjectController) Delete(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	if err := ctrl.deleteHandler.Delete(c, projectapp.DeleteProjectCommand{
		ID:      id,
		AdminID: adminID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/projects/access/join
func (ctrl *ProjectController) Join(c *gin.Context) {
	var req struct {
		Name   string `json:"name"`
		APIKey string `json:"apiKey"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	adminID := c.GetString("userID")

	project, err := ctrl.joinHandler.JoinProjectByName(
		c,
		projectapp.JoinProjectByNameCommand{
			Name:    req.Name,
			APIKey:  req.APIKey,
			AdminID: adminID,
		},
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, project.ToPrimitives())
}

// GET /api/projects/my
func (ctrl *ProjectController) GetMy(c *gin.Context) {
	adminID := c.GetString("userID")

	result, err := ctrl.getMyHandler.GetMyProjects(c, projectapp.GetMyProjectsQuery{
		AdminID: adminID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GET /api/projects/:id
func (ctrl *ProjectController) GetByID(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("userID")

	dto, err := ctrl.getByIDHandler.GetProjectByID(c, projectapp.GetProjectByIDQuery{
		ProjectID: id,
		AdminID:   adminID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto)
}
