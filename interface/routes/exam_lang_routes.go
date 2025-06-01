package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func ExamLangRoutes(router *gin.Engine, ClassController controller.ExamLangController, jwtService service.JWTService) {

	examLangPrivate := router.Group("/api/exam_lang").Use(middleware.Authenticate(jwtService))
	{
		examLangPrivate.GET("/exam/:exam_id", ClassController.GetByExamID)
		examLangPrivate.GET("/lang/:lang_id", ClassController.GetByLangID)
		examLangPrivate.POST("/create", ClassController.Create)
		examLangPrivate.POST("/create_many", ClassController.CreateMany)
		examLangPrivate.DELETE("/delete/:id", ClassController.Delete)
	}

}