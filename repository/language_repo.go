package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type LanguageRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id uint) (entity.Language, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Language, error)
	Create(ctx context.Context, tx *gorm.DB, language entity.Language) (entity.Language, error)
	Update(ctx context.Context, tx *gorm.DB, language entity.Language) (entity.Language, error)
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type languageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) LanguageRepository {
	return &languageRepository{db: db}
}

func (r *languageRepository) GetByID(ctx context.Context, tx *gorm.DB, id uint) (entity.Language, error) {
	if tx == nil {
		tx = r.db
	}
	var language entity.Language
	err := tx.WithContext(ctx).First(&language, id).Error
	return language, err
}

func (r *languageRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Language, error) {
	if tx == nil {
		tx = r.db
	}
	var languages []entity.Language
	err := tx.WithContext(ctx).Find(&languages).Error
	return languages, err
}

func (r *languageRepository) Create(ctx context.Context, tx *gorm.DB, language entity.Language) (entity.Language, error) {
	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).Create(&language).Error
	return language, err
}

func (r *languageRepository) Update(ctx context.Context, tx *gorm.DB, language entity.Language) (entity.Language, error) {
	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).Save(&language).Error
	return language, err
}

func (r *languageRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Delete(&entity.Language{}, id).Error
}
