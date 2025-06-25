package service

import (
	"context"
	"fmt"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"

	"github.com/google/uuid"
)

type (
	problemService struct {
		repo domain.ProblemRepository
		authRepo domain.AuthRepo
	}

	ProblemService interface {
		GetByID(ctx context.Context, id string, userId, problemId string) (dto.ProblemResponse, error)
		GetByExamID(ctx context.Context, userAgent, requestHash, configKeyHash, fullURL, sessionId string, examID string, userId string) ([]dto.ProblemWithStatusResponse, error)
		GetAll(ctx context.Context, userId string) ([]dto.ProblemResponse, error)
		Create(ctx context.Context, req dto.ProblemCreateRequest, userId string) (dto.ProblemResponse, error)
		Update(ctx context.Context, req dto.ProblemUpdateRequest, id string, userId string) (dto.ProblemUpdateResponse, error)
		Delete(ctx context.Context, id string, userId string) error
	}
)

func NewProblemService(repo domain.ProblemRepository, authRepo domain.AuthRepo) ProblemService {
	return &problemService{
		repo: repo,
		authRepo: authRepo,
	}
}

func (ps *problemService) GetByID(ctx context.Context, id string, userId, problemId string) (dto.ProblemResponse, error) {

	problem, err := ps.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.ProblemResponse{}, err
	}

	return dto.ProblemResponse{
		ID:      		problem.ID.String(),
		// ExamID:  		problem.ExamID,
		Title:        	problem.Title,
		Description:  	problem.Description,
		Constraints:  	problem.Constraints,
		SampleInput:  	problem.SampleInput,
		SampleOutput: 	problem.SampleOutput,
		CpuTimeLimit: problem.CpuTimeLimit,
		MemoryLimit: problem.MemoryLimit,
	}, nil
}

func (ps *problemService) GetByExamID(ctx context.Context, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examID string) ([]dto.ProblemWithStatusResponse, error) {
	// authorized, err := ps.repo.IsUserInExamClass(ctx, nil, userId, examID)
	// if err != nil {
	// 	return nil, err
	// }
	// if !authorized {
	// 	return nil, dto_error.ErrAuthorizeFor("this problem")
	// }

	check, err := ps.authRepo.CanAccessProblem(ctx, userAgent, requestHash, configKeyHash, fullURL, sessionId, userId, examID)
	if err != nil {
		return nil, err
	}


	var responses []dto.ProblemWithStatusResponse
	if check {
		examProblems, err := ps.repo.GetByExamIDStudent(ctx, nil, examID, userId)
		if err != nil {
			return nil, err
		}
		responses = examProblems
	} else {
		examProblems, err := ps.repo.GetByExamID(ctx, nil, examID)
		if err != nil {
			return nil, err
		}

		
		for _, ep := range examProblems {
			responses = append(responses, dto.ProblemWithStatusResponse{
				ID:           ep.Problem.ID.String(),
				Title:        ep.Problem.Title,
				Description:  ep.Problem.Description,
				Constraints:  ep.Problem.Constraints,
				SampleInput:  ep.Problem.SampleInput,
				SampleOutput: ep.Problem.SampleOutput,
				CreatedAt:    ep.CreatedAt.String(), // ini dari ExamProblem
			})
		}
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
			// ExamID:  problem.ExamID,
			Title:        	problem.Title,
			Description:  	problem.Description,
			Constraints:  	problem.Constraints,
			SampleInput:  	problem.SampleInput,
			SampleOutput: 	problem.SampleOutput,
			CreatedAt:      problem.CreatedAt.String(),
		})
	}

	return responses, nil
}

func (ps *problemService) Create(ctx context.Context, req dto.ProblemCreateRequest, userId string) (dto.ProblemResponse, error) {

	existingProblem, err := ps.repo.GetByTitle(ctx, nil, req.Title)
	if err == nil && existingProblem.ID != uuid.Nil {
		// Found a problem with same title
		return dto.ProblemResponse{}, fmt.Errorf("problem with title '%s' already exists", req.Title)
	}

	problem := entity.Problem{
		// ExamID:  req.ExamID,
		Title:        	req.Title,
		Description:  	req.Description,
		Constraints:  	req.Constraints,
		SampleInput:  	req.SampleInput,
		SampleOutput: 	req.SampleOutput,
		CpuTimeLimit: req.CpuTimeLimit,
		MemoryLimit: req.MemoryLimit,
	}

	createdProblem, err := ps.repo.Create(ctx, nil, problem)
	if err != nil {
		return dto.ProblemResponse{}, err
	}

	return dto.ProblemResponse{
		ID:      		createdProblem.ID.String(),
		// ExamID:  		createdProblem.ExamID,
		Title:        	createdProblem.Title,
		Description:  	createdProblem.Description,
		Constraints:  	createdProblem.Constraints,
		SampleInput:  	createdProblem.SampleInput,
		SampleOutput: 	createdProblem.SampleOutput,
	}, nil
}

func (ps *problemService) Update(ctx context.Context, req dto.ProblemUpdateRequest, id string, userId string) (dto.ProblemUpdateResponse, error) {

	// authorized, err := ps.repo.IsUserInProblemClass(ctx, nil, userId, id)
	// if err != nil {
	// 	return dto.ProblemUpdateResponse{}, err
	// }
	// if !authorized {
	// 	return dto.ProblemUpdateResponse{}, dto_error.ErrAuthorizeFor("this problem")
	// }

	problem, err := ps.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.ProblemUpdateResponse{}, dto.ErrExamNotFound
	}
	
	data := entity.Problem{
		ID:  			problem.ID,
		Title:        	req.Title,
		Description:  	req.Description,
		Constraints:  	req.Constraints,
		SampleInput:  	req.SampleInput,
		SampleOutput: 	req.SampleOutput,
		CpuTimeLimit: req.CpuTimeLimit,
		MemoryLimit: req.MemoryLimit,
	}

	fmt.Println("ini problem update:", data)

	updatedProblem, err := ps.repo.Update(ctx, nil, data)
	if err != nil {
		return dto.ProblemUpdateResponse{}, err
	}

	return dto.ProblemUpdateResponse{
		ID:      updatedProblem.ID.String(),
		Title:        	req.Title,
		Description:  	req.Description,
		Constraints:  	req.Constraints,
		SampleInput:  	req.SampleInput,
		SampleOutput: 	req.SampleOutput,
	}, nil
}

func (ps *problemService) Delete(ctx context.Context, id string, userId string) error {

	// authorized, err := ps.repo.IsUserInProblemClass(ctx, nil, userId, id)
	// if err != nil {
	// 	return err
	// }
	// if !authorized {
	// 	return dto_error.ErrAuthorizeFor("this problem")
	// }

	if err := ps.repo.Delete(ctx, nil, id); err != nil {
		return dto.ErrDeleteProblem
	}
	return nil
}
