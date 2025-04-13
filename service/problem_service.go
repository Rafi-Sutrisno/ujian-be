package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
)

type (
	problemService struct {
		repo repository.ProblemRepository
	}

	ProblemService interface {
		GetByID(ctx context.Context, id string) (dto.ProblemResponse, error)
		GetByExamID(ctx context.Context, examID string) ([]dto.ProblemResponse, error)
		GetAll(ctx context.Context) ([]dto.ProblemResponse, error)
		Create(ctx context.Context, req dto.ProblemCreateRequest) (dto.ProblemResponse, error)
		Update(ctx context.Context, req dto.ProblemUpdateRequest, id string) (dto.ProblemUpdateResponse, error)
		Delete(ctx context.Context, id string) error
	}
)

func NewProblemService(repo repository.ProblemRepository) ProblemService {
	return &problemService{
		repo: repo,
	}
}

func (ps *problemService) GetByID(ctx context.Context, id string) (dto.ProblemResponse, error) {
	problem, err := ps.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.ProblemResponse{}, err
	}

	return dto.ProblemResponse{
		ID:      		problem.ID.String(),
		ExamID:  		problem.ExamID,
		Title:        	problem.Title,
		Description:  	problem.Description,
		Constraints:  	problem.Constraints,
		SampleInput:  	problem.SampleInput,
		SampleOutput: 	problem.SampleOutput,
	}, nil
}

func (ps *problemService) GetByExamID(ctx context.Context, examID string) ([]dto.ProblemResponse, error) {
	problems, err := ps.repo.GetByExamID(ctx, nil, examID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ProblemResponse
	for _, problem := range problems {
		responses = append(responses, dto.ProblemResponse{
			ID:      		problem.ID.String(),
			ExamID:  		problem.ExamID,
			Title:        	problem.Title,
			Description:  	problem.Description,
			Constraints:  	problem.Constraints,
			SampleInput:  	problem.SampleInput,
			SampleOutput: 	problem.SampleOutput,
		})
	}

	return responses, nil
}

func (ps *problemService) GetAll(ctx context.Context) ([]dto.ProblemResponse, error) {
	problems, err := ps.repo.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	var responses []dto.ProblemResponse
	for _, problem := range problems {
		responses = append(responses, dto.ProblemResponse{
			ID:      problem.ID.String(),
			ExamID:  problem.ExamID,
			Title:        	problem.Title,
			Description:  	problem.Description,
			Constraints:  	problem.Constraints,
			SampleInput:  	problem.SampleInput,
			SampleOutput: 	problem.SampleOutput,
		})
	}

	return responses, nil
}

func (ps *problemService) Create(ctx context.Context, req dto.ProblemCreateRequest) (dto.ProblemResponse, error) {
	problem := entity.Problem{
		ExamID:  req.ExamID,
		Title:        	req.Title,
		Description:  	req.Description,
		Constraints:  	req.Constraints,
		SampleInput:  	req.SampleInput,
		SampleOutput: 	req.SampleOutput,
	}

	createdProblem, err := ps.repo.Create(ctx, nil, problem)
	if err != nil {
		return dto.ProblemResponse{}, err
	}

	return dto.ProblemResponse{
		ID:      		createdProblem.ID.String(),
		ExamID:  		createdProblem.ExamID,
		Title:        	createdProblem.Title,
		Description:  	createdProblem.Description,
		Constraints:  	createdProblem.Constraints,
		SampleInput:  	createdProblem.SampleInput,
		SampleOutput: 	createdProblem.SampleOutput,
	}, nil
}

func (ps *problemService) Update(ctx context.Context, req dto.ProblemUpdateRequest, id string) (dto.ProblemUpdateResponse, error) {

	problem, err := ps.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.ProblemUpdateResponse{}, dto.ErrExamNotFound
	}
	
	data := entity.Problem{
		ID:  			problem.ID,
		ExamID: 		problem.ExamID,
		Title:        	req.Title,
		Description:  	req.Description,
		Constraints:  	req.Constraints,
		SampleInput:  	req.SampleInput,
		SampleOutput: 	req.SampleOutput,
	}

	updatedProblem, err := ps.repo.Update(ctx, nil, data)
	if err != nil {
		return dto.ProblemUpdateResponse{}, err
	}

	return dto.ProblemUpdateResponse{
		ID:      updatedProblem.ID.String(),
		ExamID:  updatedProblem.ExamID,
		Title:        	req.Title,
		Description:  	req.Description,
		Constraints:  	req.Constraints,
		SampleInput:  	req.SampleInput,
		SampleOutput: 	req.SampleOutput,
	}, nil
}

func (ps *problemService) Delete(ctx context.Context, id string) error {
	if err := ps.repo.Delete(ctx, nil, id); err != nil {
		return dto.ErrDeleteProblem
	}
	return nil
}
