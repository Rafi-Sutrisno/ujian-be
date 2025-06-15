package repository

import (
	"context"
	"mods/domain/entity"

	"gorm.io/gorm"
)

type ProblemRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Problem, error)
	GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Problem, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Problem, error)
	Create(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error)
	Update(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error)
	Delete(ctx context.Context, tx *gorm.DB, id string) error
	GetByTitle(ctx context.Context, tx *gorm.DB, title string) (entity.Problem, error)
	// IsUserInProblemClass(ctx context.Context, tx *gorm.DB, userID, problemID string) (bool, error)
	// IsUserInExamClass(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, error)
}