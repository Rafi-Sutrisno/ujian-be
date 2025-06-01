package repository

import (
	"context"
	"mods/domain/entity"

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
)