package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	judge0 "mods/repository/external"
	"strconv"
	"strings"
	"time"
)

type SubmissionService interface {
	RunCode(ctx context.Context, req dto.Judge0Request) (dto.Judge0Response, error)
	SubmitCode(ctx context.Context, req dto.SubmissionRequest, userId string) (dto.SubmissionResponse, error)
	StartSubmissionPolling(ctx context.Context)
	CreateSubmission(ctx context.Context, request dto.SubmissionCreateRequest) (dto.SubmissionResponse, error)
	GetByID(ctx context.Context, id string) (dto.SubmissionResponse, error)
	GetByExamIDandUserID(ctx context.Context, examID string, userID string) ([]dto.SubmissionResponse, error)
	GetByExamID(ctx context.Context, examID string, userID string) ([]dto.SubmissionResponse, error)
	GetByProblemID(ctx context.Context, problemID string) ([]dto.SubmissionResponse, error)
	GetByUserID(ctx context.Context, userID string) ([]dto.SubmissionResponse, error)
}

type submissionService struct {
	submissionRepo repository.SubmissionRepository
	testcaseRepo repository.TestCaseRepository
}

func NewSubmissionService(submissionRepo repository.SubmissionRepository, testcaseRepo repository.TestCaseRepository) SubmissionService {
	return &submissionService{
		submissionRepo: submissionRepo,
		testcaseRepo: testcaseRepo,
	}
}

func (s *submissionService) RunCode(ctx context.Context, req dto.Judge0Request) (dto.Judge0Response, error) {
	return judge0.SubmitToJudge0(req)
}

func (s *submissionService) SubmitCode(ctx context.Context, req dto.SubmissionRequest, userId string) (dto.SubmissionResponse, error) {
	testCases, err := s.testcaseRepo.GetByProblemID(ctx, nil, req.ProblemID)
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	var submissions []dto.Judge0SubmissionRequest
	for _, tc := range testCases {
		submissions = append(submissions, dto.Judge0SubmissionRequest{
			LanguageID:     req.LanguageID,
			SourceCode:     base64.StdEncoding.EncodeToString([]byte(req.SourceCode)),
			Stdin:          base64.StdEncoding.EncodeToString([]byte(tc.InputData + "\n")),
			ExpectedOutput: base64.StdEncoding.EncodeToString([]byte(tc.ExpectedOutput + "\n")),
		})
	}

	batchReq := dto.Judge0BatchSubmissionRequest{Submissions: submissions}
	batchResp, err := judge0.SubmitToJudge0Batch(batchReq)
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	var tokenList []string
	for _, item := range batchResp {
		tokenList = append(tokenList, item.Token)
	}
	tokenStr := strings.Join(tokenList, ",")

	submission := entity.Submission{
		UserID:         userId,
		ExamID:         req.ExamID,
		ProblemID:      req.ProblemID,
		Code:           base64.StdEncoding.EncodeToString([]byte(req.SourceCode)),
		LangID:         "1",
		SubmissionTime: time.Now().Format(time.RFC3339),
		Status:         "in_queue",
		Judge0Token:    tokenStr,
		Time:           "",
		Memory:         "",
	}

	created, err := s.submissionRepo.Create(ctx, nil, submission)
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	return dto.SubmissionResponse{
		ID:     created.ID.String(),
		Status: created.Status,
	}, nil
}

func (s *submissionService) StartSubmissionPolling(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second) 
	go func() {
		for {
			select {
			case <-ticker.C:
				s.pollPendingSubmissions(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *submissionService) pollPendingSubmissions(ctx context.Context) {
	// 1. Get all in-queue submissions
	submissions, err := s.submissionRepo.GetPendingSubmissions(ctx)
	if err != nil {
		log.Printf("poll error: failed to get pending submissions: %v", err)
		return
	}
	log.Printf("masuk polling dengan submission queue: %v", submissions)

	for _, submission := range submissions {
		tokens := strings.Split(submission.Judge0Token, ",")
		batchResp, err := judge0.GetBatchResults(tokens)
		if err != nil {
			log.Printf("poll error: failed to get batch results for submission %s: %v", submission.ID, err)
			continue
		}

		// check if all results are done
		allDone := true
		maxTime := float64(0)
		maxMemory := 0
		status := "accepted"

		for _, res := range batchResp.Submissions {
			if res.Status.ID <= 2 { // In queue or Processing
				allDone = false
				break
			}

			// Convert string to float64
			timeVal, err := strconv.ParseFloat(res.Time, 64)
			if err != nil {
				log.Printf("poll error: invalid time format for submission %s: %v", submission.ID, err)
				continue
			}

			if timeVal > maxTime {
				maxTime = timeVal
			}
			if res.Memory > maxMemory {
				maxMemory = res.Memory
			}
			if res.Status.ID == 6 {
				status = "wrong_answer"
			} else if res.Status.ID == 11 {
				status = "compilation_error"
			}
		}

		if allDone {
			submission.Status = status
			submission.Time = fmt.Sprintf("%.2f", maxTime)
			submission.Memory = fmt.Sprintf("%d", maxMemory)

			if _, err := s.submissionRepo.Update(ctx, nil, submission); err != nil {
				log.Printf("poll error: failed to update submission %s: %v", submission.ID, err)
			}
		}
	}
}


func (s *submissionService) CreateSubmission(ctx context.Context, request dto.SubmissionCreateRequest) (dto.SubmissionResponse, error) {
	// Convert the request to the entity
	submission := entity.Submission{
		UserID:      request.UserID,
		ExamID:      request.ExamID,
		ProblemID:   request.ProblemID,
		LangID:      request.LangID,
		Code:        request.Code,
		SubmissionTime: time.Now().Format(time.RFC3339), // Set current time as submission time
		Status:      "accepted", // Default status to "accepted"
	}

	// Save to the database
	submitted, err := s.submissionRepo.Create(ctx, nil, submission)
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	// Map the created submission entity to the response DTO
	response := dto.SubmissionResponse{
		ID:            submitted.ID.String(),
		UserID:        submitted.UserID,
		ExamID:        submitted.ExamID,
		ProblemID:     submitted.ProblemID,
		LangID:        submitted.LangID,
		Code:          submitted.Code,
		SubmissionTime: submitted.SubmissionTime,
		Status:        submitted.Status,
	}

	return response, nil
}

func (s *submissionService) GetByID(ctx context.Context, id string) (dto.SubmissionResponse, error) {
	// Get submission by ID from the repository
	submission, err := s.submissionRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	// Map the submission entity to response DTO
	response := dto.SubmissionResponse{
		ID:            submission.ID.String(),
		UserID:        submission.UserID,
		ExamID:        submission.ExamID,
		ProblemID:     submission.ProblemID,
		LangID:        submission.LangID,
		Code:          submission.Code,
		SubmissionTime: submission.SubmissionTime,
		Status:        submission.Status,
	}

	return response, nil
}

func (s *submissionService) GetByExamIDandUserID(ctx context.Context, examID string, userID string) ([]dto.SubmissionResponse, error) {
	// Get submissions by Exam ID from the repository
	submissions, err := s.submissionRepo.GetByExamIDandUserID(ctx, nil, examID, userID)
	if err != nil {
		return nil, err
	}

	// Map the submissions entities to response DTOs
	var response []dto.SubmissionResponse
	for _, submission := range submissions {
		response = append(response, dto.SubmissionResponse{
			ID:            submission.ID.String(),
			UserID:        submission.UserID,
			ExamID:        submission.ExamID,
			ProblemID:     submission.ProblemID,
			LangID:        submission.LangID,
			Code:          submission.Code,
			SubmissionTime: submission.SubmissionTime,
			Status:        submission.Status,
			Time:          submission.Time,
			Memory: 	   submission.Memory,
		})
	}

	return response, nil
}

func (s *submissionService) GetByExamID(ctx context.Context, examID string, userID string) ([]dto.SubmissionResponse, error) {
	// Get submissions by Exam ID from the repository
	submissions, err := s.submissionRepo.GetByExamID(ctx, nil, examID)
	if err != nil {
		return nil, err
	}

	// Map the submissions entities to response DTOs
	var response []dto.SubmissionResponse
	for _, submission := range submissions {
		response = append(response, dto.SubmissionResponse{
			ID:            submission.ID.String(),
			UserID:        submission.UserID,
			ExamID:        submission.ExamID,
			ProblemID:     submission.ProblemID,
			LangID:        submission.LangID,
			Code:          submission.Code,
			SubmissionTime: submission.SubmissionTime,
			Status:        submission.Status,
			Time:          submission.Time,
			Memory: 	   submission.Memory,
		})
	}

	return response, nil
}

func (s *submissionService) GetByProblemID(ctx context.Context, problemID string) ([]dto.SubmissionResponse, error) {
	// Get submissions by Problem ID from the repository
	submissions, err := s.submissionRepo.GetByProblemID(ctx, nil, problemID)
	if err != nil {
		return nil, err
	}

	// Map the submissions entities to response DTOs
	var response []dto.SubmissionResponse
	for _, submission := range submissions {
		response = append(response, dto.SubmissionResponse{
			ID:            submission.ID.String(),
			UserID:        submission.UserID,
			ExamID:        submission.ExamID,
			ProblemID:     submission.ProblemID,
			LangID:        submission.LangID,
			Code:          submission.Code,
			SubmissionTime: submission.SubmissionTime,
			Status:        submission.Status,
		})
	}

	return response, nil
}

func (s *submissionService) GetByUserID(ctx context.Context, userID string) ([]dto.SubmissionResponse, error) {
	// Get submissions by User ID from the repository
	submissions, err := s.submissionRepo.GetByUserID(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	// Map the submissions entities to response DTOs
	var response []dto.SubmissionResponse
	for _, submission := range submissions {
		response = append(response, dto.SubmissionResponse{
			ID:            submission.ID.String(),
			UserID:        submission.UserID,
			ExamID:        submission.ExamID,
			ProblemID:     submission.ProblemID,
			LangID:        submission.LangID,
			Code:          submission.Code,
			SubmissionTime: submission.SubmissionTime,
			Status:        submission.Status,
		})
	}

	return response, nil
}

