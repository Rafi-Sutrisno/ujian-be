package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	ExamSessionRepository interface {
		CreateSession(ctx context.Context, tx *gorm.DB, session entity.ExamSesssion)  error
		FindByUserAndExam(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, error)
		GetByExamID(ctx context.Context, tx *gorm.DB, examId string) ([]entity.ExamSesssion, error)
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	examSessionRepository struct {
		db *gorm.DB
	}
)

func NewExamSessionRepository(db *gorm.DB) ExamSessionRepository {
	return &examSessionRepository{
		db: db,
	}
}

func (er *examSessionRepository) CreateSession(ctx context.Context, tx *gorm.DB, session entity.ExamSesssion)  error {
	if tx == nil {
		tx = er.db
	}

	if err := tx.WithContext(ctx).Create(&session).Error; err != nil {
		return err
	}

	return  nil
}

func (r *examSessionRepository) FindByUserAndExam(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.ExamSesssion{}).
		Where("user_id = ? AND exam_id = ?", userId, examId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *examSessionRepository) GetByExamID(ctx context.Context, tx *gorm.DB, examId string) ([]entity.ExamSesssion, error) {
	if tx == nil {
		tx = r.db
	}

	var sessions []entity.ExamSesssion
	err := tx.WithContext(ctx).
		Preload("User").
		Where("exam_id = ?", examId).
		Find(&sessions).Error

	if err != nil {
		return nil, err
	}

	return sessions, nil
}


func (r *examSessionRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.ExamSesssion{}).Error
}


