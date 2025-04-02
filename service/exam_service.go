package service

import (
	"context"
	"fmt"
	"mods/dto"
	"mods/entity"
	"mods/repository"
)

type (
	examService struct {
		examRepository repository.ExamRepository
		userExamService UserExamService
	}

	ExamService interface {
		CreateExam(ctx context.Context, req dto.ExamCreateRequest, userId string) (dto.ExamResponse, error)
		GetAllExamWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ExamPaginationResponse, error)
		GetExamById(ctx context.Context, examId string) (dto.ExamResponse, error)
		Update(ctx context.Context, req dto.ExamUpdateRequest, examId string) (dto.ExamUpdateResponse, error)
		Delete(ctx context.Context, examId string) error
		
	}
)

func NewExamService(er repository.ExamRepository, userExamService UserExamService) ExamService {
	return &examService{
		examRepository: er,
		userExamService: userExamService,
	}
}

func (es *examService) CreateExam(ctx context.Context, req dto.ExamCreateRequest, userId string) (dto.ExamResponse, error) {
	// 1. Create Exam entity
	exam := entity.Exam{
		Name:        req.Name,
		ShortName:   req.ShortName,
		IsPublished: req.IsPublished,
		StartTime:   req.StartTime,
		Duration:    req.Duration,
		CreatedBy:   userId,
	}

	// 2. Save Exam to DB
	examCreate, err := es.examRepository.CreateExam(ctx, nil, exam)
	if err != nil {
		return dto.ExamResponse{}, dto.ErrCreateExam
	}

	// 3. Prepare UserExamCreateRequest DTO
	userExamReq := dto.UserExamCreateRequest{
		UserID: userId,
		ExamID: examCreate.ID.String(),
		Role:   "judge",
	}

	// 4. Create UserExam via service
	_, err = es.userExamService.CreateUserExam(ctx, userExamReq)
	if err != nil {
		return dto.ExamResponse{}, dto.ErrCreateExam
	}

	// 5. Return successful response
	return dto.ExamResponse{
		ID:          examCreate.ID.String(),
		Name:        examCreate.Name,
		ShortName:   examCreate.ShortName,
		IsPublished: examCreate.IsPublished,
		StartTime:   examCreate.StartTime,
		Duration:    examCreate.Duration,
		EndTime:     examCreate.EndTime,
		CreatedBy:   examCreate.CreatedBy,
	}, nil
}


func (us *examService) GetAllExamWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ExamPaginationResponse, error) {
	dataWithPaginate, err := us.examRepository.GetAllExamWithPagination(ctx, nil, req)
	if err != nil {
		return dto.ExamPaginationResponse{}, err
	}

	var datas []dto.ExamResponse
	for _, exam := range dataWithPaginate.Exams {
		data := dto.ExamResponse{
			ID:          exam.ID.String(),
			Name:        exam.Name,
			ShortName:   exam.ShortName,
			IsPublished: exam.IsPublished,
			StartTime:   exam.StartTime,
			Duration:    exam.Duration,
			EndTime:     exam.EndTime,
			CreatedBy:   exam.CreatedBy,
		}

		datas = append(datas, data)
	}

	return dto.ExamPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (us *examService) GetExamById(ctx context.Context, examId string) (dto.ExamResponse, error) {
	fmt.Println("exam id di service:", examId)
	exam, err := us.examRepository.GetExamById(ctx, nil, examId)
	if err != nil {
		return dto.ExamResponse{}, dto.ErrExamNotFound
	}

	return dto.ExamResponse{
		ID:         	exam.ID.String(),
		Name:       	exam.Name,
		ShortName:  	exam.ShortName,
		IsPublished: 	exam.IsPublished,
		StartTime:   	exam.StartTime,
		Duration: 		exam.Duration,
		EndTime: 		exam.EndTime,
		CreatedBy: 		exam.CreatedBy,
	}, nil
}

func (us *examService) Update(ctx context.Context, req dto.ExamUpdateRequest, examId string) (dto.ExamUpdateResponse, error) {
	exam, err := us.examRepository.GetExamById(ctx, nil, examId)
	if err != nil {
		return dto.ExamUpdateResponse{}, dto.ErrExamNotFound
	}

	data := entity.Exam{
		ID:         exam.ID,
		Name:       req.Name,
		ShortName:  req.ShortName,
		IsPublished: req.IsPublished,
		StartTime:  req.StartTime,
		Duration:   req.Duration,
		EndTime:    req.StartTime,
	}

	examUpdate, err := us.examRepository.UpdateExam(ctx, nil, data)
	if err != nil {
		return dto.ExamUpdateResponse{}, dto.ErrUpdateExam
	}

	return dto.ExamUpdateResponse{
		ID:         	examUpdate.ID.String(),
		Name:       	examUpdate.Name,
		ShortName:  	examUpdate.ShortName,
		IsPublished: 	examUpdate.IsPublished,
		StartTime: 		examUpdate.StartTime,
		Duration: 		examUpdate.Duration,
		EndTime: 		examUpdate.EndTime,
		CreatedBy: 		examUpdate.CreatedBy,
	}, nil
}

func (us *examService) Delete(ctx context.Context, examId string) error {
	exam, err := us.examRepository.GetExamById(ctx, nil, examId)
	if err != nil {
		return dto.ErrExamNotFound
	}

	err = us.examRepository.DeleteExam(ctx, nil, exam.ID.String())
	if err != nil {
		return dto.ErrDeleteExam
	}

	return nil
}