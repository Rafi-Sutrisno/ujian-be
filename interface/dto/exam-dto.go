package dto

import (
	"errors"
	"mods/domain/entity"
	"time"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_EXAM           = "failed create exam"
	MESSAGE_FAILED_GET_LIST_EXAM           = "failed get list exam"
	// MESSAGE_FAILED_TOKEN_NOT_VALID         = "token not valid"
	// MESSAGE_FAILED_TOKEN_NOT_FOUND         = "token not found"
	MESSAGE_FAILED_GET_EXAM               = "failed get exam"

	MESSAGE_FAILED_UPDATE_EXAM             = "failed update exam"
	MESSAGE_FAILED_DELETE_EXAM             = "failed delete exam"
	// MESSAGE_FAILED_PROSES_REQUEST          = "failed proses request"
	// MESSAGE_FAILED_DENIED_ACCESS           = "denied access"

	// Success
	MESSAGE_SUCCESS_CREATE_EXAM           = "success create exam"
	MESSAGE_SUCCESS_GET_LIST_EXAM           = "success get list exam"
	MESSAGE_SUCCESS_GET_EXAM               = "success get exam"

	MESSAGE_SUCCESS_UPDATE_EXAM             = "success update exam"
	MESSAGE_SUCCESS_DELETE_EXAM             = "success delete exam"

)

var (
	ErrCreateExam             = errors.New("failed to create exam")
	ErrGetAllExam             = errors.New("failed to get all exam")
	ErrGetExamById            = errors.New("failed to get exam by id")
	ErrShortNameAlreadyExists     = errors.New("email already exist")
	ErrUpdateExam             = errors.New("failed to update exam")
    ErrGetAllExamsByClassId    = errors.New("failed to get all exams by class id")
	ErrExamNotFound           = errors.New("exam not found")

	ErrDeleteExam             = errors.New("failed to delete exam")
	// ErrDeniedAccess           = errors.New("denied access")
	// ErrAccountNotVerified     = errors.New("account not verified")
	// ErrTokenInvalid           = errors.New("token invalid")
	// ErrTokenExpired           = errors.New("token expired")

)

type (
	// untuk create
	ExamCreateRequest struct {
		ClassID     		string    		`json:"class_id" binding:"required"`
		Name        		string    		`json:"name" binding:"required"`
		ShortName   		string    		`json:"short_name" binding:"required"`
		IsPublished 		bool      		`json:"is_published"`
		StartTime   		time.Time 		`json:"start_time" binding:"required"`
		DurationStr 		string    	  	`json:"duration" binding:"required"` 
		Duration    		time.Duration 	`json:"-"`
		IsSEBRestricted   	bool            `json:"is_seb_restricted"`
		SEBBrowserKey      	string        	`json:"seb_browser_key"` 
		SEBConfigKey      	string        	`json:"seb_config_key"` 
		SEBQuitURL      	string        	`json:"seb_quit_url"`   
	}

	ExamResponse struct {
		ID          		string        `json:"id"`
		ClassID     		string        `json:"class_id"`
		Name        		string        `json:"name"`
		ShortName   		string        `json:"short_name"`
		IsPublished 		bool          `json:"is_published"`
		StartTime   		time.Time     `json:"start_time"`           
		Duration    		string 		  `json:"duration"`             
		EndTime     		time.Time     `json:"end_time"`
		IsSEBRestricted   	bool            `json:"is_seb_restricted"`
		SEBBrowserKey      	string        	`json:"seb_browser_key"` 
		SEBConfigKey      	string        	`json:"seb_config_key"`   
		SEBQuitURL      	string        	`json:"seb_quit_url"`   
		AllowedLanguages    []LanguageResponse          `json:"allowed_languages"`      
	}


	// untuk update
	ExamUpdateRequest struct {
		Name        string        `json:"name" `
		ShortName   string        `json:"short_name"`
		IsPublished bool          `json:"is_published"`
		StartTime   time.Time     `json:"start_time"` 
		DurationStr string        `json:"duration"`           
		Duration    time.Duration `json:"-"`
		IsSEBRestricted   	bool            `json:"is_seb_restricted"`
		SEBBrowserKey      	string        	`json:"seb_browser_key"` 
		SEBConfigKey      	string        	`json:"seb_config_key"`     
		SEBQuitURL      	string        	`json:"seb_quit_url"`                      
	}

	ExamUpdateResponse struct {
		ID          string        `json:"id"`
		ClassID     string        `json:"class_id"`
		Name        string        `json:"name"`
		ShortName   string        `json:"short_name"`
		IsPublished bool          `json:"is_published"`
		StartTime   time.Time     `json:"start_time"`           
		Duration    time.Duration `json:"duration"`             
		EndTime     time.Time     `json:"end_time"`           
		
	}

	// untuk get
	ExamPaginationResponse struct {
		Data []ExamResponse `json:"data"`
		PaginationResponse
	}

	GetAllExamRepositoryResponse struct {
		Exams []entity.Exam `json:"exams"`
		PaginationResponse
	}
)