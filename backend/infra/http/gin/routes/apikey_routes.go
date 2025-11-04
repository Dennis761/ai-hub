package routes

import (
	"ai_hub.com/app/infra/http/gin/controllers"
	"ai_hub.com/app/infra/http/gin/middlewares"
	"github.com/gin-gonic/gin"
)

func APIKeyRoutes(r gin.IRouter, controller *controllers.APIKeyController, verifier middlewares.TokenVerifier) {
	api := r.Group("/api-keys")
	api.Use(middlewares.Auth(verifier))
	{
		// Commands
		api.POST("", controller.Create)
		api.PATCH("/:id", controller.Update)
		api.DELETE("/:id", controller.Delete)

		// Queries
		api.GET("/my-keys", controller.GetMyKeys)
		api.GET("/:id", controller.GetByID)
	}
}

type APIKeyRouter struct {
	controller *controllers.APIKeyController
	verifier   middlewares.TokenVerifier
}

func NewAPIKeyRoutes(controller *controllers.APIKeyController, verifier middlewares.TokenVerifier) *APIKeyRouter {
	return &APIKeyRouter{
		controller: controller,
		verifier:   verifier,
	}
}

func (ar *APIKeyRouter) Mount(r gin.IRouter) {
	APIKeyRoutes(r, ar.controller, ar.verifier)
}
