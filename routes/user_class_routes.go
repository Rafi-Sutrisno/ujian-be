package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func UserClassRoutes(router *gin.Engine, UserClassController controller.UserClassController, jwtService service.JWTService) {

	userClassPrivate := router.Group("/api/user_class").Use(middleware.Authenticate(jwtService))
	{
		userClassPrivate.GET("/me", UserClassController.GetByToken)
		userClassPrivate.GET("/:user_id", UserClassController.GetByUserID)
		userClassPrivate.GET("/:class_id", UserClassController.GetByClassID)
		userClassPrivate.GET("/create", UserClassController.Create)
		userClassPrivate.PATCH("/delete/:id", UserClassController.Delete)
	}
	// examPrivateAdmin := router.Group("/api/class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	// {
	// 	examPrivateAdmin.POST("/add", ClassController.CreateExam)
	// 	examPrivateAdmin.GET("/all", ClassController.GetAllExam)
	// 	examPrivateAdmin.DELETE("/delete/:exam_id", ClassController.Delete)
	// }

}