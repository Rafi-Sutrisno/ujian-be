package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
	dto_error "mods/interface/dto/error"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	examService struct {
		examRepository domain.ExamRepository
		problemRepo domain.ProblemRepository
		examProblemRepo domain.ExamProblemRepository
		examLangRepo domain.ExamLangRepository
		authRepo domain.AuthRepo
	}

	ExamService interface {
		CreateExam(ctx context.Context, req dto.ExamCreateRequest, userId string) (dto.ExamResponse, error)
		UploadExamFromYaml(ctx context.Context, classID string, file multipart.File, userId string) (dto.ExamResponse, error)
		GetAllExamWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ExamPaginationResponse, error)
		GetExamById(ctx context.Context, examId string, userId string) (dto.ExamResponse, error)
		GetByClassID(ctx context.Context, classID string, userId string) ([]dto.ExamResponse, error)
		GetByUserID(ctx context.Context,  userId string) ([]dto.ExamResponse, error)
		Update(ctx context.Context, req dto.ExamUpdateRequest, examId string, userId string) (dto.ExamUpdateResponse, error)
		Delete(ctx context.Context, examId string, userId string) error	
	}
)

func NewExamService(er domain.ExamRepository,pr domain.ProblemRepository, epr domain.ExamProblemRepository, elr domain.ExamLangRepository, authRepo domain.AuthRepo) ExamService {
	return &examService{
		examRepository: er,
		problemRepo: pr,
		examProblemRepo: epr,
		examLangRepo: elr,
		authRepo: authRepo,
	}
}



func (es *examService) UploadExamFromYaml(ctx context.Context, classID string, file multipart.File, userId string) (dto.ExamResponse, error) {
	// decode file
	var yamlReq dto.ExamYamlRequest
	data, err := io.ReadAll(file)
	if err != nil {
		return dto.ExamResponse{}, err
	}
	err = yaml.Unmarshal(data, &yamlReq)
	if err != nil {
		return dto.ExamResponse{}, err
	}

	// prepare exam data
	duration, err := time.ParseDuration(yamlReq.DurationStr)
	if err != nil {
		return dto.ExamResponse{}, err
	}

	// Validate required fields manually
	if yamlReq.Name == "" ||
		yamlReq.ShortName == "" ||
		yamlReq.DurationStr == "" ||
		yamlReq.StartTime.IsZero() ||
		len(yamlReq.Problems) == 0 ||
		len(yamlReq.Languages) == 0 {
		return dto.ExamResponse{}, fmt.Errorf("required fields missing in YAML")
	}


	examReq := dto.ExamCreateRequest{
		ClassID:         classID, // from URL now
		Name:            yamlReq.Name,
		ShortName:       yamlReq.ShortName,
		IsPublished:     yamlReq.IsPublished,
		StartTime:       yamlReq.StartTime,
		DurationStr:     yamlReq.DurationStr,
		Duration:        duration,
		IsSEBRestricted: yamlReq.IsSEBRestricted,
		SEBBrowserKey:   yamlReq.SEBBrowserKey,
		SEBConfigKey:    yamlReq.SEBConfigKey,
		SEBQuitURL:      yamlReq.SEBQuitURL,
	}

	// create exam
	examResp, err := es.CreateExam(ctx, examReq, userId)
	if err != nil {
		return dto.ExamResponse{}, err
	}

	// for each problem title in YAML, get its ID
	var problemReqs []dto.ExamProblemCreateRequest
	for _, p := range yamlReq.Problems {
		prob, err := es.problemRepo.GetByTitle(ctx, nil, p.Title)
		if err != nil {
			return dto.ExamResponse{}, fmt.Errorf("problem with title '%s' not found", p.Title)
		}
		problemReqs = append(problemReqs, dto.ExamProblemCreateRequest{
			ExamID:    examResp.ID,
			ProblemID: prob.ID.String(),
		})
	}

	// assign problems
	err = es.examProblemRepo.CreateMany(ctx, nil, mapToEntity(problemReqs))
	if err != nil {
		return dto.ExamResponse{}, err
	}

	langMap := map[string]uint{
	"C":   1,
	"C++": 2,
	"C#":  3,
	}

	var examLangs []entity.ExamLang
	for _, lang := range yamlReq.Languages {
		id, ok := langMap[lang.Name]
		if !ok {
			return dto.ExamResponse{}, fmt.Errorf("language '%s' not supported", lang.Name)
		}
		examLangs = append(examLangs, entity.ExamLang{
			ExamID: examResp.ID,
			LangID: id,
		})
	}

	// save to database
	if err := es.examLangRepo.CreateMany(ctx, nil, examLangs); err != nil {
		return dto.ExamResponse{}, err
	}

	return examResp, nil
}

func mapToEntity(reqs []dto.ExamProblemCreateRequest) []entity.ExamProblem {
	var eps []entity.ExamProblem
	for _, r := range reqs {
		eps = append(eps, entity.ExamProblem{
			ExamID:    r.ExamID,
			ProblemID: r.ProblemID,
		})
	}
	return eps
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
	user, err := es.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	var check bool
	if(user.RoleID == 2){
		check = true
	}else{
		check = false
	}
	
	exams, err := es.examRepository.GetByClassID(ctx, nil, classID, check)
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
			CreatedAt: exam.CreatedAt.String(),
		})
	}

	return responses, nil
}

func (es *examService) GetByUserID(ctx context.Context, userId string) ([]dto.ExamResponse, error) {
	user, err := es.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	var check bool
	if(user.RoleID == 2){
		check = true
	}else{
		check = false
	}

	exams, err := es.examRepository.GetByUserID(ctx, nil, userId, check)
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
			CreatedAt: exam.CreatedAt.String(),
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
	user, err := us.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return  dto.ExamResponse{}, err
	}

	var check bool
	if(user.RoleID == 2){
		check = true
	}else{
		check = false
	}

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
		SEBBrowserKey:   ifThenElse(check, "", exam.SEBBrowserKey),
	    SEBConfigKey:    ifThenElse(check, "", exam.SEBConfigKey),
	    SEBQuitURL:      exam.SEBQuitURL,
		AllowedLanguages: allowedLangs,
		CreatedAt: exam.CreatedAt.String(),
	}, nil
}

func ifThenElse(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
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
	// fmt.Println("ini req update exam:", req)

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

	// fmt.Println("ini exam update:", data)

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