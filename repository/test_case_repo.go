package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type TestCaseRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TestCase, error)
	GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.TestCase, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TestCase, error)
	Create(ctx context.Context, tx *gorm.DB, testCase entity.TestCase) (entity.TestCase, error)
	Update(ctx context.Context, tx *gorm.DB, testCase entity.TestCase) (entity.TestCase, error)
	Delete(ctx context.Context, tx *gorm.DB, id string) error
}

type testCaseRepository struct {
	db *gorm.DB
}

func NewTestCaseRepository(db *gorm.DB) TestCaseRepository {
	return &testCaseRepository{
		db: db,
	}
}

func (r *testCaseRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TestCase, error) {
	if tx == nil {
		tx = r.db
	}

	var testCase entity.TestCase
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&testCase).Error; err != nil {
		return entity.TestCase{}, err
	}

	return testCase, nil
}

func (r *testCaseRepository) GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.TestCase, error) {
	if tx == nil {
		tx = r.db
	}

	var testCases []entity.TestCase
	if err := tx.WithContext(ctx).Where("problem_id = ?", problemID).Find(&testCases).Error; err != nil {
		return nil, err
	}

	return testCases, nil
}

func (r *testCaseRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TestCase, error) {
	if tx == nil {
		tx = r.db
	}

	var testCases []entity.TestCase
	if err := tx.WithContext(ctx).Find(&testCases).Error; err != nil {
		return nil, err
	}

	return testCases, nil
}

func (r *testCaseRepository) Create(ctx context.Context, tx *gorm.DB, testCase entity.TestCase) (entity.TestCase, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&testCase).Error; err != nil {
		return entity.TestCase{}, err
	}

	return testCase, nil
}

func (r *testCaseRepository) Update(ctx context.Context, tx *gorm.DB, testCase entity.TestCase) (entity.TestCase, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&testCase).Error; err != nil {
		return entity.TestCase{}, err
	}

	return testCase, nil
}

func (r *testCaseRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.TestCase{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}