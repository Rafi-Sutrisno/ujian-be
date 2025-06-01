package repository

import (
	"context"
	"mods/domain/entity"

	"gorm.io/gorm"
)

type LanguageRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id uint) (entity.Language, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Language, error)
	Create(ctx context.Context, tx *gorm.DB, language entity.Language) (entity.Language, error)
	Update(ctx context.Context, tx *gorm.DB, language entity.Language) (entity.Language, error)
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}