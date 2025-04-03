package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
)

type (
	classService struct {
		classRepository repository.ClassRepository
	}

	ClassService interface {
		GetById(ctx context.Context, classId string) (dto.ClassResponse, error)
		GetAllWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.ClassPaginationResponse, error)
		Create(ctx context.Context, req dto.ClassCreateRequest) (dto.ClassResponse, error)
		Update(ctx context.Context, req dto.ClassUpdateRequest, classId string) (dto.ClassUpdateResponse, error)
		Delete(ctx context.Context, classId string) error
		
	}
)

func NewClassService(cr repository.ClassRepository) ClassService {
	return &classService{
		classRepository: cr,
	}
}

func (cs *classService) Create(ctx context.Context, req dto.ClassCreateRequest) (dto.ClassResponse, error) {
	class := entity.Class{
		Name:        req.Name,
		ShortName:   req.ShortName,
	}

	classCreate, err := cs.classRepository.Create(ctx, nil, class)
	if err != nil {
		return dto.ClassResponse{}, dto.ErrCreateClass
	}

	return dto.ClassResponse{
		ID:          classCreate.ID.String(),
		Name:        classCreate.Name,
		ShortName:   classCreate.ShortName,
	}, nil
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
			ShortName:   class.ShortName,
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
		ShortName:  	class.ShortName,
	}, nil
}

func (cs *classService) Update(ctx context.Context, req dto.ClassUpdateRequest, classId string) (dto.ClassUpdateResponse, error) {
	class, err := cs.classRepository.GetById(ctx, nil, classId)
	if err != nil {
		return dto.ClassUpdateResponse{}, dto.ErrClassNotFound
	}

	data := entity.Class{
		ID:         class.ID,
		Name:       req.Name,
		ShortName:  req.ShortName,
	}

	classUpdate, err := cs.classRepository.Update(ctx, nil, data)
	if err != nil {
		return dto.ClassUpdateResponse{}, dto.ErrUpdateClass
	}

	return dto.ClassUpdateResponse{
		ID:         	classUpdate.ID.String(),
		Name:       	classUpdate.Name,
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