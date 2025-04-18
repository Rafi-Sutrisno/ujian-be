package controller

import (
	"fmt"
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	classController struct {
		classService service.ClassService
	}

	ClassController interface {
		
		GetById(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetAllWithPagination(ctx *gin.Context)
		Create(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}
)

func NewClassController(cs service.ClassService) ClassController {
	return &classController{
		classService: cs,
	}
}

func (cc *classController)Create(ctx *gin.Context){
	var classReq dto.ClassCreateRequest

	if err := ctx.ShouldBind(&classReq); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	createdClass, err := cc.classService.Create(ctx.Request.Context(), classReq)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_CLASS, createdClass)
	ctx.JSON(http.StatusCreated, res)
}

func (cc *classController) GetById(ctx *gin.Context) {
	classId := ctx.Param("class_id")
	fmt.Println("class id di controller:", classId)

	result, err := cc.classService.GetById(ctx.Request.Context(), classId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (cc *classController) GetAllWithPagination(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := cc.classService.GetAllWithPagination(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	resp := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_LIST_CLASS,
		Data:    result.Data,
		Meta:    result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (cc *classController) GetAll(ctx *gin.Context) {
	results, err := cc.classService.GetAll(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_CLASS, results)
	ctx.JSON(http.StatusOK, res)
}

func (cc *classController) Update(ctx *gin.Context) {
	var req dto.ClassUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	classId := ctx.Param("class_id")
	// userId := ctx.MustGet("user_id").(string)
	result, err := cc.classService.Update(ctx.Request.Context(), req, classId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_CLASS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_CLASS, result)
	ctx.JSON(http.StatusOK, res)
}

func (cc *classController) Delete(ctx *gin.Context) {
	classId := ctx.Param("class_id")
	// userId := ctx.MustGet("user_id").(string)

	if err := cc.classService.Delete(ctx.Request.Context(), classId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_CLASS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_CLASS, nil)
	ctx.JSON(http.StatusOK, res)
}