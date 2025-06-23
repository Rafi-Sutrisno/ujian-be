package repository

import (
	"context"
	"mods/domain/entity"
)

type AuthRepo interface {
	IsUserInExamClass(ctx context.Context, userId, examId string) (bool, error)
	IsExamActive(ctx context.Context, examId string) (bool, int64, error)
	HasExamSession(ctx context.Context, examId, sessionId string) (bool, error)
	CanStartExam(ctx context.Context, userAgent,requestHash, configKeyHash, fullURL, userId, examId string) (int64, error)
	CanAccessExam(ctx context.Context, userAgent,requestHash, configKeyHash, fullURL, sessionID, userId, examId string) error
	CanSeeExamResult(ctx context.Context,  userId, examId string) error
	CanAccessProblem(ctx context.Context,userAgent,requestHash, configKeyHash, fullURL, sessionID, userId, problemId string) (bool ,error)
	// ValidateSEBRequest(ginCtx *gin.Context, ctx context.Context, examId string) error
	GetUserById(ctx context.Context, userId string) (entity.User, error)
}
