package routes

import (
	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"github.com/gin-gonic/gin"
)

func TaskRoutes(r gin.IRouter, controller *controllers.TaskController, verifier middlewares.TokenVerifier) {
	group := r.Group("/tasks")
	group.Use(middlewares.Auth(verifier))
	{
		// Commands
		group.POST("", controller.Create)
		group.PATCH("/:id", controller.Update)
		group.DELETE("/:id", controller.Delete)

		group.GET("/project/:projectId", controller.GetByProject)
	}
}

type TaskRouter struct {
	controller *controllers.TaskController
	verifier   middlewares.TokenVerifier
}

func NewTaskRoutes(controller *controllers.TaskController, verifier middlewares.TokenVerifier) *TaskRouter {
	return &TaskRouter{
		controller: controller,
		verifier:   verifier,
	}
}

func (tr *TaskRouter) Mount(r gin.IRouter) {
	TaskRoutes(r, tr.controller, tr.verifier)
}
