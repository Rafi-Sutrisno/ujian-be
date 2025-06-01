package repository

import (
	"context"
	"mods/domain/entity"
	"mods/interface/dto"

	"gorm.io/gorm"
)

type (
	ExamRepository interface {
		CreateExam(ctx context.Context, tx *gorm.DB, exam entity.Exam) (entity.Exam, error)
		GetExamById(ctx context.Context, tx *gorm.DB, examId string) (entity.Exam, error)
		GetByClassID(ctx context.Context, tx *gorm.DB, classID string) ([]entity.Exam, error)
		GetAllExamWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllExamRepositoryResponse, error)
		UpdateExam(ctx context.Context, tx *gorm.DB, exam entity.Exam) (entity.Exam, error)
		DeleteExam(ctx context.Context, tx *gorm.DB, examId string) error
		IsUserInExamClass(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, error)
		IsUserInClass(ctx context.Context, tx *gorm.DB, userID, classID string) (bool, error)


	}
)
