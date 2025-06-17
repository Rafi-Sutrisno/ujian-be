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
    // log.Println("mark 0")
	existing, err := s.repo.GetByIdentifiers(ctx, req.UserID, req.ExamID, req.ProblemID, req.Language)
	// log.Println("mark 1")
	if err == nil {
		// log.Println("mark 2")
		// Record exists, so update it
		draft.ID = existing.ID
		updated, err := s.repo.Update(ctx, nil, draft)
		if err != nil {
			return dto.UserCodeDraftResponse{}, err
		}
		return dto.UserCodeDraftResponse{
			ID:        updated.ID.String(),
			UserID:    updated.UserID,
			ExamID:    updated.ExamID,
			ProblemID: updated.ProblemID,
			Language:  updated.Language,
			Code:      updated.Code,
		}, nil
	}
	// log.Println("mark 3")
	// If not found, create new
	created, err := s.repo.Create(ctx, nil, draft)
	if err != nil {
		return dto.UserCodeDraftResponse{}, err
	}
    // log.Println("mark 4")
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
