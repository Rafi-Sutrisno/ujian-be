package routes

import (
	"mods/controller"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ExamRoutes(router *gin.Engine, ExamController controller.ExamController, jwtService service.JWTService) {

	// examPublic := router.Group("/api/exam").Use(middleware.Authenticate(jwtService))
	// {
	// 	examPublic.GET("/:exam_id", ExamController.GetExamById)
	// 	examPublic.PATCH("/update/:exam_id", ExamController.Update)
	// 	examPublic.POST("/add", ExamController.CreateExam)
	// 	examPublic.GET("/all", ExamController.GetAllExam)
	// 	examPublic.DELETE("/delete/:exam_id", ExamController.Delete)
	// }
	examPublic := router.Group("/api/exam")
	{
		examPublic.GET("/:exam_id", ExamController.GetExamById)
		examPublic.GET("/byclass/:class_id", ExamController.GetByClassID)
		examPublic.PATCH("/:exam_id", ExamController.Update)
		examPublic.POST("/create", ExamController.CreateExam)
		examPublic.GET("/all", ExamController.GetAllExam)
		examPublic.DELETE("/:exam_id", ExamController.Delete)
	}
	// examPrivateAdmin := router.Group("/api/exam").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	// {
	// 	examPrivateAdmin.POST("/add", ExamController.CreateExam)
	// 	examPrivateAdmin.GET("/all", ExamController.GetAllExam)
	// 	examPrivateAdmin.DELETE("/delete/:exam_id", ExamController.Delete)
	// }

}