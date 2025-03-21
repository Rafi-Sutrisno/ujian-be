package controller

import (
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	userExamController struct {
		userExamService service.UserExamService
	}

	UserExamController interface {
		CreateUserExam(ctx *gin.Context)
	}
)

func NewUserExamController(es service.UserExamService) UserExamController {
	return &userExamController{
		userExamService: es,
	}
}

func (ec *userExamController)CreateUserExam(ctx *gin.Context){
	var userExamReq dto.UserExamCreateRequest

	if err := ctx.ShouldBind(&userExamReq); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	createdExam, err := ec.userExamService.CreateUserExam(ctx.Request.Context(), userExamReq)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_EXAM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_EXAM, createdExam)
	ctx.JSON(http.StatusCreated, res)
}