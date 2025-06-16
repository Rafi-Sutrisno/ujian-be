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
	submissionController struct {
		submissionService service.SubmissionService
	}

	SubmissionController interface {
		RunCode(ctx *gin.Context)
		SubmitCode(ctx *gin.Context)
		GetCorrectStatsByExam(ctx *gin.Context)
		GetCorrectStatsByExamandStudent(ctx *gin.Context)
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
	var combinedReq dto.CombinedRequestRun
		userId := ctx.MustGet("requester_id").(string)
		examId := ctx.Param("exam_id")
	if err := ctx.ShouldBindJSON(&combinedReq); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req := combinedReq.Judge0Request
	request := combinedReq.ExamSessionCreateRequest

	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		fmt.Println("Tidak ada cookie session_id, lanjutkan tanpa session")
		sessionID = ""
	} else {
		fmt.Println("ini session id dari cookie:", sessionID)
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

	// Call service
	result, err := sc.submissionService.RunCode(ctx.Request.Context(), req, userAgent,requestHash,configKeyHash, fullURL,  sessionID,userId, examId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_RUN_CODE, result)

	ctx.JSON(http.StatusOK, res)
}

func  (sc *submissionController) SubmitCode(ctx *gin.Context) {
	var combinedReq dto.CombinedRequestSubmit
	userId := ctx.MustGet("requester_id").(string)
	examId := ctx.Param("exam_id")

	if err := ctx.ShouldBindJSON(&combinedReq); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req := combinedReq.SubmissionRequest
	request := combinedReq.ExamSessionRequest

	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		fmt.Println("Tidak ada cookie session_id, lanjutkan tanpa session")
		sessionID = ""
	} else {
		fmt.Println("ini session id dari cookie:", sessionID)
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
	
	// Call service
	result, err := sc.submissionService.SubmitCode(ctx.Request.Context(),req, userAgent,requestHash,configKeyHash, fullURL,  sessionID,  userId, examId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_SUBMISSION, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *submissionController) GetCorrectStatsByExam(ctx *gin.Context) {
	examID := ctx.Param("exam_id")

	stats, err := sc.submissionService.GetCorrectSubmissionStatsByExam(ctx.Request.Context(), examID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SUBMISSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_STATISTICS, stats)
	ctx.JSON(http.StatusOK, res)
}

func (sc *submissionController) GetCorrectStatsByExamandStudent(ctx *gin.Context) {
	examID := ctx.Param("exam_id")
	userId := ctx.MustGet("requester_id").(string)

	stats, err := sc.submissionService.GetCorrectSubmissionStatsByExamandUser(ctx.Request.Context(), examID, userId)
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
	var request dto.ExamSessionCreateRequest
	examID := ctx.Param("exam_id")
	userId := ctx.MustGet("requester_id").(string)

	if err := ctx.ShouldBind(&request); err != nil {
        res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
        ctx.JSON(http.StatusBadRequest, res)
        return
    }

	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		fmt.Println("Tidak ada cookie session_id, lanjutkan tanpa session")
		sessionID = ""
	} else {
		fmt.Println("ini session id dari cookie:", sessionID)
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

	
	result, err := sc.submissionService.GetByExamIDandUserID(ctx.Request.Context(), userAgent,requestHash,configKeyHash, fullURL, sessionID,userId, examID )
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
