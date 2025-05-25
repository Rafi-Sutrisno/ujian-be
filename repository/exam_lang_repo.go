package repository

import (
	"context"
	"mods/entity"

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

	examLangRepository struct {
		db *gorm.DB
	}
)

func NewExamLangRepository(db *gorm.DB) ExamLangRepository {
	return &examLangRepository{
		db: db,
	}
}

func (r *examLangRepository) GetAllByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.ExamLang, error) {
	if tx == nil {
		tx = r.db
	}

	var examLangs []entity.ExamLang
	if err := tx.WithContext(ctx).Where("exam_id = ?", examID).Find(&examLangs).Error; err != nil {
		return nil, err
	}

	return examLangs, nil
}

func (r *examLangRepository) GetAllByLangID(ctx context.Context, tx *gorm.DB, langID uint) ([]entity.ExamLang, error) {
	if tx == nil {
		tx = r.db
	}

	var examLangs []entity.ExamLang
	if err := tx.WithContext(ctx).Where("lang_id = ?", langID).Find(&examLangs).Error; err != nil {
		return nil, err
	}

	return examLangs, nil
}

func (r *examLangRepository) Create(ctx context.Context, tx *gorm.DB, examLang entity.ExamLang) (entity.ExamLang, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&examLang).Error; err != nil {
		return entity.ExamLang{}, err
	}

	return examLang, nil
}

func (r *examLangRepository) CreateMany(ctx context.Context, tx *gorm.DB, examLangs []entity.ExamLang) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&examLangs).Error; err != nil {
		return err
	}

	return nil
}

func (r *examLangRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.ExamLang{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *examLangRepository) DeleteByExamID(ctx context.Context, tx *gorm.DB, examID string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("exam_id = ?", examID).Delete(&entity.ExamLang{}).Error; err != nil {
		return err
	}

	return nil
}
