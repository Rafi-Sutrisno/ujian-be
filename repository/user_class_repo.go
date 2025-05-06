package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	UserClassRepository interface {
		GetById(ctx context.Context, tx *gorm.DB, Id string) (entity.UserClass, error)
		CheckExist(ctx context.Context, tx *gorm.DB, userID string, classID string) (bool, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.UserClass, error)
		GetByClassID(ctx context.Context, tx *gorm.DB, classID string) ([]entity.UserClass, error)
		Create(ctx context.Context, tx *gorm.DB, userClass entity.UserClass) (entity.UserClass, error)
		CreateMany(ctx context.Context, tx *gorm.DB, userClasses []entity.UserClass) error
		Delete(ctx context.Context, tx *gorm.DB, id string) error
		IsUserInClass(ctx context.Context, tx *gorm.DB, userID, classID string) (bool, error)
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

func (ucr *userClassRepository) GetById(ctx context.Context, tx *gorm.DB, Id string) (entity.UserClass, error) {
	if tx == nil {
		tx = ucr.db
	}

	var userClass entity.UserClass
	if err := tx.WithContext(ctx).Where("id = ?", Id).First(&userClass).Error; err != nil {
		return entity.UserClass{}, err
	}

	return userClass, nil
}

func (ucr *userClassRepository) CheckExist(ctx context.Context, tx *gorm.DB, userID string, classID string) (bool, error) {
	if tx == nil {
		tx = ucr.db
	}

	var count int64
	if err := tx.WithContext(ctx).
		Model(&entity.UserClass{}).
		Where("user_id = ? AND class_id = ?", userID, classID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
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
	if err := tx.WithContext(ctx).Preload("User").Where("class_id = ?", classID).Find(&userClasses).Error; err != nil {
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

func (ucr *userClassRepository) IsUserInClass(ctx context.Context, tx *gorm.DB, userID, classID string) (bool, error) {
	if tx == nil {
		tx = ucr.db
	}

	var count int64
	err := tx.WithContext(ctx).Model(&entity.UserClass{}).
		Where("user_id = ? AND class_id = ?", userID, classID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

