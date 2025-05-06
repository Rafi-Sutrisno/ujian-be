package service

import (
	"context"
	"mods/dto"
	dto_error "mods/dto/error"
	"mods/entity"
	"mods/repository"
)

type (
	examProblemService struct {
		repo repository.ExamProblemRepository
	}

	ExamProblemService interface {
		GetByExamID(ctx context.Context, examID string, userId string) ([]dto.ExamProblemResponse, error)
		GetByProblemID(ctx context.Context, problemID string, userId string) ([]dto.ExamProblemResponse, error)
		GetUnassignedByExamID(ctx context.Context, examID string, userId string) ([]dto.ExamProblemResponse, error)
		Create(ctx context.Context, req dto.ExamProblemCreateRequest, userId string) (dto.ExamProblemResponse, error)
		CreateMany(ctx context.Context, reqs []dto.ExamProblemCreateRequest, userId string) error
		Delete(ctx context.Context, id string, userId string) error
	}
)

func NewExamProblemService(repo repository.ExamProblemRepository) ExamProblemService {
	return &examProblemService{
		repo: repo,
	}
}

func (ucs *examProblemService) GetByExamID(ctx context.Context, examID string, userId string) ([]dto.ExamProblemResponse, error) {
	exists, err := ucs.repo.IsUserInExam(ctx, nil, userId, examID)
	if err != nil {
		return nil, dto_error.ErrAuthorizeFor("this exam")
	}

	if !exists {
		return nil, dto_error.ErrAuthorizeFor("this exam")
	}

	examProblem, err := ucs.repo.GetByExamID(ctx, nil, examID)
	if err != nil {
		return nil, dto.ErrGetAllExamProblemByExamId
	}

	var responses []dto.ExamProblemResponse
	for _, uc := range examProblem {
		problem := &dto.ProblemResponse{
			ID: uc.ProblemID,
			Title: uc.Problem.Title,
		}

		responses = append(responses, dto.ExamProblemResponse{
			ID:        	uc.ID.String(),
			ExamID: 	uc.ExamID,
			ProblemID: 	uc.ProblemID,
			Problem: 	problem,
		})
	}

	return responses, nil
}

func (ucs *examProblemService) GetUnassignedByExamID(ctx context.Context, examID string, userId string) ([]dto.ExamProblemResponse, error) {
	exists, err := ucs.repo.IsUserInExam(ctx, nil, userId, examID)
	if err != nil {
		return nil, dto_error.ErrAuthorizeFor("this exam")
	}

	if !exists {
		return nil, dto_error.ErrAuthorizeFor("this exam")
	}

	examProblem, err := ucs.repo.GetUnassignedProblemsByExamID(ctx, nil, examID)
	if err != nil {
		return nil, dto.ErrGetAllExamProblemByExamId
	}

	var responses []dto.ExamProblemResponse
	for _, uc := range examProblem {
		problem := &dto.ProblemResponse{
			ID: uc.ID.String(),
			Title: uc.Title,
		}

		responses = append(responses, dto.ExamProblemResponse{
			ID:        	uc.ID.String(),
			Problem: 	problem,
		})
	}

	return responses, nil
}

func (ucs *examProblemService) GetByProblemID(ctx context.Context, examID string, userId string) ([]dto.ExamProblemResponse, error) {
	exists, err := ucs.repo.IsUserInExam(ctx, nil, userId, examID)
	if err != nil {
		return nil, dto_error.ErrAuthorizeFor("this exam")
	}

	if !exists {
		return nil, dto_error.ErrAuthorizeFor("this exam")
	}

	examProblem, err := ucs.repo.GetByProblemID(ctx, nil, examID)
	if err != nil {
		return nil, dto.ErrGetAllExamProblemByProblemId
	}

	var responses []dto.ExamProblemResponse
	for _, uc := range examProblem {
		exam := &dto.ExamResponse{
			ID: uc.ExamID,
			Name:       uc.Exam.Name,
		}

		responses = append(responses, dto.ExamProblemResponse{
			ID:        	uc.ID.String(),
			ExamID: 	uc.ExamID,
			ProblemID: 	uc.ProblemID,
			Exam: 	  	exam,
		})
	}

	return responses, nil
}


func (ucs *examProblemService) Create(ctx context.Context, req dto.ExamProblemCreateRequest, userId string) (dto.ExamProblemResponse, error) {
	exists, err := ucs.repo.IsUserInExam(ctx, nil, userId, req.ExamID)
	if err != nil {
		return dto.ExamProblemResponse{}, dto_error.ErrAuthorizeFor("this exam")
	}

	if !exists {
		return dto.ExamProblemResponse{}, dto_error.ErrAuthorizeFor("this exam")
	}

	examProblem := entity.ExamProblem{
		ExamID: 	req.ExamID,
		ProblemID: 	req.ProblemID,
	}

	createdExamProblem, err := ucs.repo.Create(ctx, nil, examProblem)
	if err != nil {
		return dto.ExamProblemResponse{}, dto.ErrCreateExamProblem
	}

	return dto.ExamProblemResponse{
		ID:        createdExamProblem.ID.String(),
		ExamID: 	createdExamProblem.ExamID,
		ProblemID: 	createdExamProblem.ProblemID,
		
	}, nil
}

func (ucs *examProblemService) CreateMany(ctx context.Context, reqs []dto.ExamProblemCreateRequest, userId string) error {
	exists, err := ucs.repo.IsUserInExam(ctx, nil, userId, reqs[0].ExamID)
	if err != nil {
		return dto_error.ErrAuthorizeFor("this exam")
	}

	if !exists {
		return dto_error.ErrAuthorizeFor("this exam")
	}
	
	var examProblems []entity.ExamProblem
	for _, req := range reqs {
		examProblems = append(examProblems, entity.ExamProblem{
			ExamID: 	req.ExamID,
			ProblemID: 	req.ProblemID,
		})
	}

	if err := ucs.repo.CreateMany(ctx, nil, examProblems); err != nil {
		return dto.ErrCreateExamProblem
	}

	return nil
}

func (ucs *examProblemService) Delete(ctx context.Context, id string, userId string) error {
	examProblem, err := ucs.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.ErrGetExamProblemById
	}

	exists, err := ucs.repo.IsUserInExam(ctx, nil, userId, examProblem.ExamID)
	if err != nil {
		return dto_error.ErrAuthorizeFor("this exam")
	}

	if !exists {
		return dto_error.ErrAuthorizeFor("this exam")
	}

	if err := ucs.repo.Delete(ctx, nil, examProblem.ID.String()); err != nil {
		return dto.ErrDeleteExamProblem
	}
	return nil
}
