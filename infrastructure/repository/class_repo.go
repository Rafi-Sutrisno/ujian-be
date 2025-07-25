package repository

import (
	"context"
	"math"
	"mods/domain/entity"
	domainrepo "mods/domain/repository"
	"mods/interface/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (

	classRepository struct {
		db *gorm.DB
	}
)

func NewClassRepository(db *gorm.DB) domainrepo.ClassRepository {
	return &classRepository{
		db: db,
	}
}

func (cr *classRepository) GetById(ctx context.Context, tx *gorm.DB, classId string) (entity.Class, error) {
	if tx == nil {
		tx = cr.db
	}

	var class entity.Class
	if err := tx.WithContext(ctx).Where("id = ?", classId).Preload("UserClass").Preload("Exams").First(&class).Error; err != nil {
		return entity.Class{}, err
	}

	return class, nil
}

func (cr *classRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Class, error) {
	if tx == nil {
		tx = cr.db
	}

	var classes []entity.Class
	if err := tx.WithContext(ctx).Find(&classes).Error; err != nil {
		return nil, err
	}

	return classes, nil
}

func (cr *classRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllClassRepositoryResponse, error) {
	if tx == nil {
		tx = cr.db
	}

	var classes []entity.Class
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Class{})

	if req.Search != "" {
		searchKeyword := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR short_name ILIKE ?", searchKeyword, searchKeyword)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllClassRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req.Page, req.PerPage)).Find(&classes).Error; err != nil {
		return dto.GetAllClassRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.GetAllClassRepositoryResponse{
		Classes: classes,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, nil
}

func (cr *classRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Class, error) {
	if tx == nil {
		tx = cr.db
	}

	var classes []entity.Class

	err := tx.WithContext(ctx).
		Table("classes").
		Joins("JOIN user_classes ON user_classes.class_id = classes.id").
		Where("user_classes.user_id = ?", userID).
		Find(&classes).Error

	if err != nil {
		return nil, err
	}

	return classes, nil
}

func (cr *classRepository) Create(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error) {
	if tx == nil {
		tx = cr.db
	}

	if err := tx.WithContext(ctx).Create(&class).Error; err != nil {
		return entity.Class{}, err
	}

	return class, nil
}

func (cr *classRepository) Update(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error) {
	if tx == nil {
		tx = cr.db
	}

	if err := tx.WithContext(ctx).Updates(&class).Error; err != nil {
		return entity.Class{}, err
	}

	return class, nil
}

func (cr *classRepository) Delete(ctx context.Context, tx *gorm.DB, classId string) error {
	if tx == nil {
		tx = cr.db
	}

	class := entity.Class{ID: uuid.MustParse(classId)}
	if err := tx.WithContext(ctx).Select("UserClass").Delete(class).Error; err != nil {
		return err
	}

	return nil
}
