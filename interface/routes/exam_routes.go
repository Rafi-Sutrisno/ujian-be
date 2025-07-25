package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func ExamRoutes(router *gin.Engine, ExamController controller.ExamController, jwtService service.JWTService) {

	examPrivate := router.Group("/api/exam").Use(middleware.Authenticate(jwtService))
	{
		examPrivate.GET("/:exam_id", ExamController.GetExamById)
		examPrivate.GET("/byclass/:class_id", ExamController.GetByClassID)
		examPrivate.GET("/byuser", ExamController.GetByUserID)
	}
	examPrivateAdmin := router.Group("/api/exam").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		examPrivateAdmin.PATCH("/:exam_id", ExamController.Update)
		examPrivateAdmin.POST("/", ExamController.CreateExam)
		examPrivateAdmin.POST("/yaml/:class_id", ExamController.UploadExamFromYaml)
		examPrivateAdmin.GET("/all", ExamController.GetAllExam)
		examPrivateAdmin.DELETE("/:exam_id", ExamController.Delete)
	}

}