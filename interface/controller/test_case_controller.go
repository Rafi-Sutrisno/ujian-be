package controller

import (
	"fmt"
	"mods/application/service"
	"mods/interface/dto"
	"mods/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type (
	testCaseController struct {
		testCaseService service.TestCaseService
	}

	TestCaseController interface {
		GetByID(ctx *gin.Context)
		GetByProblemID(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		Create(ctx *gin.Context)
		UploadZip(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}
)

func NewTestCaseController(ts service.TestCaseService) TestCaseController {
	return &testCaseController{
		testCaseService: ts,
	}
}

func (tc *testCaseController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := tc.testCaseService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TEST_CASE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TEST_CASE, result)
	ctx.JSON(http.StatusOK, res)
}

func (tc *testCaseController) GetByProblemID(ctx *gin.Context) {
	problemID := ctx.Param("problem_id")
	userId := ctx.MustGet("requester_id").(string)
	results, err := tc.testCaseService.GetByProblemID(ctx.Request.Context(), problemID, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TEST_CASE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TEST_CASE, results)
	ctx.JSON(http.StatusOK, res)
}

func (tc *testCaseController) GetAll(ctx *gin.Context) {
	results, err := tc.testCaseService.GetAll(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_TEST_CASE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_TEST_CASE, results)
	ctx.JSON(http.StatusOK, res)
}

func (tc *testCaseController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	var req dto.TestCaseCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := tc.testCaseService.Create(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TEST_CASE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_TEST_CASE, result)
	ctx.JSON(http.StatusCreated, res)
}

func (tc *testCaseController) UploadZip(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	problemID := ctx.Param("problem_id")

	file, err := ctx.FormFile("file")
	if err != nil {
		res := utils.BuildResponseFailed("Gagal menerima file ZIP", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Save uploaded zip to a temp file
	tempFilePath := "./tmp/" + file.Filename
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		res := utils.BuildResponseFailed("Gagal menyimpan file sementara", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	defer os.Remove(tempFilePath) // clean up after use

	// Parse and create test cases
	count, err := tc.testCaseService.CreateFromZip(ctx.Request.Context(), tempFilePath, problemID, userId)
	if err != nil {
		res := utils.BuildResponseFailed("Gagal membuat test case dari zip", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(fmt.Sprintf("Berhasil membuat %d test case", count), nil)
	ctx.JSON(http.StatusCreated, res)
}

func (tc *testCaseController) Update(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	id := ctx.Param("id")
	var req dto.TestCaseUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := tc.testCaseService.Update(ctx.Request.Context(), req, id, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_TEST_CASE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_TEST_CASE, result)
	ctx.JSON(http.StatusOK, res)
}

func (tc *testCaseController) Delete(ctx *gin.Context) {
	userId := ctx.MustGet("requester_id").(string)
	id := ctx.Param("id")
	if err := tc.testCaseService.Delete(ctx.Request.Context(), id, userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_TEST_CASE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_TEST_CASE, nil)
	ctx.JSON(http.StatusOK, res)
}
