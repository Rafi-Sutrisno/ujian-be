package controller

import (
	"fmt"
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	examController struct {
		examService service.ExamService
	}

	ExamController interface {
		CreateExam(ctx *gin.Context)
		GetExamById(ctx *gin.Context)
		GetByClassID(ctx *gin.Context)
		GetAllExam(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}
)

func NewExamController(es service.ExamService) ExamController {
	return &examController{
		examService: es,
	}
}

func (ec *examController)CreateExam(ctx *gin.Context){
	var examReq dto.ExamCreateRequest

	if err := ctx.ShouldBind(&examReq); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	parsedDuration, err := time.ParseDuration(examReq.DurationStr)
	if err != nil {
		res := utils.BuildResponseFailed("invalid duration format", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	examReq.Duration = parsedDuration

	now := time.Now()
	fmt.Println("Server Time:", now.Format(time.RFC3339), "Timezone:", now.Location())
	if examReq.StartTime.Before(now) || examReq.StartTime.Equal(now) {
		res := utils.BuildResponseFailed("invalid start time", "start time must be after current time", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	createdExam, err := ec.examService.CreateExam(ctx.Request.Context(), examReq)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_EXAM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_EXAM, createdExam)
	ctx.JSON(http.StatusCreated, res)
}

func (c *examController) GetByClassID(ctx *gin.Context) {
	classID := ctx.Param("class_id")
	result, err := c.examService.GetByClassID(ctx.Request.Context(), classID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_EXAM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_EXAM, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *examController) GetExamById(ctx *gin.Context) {
	examId := ctx.Param("exam_id")
	fmt.Println("exam id di controller:", examId)

	result, err := c.examService.GetExamById(ctx.Request.Context(), examId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_EXAM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_EXAM, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *examController) GetAllExam(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.examService.GetAllExamWithPagination(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_EXAM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	resp := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_LIST_EXAM,
		Data:    result.Data,
		Meta:    result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (uc *examController) Update(ctx *gin.Context) {
	var req dto.ExamUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if req.DurationStr != "" {
		parsedDuration, err := time.ParseDuration(req.DurationStr)
		if err != nil {
			res := utils.BuildResponseFailed("invalid duration format", err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		req.Duration = parsedDuration
	}
	now := time.Now()
	fmt.Println("Server Time:", now.Format(time.RFC3339), "Timezone:", now.Location())
	if req.StartTime.Before(now) || req.StartTime.Equal(now) {
		res := utils.BuildResponseFailed("invalid start time", "start time must be after current time", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	examId := ctx.Param("exam_id")
	// userId := ctx.MustGet("user_id").(string)
	result, err := uc.examService.Update(ctx.Request.Context(), req, examId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_EXAM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_EXAM, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *examController) Delete(ctx *gin.Context) {
	examId := ctx.Param("exam_id")
	// userId := ctx.MustGet("user_id").(string)

	if err := uc.examService.Delete(ctx.Request.Context(), examId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_EXAM, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_EXAM, nil)
	ctx.JSON(http.StatusOK, res)
}