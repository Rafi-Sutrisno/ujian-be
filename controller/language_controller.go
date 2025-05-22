package controller

import (
	"net/http"
	"strconv"

	"mods/dto"
	"mods/service"
	"mods/utils"

	"github.com/gin-gonic/gin"
)

type LanguageController interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type languageController struct {
	langService service.LanguageService
}

func NewLanguageController(s service.LanguageService) LanguageController {
	return &languageController{langService: s}
}

func (lc *languageController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	result, err := lc.langService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		res := utils.BuildResponseFailed("Failed to get language", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess("Successfully retrieved language", result)
	ctx.JSON(http.StatusOK, res)
}

func (lc *languageController) GetAll(ctx *gin.Context) {
	result, err := lc.langService.GetAll(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed("Failed to get languages", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponseSuccess("Successfully retrieved languages", result)
	ctx.JSON(http.StatusOK, res)
}

func (lc *languageController) Create(ctx *gin.Context) {
	var req dto.LanguageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed("Invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := lc.langService.Create(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to create language", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Successfully created language", result)
	ctx.JSON(http.StatusCreated, res)
}

func (lc *languageController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req dto.LanguageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed("Invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := lc.langService.Update(ctx.Request.Context(), uint(id), req)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to update language", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Successfully updated language", result)
	ctx.JSON(http.StatusOK, res)
}

func (lc *languageController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := lc.langService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		res := utils.BuildResponseFailed("Failed to delete language", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Successfully deleted language", nil)
	ctx.JSON(http.StatusOK, res)
}
