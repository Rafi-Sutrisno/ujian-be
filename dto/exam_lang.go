package dto

import (
	"errors"

	"mods/entity"
)

const (
	// Failed

	MESSAGE_FAILED_CREATE_EXAM_LANG           = "failed create exam_lang"
	MESSAGE_FAILED_GET_LIST_EXAM_LANG           = "failed get list exam_lang"
	MESSAGE_FAILED_GET_EXAM_LANG_TOKEN          = "failed get exam_lang token"

	MESSAGE_FAILED_GET_EXAM_LANG                = "failed get exam_lang"

	MESSAGE_FAILED_UPDATE_EXAM_LANG             = "failed update exam_lang"
	MESSAGE_FAILED_DELETE_EXAM_LANG             = "failed delete exam_lang"

	// Success
	MESSAGE_SUCCESS_CREATE_EXAM_LANG             = "success create exam_lang"
	MESSAGE_SUCCESS_GET_LIST_EXAM_LANG           = "success get list exam_lang"
	MESSAGE_SUCCESS_GET_EXAM_LANG                = "success get exam_lang"
	MESSAGE_SUCCESS_UPDATE_EXAM_LANG             = "success update exam_lang"
	MESSAGE_SUCCESS_DELETE_EXAM_LANG             = "success delete exam_lang"

)

var (
	ErrCreateExamLang             = errors.New("failed to create exam_lang")
	ErrGetAllExamLang             = errors.New("failed to get all exam_lang")
	ErrGetAllExamLangByExamId     = errors.New("failed to get all exam_lang by exam id")
	ErrGetAllExamLangByLangId     = errors.New("failed to get all exam_lang by lang id")
	ErrGetExamLangById            = errors.New("failed to get exam_lang by id")
	ErrGetExamLangByEmail         = errors.New("failed to get exam_lang by email")

	ErrUpdateExamLang             = errors.New("failed to update exam_lang")
	ErrExamLangNotAdmin           = errors.New("exam_lang not admin")
	ErrExamLangNotFound           = errors.New("exam_lang not found")

	ErrDeleteExamLang             = errors.New("failed to delete exam_lang")

)

type (
	ExamLangCreateRequest struct {
		ExamID       string `json:"exam_id" binding:"required"`
		LangID       uint `json:"lang_id" binding:"required"`
	}

	ExamLangResponse struct {
		ID           string `json:"id"`
		ExamID       string `json:"exam_id"`
		LangID       uint `json:"lang_id"`
	}

	ExamLangPaginationResponse struct {
		Data []ExamLangResponse `json:"data"`
		PaginationResponse
	}

	GetAllExamLangRepositoryResponse struct {
		ExamLang []entity.ExamLang `json:"exam_lang"`
		PaginationResponse
	}
)