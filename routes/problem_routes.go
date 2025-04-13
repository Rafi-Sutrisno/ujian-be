package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ProblemRoutes(router *gin.Engine, ProblemController controller.ProblemController, jwtService service.JWTService) {

	problemPrivate := router.Group("/api/problem").Use(middleware.Authenticate(jwtService))
	{
		problemPrivate.GET("/", ProblemController.GetAll)
		problemPrivate.GET("/:id", ProblemController.GetByID)
		problemPrivate.GET("/exam/:exam_id", ProblemController.GetByExamID)
		problemPrivate.POST("/", ProblemController.Create)
		problemPrivate.PATCH("/:id", ProblemController.Update)
		problemPrivate.DELETE("/:id", ProblemController.Delete)
	}
}
