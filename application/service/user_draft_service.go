package service

import (
	"context"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
)

type (
	userDraftService struct {
		repo domain.UserDraftRepository
	}

	UserDraftService interface {
		SaveDraft(ctx context.Context, req dto.UserCodeDraftRequest) (dto.UserCodeDraftResponse, error)
		GetDraft(ctx context.Context, userID, examID, problemID, language string) (dto.UserCodeDraftResponse, error)
	}
)

func NewUserDraftService(repo domain.UserDraftRepository) UserDraftService {
	return &userDraftService{
		repo: repo,
	}
}

func (s *userDraftService) SaveDraft(ctx context.Context, req dto.UserCodeDraftRequest) (dto.UserCodeDraftResponse, error) {
	draft := entity.UserCodeDraft{
		UserID:    req.UserID,
		ExamID:    req.ExamID,
		ProblemID: req.ProblemID,
		Language:  req.Language,
		Code:      req.Code,
	}

	// Optional: Check if it exists -> update, otherwise create
	existing, err := s.repo.GetByIdentifiers(ctx, req.UserID, req.ExamID, req.ProblemID, req.Language)
	if err == nil {
		draft.ID = existing.ID // Preserve ID if exists
	}

	created, err := s.repo.Create(ctx, nil, draft)
	if err != nil {
		return dto.UserCodeDraftResponse{}, err
	}

	return dto.UserCodeDraftResponse{
		ID:        created.ID.String(),
		UserID:    created.UserID,
		ExamID:    created.ExamID,
		ProblemID: created.ProblemID,
		Language:  created.Language,
		Code:      created.Code,
	}, nil
}

func (s *userDraftService) GetDraft(ctx context.Context, userID, examID, problemID, language string) (dto.UserCodeDraftResponse, error) {
	draft, err := s.repo.GetByIdentifiers(ctx, userID, examID, problemID, language)
	if err != nil {
		return dto.UserCodeDraftResponse{}, err
	}

	return dto.UserCodeDraftResponse{
		ID:        draft.ID.String(),
		UserID:    draft.UserID,
		ExamID:    draft.ExamID,
		ProblemID: draft.ProblemID,
		Language:  draft.Language,
		Code:      draft.Code,
	}, nil
}
