package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userPublic := router.Group("/api/user")
	{
		userPublic.GET("", UserController.GetAllUser)
		userPublic.POST("/login", UserController.LoginUser)
	}
	userPrivate := router.Group("/api/user").Use(middleware.Authenticate(jwtService))
	{
		userPrivate.GET("/me", UserController.Me)
		userPrivate.PATCH("/update", UserController.UpdateMe)
	}
	userPrivateAdmin := router.Group("/api/user").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		userPrivateAdmin.GET("/all", UserController.GetAllUser)
		userPrivateAdmin.POST("/add", UserController.Register)
		userPrivateAdmin.PATCH("/update/:id", UserController.Update)
		userPrivateAdmin.DELETE("/delete/:id", UserController.Delete)
		userPrivateAdmin.POST("/add/upload-yaml", UserController.RegisterYAML)
	}

}
