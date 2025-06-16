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
	examSessionController struct {
		examSessionService service.ExamSessionService
	}

	ExamSessionController interface {
		CheckSession(ctx *gin.Context)
		CreateSession(ctx *gin.Context)
        GetByExamID(ctx *gin.Context)
		GetBySessionID(ctx *gin.Context)
		FinishSession(ctx *gin.Context)
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
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		fmt.Println("Tidak ada cookie session_id, lanjutkan tanpa session")
		sessionID = ""
	} else {
		fmt.Println("ini session id dari cookie:", sessionID)
	}
	for key, values := range ctx.Request.Header {
		fmt.Printf("Header %s: %s\n", key, values)
	}

	userAgent := ctx.Request.UserAgent()
	requestHash := ctx.GetHeader("X-SafeExamBrowser-RequestHash")
	configKeyHash := ctx.GetHeader("X-Safeexambrowser-Configkeyhash")
    
    if err := ctx.ShouldBind(&request); err != nil {
        res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
        ctx.JSON(http.StatusBadRequest, res)
        return
    }
	fmt.Println("ini fe url:", request.FEURL)
	fmt.Println("ini requesthash dari body:", request.BrowserExamKey)
	fmt.Println("ini requesthash dari header:", requestHash)
	fmt.Println("ini confighash dari body:", request.ConfigKey)
	fmt.Println("ini confighash dari header:", configKeyHash)

	// fmt.Println(validateSEBHash(request.FEURL, "fc7289e2e6e68444bc89d0b1ad2b70a20cabd2a93c1f4a0d7ae0fd4d64dd7cf5", request.BrowserExamKey))
	// fmt.Println(validateSEBHash(request.FEURL, "53adc488d10da166352a79712136414ae286a0dd28fe01147598f9b9fa561bd1", request.ConfigKey))

	

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

	fmt.Println("ini final key url:", requestHash, configKeyHash, fullURL)

    newSession, newSessionID, timeleft, err := cc.examSessionService.CreateorUpdateSession(
		ctx.Request.Context(),
		request,
		sessionID,
		userId,
		ctx.ClientIP(),
		userAgent,requestHash,configKeyHash,
		fullURL,
	)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_EXAM_SESSION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}


	// isDev := os.Getenv("GIN_MODE") != "production"

	fmt.Println("ini time left:", int(timeleft), newSessionID)

	// http.SetCookie(ctx.Writer, &http.Cookie{
	// 	Name:     "session_id",
	// 	Value:    newSessionID,
	// 	Path:     "/",
	// 	Domain:   "34.128.84.215",
	// 	MaxAge:   int(timeleft),
	// 	HttpOnly: true,
	// 	Secure:   !isDev, // secure false in dev
	// 	SameSite: func() http.SameSite {
	// 		if isDev {
	// 			return http.SameSiteLaxMode
	// 		}
	// 		return http.SameSiteNoneMode
	// 	}(),
	// })

    res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_EXAM_SESSION, newSession)
    ctx.JSON(http.StatusCreated, res)
}


func (cc *examSessionController) CheckSession(ctx *gin.Context) {
	sessionID, err := ctx.Cookie("session_id")
	fmt.Println("ini session id check:", sessionID)
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

func (cc *examSessionController) GetBySessionID(ctx *gin.Context) {
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

func (cc *examSessionController) FinishSession(ctx *gin.Context) {
	ExamId := ctx.Param("exam_id")
	userId := ctx.MustGet("requester_id").(string)
	fmt.Println("start update status controller")
	err := cc.examSessionService.FinishSession(ctx.Request.Context(), userId, ExamId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_FINISHING_EXAM_SESSION, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	fmt.Println("success update status controller")
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_FINISHING_EXAM_SESSION, nil)
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

// func  validateSEBHash(url string, key string, recvHash string) bool {

//     hasher := sha256.New()

// 	hasher.Write([]byte(key))
// 	hasher.Write([]byte(url))
    
//     finalHash := hasher.Sum(nil)
//     hashHex := hex.EncodeToString(finalHash)

//     fmt.Println("Controller-BEK/ConfigKey: Expected Hash:", hashHex)
//     fmt.Println("Controller-BEK/ConfigKey: Received Hash:", recvHash)

//     return hashHex == recvHash
// }