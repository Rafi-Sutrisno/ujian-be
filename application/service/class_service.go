package service

import (
	"context"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
)

type (

	classService struct {
		classRepository domain.ClassRepository
		userClassRepository domain.UserClassRepository
	}

	ClassService interface {
		GetById(ctx context.Context, classId string) (dto.ClassResponse, error)
		GetAll(ctx context.Context) ([]dto.ClassResponse, error)
		GetByUserID(ctx context.Context, userID string) ([]dto.ClassResponse, error)
		GetAllWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ClassPaginationResponse, error)
		Create(ctx context.Context, req dto.ClassCreateRequest, userId string) (dto.ClassResponse, error)
		Update(ctx context.Context, req dto.ClassUpdateRequest, classId string) (dto.ClassUpdateResponse, error)
		Delete(ctx context.Context, classId string) error
	}
)

func NewClassService(cr domain.ClassRepository, ucr domain.UserClassRepository) ClassService {
	return &classService{
		classRepository: cr,
		userClassRepository: ucr,
	}
}

func (cs *classService) Create(ctx context.Context, req dto.ClassCreateRequest, userId string) (dto.ClassResponse, error) {
	class := entity.Class{
		Name:        	req.Name,
		Year: 			req.Year,
		Class: 			req.Class,
		ShortName:   	req.ShortName,
	}

	classCreate, err := cs.classRepository.Create(ctx, nil, class)
	if err != nil {
		return dto.ClassResponse{}, dto.ErrCreateClass
	}

	userClass := entity.UserClass{
		UserID:    userId,
		ClassID:   classCreate.ID.String(),
	}

	_, err = cs.userClassRepository.Create(ctx, nil, userClass)
	if err != nil {
		return dto.ClassResponse{}, dto.ErrCreateUserClass
	}

	return dto.ClassResponse{
		ID:          classCreate.ID.String(),
		Name:        classCreate.Name,
		Year: 		 classCreate.Year,
		Class: 		 classCreate.Class,
		ShortName:   classCreate.ShortName,
		CreatedAt:   classCreate.CreatedAt.String(),
	}, nil
}

func (cs *classService) GetAll(ctx context.Context) ([]dto.ClassResponse, error) {
	Classes, err := cs.classRepository.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	var responses []dto.ClassResponse
	for _, Class := range Classes {
		responses = append(responses, dto.ClassResponse{
			ID:            	Class.ID.String(),
			Name:        Class.Name,
			Year: 		 Class.Year,
			Class: 		 Class.Class,
			ShortName:   Class.ShortName,
			CreatedAt:   Class.CreatedAt.String(),
		})
	}

	return responses, nil
}

func (cs *classService) GetAllWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ClassPaginationResponse, error) {
	dataWithPaginate, err := cs.classRepository.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.ClassPaginationResponse{}, err
	}

	var datas []dto.ClassResponse
	for _, class := range dataWithPaginate.Classes {
		data := dto.ClassResponse{
			ID:          class.ID.String(),
			Name:        class.Name,
			Year: 		 class.Year,
			Class: 		 class.Class,
			ShortName:   class.ShortName,
			CreatedAt:   class.CreatedAt.String(),
		}

		datas = append(datas, data)
	}

	return dto.ClassPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (cs *classService) GetById(ctx context.Context, classId string) (dto.ClassResponse, error) {
	// fmt.Println("class id di service:", classId)
	class, err := cs.classRepository.GetById(ctx, nil, classId)
	if err != nil {
		return dto.ClassResponse{}, dto.ErrClassNotFound
	}

	return dto.ClassResponse{
		ID:         	class.ID.String(),
		Name:       	class.Name,
		Year: 		 class.Year,
		Class: 		 class.Class,
		ShortName:  	class.ShortName,
		CreatedAt:   class.CreatedAt.String(),
	}, nil
}

func (cs *classService) GetByUserID(ctx context.Context, userID string) ([]dto.ClassResponse, error) {
	classes, err := cs.classRepository.GetByUserID(ctx, nil, userID)
	if err != nil {
		return nil, err 
	}

	var responses []dto.ClassResponse
	for _, class := range classes {
		responses = append(responses, dto.ClassResponse{
			ID:         class.ID.String(),
			Name:       class.Name,
			Year:       class.Year,
			Class:      class.Class,
			ShortName:  class.ShortName,
			CreatedAt:  class.CreatedAt.String(),
		})
	}

	return responses, nil
}


func (cs *classService) Update(ctx context.Context, req dto.ClassUpdateRequest, classId string) (dto.ClassUpdateResponse, error) {
	class, err := cs.classRepository.GetById(ctx, nil, classId)
	if err != nil {
		return dto.ClassUpdateResponse{}, dto.ErrClassNotFound
	}

	data := entity.Class{
		ID:         class.ID,
		Name:       class.Name,
		Year: 		 class.Year,
		Class: 		 class.Class,
		ShortName:  class.ShortName,
	}

	if req.Name != "" {
		data.Name = req.Name
	}
	if req.ShortName != "" {
		data.ShortName = req.ShortName
	}
	if req.Year != "" {
		data.Year = req.Year
	}
	if req.Class != "" {
		data.Class = req.Class
	}

	classUpdate, err := cs.classRepository.Update(ctx, nil, data)
	if err != nil {
		return dto.ClassUpdateResponse{}, dto.ErrUpdateClass
	}

	return dto.ClassUpdateResponse{
		ID:         	classUpdate.ID.String(),
		Name:       	classUpdate.Name,
		Year: 		 classUpdate.Year,
		Class: 		 classUpdate.Class,
		ShortName:  	classUpdate.ShortName,
	}, nil
}

func (cs *classService) Delete(ctx context.Context, classId string) error {
	class, err := cs.classRepository.GetById(ctx, nil, classId)
	if err != nil {
		return dto.ErrClassNotFound
	}

	err = cs.classRepository.Delete(ctx, nil, class.ID.String())
	if err != nil {
		return dto.ErrDeleteClass
	}

	return nil
}