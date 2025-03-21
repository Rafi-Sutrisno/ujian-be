package service

import (
	"errors"
	"mods/entity"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type sessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) SessionService {
	return &sessionService{db: db}
}

func (s *sessionService) SetSession(ctx *gin.Context, sessionId string, data map[string]interface{}, role string) error {
	// Set expiry e.g., 1 hour
	expiry := time.Now().Add(1 * time.Hour)

	session := entity.DatabaseData{
		Id:        sessionId,
		Data:      data,
		ExpiredAt: expiry,
	}

	return s.db.Save(&session).Error
}

func (s *sessionService) GetSessionData(ctx *gin.Context, sessionId string) (map[string]interface{}, error) {
	var session entity.DatabaseData
	err := s.db.First(&session, "id = ?", sessionId).Error
	if err != nil {
		return nil, err
	}
	// Optional: Check if session expired
	if session.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}
	return session.Data, nil
}

func (s *sessionService) SetCSRFToken(sessionId, token string) error {
	return s.db.Model(&entity.DatabaseData{}).Where("id = ?", sessionId).Update("csrf_token", token).Error
}

func (s *sessionService) GetCSRFToken(sessionId string) (string, error) {
	var session entity.DatabaseData
	err := s.db.Select("csrf_token").First(&session, "id = ?", sessionId).Error
	if err != nil {
		return "", err
	}
	return session.CSRFToken, nil
}

func (s *sessionService) ClearSession(sessionId string) error {
	return s.db.Delete(&entity.DatabaseData{}, "id = ?", sessionId).Error
}
