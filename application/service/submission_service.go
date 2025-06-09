package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/infrastructure/judge0"
	"mods/interface/dto"
	"strconv"
	"strings"
	"time"
)

type SubmissionService interface {
	RunCode(ctx context.Context, req dto.Judge0Request, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examId string) (dto.Judge0Response, error)
	SubmitCode(ctx context.Context, req dto.SubmissionRequest, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examId string) (dto.SubmissionResponse, error)
	StartSubmissionPolling(ctx context.Context)
	GetCorrectSubmissionStatsByExam(ctx context.Context, examID string) ([]dto.ExamUserCorrectDTO, error)
	GetCorrectSubmissionStatsByExamandUser(ctx context.Context, examID, userID string) (dto.ExamUserCorrectDTO, error)
	GetByID(ctx context.Context, id string) (dto.SubmissionResponse, error)
	GetByExamIDandUserID(ctx context.Context, userAgent, requestHash, configKeyHash, fullURL, sessionId, userID, examID string, ) ([]dto.SubmissionResponse, error)
	GetByExamID(ctx context.Context, examID string, userID string) ([]dto.SubmissionResponse, error)
	GetByProblemID(ctx context.Context, problemID string) ([]dto.SubmissionResponse, error)
	GetByUserID(ctx context.Context, userID string) ([]dto.SubmissionResponse, error)
}

type submissionService struct {
	submissionRepo domain.SubmissionRepository
	testcaseRepo domain.TestCaseRepository
	langRepo domain.LanguageRepository
	problemRepo domain.ProblemRepository
	authRepo domain.AuthRepo
}

func NewSubmissionService(submissionRepo domain.SubmissionRepository, testcaseRepo domain.TestCaseRepository, langRepo domain.LanguageRepository, problemRepo domain.ProblemRepository, authRepo domain.AuthRepo) SubmissionService {
	return &submissionService{
		submissionRepo: submissionRepo,
		testcaseRepo: testcaseRepo,
		langRepo: langRepo,
		problemRepo: problemRepo,
		authRepo: authRepo,
	}
}

func (s *submissionService) RunCode(ctx context.Context, req dto.Judge0Request, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examId string) (dto.Judge0Response, error){
	if err := s.authRepo.CanAccessExam(ctx, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examId); err != nil {
		return dto.Judge0Response{}, err
	}
	lang, err := s.langRepo.GetByID(ctx, nil, uint(req.LanguageID))
	if err != nil{
		return dto.Judge0Response{}, err
	}
	u, err := strconv.ParseInt(lang.Code, 10, 0) 
	if err != nil {
		return dto.Judge0Response{}, err
	}

	req.LanguageID = int(u)

	return judge0.SubmitToJudge0(req)
}

func (s *submissionService) SubmitCode(ctx context.Context, req dto.SubmissionRequest, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examId string) (dto.SubmissionResponse, error) {
	if err := s.authRepo.CanAccessExam(ctx, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examId); err != nil {
		return dto.SubmissionResponse{}, err
	}
	problem, err := s.problemRepo.GetByID(ctx, nil, req.ProblemID)
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	cpuLimit := problem.CpuTimeLimit
	if cpuLimit == 0 {
		cpuLimit = 2.0
	}
	memoryLimit := problem.MemoryLimit
	if memoryLimit == 0 {
		memoryLimit = 128 * 1024
	}

	testCases, err := s.testcaseRepo.GetByProblemID(ctx, nil, req.ProblemID)
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	lang, err := s.langRepo.GetByID(ctx, nil, uint(req.LanguageID))
	if err != nil{
		return dto.SubmissionResponse{}, err
	}
	
	u, err := strconv.ParseInt(lang.Code, 10, 0) 
	if err != nil {
		return dto.SubmissionResponse{}, err
	}

	var submissions []dto.Judge0SubmissionRequest
	for _, tc := range testCases {
		submissions = append(submissions, dto.Judge0SubmissionRequest{
			LanguageID:     int(u),
			SourceCode:     base64.StdEncoding.EncodeToString([]byte(req.SourceCode)),
			Stdin:          base64.StdEncoding.EncodeToString([]byte(tc.InputData + "\n")),
			ExpectedOutput: base64.StdEncoding.EncodeToString([]byte(tc.ExpectedOutput + "\n")),
			CpuTimeLimit:   cpuLimit,
			CpuExtraTime:   0.5,
			WallTimeLimit:  cpuLimit + 2,
			MemoryLimit:    memoryLimit,
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
	fmt.Println("ini lang:", uint(req.LanguageID))

	submission := entity.Submission{
		UserID:         userId,
		ExamID:         req.ExamID,
		ProblemID:      req.ProblemID,
		Code:           base64.StdEncoding.EncodeToString([]byte(req.SourceCode)),
		LangID:         uint(req.LanguageID),
		SubmissionTime: time.Now().Format(time.RFC3339),
		StatusId:         1,
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
		Status: created.StatusId,
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
		finalStatus := 8 // default to "unknown/unset"

		for _, res := range batchResp.Submissions {
			if res.Status.ID <= 2 {
				allDone = false
				break
			}

			// Map Judge0 status to your custom status ID
			switch res.Status.ID {
			case 3: // Accepted
				if finalStatus > 2 {
					finalStatus = 2
				}
			case 4: // Wrong Answer
				if finalStatus > 3 {
					finalStatus = 3
				}
			case 6: // Compilation Error
				finalStatus = 4
			case 11: // Runtime Error (NZEC)
				if finalStatus > 5 {
					finalStatus = 5
				}
			case 5: // Time Limit Exceeded
				if finalStatus > 6 {
					finalStatus = 6
				}
			case 7: // Memory Limit Exceeded (optional Judge0 status)
				if finalStatus > 7 {
					finalStatus = 7
				}
			case 13: // Internal Error
				finalStatus = 8
			}

			// Parse time
			if res.Time != "" {
				timeVal, err := strconv.ParseFloat(res.Time, 64)
				if err != nil {
					log.Printf("poll error: invalid time format for submission %s: %v", submission.ID, err)
				} else if timeVal > maxTime {
					maxTime = timeVal
				}
			}

			// Parse memory
			if res.Memory > maxMemory {
				maxMemory = res.Memory
			}
		}

		if allDone {
			submission.StatusId = uint(finalStatus)
			submission.Time = fmt.Sprintf("%.2f", maxTime)
			submission.Memory = fmt.Sprintf("%d", maxMemory)
			log.Println("Submission final result:", submission)

			if _, err := s.submissionRepo.Update(ctx, nil, submission); err != nil {
				log.Printf("poll error: failed to update submission %s: %v", submission.ID, err)
			}
		}
	}
}



func (s *submissionService) GetCorrectSubmissionStatsByExam(ctx context.Context, examID string) ([]dto.ExamUserCorrectDTO, error) {
	// var results []dto.ExamUserCorrectDTO
	results, err := s.submissionRepo.GetCorrectSubmissionStatsByExam(ctx, examID)
	if err != nil {
		return []dto.ExamUserCorrectDTO{}, err
	}

	fmt.Println("ini hasil query:", results)
	return results, nil
}

func (s *submissionService) GetCorrectSubmissionStatsByExamandUser(ctx context.Context, examID, userID string) (dto.ExamUserCorrectDTO, error) {
	if err := s.authRepo.CanSeeExamResult(ctx, userID, examID); err != nil {
		return dto.ExamUserCorrectDTO{}, err
	}


	result, err := s.submissionRepo.GetCorrectSubmissionStatsByExamandStudent(ctx, examID, userID)
	if err != nil {
		return dto.ExamUserCorrectDTO{}, err
	}

	fmt.Println("ini hasil query:", result)
	return result, nil
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
		Status:        submission.StatusId,
	}

	return response, nil
}

func (s *submissionService) GetByExamIDandUserID(ctx context.Context, userAgent, requestHash, configKeyHash, fullURL, sessionId, userID, examID string) ([]dto.SubmissionResponse, error) {
	if err := s.authRepo.CanAccessExam(ctx, userAgent, requestHash, configKeyHash, fullURL, sessionId, userID, examID); err != nil {
		return nil, err
	}
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
			Status:        submission.StatusId,
			StatusName: submission.Status.Name,
			Time:          submission.Time,
			Memory: 	   submission.Memory,
			Problem: dto.ProblemResponse{
				Title: submission.Problem.Title,
			},
			Language: dto.LanguageResponse{
				Name: submission.Language.Name,
			},
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
			Status:        submission.StatusId,
			StatusName: submission.Status.Name,
			Time:          submission.Time,
			Memory: 	   submission.Memory,
			Problem: dto.ProblemResponse{
				Title: submission.Problem.Title,
			},
			Language: dto.LanguageResponse{
				Name: submission.Language.Name,
			},
			User: dto.UserResponse{
				Name: submission.User.Name,
				Noid: submission.User.Noid,
			},
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
			Status:        submission.StatusId,
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
			Status:        submission.StatusId,
		})
	}

	return response, nil
}

