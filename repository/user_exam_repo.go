package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	UserExamRepository interface {
		CreateUserExam(ctx context.Context, tx *gorm.DB, UserExam entity.UserExam) (entity.UserExam, error)
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
