package service

import (
	"context"
	"mods/dto"
	dto_error "mods/dto/error"
	"mods/entity"
	"mods/repository"
)

type (
	problemService struct {
		repo repository.ProblemRepository
	}

	ProblemService interface {
		GetByID(ctx context.Context, id string, userId string) (dto.ProblemResponse, error)
		GetByExamID(ctx context.Context, examID string, userId string) ([]dto.ProblemResponse, error)
		GetAll(ctx context.Context, userId string) ([]dto.ProblemResponse, error)
		Create(ctx context.Context, req dto.ProblemCreateRequest, userId string) (dto.ProblemResponse, error)
		Update(ctx context.Context, req dto.ProblemUpdateRequest, id string, userId string) (dto.ProblemUpdateResponse, error)
		Delete(ctx context.Context, id string, userId string) error
	}
)

func NewProblemService(repo repository.ProblemRepository) ProblemService {
	return &problemService{
		repo: repo,
	}
}

func (ps *problemService) GetByID(ctx context.Context, id string, userId string) (dto.ProblemResponse, error) {
	authorized, err := ps.repo.IsUserInProblemClass(ctx, nil, userId, id)
	if err != nil {
		return dto.ProblemResponse{}, err
	}
	if !authorized {
		return dto.ProblemResponse{}, dto_error.ErrAuthorizeFor("this problem")
	}

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

func (ps *problemService) GetByExamID(ctx context.Context, examID string, userId string) ([]dto.ProblemResponse, error) {
	authorized, err := ps.repo.IsUserInExamClass(ctx, nil, userId, examID)
	if err != nil {
		return nil, err
	}
	if !authorized {
		return nil, dto_error.ErrAuthorizeFor("this problem")
	}

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

func (ps *problemService) GetAll(ctx context.Context, userId string) ([]dto.ProblemResponse, error) {
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

func (ps *problemService) Create(ctx context.Context, req dto.ProblemCreateRequest, userId string) (dto.ProblemResponse, error) {
	authorized, err := ps.repo.IsUserInExamClass(ctx, nil, userId, req.ExamID)
	if err != nil {
		return dto.ProblemResponse{}, err
	}
	if !authorized {
		return dto.ProblemResponse{}, dto_error.ErrAuthorizeFor("this problem")
	}

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

func (ps *problemService) Update(ctx context.Context, req dto.ProblemUpdateRequest, id string, userId string) (dto.ProblemUpdateResponse, error) {

	authorized, err := ps.repo.IsUserInProblemClass(ctx, nil, userId, id)
	if err != nil {
		return dto.ProblemUpdateResponse{}, err
	}
	if !authorized {
		return dto.ProblemUpdateResponse{}, dto_error.ErrAuthorizeFor("this problem")
	}

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

func (ps *problemService) Delete(ctx context.Context, id string, userId string) error {

	authorized, err := ps.repo.IsUserInProblemClass(ctx, nil, userId, id)
	if err != nil {
		return err
	}
	if !authorized {
		return dto_error.ErrAuthorizeFor("this problem")
	}

	if err := ps.repo.Delete(ctx, nil, id); err != nil {
		return dto.ErrDeleteProblem
	}
	return nil
}
