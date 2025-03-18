package service

import "github.com/gin-gonic/gin"

type SessionService interface {
	SetSession(ctx *gin.Context, sessionId string, data map[string]interface{}, role string) error
	GetSessionData(ctx *gin.Context, sessionId string) (map[string]interface{}, error)
	SetCSRFToken(sessionId, token string) error
	GetCSRFToken(sessionId string) (string, error)
	ClearSession(sessionId string) error
}
