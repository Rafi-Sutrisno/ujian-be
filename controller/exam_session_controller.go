package controller

import (
	"fmt"
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	examSessionController struct {
		examSessionService service.ExamSessionService
	}

	ExamSessionController interface {
		CheckSession(ctx *gin.Context)
		CreateSession(ctx *gin.Context)
        GetByExamID(ctx *gin.Context)
		DeleteByID(ctx *gin.Context)
	}
)

func NewExamSessionController(es service.ExamSessionService) ExamSessionController {
	return &examSessionController{
		examSessionService: es,
	}
}



func (cc *examSessionController) CreateSession(ctx *gin.Context) {
    var request dto.ExamSessionCreateRequest
    userId := ctx.MustGet("requester_id").(string)
    
    if err := ctx.ShouldBind(&request); err != nil {
        res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
        ctx.JSON(http.StatusBadRequest, res)
        return
    }

	ipAddress := ctx.ClientIP()
    userAgent := ctx.Request.UserAgent()
	fmt.Println("User-Agent:", userAgent)

	if strings.Contains(userAgent, "SEB") {
		fmt.Println("Request is from Safe Exam Browser")
	} else {
		fmt.Println("Request is NOT from Safe Exam Browser")
	}

	requestHash := ctx.Request.Header.Get("X-SafeExamBrowser-RequestHash")
	configKeyHash := ctx.Request.Header.Get("X-Safeexambrowser-Configkeyhash")

	fmt.Println("X-SafeExamBrowser-RequestHash:", requestHash)
	fmt.Println("X-Safeexambrowser-Configkeyhash:", configKeyHash)

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, ctx.Request.Host, ctx.Request.RequestURI)
	fmt.Println("ini full url: ", fullURL)

    newSession, sessionID, err := cc.examSessionService.CreateorUpdateSession(ctx.Request.Context(), request, userId, ipAddress, userAgent, requestHash, configKeyHash, fullURL)
    if err != nil {
        res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_EXAM_SESSION, err.Error(), nil)
        ctx.JSON(http.StatusBadRequest, res)
        return
    }

	isDev := os.Getenv("GIN_MODE") != "release"

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   !isDev, // secure false in dev
		SameSite: func() http.SameSite {
			if isDev {
				return http.SameSiteLaxMode
			}
			return http.SameSiteNoneMode
		}(),
	})

    res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_EXAM_SESSION, newSession)
    ctx.JSON(http.StatusCreated, res)
}


func (cc *examSessionController) CheckSession(ctx *gin.Context) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil || sessionID == "" {
		res := utils.BuildResponseFailed("Session not found", "No session_id cookie provided", nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	session, err := cc.examSessionService.GetBySessionID(ctx.Request.Context(), sessionID)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid session", err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	res := utils.BuildResponseSuccess("Session is valid", session)
	ctx.JSON(http.StatusOK, res)
}


func (cc *examSessionController) GetByExamID(ctx *gin.Context) {
	examId := ctx.Param("exam_id")

	sessions, err := cc.examSessionService.GetByExamID(ctx.Request.Context(), examId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_EXAM_SESSION, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_EXAM_SESSION, sessions)
	ctx.JSON(http.StatusOK, res)
}

func (cc *examSessionController) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")

	err := cc.examSessionService.DeleteByID(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_EXAM_SESSION, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_EXAM_SESSION, nil)
	ctx.JSON(http.StatusOK, res)
}

