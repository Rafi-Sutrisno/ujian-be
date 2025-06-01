package routes

import (
	"mods/interface/controller"
	"mods/interface/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func LanguageRoutes(router *gin.Engine, LanguageController controller.LanguageController, jwtService service.JWTService) {
	languagePrivateAdmin := router.Group("/api/language").Use(middleware.Authenticate(jwtService)).Use(middleware.Authorize("admin"))
	{
		// languagePrivateAdmin.POST("/create", LanguageController.Create)
		// languagePrivateAdmin.GET("/all/paginate", LanguageController.GetAllWithPagination)
		languagePrivateAdmin.GET("/all", LanguageController.GetAll)
		// languagePrivateAdmin.PATCH("/update/:class_id", LanguageController.Update)
		// languagePrivateAdmin.DELETE("/delete/:class_id", LanguageController.Delete)
	}
}