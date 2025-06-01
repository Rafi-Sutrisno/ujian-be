package controller

import (
	"fmt"
	"mods/interface/dto"
	"mods/service"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	examProblemController struct {
		service service.ExamProblemService
	}

	ExamProblemController interface {
		GetByExamID(ctx *gin.Context)
		GetByProblemID(ctx *gin.Context)
		GetUnassignedByExamID(ctx *gin.Context)
		Create(ctx *gin.Context)
		CreateMany(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}
)

func NewExamProblemController(ucs service.ExamProblemService) ExamProblemController {
	return &examProblemController{
		service: ucs,
	}
}

func (ucc *examProblemController) GetByExamID(ctx *gin.Context) {
	examID := ctx.Param("exam_id")
	userID := ctx.MustGet("requester_id").(string)
	result, err := ucc.service.GetByExamID(ctx.Request.Context(),examID, userID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ucc *examProblemController) GetByProblemID(ctx *gin.Context) {
	problemID := ctx.Param("problem_id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := ucc.service.GetByProblemID(ctx.Request.Context(), problemID, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ucc *examProblemController) GetUnassignedByExamID(ctx *gin.Context) {
	examID := ctx.Param("exam_id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := ucc.service.GetUnassignedByExamID(ctx.Request.Context(), examID, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ucc *examProblemController) Create(ctx *gin.Context) {
	var req dto.ExamProblemCreateRequest
	userId := ctx.MustGet("requester_id").(string)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	createdUserClass, err := ucc.service.Create(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_USER_CLASS, createdUserClass)
	ctx.JSON(http.StatusCreated, res)
}

func (ucc *examProblemController) CreateMany(ctx *gin.Context) {
	var reqs []dto.ExamProblemCreateRequest
	userId := ctx.MustGet("requester_id").(string)
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if err := ucc.service.CreateMany(ctx.Request.Context(), reqs, userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_USER_CLASS, nil)
	ctx.JSON(http.StatusCreated, res)
}

func (ucc *examProblemController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println("ini id:", id)
	userId := ctx.MustGet("requester_id").(string)
	if err := ucc.service.Delete(ctx.Request.Context(), id, userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER_CLASS, nil)
	ctx.JSON(http.StatusOK, res)
}
