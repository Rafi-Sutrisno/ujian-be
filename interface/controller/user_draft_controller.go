package controller

import (
	"fmt"
	"mods/application/service"
	"mods/interface/dto"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	userDraftController struct {
		service service.UserDraftService
	}

	UserDraftController interface {
		SaveDraft(ctx *gin.Context)
		GetDraft(ctx *gin.Context)
	}
)

func NewUserDraftController(s service.UserDraftService) UserDraftController {
	return &userDraftController{
		service: s,
	}
}

func (c *userDraftController) SaveDraft(ctx *gin.Context) {
	var req dto.UserCodeDraftRequest
	userId := ctx.MustGet("requester_id").(string)
	fmt.Println("masuk controller")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req.UserID = userId

	draft, err := c.service.SaveDraft(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to save draft", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Draft saved successfully", draft)
	ctx.JSON(http.StatusOK, res)
}

func (c *userDraftController) GetDraft(ctx *gin.Context) {
	var req dto.UserCodeDraftRequest
	userId := ctx.MustGet("requester_id").(string)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req.UserID = userId

	if req.ExamID == "" || req.ProblemID == "" || req.Language == "" {
		res := utils.BuildResponseFailed("Missing required query parameters", "", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	fmt.Println("ini req get draft:", req)

	draft, err := c.service.GetDraft(ctx.Request.Context(), userId, req.ExamID, req.ProblemID, req.Language)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to retrieve draft", err.Error(), nil)
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess("Draft retrieved successfully", draft)
	ctx.JSON(http.StatusOK, res)
}
