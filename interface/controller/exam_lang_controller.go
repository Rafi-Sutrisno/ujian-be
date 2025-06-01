package controller

import (
	"mods/interface/dto"
	"mods/service"
	"mods/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	examLangController struct {
		examLangService service.ExamLangService
	}

	ExamLangController interface {
		GetByExamID(ctx *gin.Context)
		GetByLangID(ctx *gin.Context)
		Create(ctx *gin.Context)
		CreateMany(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}
)

func NewExamLangController(els service.ExamLangService) ExamLangController {
	return &examLangController{
		examLangService: els,
	}
}

func (elc *examLangController) GetByExamID(ctx *gin.Context) {
	examID := ctx.Param("exam_id")
	result, err := elc.examLangService.GetAllByExamID(ctx.Request.Context(), examID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_EXAM_LANG, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_EXAM_LANG, result)
	ctx.JSON(http.StatusOK, res)
}

func (elc *examLangController) GetByLangID(ctx *gin.Context) {
	langIDStr := ctx.Param("lang_id")
	langIDUint64, err := strconv.ParseUint(langIDStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_EXAM_LANG, "Invalid lang_id", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	langID := uint(langIDUint64)
	result, err := elc.examLangService.GetAllByLangID(ctx.Request.Context(), langID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_EXAM_LANG, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_EXAM_LANG, result)
	ctx.JSON(http.StatusOK, res)
}

func (elc *examLangController) Create(ctx *gin.Context) {
	var req dto.ExamLangCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	createdExamLang, err := elc.examLangService.Create(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_EXAM_LANG, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_EXAM_LANG, createdExamLang)
	ctx.JSON(http.StatusCreated, res)
}

func (elc *examLangController) CreateMany(ctx *gin.Context) {
	var reqs []dto.ExamLangCreateRequest
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if err := elc.examLangService.CreateMany(ctx.Request.Context(), reqs); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_EXAM_LANG, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_EXAM_LANG, nil)
	ctx.JSON(http.StatusCreated, res)
}

func (elc *examLangController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := elc.examLangService.Delete(ctx.Request.Context(), id); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_EXAM_LANG, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_EXAM_LANG, nil)
	ctx.JSON(http.StatusOK, res)
}
