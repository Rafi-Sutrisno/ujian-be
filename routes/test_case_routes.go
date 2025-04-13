package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func TestCaseRoutes(router *gin.Engine, TestCaseController controller.TestCaseController, jwtService service.JWTService) {

	testCasePrivate := router.Group("/api/testcase").Use(middleware.Authenticate(jwtService))
	{
		testCasePrivate.GET("/", TestCaseController.GetAll)
		testCasePrivate.GET("/:id", TestCaseController.GetByID)
		testCasePrivate.GET("/problem/:problem_id", TestCaseController.GetByProblemID)
		testCasePrivate.POST("/", TestCaseController.Create)
		testCasePrivate.PATCH("/:id", TestCaseController.Update)
		testCasePrivate.DELETE("/:id", TestCaseController.Delete)
	}
}
