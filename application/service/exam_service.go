package service

import (
	"context"
	"fmt"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
	dto_error "mods/interface/dto/error"
)

type (
	examService struct {
		examRepository domain.ExamRepository
	}

	ExamService interface {
		CreateExam(ctx context.Context, req dto.ExamCreateRequest, userId string) (dto.ExamResponse, error)
		GetAllExamWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ExamPaginationResponse, error)
		GetExamById(ctx context.Context, examId string, userId string) (dto.ExamResponse, error)
		GetByClassID(ctx context.Context, classID string, userId string) ([]dto.ExamResponse, error)
		GetByUserID(ctx context.Context,  userId string) ([]dto.ExamResponse, error)
		Update(ctx context.Context, req dto.ExamUpdateRequest, examId string, userId string) (dto.ExamUpdateResponse, error)
		Delete(ctx context.Context, examId string, userId string) error	
	}
)

func NewExamService(er domain.ExamRepository) ExamService {
	return &examService{
		examRepository: er,
	}
}




func (es *examService) CreateExam(ctx context.Context, req dto.ExamCreateRequest, userId string) (dto.ExamResponse, error) {
	exists, err := es.examRepository.IsUserInClass(ctx, nil, userId, req.ClassID)
	if err != nil {
		return dto.ExamResponse{}, dto.ErrAuthorize
	}

	if !exists {
		return dto.ExamResponse{}, dto.ErrAuthorize 
	}
	// 1. Create Exam entity
	exam := entity.Exam{
		ClassID:     req.ClassID,
		Name:        req.Name,
		ShortName:   req.ShortName,
		IsPublished: req.IsPublished,
		StartTime:   req.StartTime,
		Duration:    req.Duration,
		IsSEBRestricted: req.IsSEBRestricted,
		SEBBrowserKey: req.SEBBrowserKey,
		SEBConfigKey: req.SEBConfigKey,
		SEBQuitURL: req.SEBQuitURL,
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

func (es *examService) GetByClassID(ctx context.Context, classID string, userId string) ([]dto.ExamResponse, error) {
	
	
	exams, err := es.examRepository.GetByClassID(ctx, nil, classID)
	if err != nil {
		return nil, dto.ErrGetAllExamsByClassId
	}

	if len(exams) == 0 {
		// return empty slice, not nil, to indicate success with no data
		return []dto.ExamResponse{}, nil
	}

	exists, err := es.examRepository.IsUserInExamClass(ctx, nil, userId, exams[0].ID.String())
	if err != nil {
		return  nil,err
	}

	if !exists {
		return  nil, dto_error.ErrAuthorizeFor("this exam")
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

func (es *examService) GetByUserID(ctx context.Context, userId string) ([]dto.ExamResponse, error) {
	exams, err := es.examRepository.GetByUserID(ctx, nil, userId)
	if err != nil {
		return nil, dto.ErrGetAllExamsByClassId
	}

	if len(exams) == 0 {
		// return empty slice, not nil, to indicate success with no data
		return []dto.ExamResponse{}, nil
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

func (us *examService) GetExamById(ctx context.Context, examId string, userId string) (dto.ExamResponse, error) {
	exists, err := us.examRepository.IsUserInExamClass(ctx, nil, userId, examId)
	if err != nil {
		return  dto.ExamResponse{},err
	}

	if !exists {
		return  dto.ExamResponse{}, dto_error.ErrAuthorizeFor("this exam")
	}
	exam, err := us.examRepository.GetExamById(ctx, nil, examId)
	if err != nil {
		return dto.ExamResponse{}, dto.ErrExamNotFound
	}

	allowedLangs := make([]dto.LanguageResponse, len(exam.ExamLang))
	for i, el := range exam.ExamLang {
		allowedLangs[i] = dto.LanguageResponse{
			ID:   el.Language.ID,
			Name: el.Language.Name,
			Code: el.Language.Code,
		}
	}


	return dto.ExamResponse{
		ID:         	exam.ID.String(),
		ClassID: 		exam.ClassID,
		Name:       	exam.Name,
		ShortName:  	exam.ShortName,
		IsPublished: 	exam.IsPublished,
		StartTime:   	exam.StartTime,
		Duration: 		exam.Duration.String(),
		EndTime: 		exam.EndTime,
		IsSEBRestricted:    exam.IsSEBRestricted,
		SEBBrowserKey: 		exam.SEBBrowserKey,
		SEBConfigKey: 		exam.SEBConfigKey,
		SEBQuitURL:         exam.SEBQuitURL,
		AllowedLanguages: allowedLangs,
	}, nil
}

func (us *examService) Update(ctx context.Context, req dto.ExamUpdateRequest, examId string, userId string) (dto.ExamUpdateResponse, error) {
	exists, err := us.examRepository.IsUserInExamClass(ctx, nil, userId, examId)
	if err != nil {
		return  dto.ExamUpdateResponse{},err
	}

	if !exists {
		return  dto.ExamUpdateResponse{}, dto_error.ErrAuthorizeFor("this exam")
	}

	exam, err := us.examRepository.GetExamById(ctx, nil, examId)
	if err != nil {
		return dto.ExamUpdateResponse{}, dto.ErrExamNotFound
	}

	data := entity.Exam{
		ID:         		exam.ID,
		ClassID: 			exam.ClassID,
		Name:       		req.Name,
		ShortName:  		req.ShortName,
		IsPublished: 		req.IsPublished,
		StartTime:  		req.StartTime,
		Duration:   		req.Duration,
		EndTime:    		req.StartTime.Add(req.Duration),
		IsSEBRestricted: 	req.IsSEBRestricted,
		SEBBrowserKey: 		req.SEBBrowserKey,
		SEBConfigKey: 		req.SEBConfigKey,
		SEBQuitURL:         req.SEBQuitURL,
	}

	fmt.Println("ini exam update:", data)

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

func (us *examService) Delete(ctx context.Context, examId string, userId string) error {

	exists, err := us.examRepository.IsUserInExamClass(ctx, nil, userId, examId)
	if err != nil {
		return  err
	}

	if !exists {
		return  dto_error.ErrAuthorizeFor("this exam")
	}

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