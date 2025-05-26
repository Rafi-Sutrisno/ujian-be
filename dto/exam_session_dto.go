package dto

import (
	"errors"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_EXAM_SESSION           = "failed create exam session"
	MESSAGE_FAILED_GET_LIST_EXAM_SESSION           = "failed get list  exam session"
	// MESSAGE_FAILED_TOKEN_NOT_VALID         = "token not valid"
	// MESSAGE_FAILED_TOKEN_NOT_FOUND         = "token not found"
	MESSAGE_FAILED_GET_EXAM_SESSION               = "failed get  exam session"
	MESSAGE_FAILED_DELETE_EXAM_SESSION             = "failed delete  exam session"
	MESSAGE_FAILED_FINISHING_EXAM_SESSION             = "failed finishing  exam session"
	// MESSAGE_FAILED_PROSES_REQUEST          = "failed proses request"
	// MESSAGE_FAILED_DENIED_ACCESS           = "denied access"

	// Success
	MESSAGE_SUCCESS_FINISHING_EXAM_SESSION             = "success finishing  exam session"
	MESSAGE_SUCCESS_CREATE_EXAM_SESSION           = "success create  exam session"
	MESSAGE_SUCCESS_GET_LIST_EXAM_SESSION           = "success get list  exam session"
	MESSAGE_SUCCESS_GET_EXAM_SESSION               = "success get  exam session"
	MESSAGE_SUCCESS_DELETE_EXAM_SESSION             = "success delete  exam session"

)

var (
	ErrCreateExamSession             = errors.New("failed to create  exam session")
	ErrGetAllExamSession             = errors.New("failed to get all  exam session")
	ErrGetExamSessionById            = errors.New("failed to get  exam session by id")
	ErrUpdateExamSession             = errors.New("failed to update  exam session")
    ErrGetAllExamSessionsByClassId    = errors.New("failed to get all  exam session by class id")
	ErrExamSessionNotFound           = errors.New(" exam session not found")

	ErrDeleteExamSession             = errors.New("failed to delete  exam session")


)

type (

	ExamSessionCreateRequest struct {
		ExamID     string     `json:"exam_id" binding:"required"`
	}

	ExamSessionCreateResponse struct {
		UserID          	string        	`json:"user_id"`
		ExamID     			string     		`json:"exam_id"`
	}

	ExamSessionGetResponse struct {
		ID					string 			`json:"id"`
		UserID          	string        	`json:"user_id"`
		ExamID     			string     		`json:"exam_id"`
		IpAddress			string 			`json:"ip_address"`
		UserAgent			string 			`json:"user_agent"`
		Device				string 			`json:"device"`
		Status              uint			`json:"status"`
		User    			*UserResponse 	`json:"user"`
	}
	
)