package routes

import (
	"mods/application/service"
	"mods/interface/controller"
	"mods/interface/middleware"

	"github.com/gin-gonic/gin"
)

func UserDraftRoutes(router *gin.Engine, UserDraftController controller.UserDraftController, jwtService service.JWTService) {
	userPrivate := router.Group("/api/user/draft").Use(middleware.Authenticate(jwtService))
	{
		userPrivate.POST("/load", UserDraftController.GetDraft)
		userPrivate.POST("/save", UserDraftController.SaveDraft)
	}
}
