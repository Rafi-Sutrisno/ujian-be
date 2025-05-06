package controller

import (
	"fmt"
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"
	"path/filepath"
	"strings"

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
		GetUnassigned(ctx *gin.Context)
		AssignFile(ctx *gin.Context)
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

func (uc *userClassController) AssignFile(ctx *gin.Context) {
	var req dto.UserFileUploadRequest

	fmt.Println("masuk controller assign file")
	class_id := ctx.Param("class_id")

	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	filename := req.File.Filename
	ext := strings.ToLower(filepath.Ext(filename))

	var result map[string]interface{}
	var err error

	switch ext {
	case ".yaml", ".yml":
		result, err = uc.userClassService.AssignUsersFromYAML(ctx.Request.Context(), class_id, req.File)
	case ".csv":
		result, err = uc.userClassService.AssignUsersFromCSV(ctx.Request.Context(),class_id, req.File)
	default:
		err = fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result)
	ctx.JSON(http.StatusOK, res)
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
	userId := ctx.MustGet("requester_id").(string)
	result, err := ucc.userClassService.GetByClassID(ctx.Request.Context(), classID, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ucc *userClassController) GetUnassigned(ctx *gin.Context) {
	classID := ctx.Param("class_id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := ucc.userClassService.GetUnassignedUsersByClassID(ctx.Request.Context(), classID, userId)
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
	userId := ctx.MustGet("requester_id").(string)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	createdUserClass, err := ucc.userClassService.Create(ctx.Request.Context(), req, userId)
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
	userId := ctx.MustGet("requester_id").(string)
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if err := ucc.userClassService.CreateMany(ctx.Request.Context(), reqs, userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_USER_CLASS, nil)
	ctx.JSON(http.StatusCreated, res)
}

func (ucc *userClassController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.MustGet("requester_id").(string)
	if err := ucc.userClassService.Delete(ctx.Request.Context(), id, userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER_CLASS, nil)
	ctx.JSON(http.StatusOK, res)
}
