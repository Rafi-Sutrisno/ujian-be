package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ExamRoutes(router *gin.Engine, ExamController controller.ExamController, jwtService service.JWTService) {

	// examPrivate := router.Group("/api/exam").Use(middleware.Authenticate(jwtService))
	// {
	// 	examPrivate.GET("/:exam_id", ExamController.GetExamById)
	// 	examPrivate.PATCH("/update/:exam_id", ExamController.Update)
	// 	examPrivate.POST("/add", ExamController.CreateExam)
	// 	examPrivate.GET("/all", ExamController.GetAllExam)
	// 	examPrivate.DELETE("/delete/:exam_id", ExamController.Delete)
	// }
	examPrivate := router.Group("/api/exam").Use(middleware.Authenticate(jwtService))
	{
		examPrivate.GET("/:exam_id", ExamController.GetExamById)
		examPrivate.GET("/byclass/:class_id", ExamController.GetByClassID)
	}
	examPrivateAdmin := router.Group("/api/exam").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		examPrivateAdmin.PATCH("/:exam_id", ExamController.Update)
		examPrivateAdmin.POST("/create", ExamController.CreateExam)
		examPrivateAdmin.GET("/all", ExamController.GetAllExam)
		examPrivateAdmin.DELETE("/:exam_id", ExamController.Delete)
	}

}