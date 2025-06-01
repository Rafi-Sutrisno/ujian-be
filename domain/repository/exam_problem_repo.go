package repository

import (
	"context"
	"mods/domain/entity"

	"gorm.io/gorm"
)


type (
	ExamProblemRepository interface {
		GetById(ctx context.Context, tx *gorm.DB, Id string) (entity.ExamProblem, error)
		GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.ExamProblem, error)
		GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.ExamProblem, error)
		GetUnassignedProblemsByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Problem, error)
		Create(ctx context.Context, tx *gorm.DB, examProblem entity.ExamProblem) (entity.ExamProblem, error)
		CreateMany(ctx context.Context, tx *gorm.DB, examProblem []entity.ExamProblem) error
		Delete(ctx context.Context, tx *gorm.DB, id string) error
		IsUserInExam(ctx context.Context, tx *gorm.DB, userID, examID string) (bool, error)
	}
)