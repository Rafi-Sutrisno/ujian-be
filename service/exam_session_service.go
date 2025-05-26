package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	"strings"
	"time"
)

type (
	examSessionService struct {
		examSessionRepository repository.ExamSessionRepository
	}

	ExamSessionService interface {
		CreateorUpdateSession(ctx context.Context, req dto.ExamSessionCreateRequest, sessionId string, userId string, ipAddress string,
			userAgent string, SEBRequestHash string, configKeyHash string, url string) (dto.ExamSessionCreateResponse, string, error)
		GetBySessionID(ctx context.Context, sessionID string) (*dto.ExamSessionGetResponse, error)
		GetByExamID(ctx context.Context, examId string) ([]dto.ExamSessionGetResponse, error)
		FinishSession(ctx context.Context, UserId string, ExamId string) error
		DeleteByID(ctx context.Context, id string) error
	}
)

func NewExamSessionService(er repository.ExamSessionRepository) ExamSessionService {
	return &examSessionService{
		examSessionRepository: er,
	}
}

func (s *examSessionService) CreateorUpdateSession(ctx context.Context, req dto.ExamSessionCreateRequest, sessionId string, userId string, 	ipAddress string, userAgent string, SEBRequestHash string, configKeyHash string, url string) (dto.ExamSessionCreateResponse, string, error) {
	exam, err := s.examSessionRepository.GetSEBkey(ctx, nil, req.ExamID)
	if err != nil {
		return dto.ExamSessionCreateResponse{}, "", dto.ErrExamNotFound
	}
	fmt.Println("check 1")
	if(exam.IsSEBRestricted){
		fmt.Println("check 2")
		// validateSEBRequest(url, exam.SEBKey, SEBRequestHash)
		if(exam.SEBBrowserKey != "" ){
			fmt.Println("check 3")
			if !validateSEBRequest(url, exam.SEBBrowserKey, SEBRequestHash) {
				fmt.Println("unauthorized SEB request: browser exam key hash mismatch")
				return dto.ExamSessionCreateResponse{}, "", errors.New("unauthorized SEB request: hash mismatch")
			} else {
				fmt.Println("Request is from Safe Exam Browser and correct bek key")
			}
		}
		if(exam.SEBConfigKey != "" ){
			if !validateSEBRequest(url, exam.SEBConfigKey, configKeyHash) {
				fmt.Println("unauthorized SEB request: configuration key hash mismatch")
				return dto.ExamSessionCreateResponse{}, "", errors.New("unauthorized SEB request: hash mismatch")
			} else {
				fmt.Println("Request is from Safe Exam Browser and correct config key")
			}
		}
		if(exam.SEBBrowserKey == "" && exam.SEBConfigKey == "") {
			if strings.Contains(userAgent, "SEB") {
				fmt.Println("Request is from Safe Exam Browser")
			} else {
				return dto.ExamSessionCreateResponse{}, "", errors.New("unauthorized SEB request: user agent mismatch")
			}
		}
	}
	

	exists, row, err := s.examSessionRepository.FindByUserAndExam(ctx, nil, userId, req.ExamID)
	if err != nil {
		return dto.ExamSessionCreateResponse{}, "", err
	}
	if exists {
		fmt.Println("ini status session:", row.Status)
		if(row.Status != 0){
			return dto.ExamSessionCreateResponse{}, "", errors.New("you have already finished this exam")
		}

		if(sessionId == row.SessionID){
			return dto.ExamSessionCreateResponse{
				UserID:          userId,
				ExamID: 		req.ExamID,
			}, sessionId, nil
		}

		sessionID, err := generateRandomToken(32)
		device := detectDevice(userAgent)
		if err != nil {
			return dto.ExamSessionCreateResponse{}, "", err
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
			return dto.ExamSessionCreateResponse{}, "", err
		}

		return dto.ExamSessionCreateResponse{
			UserID:          userId,
			ExamID: 		req.ExamID,
		}, sessionID, nil
	}else {
		sessionID, err := generateRandomToken(32) // generate a 32-byte random token
		if err != nil {
			return dto.ExamSessionCreateResponse{}, "", err
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
			return dto.ExamSessionCreateResponse{}, "", err
		}
	
		return dto.ExamSessionCreateResponse{
			UserID:          userId,
			ExamID: 		req.ExamID,
		}, sessionID, nil
	}
}

func validateSEBRequest(url string, key string, recvHash string) bool {

    hasher := sha256.New()

	hasher.Write([]byte(url))
	hasher.Write([]byte(key))
    
    finalHash := hasher.Sum(nil)
    hashHex := hex.EncodeToString(finalHash)

    fmt.Println("BEK/ConfigKey: Expected Hash:", hashHex)
    fmt.Println("BEK/ConfigKey: Received Hash:", recvHash)

    return hashHex == recvHash
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
			User: &dto.UserResponse{
				ID:       	session.User.ID.String(),
				Name: 		session.User.Name,
				Noid:     	session.User.Noid,
			},
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
