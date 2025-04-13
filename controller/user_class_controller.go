package controller

import (
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	userClassController struct {
		userClassService service.UserClassService
	}

	UserClassController interface {
		GetByToken(ctx *gin.Context)
		GetByUserID(ctx *gin.Context)
		GetByClassID(ctx *gin.Context)
		Create(ctx *gin.Context)
		CreateMany(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}
)

func NewUserClassController(ucs service.UserClassService) UserClassController {
	return &userClassController{
		userClassService: ucs,
	}
}

func (ucc *userClassController) GetByToken(ctx *gin.Context) {
	// userID := ctx.Param("user_id")
	userID := ctx.MustGet("requester_id").(string)
	result, err := ucc.userClassService.GetByUserID(ctx.Request.Context(), userID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ucc *userClassController) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	// userID := ctx.MustGet("user_id").(string)
	result, err := ucc.userClassService.GetByUserID(ctx.Request.Context(), userID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ucc *userClassController) GetByClassID(ctx *gin.Context) {
	classID := ctx.Param("class_id")
	result, err := ucc.userClassService.GetByClassID(ctx.Request.Context(), classID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ucc *userClassController) Create(ctx *gin.Context) {
	var req dto.UserClassCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	createdUserClass, err := ucc.userClassService.Create(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_USER_CLASS, createdUserClass)
	ctx.JSON(http.StatusCreated, res)
}

func (ucc *userClassController) CreateMany(ctx *gin.Context) {
	var reqs []dto.UserClassCreateRequest
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if err := ucc.userClassService.CreateMany(ctx.Request.Context(), reqs); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_USER_CLASS, nil)
	ctx.JSON(http.StatusCreated, res)
}

func (ucc *userClassController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ucc.userClassService.Delete(ctx.Request.Context(), id); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER_CLASS, nil)
	ctx.JSON(http.StatusOK, res)
}
