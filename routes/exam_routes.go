package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ExamRoutes(router *gin.Engine, ExamController controller.ExamController, jwtService service.JWTService) {

	examPrivate := router.Group("/api/exam").Use(middleware.Authenticate(jwtService))
	{
		examPrivate.GET("/byuser/:id", ExamController.CreateExam)
		examPrivate.GET("/:exam_id", ExamController.GetExamById)
		examPrivate.PATCH("/update/:exam_id", ExamController.Update)
	}
	examPrivateAdmin := router.Group("/api/exam").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		examPrivateAdmin.POST("/add", ExamController.CreateExam)
		examPrivateAdmin.GET("/all", ExamController.GetAllExam)
		examPrivateAdmin.DELETE("/delete/:exam_id", ExamController.Delete)
	}

}