package repository

import (
	"context"
	"mods/domain/entity"
	domain "mods/domain/repository"

	"gorm.io/gorm"
)



type testCaseRepository struct {
	db *gorm.DB
}

func NewTestCaseRepository(db *gorm.DB) domain.TestCaseRepository {
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

	if err := tx.WithContext(ctx).Updates(&testCase).Error; err != nil {
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

func (pr *testCaseRepository)  CheckUserInTestCaseClass(ctx context.Context, tx *gorm.DB, userID, testCaseID string) (bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var count int64

	err := tx.WithContext(ctx).
		Table("user_classes").
		Joins("JOIN classes ON user_classes.class_id = classes.id").
		Joins("JOIN exams ON exams.class_id = classes.id").
		Joins("JOIN problems ON problems.exam_id = exams.id").
		Joins("JOIN test_cases ON test_cases.problem_id = problems.id").
		Where("user_classes.user_id = ? AND test_cases.id = ?", userID, testCaseID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (pr *testCaseRepository) IsUserInProblemClass(ctx context.Context, tx *gorm.DB, userID, problemID string) (bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var count int64
	err := tx.WithContext(ctx).
		Table("user_classes").
		Joins("JOIN exams ON user_classes.class_id = exams.class_id").
		Joins("JOIN problems ON exams.id = problems.exam_id").
		Where("user_classes.user_id = ? AND problems.id = ?", userID, problemID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}