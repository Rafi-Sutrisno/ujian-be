package repository

import (
	"context"
	"mods/domain/entity"
	domain "mods/domain/repository"

	"gorm.io/gorm"
)


type problemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) domain.ProblemRepository {
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
	err := tx.WithContext(ctx).
		Model(&entity.Problem{}).
		Joins("JOIN exam_problems ON exam_problems.problem_id = problems.id").
		Where("exam_problems.exam_id = ?", examID).
		Find(&problems).Error

	if err != nil {
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

func (pr *problemRepository) IsUserInProblemClass(ctx context.Context, tx *gorm.DB, userID, problemID string) (bool, error) {
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


func (pr *problemRepository) IsUserInExamClass(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, error) {
	if tx == nil {
		tx = pr.db
	}

	// First, get the class_id from the exam
	var exam entity.Exam
	if err := tx.WithContext(ctx).Select("class_id").Where("id = ?", examId).First(&exam).Error; err != nil {
		return false, err
	}

	// Now check if user exists in that class
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.UserClass{}).
		Where("user_id = ? AND class_id = ?", userId, exam.ClassID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

