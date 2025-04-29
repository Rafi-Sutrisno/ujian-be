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
		UserRepo repository.UserRepository
	}

	UserClassService interface {
		GetByUserID(ctx context.Context, userID string) ([]dto.UserClassResponse, error)
		GetByClassID(ctx context.Context, classID string, userId string) ([]dto.UserClassResponse, error)
		GetUnassignedUsersByClassID(ctx context.Context, classID string, userId string) ([]dto.UserResponse, error)
		Create(ctx context.Context, req dto.UserClassCreateRequest, userId string) (dto.UserClassResponse, error)
		CreateMany(ctx context.Context, reqs []dto.UserClassCreateRequest, userId string) error
		Delete(ctx context.Context, id string, userId string) error
	}
)

func NewUserClassService(repo repository.UserClassRepository, UserRepo repository.UserRepository) UserClassService {
	return &userClassService{
		repo: repo,
		UserRepo:  UserRepo,
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
			
		})
	}

	return responses, nil
}

func (ucs *userClassService) GetByClassID(ctx context.Context, classID string, userId string) ([]dto.UserClassResponse, error) {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, classID)
	if err != nil {
		return nil, dto.ErrAuthorize
	}

	if !exists {
		return nil, dto.ErrAuthorize 
	}

	userClasses, err := ucs.repo.GetByClassID(ctx, nil, classID)
	if err != nil {
		return nil, dto.ErrGetAllUserClassByClassId
	}

	var responses []dto.UserClassResponse
	for _, uc := range userClasses {
		user := &dto.UserResponse{
			Name:       uc.User.Name,
			Noid: 		uc.User.Noid,
			RoleID:     uc.User.RoleID,
			Email:      uc.User.Email,
		}

		responses = append(responses, dto.UserClassResponse{
			ID:        uc.ID.String(),
			UserID:    uc.UserID,
			ClassID:   uc.ClassID,
			User:      user,
		})
	}

	return responses, nil
}

func (ucs *userClassService) GetUnassignedUsersByClassID(ctx context.Context, classID string, userId string) ([]dto.UserResponse, error) {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, classID)
	if err != nil {
		return nil, dto.ErrAuthorize
	}

	if !exists {
		return nil, dto.ErrAuthorize 
	}

	allStudents, err := ucs.UserRepo.GetAllStudents(ctx, nil)
	if err != nil {
		return nil, err
	}

	assignedUserClasses, err := ucs.repo.GetByClassID(ctx, nil, classID)
	if err != nil {
		return nil, err
	}

	assignedMap := make(map[string]bool)
	for _, uc := range assignedUserClasses {
		assignedMap[uc.UserID] = true
	}

	var unassignedUsers []dto.UserResponse
	for _, student := range allStudents {
		if !assignedMap[student.ID.String()] {
			unassignedUsers = append(unassignedUsers, dto.UserResponse{
				ID:    student.ID.String(),
				Name:  student.Name,
				Noid: 		student.Noid,
				RoleID:       student.RoleID,
				Email: student.Email,
				
			})
		}
	}

	return unassignedUsers, nil
}


func (ucs *userClassService) Create(ctx context.Context, req dto.UserClassCreateRequest, userId string) (dto.UserClassResponse, error) {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, req.ClassID)
	if err != nil {
		return dto.UserClassResponse{}, dto.ErrAuthorize
	}

	if !exists {
		return dto.UserClassResponse{}, dto.ErrAuthorize // or any custom error you want
	}

	userClass := entity.UserClass{
		UserID:    req.UserID,
		ClassID:   req.ClassID,
		
	}

	createdUserClass, err := ucs.repo.Create(ctx, nil, userClass)
	if err != nil {
		return dto.UserClassResponse{}, dto.ErrCreateUserClass
	}

	return dto.UserClassResponse{
		ID:        createdUserClass.ID.String(),
		UserID:    createdUserClass.UserID,
		ClassID:   createdUserClass.ClassID,
		
	}, nil
}

func (ucs *userClassService) CreateMany(ctx context.Context, reqs []dto.UserClassCreateRequest, userId string) error {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, reqs[0].ClassID)
	if err != nil {
		return dto.ErrAuthorize
	}

	if !exists {
		return dto.ErrAuthorize // or any custom error you want
	}
	
	var userClasses []entity.UserClass
	for _, req := range reqs {
		userClasses = append(userClasses, entity.UserClass{
			UserID:    req.UserID,
			ClassID:   req.ClassID,
			
		})
	}

	if err := ucs.repo.CreateMany(ctx, nil, userClasses); err != nil {
		return dto.ErrCreateUserClass
	}

	return nil
}

func (ucs *userClassService) Delete(ctx context.Context, id string, userId string) error {
	userClasses, err := ucs.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.ErrGetAllUserClass
	}

	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, userClasses.ClassID)
	if err != nil {
		return  dto.ErrAuthorize
	}

	if !exists {
		return  dto.ErrAuthorize // or any custom error you want
	}
	userClass, err := ucs.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.ErrUserClassNotFound
	}

	if err := ucs.repo.Delete(ctx, nil, userClass.ID.String()); err != nil {
		return dto.ErrDeleteUserClass
	}
	return nil
}
