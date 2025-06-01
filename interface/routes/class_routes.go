package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

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

	classPrivate := router.Group("/api/class").Use(middleware.Authenticate(jwtService))
	{
		classPrivate.GET("/:class_id", ClassController.GetById)
		classPrivate.GET("/user", ClassController.GetByUserID) // ganti ke by token
	}

	classPrivateAdmin := router.Group("/api/class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		classPrivateAdmin.POST("/create", ClassController.Create)
		classPrivateAdmin.GET("/all/paginate", ClassController.GetAllWithPagination)
		classPrivateAdmin.GET("/all", ClassController.GetAll)
		classPrivateAdmin.PATCH("/update/:class_id", ClassController.Update)
		classPrivateAdmin.DELETE("/delete/:class_id", ClassController.Delete)
	}

	// examPrivateAdmin := router.Group("/api/class").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	// {
	// 	examPrivateAdmin.POST("/add", ClassController.CreateExam)
	// 	examPrivateAdmin.GET("/all", ClassController.GetAllExam)
	// 	examPrivateAdmin.DELETE("/delete/:exam_id", ClassController.Delete)
	// }

}