package dto

import (
	"errors"
	"time"

	"mods/domain/entity"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_SUBMISSION           = "failed create submission"
	MESSAGE_FAILED_GET_LIST_SUBMISSION           = "failed get list submission"

	MESSAGE_FAILED_GET_SUBMISSION                = "failed get submission"

	MESSAGE_FAILED_UPDATE_SUBMISSION             = "failed update submission"
	MESSAGE_FAILED_DELETE_SUBMISSION             = "failed delete submission"

	// Success
	MESSAGE_SUCCESS_CREATE_SUBMISSION             = "success create submission"
	MESSAGE_SUCCESS_RUN_CODE             		  = "success run code"
	MESSAGE_SUCCESS_GET_LIST_SUBMISSION           = "success get list submission"
	MESSAGE_SUCCESS_GET_SUBMISSION                = "success get submission"
	MESSAGE_SUCCESS_GET_STATISTICS                = "success get statistics"
	MESSAGE_SUCCESS_UPDATE_SUBMISSION             = "success update submission"
	MESSAGE_SUCCESS_DELETE_SUBMISSION             = "success delete submission"

)

var (
	ErrCreateSubmission             = errors.New("failed to create submission")
	ErrGetAllSubmission             = errors.New("failed to get all submission")
	ErrGetAllSubmissionByUserId     = errors.New("failed to get all submission by user id")
	ErrGetAllSubmissionByExamId     = errors.New("failed to get all submission by exam id")
	ErrGetAllSubmissionByProblemId  = errors.New("failed to get all submission by problem id")
	ErrGetSubmissionById            = errors.New("failed to get submission by id")
	ErrGetSubmissionByEmail         = errors.New("failed to get submission by email")

	ErrUpdateSubmission             = errors.New("failed to update submission")
	ErrSubmissionNotAdmin           = errors.New("submission not admin")
	ErrSubmissionNotFound           = errors.New("submission not found")

	ErrDeleteSubmission             = errors.New("failed to delete submission")

)

type (
	SubmissionCreateRequest struct {
		UserID      		string      `json:"user_id" binding:"required"`
		ExamID      		string      `json:"exam_id" binding:"required"`
		ProblemID   		string      `json:"problem_id" binding:"required"`
		LangID 				string      `json:"lang_id" binding:"required"`
		Code     			string      `json:"code" binding:"required"`
	}

	SubmissionResponse struct {
		ID         			string 		`json:"id"`
		UserID      		string      `json:"user_id"`
		ExamID      		string      `json:"exam_id"`
		ProblemID   		string      `json:"problem_id"`
		LangID 				uint      `json:"lang_id"`
		Code     			string      `json:"code"`
		SubmissionTime    	string      `json:"submission_time"`
		Status     			uint      `json:"status"`
		Time 				string		`json:"time"`
		Memory				string		`json:"memory"`
		Problem				ProblemResponse   `json:"problem"`
		Language			LanguageResponse  `json:"lang"`
		User				UserResponse  `json:"user"`
		StatusName          string		  `json:"status_name"`
	}

	SubmissionPaginationResponse struct {
		Data []SubmissionResponse `json:"data"`
		PaginationResponse
	}

	GetAllSubmissionRepositoryResponse struct {
		Submission []entity.Submission `json:"submission"`
		PaginationResponse
	}

	ExamUserCorrectDTO struct {
		UserID       string    `json:"user_id"`
		UserName     string    `json:"user_name"`
		UserNoID     string    `json:"user_no_id"`
		TotalCorrect int       `json:"total_correct"`
		TotalProblem int       `json:"total_problem"`
		Status       uint      `json:"status"`
		FinishedAt   time.Time `json:"finished_at"`

		AcceptedProblems       string `json:"accepted_problems"`       // comma-separated
		WrongProblems          string `json:"wrong_problems"`          // comma-separated
		NoSubmissionProblems   string `json:"no_submission_problems"`  // comma-separated
	}

)