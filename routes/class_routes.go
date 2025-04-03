package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ClassRoutes(router *gin.Engine, ClassController controller.ClassController, jwtService service.JWTService) {

	classPrivate := router.Group("/api/class").Use(middleware.Authenticate(jwtService))
	{
		classPrivate.GET("/create", ClassController.Create)
		classPrivate.GET("/:class_id", ClassController.GetById)
		classPrivate.GET("/all", ClassController.GetAllWithPagination)
		classPrivate.PATCH("/update/:class_id", ClassController.Update)
		classPrivate.PATCH("/delete/:class_id", ClassController.Delete)
	}
	// examPrivateAdmin := router.Group("/api/class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	// {
	// 	examPrivateAdmin.POST("/add", ClassController.CreateExam)
	// 	examPrivateAdmin.GET("/all", ClassController.GetAllExam)
	// 	examPrivateAdmin.DELETE("/delete/:exam_id", ClassController.Delete)
	// }

}