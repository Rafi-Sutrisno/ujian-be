package service

import (
	"context"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
)

type LanguageService interface {
	GetByID(ctx context.Context, id uint) (dto.LanguageResponse, error)
	GetAll(ctx context.Context) ([]dto.LanguageResponse, error)
	Create(ctx context.Context, req dto.LanguageRequest) (dto.LanguageResponse, error)
	Update(ctx context.Context, id uint, req dto.LanguageRequest) (dto.LanguageResponse, error)
	Delete(ctx context.Context, id uint) error
}

type languageService struct {
	repo domain.LanguageRepository
}

func NewLanguageService(repo domain.LanguageRepository) LanguageService {
	return &languageService{repo: repo}
}

func (s *languageService) GetByID(ctx context.Context, id uint) (dto.LanguageResponse, error) {
	lang, err := s.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.LanguageResponse{}, err
	}
	return dto.LanguageResponse{ID: lang.ID, Name: lang.Name, Code: lang.Code}, nil
}

func (s *languageService) GetAll(ctx context.Context) ([]dto.LanguageResponse, error) {
	langs, err := s.repo.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	var res []dto.LanguageResponse
	for _, lang := range langs {
		res = append(res, dto.LanguageResponse{ID: lang.ID, Name: lang.Name, Code: lang.Code})
	}
	return res, nil
}

func (s *languageService) Create(ctx context.Context, req dto.LanguageRequest) (dto.LanguageResponse, error) {
	lang, err := s.repo.Create(ctx, nil, entity.Language{Name: req.Name, Code: req.Code})
	if err != nil {
		return dto.LanguageResponse{}, err
	}
	return dto.LanguageResponse{ID: lang.ID, Name: lang.Name, Code: lang.Code}, nil
}

func (s *languageService) Update(ctx context.Context, id uint, req dto.LanguageRequest) (dto.LanguageResponse, error) {
	lang, err := s.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.LanguageResponse{}, err
	}
	lang.Name = req.Name
	lang.Code = req.Code

	updated, err := s.repo.Update(ctx, nil, lang)
	if err != nil {
		return dto.LanguageResponse{}, err
	}
	return dto.LanguageResponse{ID: updated.ID, Name: updated.Name, Code: updated.Code}, nil
}

func (s *languageService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, nil, id)
}
