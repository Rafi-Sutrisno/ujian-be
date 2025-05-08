package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Submission, error)
	GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Submission, error)
	GetByExamIDandUserID(ctx context.Context, tx *gorm.DB, examID string, userID string) ([]entity.Submission, error)
	GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Submission, error)
	GetPendingSubmissions(ctx context.Context) ([]entity.Submission, error)
	Update(ctx context.Context, tx *gorm.DB, sub entity.Submission) (entity.Submission, error)
	GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.Submission, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Submission, error)
	Create(ctx context.Context, tx *gorm.DB, submission entity.Submission) (entity.Submission, error)
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{
		db: db,
	}
}

func (r *submissionRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submission entity.Submission
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&submission).Error; err != nil {
		return entity.Submission{}, err
	}

	return submission, nil
}

func (r *submissionRepository) Update(ctx context.Context, tx *gorm.DB, sub entity.Submission) (entity.Submission, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).Model(&sub).Updates(map[string]interface{}{
		"status":  sub.Status,
		"time":    sub.Time,
		"memory":  sub.Memory,
	}).Error

	return sub, err
}

func (r *submissionRepository) GetPendingSubmissions(ctx context.Context) ([]entity.Submission, error) {
	var subs []entity.Submission
	err := r.db.WithContext(ctx).Where("status = ?", "in_queue").Find(&subs).Error
	return subs, err
}


func (r *submissionRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) GetByExamIDandUserID(ctx context.Context, tx *gorm.DB, examID string, userID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("exam_id = ? AND user_id = ?", examID, userID).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("exam_id = ?", examID).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("problem_id = ?", problemID).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) Create(ctx context.Context, tx *gorm.DB, submission entity.Submission) (entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&submission).Error; err != nil {
		return entity.Submission{}, err
	}

	return submission, nil
}