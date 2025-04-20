package routes

import (
	"mods/controller"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func ProblemRoutes(router *gin.Engine, ProblemController controller.ProblemController, jwtService service.JWTService) {

	// problemPrivate := router.Group("/api/problem").Use(middleware.Authenticate(jwtService))
	// {
	// 	problemPrivate.GET("/", ProblemController.GetAll)
	// 	problemPrivate.GET("/:id", ProblemController.GetByID)
	// 	problemPrivate.GET("/exam/:exam_id", ProblemController.GetByExamID)
	// 	problemPrivate.POST("/", ProblemController.Create)
	// 	problemPrivate.PATCH("/:id", ProblemController.Update)
	// 	problemPrivate.DELETE("/:id", ProblemController.Delete)
	// }

	problemPublic := router.Group("/api/problem")
	{
		problemPublic.GET("/", ProblemController.GetAll)
		problemPublic.GET("/:id", ProblemController.GetByID)
		problemPublic.GET("/exam/:exam_id", ProblemController.GetByExamID)
		problemPublic.POST("/", ProblemController.Create)
		problemPublic.PATCH("/:id", ProblemController.Update)
		problemPublic.DELETE("/:id", ProblemController.Delete)
	}
}
