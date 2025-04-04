package dto

import (
	"errors"
	"mime/multipart"

	"mods/entity"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_CLASS           = "failed create class"
	MESSAGE_FAILED_GET_LIST_CLASS           = "failed get list class"
	MESSAGE_FAILED_GET_CLASS_TOKEN          = "failed get class token"

	MESSAGE_FAILED_GET_CLASS                = "failed get class"

	MESSAGE_FAILED_UPDATE_CLASS             = "failed update class"
	MESSAGE_FAILED_DELETE_CLASS             = "failed delete class"

	// Success
	MESSAGE_SUCCESS_CREATE_CLASS             = "success create class"
	MESSAGE_SUCCESS_GET_LIST_CLASS           = "success get list class"
	MESSAGE_SUCCESS_GET_CLASS                = "success get class"
	MESSAGE_SUCCESS_UPDATE_CLASS             = "success update class"
	MESSAGE_SUCCESS_DELETE_CLASS             = "success delete class"

)

var (
	ErrCreateClass             = errors.New("failed to create class")
	ErrGetAllClass             = errors.New("failed to get all class")
	ErrGetClassById            = errors.New("failed to get class by id")
	ErrGetClassByEmail         = errors.New("failed to get class by email")

	ErrUpdateClass             = errors.New("failed to update class")
	ErrClassNotAdmin           = errors.New("class not admin")
	ErrClassNotFound           = errors.New("class not found")

	ErrDeleteClass             = errors.New("failed to delete class")

)

type (
	ClassCreateRequest struct {
		Name       string `json:"name" binding:"required"`
		ShortName  string `json:"short_name" binding:"required"`
	}

	ClassResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name" `
		ShortName  string `json:"short_name"`
	}

	ClassYAMLUploadRequest struct {
		File *multipart.FileHeader `form:"yaml_file" binding:"required"`
	}

	 ClassYAML struct {
		Name       string `yaml:"name"`
		Email      string `yaml:"email"`
		Noid 	   string `yaml:"noid"`
		Password   string `yaml:"password"`  
	}
	
	 ClassYAMLList struct {
		Class []ClassYAML `yaml:"users"`
	}

	FailedClassResponse struct {
		Noid   string `json:"noid"`
		Email  string `json:"email"`
		Reason string `json:"reason"`
	}

	ClassPaginationResponse struct {
		Data []ClassResponse `json:"data"`
		PaginationResponse
	}

	GetAllClassRepositoryResponse struct {
		Classes []entity.Class `json:"classes"`
		PaginationResponse
	}

	ClassUpdateRequest struct {
		Name       string `json:"name" form:"name"`
		ShortName      string `json:"short_name" form:"short_name"`
	}

	ClassUpdateResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		ShortName      string `json:"short_name"`
	}
)