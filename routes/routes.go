package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userPublic := router.Group("/api/user")
	{
		// public can access
		
		userPublic.GET("", UserController.GetAllUser)
		userPublic.POST("/login", UserController.LoginUser)

	}
	userPrivate := router.Group("/api/user").Use(middleware.Authenticate(jwtService))
	{
		userPrivate.POST("/add", UserController.AddUser)
		// userPrivate.GET("/me", userController.Me)
		
	}

}
