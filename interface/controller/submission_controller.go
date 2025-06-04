package controller

import (
	"mods/application/service"
	"mods/interface/dto"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	submissionController struct {
		submissionService service.SubmissionService
	}

	SubmissionController interface {
		RunCode(ctx *gin.Context)
		SubmitCode(ctx *gin.Context)
		Create(ctx *gin.Context)
		GetCorrectStatsByExam(ctx *gin.Context)
		GetByID(ctx *gin.Context)
		GetByExamIDandUserID(ctx *gin.Context)
		GetByExamID(ctx *gin.Context)
		GetByProblemID(ctx *gin.Context)
		GetByUserID(ctx *gin.Context)
	}
)

func NewSubmissionController(ss service.SubmissionService) SubmissionController {
	return &submissionController{
		submissionService: ss,
	}
}

func  (sc *submissionController) RunCode(ctx *gin.Context) {
	var req dto.Judge0Request
    userId := ctx.MustGet("requester_id").(string)
	examId := ctx.Param("exam_id")

	// Bind JSON body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Call service
	result, err := sc.submissionService.RunCode(ctx.Request.Context(), ctx, req, userId, examId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func  (sc *submissionController) SubmitCode(ctx *gin.Context) {
	var req dto.SubmissionRequest
	userId := ctx.MustGet("requester_id").(string)
	examId := ctx.Param("exam_id")
	// Bind JSON body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Call service
	result, err := sc.submissionService.SubmitCode(ctx.Request.Context(), ctx, req, userId, examId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (sc *submissionController) Create(ctx *gin.Context) {
	var req dto.SubmissionCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := sc.submissionService.CreateSubmission(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_SUBMISSION, result)
	ctx.JSON(http.StatusCreated, res)
}

func (sc *submissionController) GetCorrectStatsByExam(ctx *gin.Context) {
	examID := ctx.Param("exam_id")

	stats, err := sc.submissionService.GetCorrectSubmissionStatsByExam(ctx.Request.Context(), examID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_SUBMISSION, stats)
	ctx.JSON(http.StatusOK, res)
}


func (sc *submissionController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := sc.submissionService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_SUBMISSION, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *submissionController) GetByExamIDandUserID(ctx *gin.Context) {
	examID := ctx.Param("exam_id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := sc.submissionService.GetByExamIDandUserID(ctx.Request.Context(), examID, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_SUBMISSION, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *submissionController) GetByExamID(ctx *gin.Context) {
	examID := ctx.Param("exam_id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := sc.submissionService.GetByExamID(ctx.Request.Context(), examID, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_SUBMISSION, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *submissionController) GetByProblemID(ctx *gin.Context) {
	problemID := ctx.Param("problem_id")

	result, err := sc.submissionService.GetByProblemID(ctx.Request.Context(), problemID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_SUBMISSION, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *submissionController) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	result, err := sc.submissionService.GetByUserID(ctx.Request.Context(), userID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_SUBMISSION, result)
	ctx.JSON(http.StatusOK, res)
}
