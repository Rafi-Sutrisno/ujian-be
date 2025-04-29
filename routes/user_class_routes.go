package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func UserClassRoutes(router *gin.Engine, UserClassController controller.UserClassController, jwtService service.JWTService) {

	// userClassPrivate := router.Group("/api/user_class").Use(middleware.Authenticate(jwtService))
	// {
	// 	userClassPrivate.GET("/me", UserClassController.GetByToken)
	// 	userClassPrivate.GET("/user/:user_id", UserClassController.GetByUserID)
	// 	userClassPrivate.GET("/class/:class_id", UserClassController.GetByClassID)
	// 	userClassPrivate.POST("/create", UserClassController.Create)
	// 	userClassPrivate.POST("/create_many", UserClassController.CreateMany)
	// 	userClassPrivate.DELETE("/delete/:id", UserClassController.Delete)
	// }

	userClassPrivate := router.Group("/api/user_class").Use(middleware.Authenticate(jwtService))
	{
		// userClassPublic.GET("/me", UserClassController.GetByToken)
		// userClassPublic.GET("/user/:user_id", UserClassController.GetByUserID)
		userClassPrivate.GET("/class/:class_id", UserClassController.GetByClassID)
	}

	userClassPrivateAdmin := router.Group("/api/user_class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		// userClassPrivateAdmin.GET("/user/:user_id", UserClassController.GetByUserID)
		userClassPrivateAdmin.GET("/class/unassigned/:class_id", UserClassController.GetUnassigned)
		userClassPrivateAdmin.POST("/create", UserClassController.Create)
		userClassPrivateAdmin.POST("/create_many", UserClassController.CreateMany)
		userClassPrivateAdmin.DELETE("/delete/:id", UserClassController.Delete)
	}

	// examPrivateAdmin := router.Group("/api/class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	// {
	// 	examPrivateAdmin.POST("/add", ClassController.CreateExam)
	// 	examPrivateAdmin.GET("/all", ClassController.GetAllExam)
	// 	examPrivateAdmin.DELETE("/delete/:exam_id", ClassController.Delete)
	// }

}