package routes

import (
	"mods/controller"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ClassRoutes(router *gin.Engine, ClassController controller.ClassController, jwtService service.JWTService) {

	// classPrivate := router.Group("/api/class").Use(middleware.Authenticate(jwtService))
	// {
	// 	classPrivate.POST("/create", ClassController.Create)
	// 	classPrivate.GET("/:class_id", ClassController.GetById)
	// 	classPrivate.GET("/all", ClassController.GetAllWithPagination)
	// 	classPrivate.PATCH("/update/:class_id", ClassController.Update)
	// 	classPrivate.DELETE("/delete/:class_id", ClassController.Delete)
	// }

	classPublic := router.Group("/api/class")
	{
		classPublic.POST("/create", ClassController.Create)
		classPublic.GET("/:class_id", ClassController.GetById)
		classPublic.GET("/all/paginate", ClassController.GetAllWithPagination)
		classPublic.GET("/all", ClassController.GetAll)
		classPublic.PATCH("/update/:class_id", ClassController.Update)
		classPublic.DELETE("/delete/:class_id", ClassController.Delete)
	}
	// examPrivateAdmin := router.Group("/api/class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	// {
	// 	examPrivateAdmin.POST("/add", ClassController.CreateExam)
	// 	examPrivateAdmin.GET("/all", ClassController.GetAllExam)
	// 	examPrivateAdmin.DELETE("/delete/:exam_id", ClassController.Delete)
	// }

}