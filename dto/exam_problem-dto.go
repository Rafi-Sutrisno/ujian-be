package dto

import (
	"errors"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_EXAM_PROBLEM           = "failed create exam problem"
	MESSAGE_FAILED_GET_LIST_EXAM_PROBLEM           = "failed get list exam problem"
	MESSAGE_FAILED_GET_EXAM_PROBLEM_TOKEN          = "failed get exam problem token"

	MESSAGE_FAILED_GET_EXAM_PROBLEM                = "failed get exam problem"

	MESSAGE_FAILED_UPDATE_EXAM_PROBLEM             = "failed update exam problem"
	MESSAGE_FAILED_DELETE_EXAM_PROBLEM             = "failed delete exam problem"

	// Success
	MESSAGE_SUCCESS_CREATE_EXAM_PROBLEM           = "success create exam problem"
	MESSAGE_SUCCESS_GET_LIST_EXAM_PROBLEM           = "success get list exam problem"
	MESSAGE_SUCCESS_GET_EXAM_PROBLEM                = "success get exam problem"
	MESSAGE_SUCCESS_UPDATE_EXAM_PROBLEM             = "success update exam problem"
	MESSAGE_SUCCESS_DELETE_EXAM_PROBLEM             = "success delete exam problem"

)

var (
	ErrCreateExamProblem             = errors.New("failed to create exam problem")
	ErrGetAllExamProblem             = errors.New("failed to get all exam problem")
	ErrGetAllExamProblemByExamId     = errors.New("failed to get all exam problem by exam id")
	ErrGetAllExamProblemByProblemId    = errors.New("failed to get all exam problem by problem id")
	ErrGetExamProblemById            = errors.New("failed to get exam problem by id")
	ErrGetExamProblemByEmail         = errors.New("failed to get exam problem by email")

	ErrUpdateExamProblem             = errors.New("failed to update exam problem")
	ErrExamProblemNotAdmin           = errors.New("exam problem not admin")
	ErrExamProblemNotFound           = errors.New("exam problem not found")

	ErrDeleteExamProblem             = errors.New("failed to delete exam problem")

)

type (
	ExamProblemCreateRequest struct {
		ExamID  string `json:"exam_id" binding:"required"`
		ProblemID string `json:"problem_id" binding:"required"`
	
	}

	ExamProblemResponse struct {
		ID      string `json:"id"`
		ExamID  string `json:"exam_id" `
		ProblemID string `json:"problem_id" `
		Problem    *ProblemResponse `json:"problem"`
		Exam    *ExamResponse `json:"exam"`
	}
)