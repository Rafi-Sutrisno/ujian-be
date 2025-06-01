package repository

import (
	"context"
	"mods/domain/entity"
	"mods/interface/dto"

	"gorm.io/gorm"
)

type (
	ClassRepository interface {
		GetById(ctx context.Context, tx *gorm.DB, classId string) (entity.Class, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllClassRepositoryResponse, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Class, error)

		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Class, error)
		Create(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error)
		Update(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error)
		Delete(ctx context.Context, tx *gorm.DB, classId string) error
		
	}

)