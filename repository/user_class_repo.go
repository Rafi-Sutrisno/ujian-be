package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	UserClassRepository interface {
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.UserClass, error)
		GetByClassID(ctx context.Context, tx *gorm.DB, classID string) ([]entity.UserClass, error)
		Create(ctx context.Context, tx *gorm.DB, userClass entity.UserClass) (entity.UserClass, error)
		CreateMany(ctx context.Context, tx *gorm.DB, userClasses []entity.UserClass) error
		Delete(ctx context.Context, tx *gorm.DB, id string) error
	}

	userClassRepository struct {
		db *gorm.DB
	}
)

func NewUserClassRepository(db *gorm.DB) UserClassRepository {
	return &userClassRepository{
		db: db,
	}
}

func (ucr *userClassRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.UserClass, error) {
	if tx == nil {
		tx = ucr.db
	}

	var userClasses []entity.UserClass
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Find(&userClasses).Error; err != nil {
		return nil, err
	}

	return userClasses, nil
}

func (ucr *userClassRepository) GetByClassID(ctx context.Context, tx *gorm.DB, classID string) ([]entity.UserClass, error) {
	if tx == nil {
		tx = ucr.db
	}

	var userClasses []entity.UserClass
	if err := tx.WithContext(ctx).Where("class_id = ?", classID).Find(&userClasses).Error; err != nil {
		return nil, err
	}

	return userClasses, nil
}

func (ucr *userClassRepository) Create(ctx context.Context, tx *gorm.DB, userClass entity.UserClass) (entity.UserClass, error) {
	if tx == nil {
		tx = ucr.db
	}

	if err := tx.WithContext(ctx).Create(&userClass).Error; err != nil {
		return entity.UserClass{}, err
	}

	return userClass, nil
}

func (ucr *userClassRepository) CreateMany(ctx context.Context, tx *gorm.DB, userClasses []entity.UserClass) error {
	if tx == nil {
		tx = ucr.db
	}

	if err := tx.WithContext(ctx).Create(&userClasses).Error; err != nil {
		return err
	}

	return nil
}

func (ucr *userClassRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ucr.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.UserClass{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
