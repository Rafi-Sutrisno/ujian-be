package repository

import (
	"context"
	"mods/domain/entity"

	"gorm.io/gorm"
)

type (
	ExamLangRepository interface {
		GetAllByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.ExamLang, error)
		GetAllByLangID(ctx context.Context, tx *gorm.DB, langID uint) ([]entity.ExamLang, error)
		Create(ctx context.Context, tx *gorm.DB, examLang entity.ExamLang) (entity.ExamLang, error)
		CreateMany(ctx context.Context, tx *gorm.DB, examLangs []entity.ExamLang) error
		Delete(ctx context.Context, tx *gorm.DB, id string) error
		DeleteByExamID(ctx context.Context, tx *gorm.DB, examID string) error

	}
)