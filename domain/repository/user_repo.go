package repository

import (
	"context"
	"mods/domain/entity"
	"mods/interface/dto"

	"gorm.io/gorm"
)

type (
	UserRepository interface {
		RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		GetAllUsers(ctx context.Context, tx *gorm.DB) ([]entity.User, error)
		GetAllStudents(ctx context.Context, tx *gorm.DB) ([]entity.User, error)
		GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error)
		CheckUsername(ctx context.Context, tx *gorm.DB, Username string) (entity.User, bool, error)
		CheckNoid(ctx context.Context, tx *gorm.DB, Noid string) (entity.User, bool, error)
		GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error)
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
		UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error
	}
)
