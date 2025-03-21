package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func UserExamRoutes(router *gin.Engine, UserExamController controller.UserExamController, jwtService service.JWTService) {

	// userPrivate := router.Group("/api/user").Use(middleware.Authenticate(jwtService))
	// {
	// 	userPrivate.GET("/me", ExamController.Me)
	// 	userPrivate.PATCH("/update", ExamController.UpdateMe)
	// }
	userExamPrivateAdmin := router.Group("/api/exam/user").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		userExamPrivateAdmin.POST("/add", UserExamController.CreateUserExam)
	}

}