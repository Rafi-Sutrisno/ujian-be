package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func ProblemRoutes(router *gin.Engine, ProblemController controller.ProblemController, jwtService service.JWTService) {

	problemPrivate := router.Group("/api/problem").Use(middleware.Authenticate(jwtService))
	{
		problemPrivate.GET("/:id", ProblemController.GetByID)
		problemPrivate.GET("/exam/:exam_id", ProblemController.GetByExamID)
	}

	problemPrivateAdmin := router.Group("/api/problem").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		problemPrivateAdmin.GET("/", ProblemController.GetAll)
		problemPrivateAdmin.POST("/", ProblemController.Create)
		problemPrivateAdmin.PATCH("/:id", ProblemController.Update)
		problemPrivateAdmin.DELETE("/:id", ProblemController.Delete)
	}
}
