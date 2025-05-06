package dto

import (
	"errors"
	"mime/multipart"

	"mods/entity"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_USER_CLASS           = "failed create user class"
	MESSAGE_FAILED_GET_LIST_USER_CLASS           = "failed get list user class"
	MESSAGE_FAILED_GET_USER_CLASS_TOKEN          = "failed get user class token"

	MESSAGE_FAILED_GET_USER_CLASS                = "failed get user class"

	MESSAGE_FAILED_UPDATE_USER_CLASS             = "failed update user class"
	MESSAGE_FAILED_DELETE_USER_CLASS             = "failed delete user class"

	// Success
	MESSAGE_SUCCESS_CREATE_USER_CLASS           = "success create user class"
	MESSAGE_SUCCESS_GET_LIST_USER_CLASS           = "success get list user class"
	MESSAGE_SUCCESS_GET_USER_CLASS                = "success get user class"
	MESSAGE_SUCCESS_UPDATE_USER_CLASS             = "success update user class"
	MESSAGE_SUCCESS_DELETE_USER_CLASS             = "success delete user class"

)

var (
	ErrCreateUserClass             = errors.New("failed to create user class")
	ErrGetAllUserClass             = errors.New("failed to get all user class")
	ErrGetAllUserClassByUserId     = errors.New("failed to get all user class by user id")
	ErrAuthorize                   = errors.New("you are not authorized for this class")
	ErrGetAllUserClassByClassId    = errors.New("failed to get all user class by class id")
	ErrGetUserClassById            = errors.New("failed to get user class by id")
	ErrGetUserClassByEmail         = errors.New("failed to get user class by email")

	ErrUpdateUserClass             = errors.New("failed to update user class")
	ErrUserClassNotAdmin           = errors.New("user class not admin")
	ErrUserClassNotFound           = errors.New("user class not found")

	ErrDeleteUserClass             = errors.New("failed to delete user class")

)

type (
	UserClassCreateRequest struct {
		UserID  string `json:"user_id" binding:"required"`
		ClassID string `json:"class_id" binding:"required"`
	
	}

	UserClassResponse struct {
		ID      string `json:"id"`
		UserID  string `json:"user_id" `
		ClassID string `json:"class_id" `
		User    *UserResponse `json:"user"`
	}
	

	UserClassYAMLUploadRequest struct {
		File *multipart.FileHeader `form:"yaml_file" binding:"required"`
	}

	 UserClassYAML struct {
		Name       string `yaml:"name"`
		Email      string `yaml:"email"`
		Noid 	   string `yaml:"noid"`
		Password   string `yaml:"password"`  
	}
	
	 UserClassYAMLList struct {
		UserClass []UserClassYAML `yaml:"users"`
	}

	FailedUserClassResponse struct {
		Noid   string `json:"noid"`
		Email  string `json:"email"`
		Reason string `json:"reason"`
	}

	UserClassPaginationResponse struct {
		Data []UserClassResponse `json:"data"`
		PaginationResponse
	}

	GetAllUserClassRepositoryResponse struct {
		UserClasses []entity.UserClass `json:"user_classes"`
		PaginationResponse
	}

	UserClassUpdateRequest struct {
		Name       string `json:"name" form:"name"`
		ShortName      string `json:"short_name" form:"short_name"`
	}

	UserClassUpdateResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		ShortName      string `json:"short_name"`
	}
)