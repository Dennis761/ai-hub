package routes

import (
	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"github.com/gin-gonic/gin"
)

func ProjectRoutes(r gin.IRouter, controller *controllers.ProjectController, verifier middlewares.TokenVerifier) {
	group := r.Group("/projects")
	group.Use(middlewares.Auth(verifier))
	{
		// Commands
		group.POST("", controller.Create)
		group.PATCH("/:id", controller.Update)
		group.DELETE("/:id", controller.Delete)
		group.POST("/join", controller.Join)

		// Queries
		group.GET("/my-projects", controller.GetMy)
		group.GET("/:id", controller.GetByID)
	}
}

type ProjectRouter struct {
	controller *controllers.ProjectController
	verifier   middlewares.TokenVerifier
}

func NewProjectRoutes(controller *controllers.ProjectController, verifier middlewares.TokenVerifier) *ProjectRouter {
	return &ProjectRouter{
		controller: controller,
		verifier:   verifier,
	}
}

func (pr *ProjectRouter) Mount(r gin.IRouter) {
	ProjectRoutes(r, pr.controller, pr.verifier)
}
