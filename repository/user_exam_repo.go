package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	UserExamRepository interface {
		CreateUserExam(ctx context.Context, tx *gorm.DB, UserExam entity.UserExam) (entity.UserExam, error)
		CheckIsJudge(ctx context.Context, tx *gorm.DB, userId string, examId string) (bool, error)
		GetByExamId(ctx context.Context, tx *gorm.DB, examId string) ([]entity.UserExam, error)
	}

	userExamRepository struct {
		db *gorm.DB
	}
)

func NewUserExamRepository(db *gorm.DB) UserExamRepository {
	return &userExamRepository{
		db: db,
	}
}

func (er *userExamRepository) CreateUserExam(ctx context.Context,tx *gorm.DB, UserExam entity.UserExam) (entity.UserExam, error) {
	if tx == nil {
		tx = er.db
	}

	if err := tx.WithContext(ctx).Create(&UserExam).Error; err != nil {
		return entity.UserExam{}, err
	}

	return UserExam, nil
}

func (er *userExamRepository) CheckIsJudge(ctx context.Context, tx *gorm.DB, userId string, examId string) (bool, error) {
	if tx == nil {
		tx = er.db
	}

	var userExam entity.UserExam
	err := tx.WithContext(ctx).
		Where("user_id = ? AND exam_id = ? AND role = ?", userId, examId, "judge").
		First(&userExam).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // Not a judge
		}
		return false, err // DB error
	}

	return true, nil // Is a judge
}

func (er *userExamRepository) GetByExamId(ctx context.Context, tx *gorm.DB, examId string) ([]entity.UserExam, error) {
	if tx == nil {
		tx = er.db
	}

	var userExams []entity.UserExam
	if err := tx.WithContext(ctx).Where("exam_id = ?", examId).Find(&userExams).Error; err != nil {
		return []entity.UserExam{}, err
	}

	return userExams, nil
}
