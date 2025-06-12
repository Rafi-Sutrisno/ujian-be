package repository

import (
	"context"
	"math"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (

	examRepository struct {
		db *gorm.DB
	}
)

func NewExamRepository(db *gorm.DB) domain.ExamRepository {
	return &examRepository{
		db: db,
	}
}




func (er *examRepository) CreateExam(ctx context.Context, tx *gorm.DB, exam entity.Exam) (entity.Exam, error) {
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

	var exam entity.Exam
	if err := tx.WithContext(ctx).
		Preload("ExamLang").
		Preload("ExamLang.Language").
		Where("id = ?", examId).
		Take(&exam).Error; err != nil {
		return entity.Exam{}, err
	}

	return exam, nil
}


func (ur *examRepository) GetByClassID(ctx context.Context, tx *gorm.DB, classID string, check bool) ([]entity.Exam, error) {
	if tx == nil {
		tx = ur.db
	}

	var exams []entity.Exam

	query := tx.WithContext(ctx).Where("class_id = ?", classID)
	if check {
		query = query.Where("is_published = ?", true)
	}

	if err := query.Find(&exams).Error; err != nil {
		return nil, err
	}

	return exams, nil
}


func (ur *examRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string, check bool) ([]entity.Exam, error) {
	if tx == nil {
		tx = ur.db
	}

	var exams []entity.Exam

	query := tx.WithContext(ctx).
		Joins("JOIN classes ON classes.id = exams.class_id").
		Joins("JOIN user_classes ON user_classes.class_id = classes.id").
		Where("user_classes.user_id = ?", userID)

	if check {
		query = query.Where("exams.is_published = ?", true)
	}

	if err := query.Find(&exams).Error; err != nil {
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

	if err := tx.WithContext(ctx).Model(&entity.Exam{ID: exam.ID}).Updates(map[string]interface{}{
		"name":              exam.Name,
		"short_name":        exam.ShortName,
		"is_published":      exam.IsPublished,
		"start_time":        exam.StartTime,
		"duration":          exam.Duration,
		"end_time":          exam.EndTime,
		"is_seb_restricted": exam.IsSEBRestricted,
		"seb_browser_key":   exam.SEBBrowserKey,
		"seb_config_key":    exam.SEBConfigKey,
		"seb_quit_url":      exam.SEBQuitURL,
	}).Error; err != nil {
		return entity.Exam{}, err
	}


	return exam, nil
}

func (ur *examRepository) DeleteExam(ctx context.Context, tx *gorm.DB, examId string) error {
	if tx == nil {
		tx = ur.db
	}

	// Delete exam along with its related ExamLang entries
	exam := entity.Exam{ID: uuid.MustParse(examId)}
	if err := tx.WithContext(ctx).Select("ExamLang").Delete(&exam).Error; err != nil {
		return err
	}

	return nil
}

func (er *examRepository) IsUserInExamClass(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, error) {
	if tx == nil {
		tx = er.db
	}

	// First, get the class_id from the exam
	var exam entity.Exam
	if err := tx.WithContext(ctx).Select("class_id").Where("id = ?", examId).First(&exam).Error; err != nil {
		return false, err
	}

	// Now check if user exists in that class
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.UserClass{}).
		Where("user_id = ? AND class_id = ?", userId, exam.ClassID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (er *examRepository) IsUserInClass(ctx context.Context, tx *gorm.DB, userID, classID string) (bool, error) {
	if tx == nil {
		tx = er.db
	}

	var count int64
	err := tx.WithContext(ctx).Model(&entity.UserClass{}).
		Where("user_id = ? AND class_id = ?", userID, classID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}