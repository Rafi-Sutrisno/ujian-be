package repository

import (
	"context"
	"math"
	"mods/dto"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	UserRepository interface {
		RegisterUser(ctx context.Context, tx *gorm.DB,  user entity.User) (entity.User, error)
		GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error)
		CheckNoId(ctx context.Context,tx *gorm.DB, NoId string) (entity.User, bool, error)
		GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error)
		UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error
	}

	userRepository struct {
		db *gorm.DB
	}
)


func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) RegisterUser(ctx context.Context,tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = ur.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (ur *userRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = ur.db
	}

	var users []entity.User
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.User{})

	if req.Search != "" {
		searchKeyword := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR noid ILIKE ?", searchKeyword, searchKeyword, searchKeyword)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req.Page, req.PerPage)).Find(&users).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.GetAllUserRepositoryResponse{
		Users:     users,
		PaginationResponse: dto.PaginationResponse{
			Page: 		 req.Page,
			PerPage: 	 req.PerPage,
			Count: 		 count,
			MaxPage: 	 totalPage,
		},
	}, err
}


func (ur *userRepository) CheckNoId(ctx context.Context, tx *gorm.DB, NoId string) (entity.User, bool, error) {

	if tx == nil {
		tx = ur.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("noid = ?", NoId).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

func (ur *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error) {
	if tx == nil {
		tx = ur.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("id = ?", userId).Take(&user).Error; err != nil {
		return entity.User{}, err
	}


	return user, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = ur.db
	}

	if err := tx.WithContext(ctx).Updates(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error {
	if tx == nil {
		tx = ur.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.User{}, "id = ?", userId).Error; err != nil {
		return err
	}

	return nil
}
