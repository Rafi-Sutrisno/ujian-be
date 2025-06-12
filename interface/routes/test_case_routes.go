package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func TestCaseRoutes(router *gin.Engine, TestCaseController controller.TestCaseController, jwtService service.JWTService) {


	testCasePrivate := router.Group("/api/testcase").Use(middleware.Authenticate(jwtService))
	{
		testCasePrivate.GET("/:id", TestCaseController.GetByID)
		testCasePrivate.GET("/problem/:problem_id", TestCaseController.GetByProblemID)
	}
	
	testCasePrivateAdmin := router.Group("/api/testcase").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		testCasePrivateAdmin.GET("/", TestCaseController.GetAll)
		testCasePrivateAdmin.POST("/", TestCaseController.Create)
		testCasePrivateAdmin.PATCH("/:id", TestCaseController.Update)
		testCasePrivateAdmin.DELETE("/:id", TestCaseController.Delete)
	}
}
