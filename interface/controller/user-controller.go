package controller

import (
	// "encoding/json"
	"fmt"
	"mods/application/service"
	"mods/interface/dto"
	"mods/utils"
	"net/http"
	"path/filepath"
	"strings"

	// "os/exec"
	// "strconv"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
	
}

type UserController interface {
	// regist login
	Register(ctx *gin.Context)
	GetAllUser(ctx *gin.Context)
	GetAllUserPaginate(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	Me(ctx *gin.Context)
	GetByUserId(ctx *gin.Context)
	UpdateEmailMe(ctx *gin.Context)
	UpdatePassMe(ctx *gin.Context)
	SendResetPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	RegisterFile(ctx *gin.Context)
}

func NewUserController(us service.UserService) UserController {
	return &userController{
		userService: us,

	}
}

func (c *userController) Me(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)

	result, err := c.userService.GetUserById(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	
	fmt.Println("ini result controller:", result)

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) GetByUserId(ctx *gin.Context) {
	userId := ctx.Param("id")

	result, err := c.userService.GetUserById(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	
	fmt.Println("ini result controller:", result)

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) Register(ctx *gin.Context) {
	var user dto.UserCreateRequest
	if err := ctx.ShouldBind(&user); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.Register(ctx.Request.Context(), user)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) GetAllUser(ctx *gin.Context) {
	results, err := uc.userService.GetAllUsers(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER, results)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) GetAllUserPaginate(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.GetAllUserWithPagination(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	resp := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_LIST_USER,
		Data:    result.Data,
		Meta:    result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (uc *userController) LoginUser(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result, err := uc.userService.Verify(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGIN, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) Update(ctx *gin.Context) {
	var req dto.UserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.Param("id")
	// userId := ctx.MustGet("user_id").(string)
	result, err := uc.userService.Update(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateEmailMe(ctx *gin.Context) {
	var req dto.UserUpdateEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}


	userId := ctx.MustGet("requester_id").(string)
	result, err := uc.userService.UpdateMe(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) UpdatePassMe(ctx *gin.Context) {
	var req dto.UserUpdatePasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}


	userId := ctx.MustGet("requester_id").(string)
	err := uc.userService.UpdatePassMe(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) Delete(ctx *gin.Context) {
	userId := ctx.Param("id")
	// userId := ctx.MustGet("user_id").(string)

	if err := uc.userService.Delete(ctx.Request.Context(), userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) RegisterFile(ctx *gin.Context) {
	var req dto.UserFileUploadRequest

	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	filename := req.File.Filename
	ext := strings.ToLower(filepath.Ext(filename))

	var result map[string]interface{}
	var err error

	switch ext {
	case ".yaml", ".yml":
		result, err = uc.userService.RegisterUsersFromYAML(ctx.Request.Context(), req.File)
	case ".csv":
		result, err = uc.userService.RegisterUsersFromCSV(ctx.Request.Context(), req.File)
	default:
		err = fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) SendResetPassword(ctx *gin.Context) {
	var req dto.SendResetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		fmt.Println("error binding:", err.Error()) // 
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.userService.SendForgotPasswordEmail(ctx.Request.Context(), req)
	if err != nil {
		fmt.Println("error service:", err.Error()) //
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SEND_FORGOT_PASSWORD_SUCCESS, nil)
	ctx.JSON(http.StatusOK, res)
}


func (c *userController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_EMAIL, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_RESET_PASSWORD, result)
	ctx.JSON(http.StatusOK, res)
}