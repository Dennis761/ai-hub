package routes

import (
	"ai_hub.com/app/infra/http/gin/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r gin.IRouter, controller *controllers.AdminController) {
	admin := r.Group("/admin")
	{
		admin.POST("/login", controller.Login)
		admin.POST("/register", controller.Register)
		admin.POST("/:id/verify", controller.Verify)
		admin.POST("/reset/start", controller.StartReset)
		admin.POST("/reset/confirm", controller.ConfirmReset)
		admin.POST("/reset/change", controller.ChangePassword)
		admin.PATCH("/:id/rename", controller.Rename)
		admin.DELETE("/:id", controller.Delete)
	}
}

type AdminRouter struct {
	controller *controllers.AdminController
}

func NewAdminRoutes(controller *controllers.AdminController) *AdminRouter {
	return &AdminRouter{controller: controller}
}

func (ar *AdminRouter) Mount(r gin.IRouter) {
	AdminRoutes(r, ar.controller)
}
