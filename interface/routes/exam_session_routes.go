package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func ExamSessionRoutes(router *gin.Engine, ExamSessionController controller.ExamSessionController, jwtService service.JWTService) {
	examSessionPrivate := router.Group("/api/exam_session").Use(middleware.Authenticate(jwtService))
	{
		examSessionPrivate.GET("/check_session", ExamSessionController.CheckSession)
		examSessionPrivate.POST("/start_exam", ExamSessionController.CreateSession)
		examSessionPrivate.POST("/finish_exam/:exam_id", ExamSessionController.FinishSession)
	}
	examSessionPrivateAdmin := router.Group("/api/exam_session").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		examSessionPrivateAdmin.GET("/byexamid/:exam_id", ExamSessionController.GetByExamID)
		// examSessionPrivateAdmin.DELETE("/:id", ExamSessionController.DeleteByID)
	}

}