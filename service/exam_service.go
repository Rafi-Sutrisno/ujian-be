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
	}

	ExamService interface {
		CreateExam(ctx context.Context, req dto.ExamCreateRequest) (dto.ExamResponse, error)
		GetAllExamWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ExamPaginationResponse, error)
		GetExamById(ctx context.Context, examId string) (dto.ExamResponse, error)
		GetByClassID(ctx context.Context, classID string) ([]dto.ExamResponse, error)
		Update(ctx context.Context, req dto.ExamUpdateRequest, examId string) (dto.ExamUpdateResponse, error)
		Delete(ctx context.Context, examId string) error
		
	}
)

func NewExamService(er repository.ExamRepository) ExamService {
	return &examService{
		examRepository: er,
	}
}

func (es *examService) CreateExam(ctx context.Context, req dto.ExamCreateRequest) (dto.ExamResponse, error) {
	// 1. Create Exam entity
	exam := entity.Exam{
		ClassID:     req.ClassID,
		Name:        req.Name,
		ShortName:   req.ShortName,
		IsPublished: req.IsPublished,
		StartTime:   req.StartTime,
		Duration:    req.Duration,
		
	}

	// 2. Save Exam to DB
	examCreate, err := es.examRepository.CreateExam(ctx, nil, exam)
	if err != nil {
		return dto.ExamResponse{}, err
	}

	// 5. Return successful response
	return dto.ExamResponse{
		ID:          examCreate.ID.String(),
		ClassID:     examCreate.ClassID,
		Name:        examCreate.Name,
		ShortName:   examCreate.ShortName,
		IsPublished: examCreate.IsPublished,
		StartTime:   examCreate.StartTime,
		Duration:    examCreate.Duration.String(),
		EndTime:     examCreate.EndTime,
	}, nil
}

func (es *examService) GetByClassID(ctx context.Context, classID string) ([]dto.ExamResponse, error) {
	exams, err := es.examRepository.GetByClassID(ctx, nil, classID)
	if err != nil {
		return nil, dto.ErrGetAllExamsByClassId
	}

	var responses []dto.ExamResponse
	for _, exam := range exams {
		

		responses = append(responses, dto.ExamResponse{
			ID:        exam.ID.String(),
			ClassID:   exam.ClassID,
			Name:  exam.Name,
			ShortName: exam.ShortName,
			IsPublished: exam.IsPublished,
			StartTime: exam.StartTime,
			Duration: exam.Duration.String(),
			EndTime: exam.EndTime,

		})
	}

	return responses, nil
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
			Duration:    exam.Duration.String(),
			EndTime:     exam.EndTime,
			
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
		ClassID: exam.ClassID,
		Name:       	exam.Name,
		ShortName:  	exam.ShortName,
		IsPublished: 	exam.IsPublished,
		StartTime:   	exam.StartTime,
		Duration: 		exam.Duration.String(),
		EndTime: 		exam.EndTime,
		
	}, nil
}

func (us *examService) Update(ctx context.Context, req dto.ExamUpdateRequest, examId string) (dto.ExamUpdateResponse, error) {
	exam, err := us.examRepository.GetExamById(ctx, nil, examId)
	if err != nil {
		return dto.ExamUpdateResponse{}, dto.ErrExamNotFound
	}

	data := entity.Exam{
		ID:         exam.ID,
		ClassID: 	exam.ClassID,
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
		ClassID:  		examUpdate.ClassID,	
		Name:       	examUpdate.Name,
		ShortName:  	examUpdate.ShortName,
		IsPublished: 	examUpdate.IsPublished,
		StartTime: 		examUpdate.StartTime,
		Duration: 		examUpdate.Duration,
		EndTime: 		examUpdate.EndTime,
		
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