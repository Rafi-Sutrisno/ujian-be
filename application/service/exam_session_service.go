package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
	"strings"
	"time"
)

type (
	examSessionService struct {
		examSessionRepository domain.ExamSessionRepository
		authRepo domain.AuthRepo
	}

	ExamSessionService interface {
		CreateorUpdateSession(
			ctx context.Context,
			req dto.ExamSessionCreateRequest,
			sessionId string,
			userId string,
			ipAddress string,
			userAgent string,
			requestHash, configKeyHash,
			fullURL string,
		) (dto.ExamSessionCreateResponse, string, int64, error)
		GetBySessionID(ctx context.Context, sessionID string) (*dto.ExamSessionGetResponse, error)
		GetByExamID(ctx context.Context, examId string) ([]dto.ExamSessionGetResponse, error)
		FinishSession(ctx context.Context, UserId string, ExamId string) error
		DeleteByID(ctx context.Context, id string) error
	}
)

func NewExamSessionService(er domain.ExamSessionRepository, authRepo domain.AuthRepo) ExamSessionService {
	return &examSessionService{
		examSessionRepository: er,
		authRepo: authRepo,
	}
}

func (s *examSessionService) CreateorUpdateSession(
	ctx context.Context,
	req dto.ExamSessionCreateRequest,
	sessionId string,
	userId string,
	ipAddress string,
	userAgent string,
	requestHash, configKeyHash,
	fullURL string,
) (dto.ExamSessionCreateResponse, string, int64, error) {
	timeleft, err := s.authRepo.CanStartExam(ctx, userAgent, requestHash, configKeyHash, fullURL, userId, req.ExamID);
	if  err != nil {
		return dto.ExamSessionCreateResponse{},"", 0, err
	}
	

	exists, row, err := s.examSessionRepository.FindByUserAndExam(ctx, nil, userId, req.ExamID)
	if err != nil {
		return dto.ExamSessionCreateResponse{}, "", 0, err
	}
	if exists {
		fmt.Println("ini status session:", row.Status)
		if(row.Status != 0){
			return dto.ExamSessionCreateResponse{}, "", 0, errors.New("you have already finished this exam")
		}

		if(sessionId == row.SessionID){
			return dto.ExamSessionCreateResponse{
				UserID:          userId,
				ExamID: 		 req.ExamID,
			}, sessionId, timeleft, nil
		}

		sessionID, err := generateRandomToken(32)
		device := detectDevice(userAgent)
		if err != nil {
			return dto.ExamSessionCreateResponse{}, "", 0, err
		}

		newSession := entity.ExamSesssion{
			UserID:    userId,
			ExamID:    req.ExamID,
			SessionID: sessionID,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Device:    device,
			Timestamp: entity.Timestamp{
				UpdatedAt: time.Now(),
			},
		}
	
		_, err = s.examSessionRepository.UpdateSession(ctx, nil, newSession)
		if err != nil {
			return dto.ExamSessionCreateResponse{}, "", 0, err
		}

		return dto.ExamSessionCreateResponse{
			UserID:          userId,
			ExamID: 		req.ExamID,
		}, sessionID, timeleft, nil
	}else {
		sessionID, err := generateRandomToken(32) // generate a 32-byte random token
		if err != nil {
			return dto.ExamSessionCreateResponse{}, "", 0, err
		}
	
		device := detectDevice(userAgent)
	
		newSession := entity.ExamSesssion{
			UserID:    userId,
			ExamID:    req.ExamID,
			SessionID: sessionID,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Device:    device,
			Status: 0,
			Timestamp: entity.Timestamp{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
	
		err = s.examSessionRepository.CreateSession(ctx, nil, newSession)
		if err != nil {
			return dto.ExamSessionCreateResponse{}, "", 0, err
		}
	
		return dto.ExamSessionCreateResponse{
			UserID:          userId,
			ExamID: 		req.ExamID,
		}, sessionID, timeleft, nil
	}
}



func detectDevice(userAgent string) string {
	userAgent = strings.ToLower(userAgent)
	if strings.Contains(userAgent, "mobile") {
		return "Mobile"
	}
	return "Desktop"
}

func (s *examSessionService) GetBySessionID(ctx context.Context, sessionID string) (*dto.ExamSessionGetResponse, error) {
	session, err := s.examSessionRepository.GetBySessionID(ctx, nil, sessionID)
	if err != nil {
		return nil, err
	}

	return &dto.ExamSessionGetResponse{
		ID:     session.ID.String(),
		UserID: session.UserID,
		ExamID: session.ExamID,
	}, nil
}


func (s *examSessionService) GetByExamID(ctx context.Context, examId string) ([]dto.ExamSessionGetResponse, error) {
	sessions, err := s.examSessionRepository.GetByExamID(ctx, nil, examId)
	if err != nil {
		return nil, err
	}

	var responses []dto.ExamSessionGetResponse
	for _, session := range sessions {
		res := dto.ExamSessionGetResponse{
			ID:     session.ID.String(),
			UserID: session.UserID,
			ExamID: session.ExamID,
			IpAddress: session.IPAddress,
			UserAgent: session.UserAgent,
			Device: session.Device,
			Status: session.Status,
			FinishedAt: session.FinishedAt,
			User: &dto.UserResponse{
				ID:       	session.User.ID.String(),
				Name: 		session.User.Name,
				Noid:     	session.User.Noid,
			},
			CreatedAt: session.CreatedAt.String(),
		}
		responses = append(responses, res)
	}

	return responses, nil
}

func (s *examSessionService) FinishSession(ctx context.Context, UserId string, ExamId string) error {
	return s.examSessionRepository.FinishSession(ctx, nil, UserId, ExamId)
}


func (s *examSessionService) DeleteByID(ctx context.Context, id string) error {
	return s.examSessionRepository.DeleteByID(ctx, nil, id)
}


// generateRandomToken creates a base64-encoded secure random token
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
