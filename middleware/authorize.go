package middleware

import (
	"mods/dto"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleVal, exists := ctx.Get("requester_role")
		if !exists {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, "Role not found", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role, ok := roleVal.(string)
		if !ok || role != requiredRole {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, "Access denied: insufficient role", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Next()
	}
}
