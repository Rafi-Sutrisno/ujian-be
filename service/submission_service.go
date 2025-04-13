package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	"time"
)

type SubmissionService interface {
	CreateSubmission(ctx context.Context, request dto.SubmissionCreateRequest) (dto.SubmissionResponse, error)
	GetByID(ctx context.Context, id string) (dto.SubmissionResponse, error)
	GetByExamID(ctx context.Context, examID string) ([]dto.SubmissionResponse, error)
	GetByProblemID(ctx context.Context, problemID string) ([]dto.SubmissionResponse, error)
	GetByUserID(ctx context.Context, userID string) ([]dto.SubmissionResponse, error)
}

type submissionService struct {
	submissionRepo repository.SubmissionRepository
}

func NewSubmissionService(submissionRepo repository.SubmissionRepository) SubmissionService {
	return &submissionService{
		submissionRepo: submissionRepo,
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

func (s *submissionService) GetByExamID(ctx context.Context, examID string) ([]dto.SubmissionResponse, error) {
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
