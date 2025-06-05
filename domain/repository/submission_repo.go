package repository

import (
	"context"
	"mods/domain/entity"
	"mods/interface/dto"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Submission, error)
	GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Submission, error)
	GetByExamIDandUserID(ctx context.Context, tx *gorm.DB, examID string, userID string) ([]entity.Submission, error)
	GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Submission, error)
	GetPendingSubmissions(ctx context.Context) ([]entity.Submission, error)
	Update(ctx context.Context, tx *gorm.DB, sub entity.Submission) (entity.Submission, error)
	GetCorrectSubmissionStatsByExam(ctx context.Context, examID string) ([]dto.ExamUserCorrectDTO, error)
	GetCorrectSubmissionStatsByExamandStudent(ctx context.Context, examID, userID string) (dto.ExamUserCorrectDTO, error)
	GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.Submission, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Submission, error)
	Create(ctx context.Context, tx *gorm.DB, submission entity.Submission) (entity.Submission, error)
}