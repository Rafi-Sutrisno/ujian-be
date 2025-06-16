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
	problem_id := ctx.Param("problem_id")
	userId := ctx.MustGet("requester_id").(string)
	result, err := pc.problemService.GetByID(ctx.Request.Context(), problem_id, userId, problem_id)
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
    var request dto.ExamSessionCreateRequest

	if err := ctx.ShouldBind(&request); err != nil {
        res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
        ctx.JSON(http.StatusBadRequest, res)
        return
    }
	
	userAgent := ctx.Request.UserAgent()
	requestHash := ctx.GetHeader("X-SafeExamBrowser-RequestHash")
	configKeyHash := ctx.GetHeader("X-Safeexambrowser-Configkeyhash")

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, ctx.Request.Host, ctx.Request.RequestURI)
	fmt.Println("ini full url: ", fullURL)

	if(requestHash == ""){
		requestHash=request.BrowserExamKey
		configKeyHash=request.ConfigKey
		fullURL=request.FEURL
	}

	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		fmt.Println("Tidak ada cookie session_id, lanjutkan tanpa session")
		sessionID = ""
	} else {
		fmt.Println("ini session id dari cookie:", sessionID)
	}

	result, err := pc.problemService.GetByExamID(ctx.Request.Context(), userAgent,requestHash, configKeyHash, fullURL, sessionID, userId, examID)
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
