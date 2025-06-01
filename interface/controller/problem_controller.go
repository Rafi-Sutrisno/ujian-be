package controller

import (
	"mods/interface/dto"
	"mods/service"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	problemController struct {
		problemService service.ProblemService
	}

	ProblemController interface {
		GetByID(ctx *gin.Context)
		GetByExamID(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		Create(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}
)

func NewProblemController(ps service.ProblemService) ProblemController {
	return &problemController{
		problemService: ps,
	}
}

func (pc *problemController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := pc.problemService.GetByID(ctx.Request.Context(), id, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PROBLEM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PROBLEM, result)
	ctx.JSON(http.StatusOK, res)
}

func (pc *problemController) GetByExamID(ctx *gin.Context) {
	examID := ctx.Param("exam_id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := pc.problemService.GetByExamID(ctx.Request.Context(), examID, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_PROBLEM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_PROBLEM, result)
	ctx.JSON(http.StatusOK, res)
}

func (pc *problemController) GetAll(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	result, err := pc.problemService.GetAll(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_PROBLEM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_PROBLEM, result)
	ctx.JSON(http.StatusOK, res)
}

func (pc *problemController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	var req dto.ProblemCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	createdProblem, err := pc.problemService.Create(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PROBLEM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_PROBLEM, createdProblem)
	ctx.JSON(http.StatusCreated, res)
}

func (pc *problemController) Update(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	id := ctx.Param("id")
	var req dto.ProblemUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	updatedProblem, err := pc.problemService.Update(ctx.Request.Context(), req, id, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PROBLEM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_PROBLEM, updatedProblem)
	ctx.JSON(http.StatusOK, res)
}

func (pc *problemController) Delete(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	id := ctx.Param("id")
	if err := pc.problemService.Delete(ctx.Request.Context(), id, userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PROBLEM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_PROBLEM, nil)
	ctx.JSON(http.StatusOK, res)
}
