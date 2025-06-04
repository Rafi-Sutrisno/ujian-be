package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func UserClassRoutes(router *gin.Engine, UserClassController controller.UserClassController, jwtService service.JWTService) {

	userClassPrivate := router.Group("/api/user_class").Use(middleware.Authenticate(jwtService))
	{
		userClassPrivate.GET("/class/:class_id", UserClassController.GetByClassID)
	}

	userClassPrivateAdmin := router.Group("/api/user_class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		userClassPrivateAdmin.GET("/class/unassigned/:class_id", UserClassController.GetUnassigned)
		userClassPrivateAdmin.POST("/upload-file/:class_id", UserClassController.AssignFile)
		userClassPrivateAdmin.POST("/", UserClassController.Create)
		userClassPrivateAdmin.POST("/create_many", UserClassController.CreateMany)
		userClassPrivateAdmin.DELETE("/:id", UserClassController.Delete)
	}

}