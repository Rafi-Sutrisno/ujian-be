package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func SubmissionRoutes(router *gin.Engine, SubmissionController controller.SubmissionController, jwtService service.JWTService) {

	submissionPrivate := router.Group("/api/submission").Use(middleware.Authenticate(jwtService))
	{
		submissionPrivate.POST("/submit", SubmissionController.SubmitCode)
		submissionPrivate.POST("/run", SubmissionController.RunCode)
		submissionPrivate.POST("/", SubmissionController.Create)
		submissionPrivate.GET("/:id", SubmissionController.GetByID)
		submissionPrivate.GET("/exam/:exam_id", SubmissionController.GetByExamID)
		submissionPrivate.GET("/problem/:problem_id", SubmissionController.GetByProblemID)
		submissionPrivate.GET("/user/:user_id", SubmissionController.GetByUserID)
	}
}
