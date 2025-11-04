package routes

import (
	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"github.com/gin-gonic/gin"
)

func PromptRoutes(r gin.IRouter, controller *controllers.PromptController, verifier middlewares.TokenVerifier) {
	group := r.Group("/prompts")
	group.Use(middlewares.Auth(verifier))
	{
		// Commands
		group.POST("", controller.Create)
		group.PATCH("/:id", controller.Update)
		group.DELETE("/:id", controller.Delete)
		group.POST("/:id/rollback", controller.Rollback)
		group.POST("/reorder", controller.Reorder)
		group.POST("/:id/run", controller.Run)

		// Queries
		group.GET("/task/:taskId", controller.GetByTask)
		group.GET("/:id", controller.GetByID)
	}
}

type PromptRouter struct {
	controller *controllers.PromptController
	verifier   middlewares.TokenVerifier
}

func NewPromptRoutes(controller *controllers.PromptController, verifier middlewares.TokenVerifier) *PromptRouter {
	return &PromptRouter{
		controller: controller,
		verifier:   verifier,
	}
}

func (pr *PromptRouter) Mount(r gin.IRouter) {
	PromptRoutes(r, pr.controller, pr.verifier)
}
