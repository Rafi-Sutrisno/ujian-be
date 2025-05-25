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
		userPublic.POST("/login", UserController.LoginUser)
		userPublic.POST("/forgot_password", UserController.SendResetPassword)
		userPublic.POST("/reset_password", UserController.ResetPassword)
	}

	userPrivate := router.Group("/api/user").Use(middleware.Authenticate(jwtService))
	{
		userPrivate.POST("/login/dummy", UserController.LoginUser)
		userPrivate.GET("/me", UserController.Me)
		userPrivate.PATCH("/update/me", UserController.UpdateEmailMe)
	}

	userPrivateAdmin := router.Group("/api/user").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		userPrivateAdmin.GET("/:id", UserController.GetByUserId)
		userPrivateAdmin.GET("/all/paginate", UserController.GetAllUserPaginate)
		userPrivateAdmin.GET("/all", UserController.GetAllUser)
		userPrivateAdmin.POST("/add", UserController.Register)
		userPrivateAdmin.PATCH("/update/:id", UserController.Update)
		userPrivateAdmin.DELETE("/delete/:id", UserController.Delete)
		userPrivateAdmin.POST("/add/upload-file", UserController.RegisterFile)
	}

	// userPublic := router.Group("/api/user")
	// {
	// 	userPublic.GET("/all/paginate", UserController.GetAllUserPaginate)
	// 	userPublic.GET("/:id", UserController.Me)
	// 	userPublic.GET("/all", UserController.GetAllUser)
	// 	userPublic.POST("/login", UserController.LoginUser)
	// 	userPublic.POST("/add", UserController.Register)
	// 	userPublic.PATCH("/update/:id", UserController.Update)
	// 	userPublic.DELETE("/delete/:id", UserController.Delete)
	// 	userPublic.POST("/add/upload-yaml", UserController.RegisterYAML)
	// }

}
