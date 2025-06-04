package repository

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthRepo interface {
	IsUserInExamClass(ctx context.Context, userId, examId string) (bool, error)
	IsExamActive(ctx context.Context, examId string) (bool, int64, error)
	HasExamSession(ctx context.Context, examId, sessionId string) (bool, error)
	CanStartExam(ctx context.Context, userId, examId string) (int64, error)
	CanAccessExam(ctx context.Context, ginCtx *gin.Context, userId, examId string) error
	CanAccessProblem(ctx context.Context, ginCtx *gin.Context, userId, problemId string) error
	ValidateSEBRequest(ginCtx *gin.Context, ctx context.Context, examId string) error
}
