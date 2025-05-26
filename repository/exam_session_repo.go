package repository

import (
	"context"
	"errors"
	"fmt"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	ExamSessionRepository interface {
		CreateSession(ctx context.Context, tx *gorm.DB, session entity.ExamSesssion)  error
		GetBySessionID(ctx context.Context, tx *gorm.DB, sessionID string) (*entity.ExamSesssion, error)
		FindByUserAndExam(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, entity.ExamSesssion, error)
		UpdateSession(ctx context.Context, tx *gorm.DB, session entity.ExamSesssion) (entity.ExamSesssion, error)
		GetByExamID(ctx context.Context, tx *gorm.DB, examId string) ([]entity.ExamSesssion, error)
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
		GetSEBkey(ctx context.Context, tx *gorm.DB, examId string) (entity.Exam, error)
		FinishSession(ctx context.Context, tx *gorm.DB, UserId string, ExamId string) ( error)
	}

	examSessionRepository struct {
		db *gorm.DB
	}
)

func NewExamSessionRepository(db *gorm.DB) ExamSessionRepository {
	return &examSessionRepository{
		db: db,
	}
}

func (er *examSessionRepository) CreateSession(ctx context.Context, tx *gorm.DB, session entity.ExamSesssion)  error {
	if tx == nil {
		tx = er.db
	}

	if err := tx.WithContext(ctx).Create(&session).Error; err != nil {
		return err
	}

	return  nil
}

func (r *examSessionRepository) GetBySessionID(ctx context.Context, tx *gorm.DB, sessionID string) (*entity.ExamSesssion, error) {
	if tx == nil {
		tx = r.db
	}

	var session entity.ExamSesssion
	err := tx.WithContext(ctx).
		Preload("User").
		Where("session_id = ?", sessionID).
		First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}


func (r *examSessionRepository) FindByUserAndExam(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, entity.ExamSesssion, error) {
	var session entity.ExamSesssion
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND exam_id = ?", userId, examId).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, entity.ExamSesssion{}, nil
		}
		return false, entity.ExamSesssion{}, err
	}

	return true, session, nil
}

func (r *examSessionRepository) GetByExamID(ctx context.Context, tx *gorm.DB, examId string) ([]entity.ExamSesssion, error) {
	if tx == nil {
		tx = r.db
	}

	var sessions []entity.ExamSesssion
	err := tx.WithContext(ctx).
		Preload("User").
		Where("exam_id = ?", examId).
		Find(&sessions).Error

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *examSessionRepository) UpdateSession(ctx context.Context, tx *gorm.DB, session entity.ExamSesssion) (entity.ExamSesssion, error) {
	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).
		Model(&entity.ExamSesssion{}).
		Where("user_id = ? AND exam_id = ?", session.UserID, session.ExamID).
		Updates(map[string]interface{}{
			"session_id": session.SessionID,
			"ip_address": session.IPAddress,
			"user_agent": session.UserAgent,
			"device":     session.Device,
			"updated_at": session.Timestamp.UpdatedAt,
		}).Error

	if err != nil {
		return entity.ExamSesssion{}, err
	}

	return session, nil
}

func (r *examSessionRepository) FinishSession(ctx context.Context, tx *gorm.DB, UserId string, ExamId string) ( error) {
	if tx == nil {
		tx = r.db
	}

	// Only update relevant fields, not the entire struct
	err := tx.WithContext(ctx).
		Model(&entity.ExamSesssion{}).
		Where("user_id = ? AND exam_id = ?", UserId, ExamId).
		Updates(map[string]interface{}{
			"status": 1,
		}).Error

	if err != nil {
		return  err
	}

	fmt.Println("success update status repo")
	return  nil
}

func (r *examSessionRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.ExamSesssion{}).Error
}

func (ur *examSessionRepository) GetSEBkey(ctx context.Context, tx *gorm.DB, examId string) (entity.Exam, error) {
	if tx == nil {
		tx = ur.db
	}

	var exam entity.Exam
	if err := tx.WithContext(ctx).Where("id = ?", examId).Take(&exam).Error; err != nil {
		return entity.Exam{}, err
	}

	return exam, nil
}

