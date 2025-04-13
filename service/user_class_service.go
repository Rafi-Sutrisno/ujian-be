package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
)

type (
	userClassService struct {
		repo repository.UserClassRepository
	}

	UserClassService interface {
		GetByUserID(ctx context.Context, userID string) ([]dto.UserClassResponse, error)
		GetByClassID(ctx context.Context, classID string) ([]dto.UserClassResponse, error)
		Create(ctx context.Context, req dto.UserClassCreateRequest) (dto.UserClassResponse, error)
		CreateMany(ctx context.Context, reqs []dto.UserClassCreateRequest) error
		Delete(ctx context.Context, id string) error
	}
)

func NewUserClassService(repo repository.UserClassRepository) UserClassService {
	return &userClassService{
		repo: repo,
	}
}

func (ucs *userClassService) GetByUserID(ctx context.Context, userID string) ([]dto.UserClassResponse, error) {
	userClasses, err := ucs.repo.GetByUserID(ctx, nil, userID)
	if err != nil {
		return nil, dto.ErrGetAllUserClassByUserId
	}

	var responses []dto.UserClassResponse
	for _, uc := range userClasses {
		responses = append(responses, dto.UserClassResponse{
			ID:        uc.ID.String(),
			UserID:    uc.UserID,
			ClassID:   uc.ClassID,
			RoleID:    uc.RoleID,
		})
	}

	return responses, nil
}

func (ucs *userClassService) GetByClassID(ctx context.Context, classID string) ([]dto.UserClassResponse, error) {
	userClasses, err := ucs.repo.GetByClassID(ctx, nil, classID)
	if err != nil {
		return nil, dto.ErrGetAllUserClassByClassId
	}

	var responses []dto.UserClassResponse
	for _, uc := range userClasses {
		responses = append(responses, dto.UserClassResponse{
			ID:        uc.ID.String(),
			UserID:    uc.UserID,
			ClassID:   uc.ClassID,
			RoleID:    uc.RoleID,
		})
	}

	return responses, nil
}

func (ucs *userClassService) Create(ctx context.Context, req dto.UserClassCreateRequest) (dto.UserClassResponse, error) {
	userClass := entity.UserClass{
		UserID:    req.UserID,
		ClassID:   req.ClassID,
		RoleID:    req.RoleID,
	}

	createdUserClass, err := ucs.repo.Create(ctx, nil, userClass)
	if err != nil {
		return dto.UserClassResponse{}, dto.ErrCreateUserClass
	}

	return dto.UserClassResponse{
		ID:        createdUserClass.ID.String(),
		UserID:    createdUserClass.UserID,
		ClassID:   createdUserClass.ClassID,
		RoleID:    createdUserClass.RoleID,
	}, nil
}

func (ucs *userClassService) CreateMany(ctx context.Context, reqs []dto.UserClassCreateRequest) error {
	var userClasses []entity.UserClass
	for _, req := range reqs {
		userClasses = append(userClasses, entity.UserClass{
			UserID:    req.UserID,
			ClassID:   req.ClassID,
			RoleID:    req.RoleID,
		})
	}

	if err := ucs.repo.CreateMany(ctx, nil, userClasses); err != nil {
		return dto.ErrCreateUserClass
	}

	return nil
}

func (ucs *userClassService) Delete(ctx context.Context, id string) error {
	userClass, err := ucs.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.ErrUserClassNotFound
	}

	if err := ucs.repo.Delete(ctx, nil, userClass.ID.String()); err != nil {
		return dto.ErrDeleteUserClass
	}
	return nil
}
