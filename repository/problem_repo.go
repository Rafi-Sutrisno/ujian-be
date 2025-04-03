package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type ProblemRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Problem, error)
	GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Problem, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Problem, error)
	Create(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error)
	Update(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error)
	Delete(ctx context.Context, tx *gorm.DB, id string) error
}

type problemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) ProblemRepository {
	return &problemRepository{
		db: db,
	}
}

func (pr *problemRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	var problem entity.Problem
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&problem).Error; err != nil {
		return entity.Problem{}, err
	}

	return problem, nil
}

func (pr *problemRepository) GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	var problems []entity.Problem
	if err := tx.WithContext(ctx).Where("exam_id = ?", examID).Find(&problems).Error; err != nil {
		return nil, err
	}

	return problems, nil
}

func (pr *problemRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	var problems []entity.Problem
	if err := tx.WithContext(ctx).Find(&problems).Error; err != nil {
		return nil, err
	}

	return problems, nil
}

func (pr *problemRepository) Create(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Create(&problem).Error; err != nil {
		return entity.Problem{}, err
	}

	return problem, nil
}

func (pr *problemRepository) Update(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Save(&problem).Error; err != nil {
		return entity.Problem{}, err
	}

	return problem, nil
}

func (pr *problemRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.Problem{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
