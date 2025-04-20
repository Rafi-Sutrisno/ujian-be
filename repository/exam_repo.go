package repository

import (
	"context"
	"fmt"
	"math"
	"mods/dto"
	"mods/entity"

	"gorm.io/gorm"
)

type (
	ExamRepository interface {
		CreateExam(ctx context.Context, tx *gorm.DB, exam entity.Exam) (entity.Exam, error)
		GetExamById(ctx context.Context, tx *gorm.DB, examId string) (entity.Exam, error)
		GetByClassID(ctx context.Context, tx *gorm.DB, classID string) ([]entity.Exam, error)
		GetAllExamWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllExamRepositoryResponse, error)
		UpdateExam(ctx context.Context, tx *gorm.DB, exam entity.Exam) (entity.Exam, error)
		DeleteExam(ctx context.Context, tx *gorm.DB, examId string) error
	}

	examRepository struct {
		db *gorm.DB
	}
)

func NewExamRepository(db *gorm.DB) ExamRepository {
	return &examRepository{
		db: db,
	}
}

func (er *examRepository) CreateExam(ctx context.Context,tx *gorm.DB, exam entity.Exam) (entity.Exam, error) {
	if tx == nil {
		tx = er.db
	}

	if err := tx.WithContext(ctx).Create(&exam).Error; err != nil {
		return entity.Exam{}, err
	}

	return exam, nil
}

func (ur *examRepository) GetExamById(ctx context.Context, tx *gorm.DB, examId string) (entity.Exam, error) {
	if tx == nil {
		tx = ur.db
	}

	fmt.Println("exam id di Repo:", examId)

	var exam entity.Exam
	if err := tx.WithContext(ctx).Where("id = ?", examId).Take(&exam).Error; err != nil {
		return entity.Exam{}, err
	}

	return exam, nil
}

func (ur *examRepository) GetByClassID(ctx context.Context, tx *gorm.DB, classID string) ([]entity.Exam, error) {
	if tx == nil {
		tx = ur.db
	}

	var exams []entity.Exam
	if err := tx.WithContext(ctx).Where("class_id = ?", classID).Find(&exams).Error; err != nil {
		return nil, err
	}

	return exams, nil
}

func (ur *examRepository) GetAllExamWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllExamRepositoryResponse, error) {
	if tx == nil {
		tx = ur.db
	}

	var exams []entity.Exam
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Exam{})

	if req.Search != "" {
		searchKeyword := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR short_name ILIKE ?", searchKeyword, searchKeyword)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllExamRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req.Page, req.PerPage)).Find(&exams).Error; err != nil {
		return dto.GetAllExamRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.GetAllExamRepositoryResponse{
		Exams:     exams,
		PaginationResponse: dto.PaginationResponse{
			Page: 		 req.Page,
			PerPage: 	 req.PerPage,
			Count: 		 count,
			MaxPage: 	 totalPage,
		},
	}, err
}

func (ur *examRepository) UpdateExam(ctx context.Context, tx *gorm.DB, exam entity.Exam) (entity.Exam, error) {
	if tx == nil {
		tx = ur.db
	}

	if err := tx.WithContext(ctx).Updates(&exam).Error; err != nil {
		return entity.Exam{}, err
	}

	return exam, nil
}

func (ur *examRepository) DeleteExam(ctx context.Context, tx *gorm.DB, examId string) error {
	if tx == nil {
		tx = ur.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.Exam{}, "id = ?", examId).Error; err != nil {
		return err
	}

	return nil
}
