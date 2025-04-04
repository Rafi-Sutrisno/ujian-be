package dto

import (
	"errors"

	"mods/entity"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_TEST_CASE           = "failed create test_case"
	MESSAGE_FAILED_GET_LIST_TEST_CASE           = "failed get list test_case"

	MESSAGE_FAILED_GET_TEST_CASE                = "failed get test_case"

	MESSAGE_FAILED_UPDATE_TEST_CASE             = "failed update test_case"
	MESSAGE_FAILED_DELETE_TEST_CASE             = "failed delete test_case"

	// Success
	MESSAGE_SUCCESS_CREATE_TEST_CASE             = "success create test_case"
	MESSAGE_SUCCESS_GET_LIST_TEST_CASE           = "success get list test_case"
	MESSAGE_SUCCESS_GET_TEST_CASE                = "success get test_case"
	MESSAGE_SUCCESS_UPDATE_TEST_CASE             = "success update test_case"
	MESSAGE_SUCCESS_DELETE_TEST_CASE             = "success delete test_case"

)

var (
	ErrCreateTestCase             = errors.New("failed to create test_case")
	ErrGetAllTestCase             = errors.New("failed to get all test_case")
	ErrGetAllTestCaseByProblemId  = errors.New("failed to get all test_case by problem id")
	ErrGetTestCaseById            = errors.New("failed to get test_case by id")
	ErrGetTestCaseByEmail         = errors.New("failed to get test_case by email")

	ErrUpdateTestCase             = errors.New("failed to update test_case")
	ErrTestCaseNotAdmin           = errors.New("test_case not admin")
	ErrTestCaseNotFound           = errors.New("test_case not found")

	ErrDeleteTestCase             = errors.New("failed to delete test_case")

)

type (
	TestCaseCreateRequest struct {
		ProblemID       string    `json:"problem_id" binding:"required"`
		InputData       string    `json:"input_data" binding:"required"`
		ExpectedOutput  string    `json:"expected_output" binding:"required"`
	}

	TestCaseResponse struct {
		ID         		string    `json:"id"`
		ProblemID       string    `json:"problem_id"`
		InputData       string    `json:"input_data"`
		ExpectedOutput  string    `json:"expected_output"`
	}

	TestCasePaginationResponse struct {
		Data []TestCaseResponse `json:"data"`
		PaginationResponse
	}

	GetAllTestCaseRepositoryResponse struct {
		TestCase []entity.TestCase `json:"test_case"`
		PaginationResponse
	}

	TestCaseUpdateRequest struct {
		InputData       string    `json:"input_data" form:"input_data"`
		ExpectedOutput  string    `json:"expected_output" form:"expected_output"` 
	}

	TestCaseUpdateResponse struct {
		ID         		string 	  `json:"id"`
		ProblemID       string    `json:"problem_id"`
		InputData       string    `json:"input_data"`
		ExpectedOutput  string    `json:"expected_output"`
	}
)