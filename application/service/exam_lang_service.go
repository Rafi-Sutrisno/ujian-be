package service

import (
	"context"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
)

type (
	examLangService struct {
		repo domain.ExamLangRepository
	}

	ExamLangService interface {
		GetAllByExamID(ctx context.Context, examID string) ([]dto.ExamLangResponse, error)
		GetAllByLangID(ctx context.Context, langID uint) ([]dto.ExamLangResponse, error)
		Create(ctx context.Context, req dto.ExamLangCreateRequest) (dto.ExamLangResponse, error)
		CreateMany(ctx context.Context, reqs []dto.ExamLangCreateRequest) error
		Delete(ctx context.Context, id string) error
	}
)

func NewExamLangService(repo domain.ExamLangRepository) ExamLangService {
	return &examLangService{
		repo: repo,
	}
}

func (els *examLangService) GetAllByExamID(ctx context.Context, examID string) ([]dto.ExamLangResponse, error) {
	examLangs, err := els.repo.GetAllByExamID(ctx, nil, examID)
	if err != nil {
		return nil, dto.ErrGetAllExamLangByExamId
	}

	var responses []dto.ExamLangResponse
	for _, el := range examLangs {
		responses = append(responses, dto.ExamLangResponse{
			ID:     el.ID.String(),
			ExamID: el.ExamID,
			LangID: el.LangID,
		})
	}

	return responses, nil
}

func (els *examLangService) GetAllByLangID(ctx context.Context, langID uint) ([]dto.ExamLangResponse, error) {
	examLangs, err := els.repo.GetAllByLangID(ctx, nil, langID)
	if err != nil {
		return nil, dto.ErrGetAllExamLangByLangId
	}

	var responses []dto.ExamLangResponse
	for _, el := range examLangs {
		responses = append(responses, dto.ExamLangResponse{
			ID:     el.ID.String(),
			ExamID: el.ExamID,
			LangID: el.LangID,
		})
	}

	return responses, nil
}

func (els *examLangService) Create(ctx context.Context, req dto.ExamLangCreateRequest) (dto.ExamLangResponse, error) {
	examLang := entity.ExamLang{
		ExamID: req.ExamID,
		LangID: req.LangID,
	}

	createdExamLang, err := els.repo.Create(ctx, nil, examLang)
	if err != nil {
		return dto.ExamLangResponse{}, dto.ErrCreateExamLang
	}

	return dto.ExamLangResponse{
		ID:     createdExamLang.ID.String(),
		ExamID: createdExamLang.ExamID,
		LangID: createdExamLang.LangID,
	}, nil
}

func (els *examLangService) CreateMany(ctx context.Context, reqs []dto.ExamLangCreateRequest) error {
	if len(reqs) == 0 {
		return nil // nothing to process
	}

	// Delete all existing entries for the exam ID
	err := els.repo.DeleteByExamID(ctx, nil, reqs[0].ExamID)
	if err != nil {
		return err // optionally wrap this with a custom error if needed
	}

	// Prepare the new records
	var examLangs []entity.ExamLang
	for _, req := range reqs {
		examLangs = append(examLangs, entity.ExamLang{
			ExamID: req.ExamID,
			LangID: req.LangID,
		})
	}

	// Create new entries
	if err := els.repo.CreateMany(ctx, nil, examLangs); err != nil {
		return dto.ErrCreateExamLang
	}

	return nil
}


func (els *examLangService) Delete(ctx context.Context, id string) error {
	
	if err := els.repo.Delete(ctx, nil, id); err != nil {
		return dto.ErrDeleteExamLang
	}
	return nil
}