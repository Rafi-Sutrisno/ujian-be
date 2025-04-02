package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
)

type (
	userExamService struct {
		userExamRepository repository.UserExamRepository
	}

	UserExamService interface {
		CreateUserExam(ctx context.Context, req dto.UserExamCreateRequest) (dto.UserExamResponse, error)
		GetByExamId(ctx context.Context,  examId string) ([]dto.UserExamResponse, error)
	}
)

func NewUserExamService(er repository.UserExamRepository) UserExamService {
	return &userExamService{
		userExamRepository: er,
	}
}

func (es * userExamService) CreateUserExam(ctx context.Context, req dto.UserExamCreateRequest) (dto.UserExamResponse, error){
	userExam := entity.UserExam{
		UserID: req.UserID,
		ExamID: req.ExamID,
		Role: req.Role,
	}

	userExamCreate, err := es.userExamRepository.CreateUserExam(ctx, nil, userExam)

	if err != nil {
		return dto.UserExamResponse{}, dto.ErrCreateExam
	}

	return dto.UserExamResponse{
		ID: 			userExamCreate.ID.String(),
		UserID: 		userExamCreate.UserID,
		ExamID: 		userExamCreate.ExamID,
		Role: 			userExamCreate.Role,
	}, nil
}

func (us *userExamService) GetByExamId(ctx context.Context, examId string) ([]dto.UserExamResponse, error) {
	userExams, err := us.userExamRepository.GetByExamId(ctx, nil, examId)
	if err != nil {
		return nil, err
	}

	var responses []dto.UserExamResponse
	for _, ue := range userExams {
		responses = append(responses, dto.UserExamResponse{
			ID:     ue.ID.String(),
			UserID: ue.UserID,
			ExamID: ue.ExamID,
			Role:   ue.Role,
		})
	}

	return responses, nil
}