package dto

import (
	"errors"

	"mods/domain/entity"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_PROBLEM           = "failed create problem"
	MESSAGE_FAILED_GET_LIST_PROBLEM           = "failed get list problem"

	MESSAGE_FAILED_GET_PROBLEM                = "failed get problem"

	MESSAGE_FAILED_UPDATE_PROBLEM             = "failed update problem"
	MESSAGE_FAILED_DELETE_PROBLEM             = "failed delete problem"

	// Success
	MESSAGE_SUCCESS_CREATE_PROBLEM             = "success create problem"
	MESSAGE_SUCCESS_GET_LIST_PROBLEM           = "success get list problem"
	MESSAGE_SUCCESS_GET_PROBLEM                = "success get problem"
	MESSAGE_SUCCESS_UPDATE_PROBLEM             = "success update problem"
	MESSAGE_SUCCESS_DELETE_PROBLEM             = "success delete problem"

)

var (
	ErrCreateProblem             = errors.New("failed to create problem")
	ErrGetAllProblem             = errors.New("failed to get all problem")
	ErrGetAllProblemByExamId     = errors.New("failed to get all problem by exam id")
	ErrGetProblemById            = errors.New("failed to get problem by id")
	ErrGetProblemByEmail         = errors.New("failed to get problem by email")

	ErrUpdateProblem             = errors.New("failed to update problem")
	ErrProblemNotAdmin           = errors.New("problem not admin")
	ErrProblemNotFound           = errors.New("problem not found")

	ErrDeleteProblem             = errors.New("failed to delete problem")

)

type (
	ProblemCreateRequest struct {
		Title         string    `json:"title" binding:"required"`
		Description   string    `json:"description" binding:"required"`
		Constraints   string    `json:"constraints" binding:"required"`
		SampleInput   string    `json:"sample_input" binding:"required"`
		SampleOutput  string    `json:"sample_output" binding:"required"`
		CpuTimeLimit  float64 `json:"cpu_time_limit,omitempty"`
    	MemoryLimit   int     `json:"memory_limit,omitempty"`
	}

	ProblemResponse struct {
		ID         	  string    `json:"id"`
		Title         string    `json:"title"`
		Description   string    `json:"description" `
		Constraints   string    `json:"constraints" `
		SampleInput   string    `json:"sample_input" `
		SampleOutput  string    `json:"sample_output" `
		CreatedAt     string	`json:"created_at" `
		CpuTimeLimit  float64 `json:"cpu_time_limit"`
    	MemoryLimit   int     `json:"memory_limit"`
	}

	ProblemPaginationResponse struct {
		Data []ProblemResponse `json:"data"`
		PaginationResponse
	}

	GetAllProblemRepositoryResponse struct {
		Problem []entity.Problem `json:"problem"`
		PaginationResponse
	}

	ProblemUpdateRequest struct {
		Title         string    `json:"title" form:"title"`
		Description   string    `json:"description" form:"description"`
		Constraints   string    `json:"constraints" form:"constraints"`
		SampleInput   string    `json:"sample_input" form:"sample_input"`
		SampleOutput  string    `json:"sample_output" form:"sample_output"`
		CpuTimeLimit  float64 `json:"cpu_time_limit"`
		MemoryLimit   int     `json:"memory_limit"`
	}

	ProblemUpdateResponse struct {
		ID            string    `json:"id"`
		Title         string    `json:"title"`
		Description   string    `json:"description" `
		Constraints   string    `json:"constraints" `
		SampleInput   string    `json:"sample_input" `
		SampleOutput  string    `json:"sample_output" `
	}
)