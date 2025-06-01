package repository

import (
	"context"
	"mods/domain/entity"

	"gorm.io/gorm"
)

type TestCaseRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TestCase, error)
	GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.TestCase, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TestCase, error)
	Create(ctx context.Context, tx *gorm.DB, testCase entity.TestCase) (entity.TestCase, error)
	Update(ctx context.Context, tx *gorm.DB, testCase entity.TestCase) (entity.TestCase, error)
	Delete(ctx context.Context, tx *gorm.DB, id string) error
}