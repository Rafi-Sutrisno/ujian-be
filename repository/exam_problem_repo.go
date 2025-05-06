package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	ExamProblemRepository interface {
		GetById(ctx context.Context, tx *gorm.DB, Id string) (entity.ExamProblem, error)
		GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.ExamProblem, error)
		GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.ExamProblem, error)
		GetUnassignedProblemsByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Problem, error)
		Create(ctx context.Context, tx *gorm.DB, examProblem entity.ExamProblem) (entity.ExamProblem, error)
		CreateMany(ctx context.Context, tx *gorm.DB, examProblem []entity.ExamProblem) error
		Delete(ctx context.Context, tx *gorm.DB, id string) error
		IsUserInExam(ctx context.Context, tx *gorm.DB, userID, examID string) (bool, error)
	}

	examProblemRepository struct {
		db *gorm.DB
	}
)

func NewExamProblemRepository(db *gorm.DB) ExamProblemRepository {
	return &examProblemRepository{
		db: db,
	}
}

func (ucr *examProblemRepository) GetById(ctx context.Context, tx *gorm.DB, Id string) (entity.ExamProblem, error) {
	if tx == nil {
		tx = ucr.db
	}

	var examProblem entity.ExamProblem
	if err := tx.WithContext(ctx).Where("id = ?", Id).First(&examProblem).Error; err != nil {
		return entity.ExamProblem{}, err
	}

	return examProblem, nil
}

func (ucr *examProblemRepository) GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.ExamProblem, error) {
	if tx == nil {
		tx = ucr.db
	}

	var examProblem []entity.ExamProblem
	if err := tx.WithContext(ctx).Preload("Problem").Where("exam_id = ?", examID).Find(&examProblem).Error; err != nil {
		return nil, err
	}

	return examProblem, nil
}

func (ucr *examProblemRepository) GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.ExamProblem, error) {
	if tx == nil {
		tx = ucr.db
	}

	var examProblem []entity.ExamProblem
	if err := tx.WithContext(ctx).Preload("Exam").Where("problem_id = ?", problemID).Find(&examProblem).Error; err != nil {
		return nil, err
	}

	return examProblem, nil
}

func (ucr *examProblemRepository) GetUnassignedProblemsByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Problem, error) {
	if tx == nil {
		tx = ucr.db
	}

	var problems []entity.Problem
	subQuery := tx.Model(&entity.ExamProblem{}).
		Select("problem_id").
		Where("exam_id = ?", examID)

	err := tx.WithContext(ctx).
		Where("id NOT IN (?)", subQuery).
		Find(&problems).Error

	if err != nil {
		return nil, err
	}

	return problems, nil
}


func (ucr *examProblemRepository) Create(ctx context.Context, tx *gorm.DB, examProblem entity.ExamProblem) (entity.ExamProblem, error) {
	if tx == nil {
		tx = ucr.db
	}

	if err := tx.WithContext(ctx).Create(&examProblem).Error; err != nil {
		return entity.ExamProblem{}, err
	}

	return examProblem, nil
}

func (ucr *examProblemRepository) CreateMany(ctx context.Context, tx *gorm.DB, examProblems []entity.ExamProblem) error {
	if tx == nil {
		tx = ucr.db
	}

	if err := tx.WithContext(ctx).Create(&examProblems).Error; err != nil {
		return err
	}

	return nil
}

func (ucr *examProblemRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ucr.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.ExamProblem{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (ucr *examProblemRepository) IsUserInExam(ctx context.Context, tx *gorm.DB, userID, examID string) (bool, error) {
	if tx == nil {
		tx = ucr.db
	}

	var count int64
	err := tx.WithContext(ctx).
		Table("user_classes").
		Joins("JOIN exams ON exams.class_id = user_classes.class_id").
		Where("user_classes.user_id = ? AND exams.id = ?", userID, examID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}


