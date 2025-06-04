package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func SubmissionRoutes(router *gin.Engine, SubmissionController controller.SubmissionController, jwtService service.JWTService) {

	submissionPrivate := router.Group("/api/submission").Use(middleware.Authenticate(jwtService))
	{
		submissionPrivate.POST("/submit/:exam_id", SubmissionController.SubmitCode)
		submissionPrivate.POST("/run/:exam_id", SubmissionController.RunCode)
		// submissionPrivate.POST("/", SubmissionController.Create)
		submissionPrivate.GET("/:id", SubmissionController.GetByID)
		submissionPrivate.GET("/exam/student/:exam_id", SubmissionController.GetByExamIDandUserID)
		// submissionPrivate.GET("/problem/:problem_id", SubmissionController.GetByProblemID)
		// submissionPrivate.GET("/user/:user_id", SubmissionController.GetByUserID)
	}

	submissionPrivateAdmin := router.Group("/api/submission").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		submissionPrivateAdmin.GET("/exam/:exam_id", SubmissionController.GetByExamID)
		submissionPrivateAdmin.GET("/stats/exam/:exam_id", SubmissionController.GetCorrectStatsByExam)
	}
}
