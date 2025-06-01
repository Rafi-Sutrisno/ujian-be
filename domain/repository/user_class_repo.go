package repository

import (
	"context"
	"mods/domain/entity"

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

)