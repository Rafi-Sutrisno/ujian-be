package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ExamProblemRoutes(router *gin.Engine, ExamProblemController controller.ExamProblemController, jwtService service.JWTService) {


	examProblemPrivate := router.Group("/api/exam_problem").Use(middleware.Authenticate(jwtService))
	{
		examProblemPrivate.GET("/exam/:exam_id", ExamProblemController.GetByExamID)
		examProblemPrivate.GET("/problem/:problem_id", ExamProblemController.GetByProblemID)
	}

	examProblemPrivateAdmin := router.Group("/api/exam_problem").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		examProblemPrivateAdmin.POST("/create", ExamProblemController.Create)
		examProblemPrivateAdmin.POST("/create_many", ExamProblemController.CreateMany)
		examProblemPrivateAdmin.DELETE("/delete/:id", ExamProblemController.Delete)
	}

}